package main

 import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	cloudevents "github.com/cloudevents/sdk-go/v2" // make sure to use v2 cloudevents here
	"github.com/kelseyhightower/envconfig"

	keptn "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

var keptnOptions = keptn.KeptnOpts{}

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int `envconfig:"RCV_PORT" default:"8080"`
	// Path to which cloudevents are sent
	Path string `envconfig:"RCV_PATH" default:"/"`
	// Whether we are running locally (e.g., for testing) or on production
	Env string `envconfig:"ENV" default:"local"`
	// URL of the Keptn configuration service (this is where we can fetch files from the config repo)
	ConfigurationServiceUrl string `envconfig:"CONFIGURATION_SERVICE" default:""`
}

type JiraDetails struct {
	BaseURL              string
	Username             string
	APIToken             string
	AssigneeID           string
	ReporterID           string
	ProjectKey           string
	IssueType            string
	TicketForProblems    bool
	TicketForEvaluations bool
}

type KeptnDetails struct {
	KeptnDomain    string
	KeptnBridgeURL string
}

var JIRA_DETAILS JiraDetails
var KEPTN_DETAILS KeptnDetails

// ServiceName specifies the current services name (e.g., used as source when sending CloudEvents)
const ServiceName = "jira-service"

// This method gets called when a new event is received from the Keptn Event Distributor
func processKeptnCloudEvent(ctx context.Context, event cloudevents.Event) error {

	// create keptn handler
	log.Printf("[main.go] Initializing Keptn Handler")
	myKeptn, err := keptnv2.NewKeptn(&event, keptnOptions)
	if err != nil {
		return errors.New("Could not create Keptn Handler: " + err.Error())
	}

	setupAndDebug(myKeptn, event)

	switch event.Type() {

	case keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName): // sh.keptn.event.evaluation.triggered
		log.Printf("Processing Evaluation.Triggered Event")

		eventData := &keptnv2.EvaluationTriggeredEventData{}
		parseKeptnCloudEventPayload(event, eventData)

		return HandleEvaluationTriggeredEvent(myKeptn, event, eventData)

	case keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName): // sk.keptn.event.evaluation.finished
		log.Printf("Processing Evaluation.Finished Event")

		eventData := &keptnv2.EvaluationFinishedEventData{}
		parseKeptnCloudEventPayload(event, eventData)

		// Handle evaluation.finished and return any errors
		return HandleEvaluationFinishedEvent(myKeptn, event, eventData)
	}

	return nil

}

/**
 * Usage: ./main
 * no args: starts listening for cloudnative events on localhost:port/path
 *
 * Environment Variables
 * env=runlocal   -> will fetch resources from local drive instead of configuration service
 */
func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("[main.go] Failed to process env var: %s", err)
	}

	os.Exit(_main(os.Args[1:], env))
}

/**
 * Opens up a listener on localhost:port/path and passes incoming requets to gotEvent
 */
func _main(args []string, env envConfig) int {
	// configure keptn options
	if env.Env == "local" {
		log.Println("[main.go] env=local: Running with local filesystem to fetch resources")
		keptnOptions.UseLocalFileSystem = true
	}

	keptnOptions.ConfigurationServiceURL = env.ConfigurationServiceUrl

	log.Printf("[main.go] 19:34 build")
	log.Printf("[main.go] Starting %s...", ServiceName)
	log.Printf("[main.go]     on Port = %d; Path=%s", env.Port, env.Path)

	ctx := context.Background()
	ctx = cloudevents.WithEncodingStructured(ctx)

	log.Printf("[main.go] Creating new http handler")

	// configure http server to receive cloudevents
	p, err := cloudevents.NewHTTP(cloudevents.WithPath(env.Path), cloudevents.WithPort(env.Port))

	if err != nil {
		log.Fatalf("[main.go] failed to create client, %v", err)
	}
	c, err := cloudevents.NewClient(p)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	log.Printf("[main.go] Starting receiver")
	log.Fatal(c.StartReceiver(ctx, processKeptnCloudEvent))

	return 0
}

/**
 * Parses a Keptn Cloud Event payload (data attribute)
 */
func parseKeptnCloudEventPayload(event cloudevents.Event, data interface{}) error {
	err := event.DataAs(data)
	if err != nil {
		log.Fatalf("Got Data Error: %s", err.Error())
		return err
	}
	return nil
}

func setJIRADetails() {
	JIRA_DETAILS.BaseURL = os.Getenv("JIRA_BASE_URL")
	JIRA_DETAILS.Username = os.Getenv("JIRA_USERNAME")
	JIRA_DETAILS.AssigneeID = os.Getenv("JIRA_ASSIGNEE_ID")
	JIRA_DETAILS.ReporterID = os.Getenv("JIRA_REPORTER_ID")
	JIRA_DETAILS.APIToken = os.Getenv("JIRA_API_TOKEN")
	JIRA_DETAILS.ProjectKey = os.Getenv("JIRA_PROJECT_KEY")
	JIRA_DETAILS.IssueType = os.Getenv("JIRA_ISSUE_TYPE")
	JIRA_DETAILS.TicketForProblems, _ = strconv.ParseBool(os.Getenv("JIRA_TICKET_FOR_PROBLEMS"))
	JIRA_DETAILS.TicketForEvaluations, _ = strconv.ParseBool(os.Getenv("JIRA_TICKET_FOR_EVALUATIONS"))
}

func setKeptnDetails() {
	KEPTN_DETAILS.KeptnDomain = os.Getenv("KEPTN_DOMAIN")

	// If Bridge URL isn't set in YAML file, default to the KEPTN_DOMAIN which is mandatory
	if os.Getenv("KEPTN_BRIDGE_URL") == "" {
		KEPTN_DETAILS.KeptnBridgeURL = os.Getenv("KEPTN_DOMAIN")
	} else {
		KEPTN_DETAILS.KeptnBridgeURL = os.Getenv("KEPTN_BRIDGE_URL")
	}
}

func setupAndDebug(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event) {
	log.Printf("[main.go] gotEvent(%s): %s - %s", incomingEvent.Type(), myKeptn.KeptnContext, incomingEvent.Context.GetID())

	// Get Debug Mode
	// This is set in the service.yaml as DEBUG "true"
	DEBUG, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	log.Printf("[main.go] Debug Mode: %v \n", DEBUG)

	// Set JIRA Details
	setJIRADetails()

	// Get Dynatrace Tenant
	dynaTraceTenant := os.Getenv("DT_TENANT")

	// KEPTN_DOMAIN must be set but KEPTN_BRIDGE_URL is optional in jira-service deployment.yaml file
	setKeptnDetails()

	if JIRA_DETAILS.BaseURL == "" ||
		JIRA_DETAILS.Username == "" ||
		JIRA_DETAILS.APIToken == "" ||
		JIRA_DETAILS.ProjectKey == "" ||
		JIRA_DETAILS.IssueType == "" ||
		KEPTN_DETAILS.KeptnDomain == "" {
		fmt.Println("[main.go] Missing mandatory input parameters JIRA_BASE_URL and / or JIRA_USERNAME and / or JIRA_API_TOKEN and / or JIRA_PROJECT_KEY and / or JIRA_ISSUE_TYPE and / or KEPTN_DOMAIN.")
	}

	if DEBUG {
		fmt.Println("[main.go] --- Printing JIRA Input Details ---")
		fmt.Printf("[main.go] Base URL: %s \n", JIRA_DETAILS.BaseURL)
		fmt.Printf("[main.go] Username: %s \n", JIRA_DETAILS.Username)
		fmt.Printf("[main.go] Assignee ID: %s \n", JIRA_DETAILS.AssigneeID)
		fmt.Printf("[main.go] Reporter ID: %s \n", JIRA_DETAILS.ReporterID)
		fmt.Printf("[main.go] API Token: %s \n", JIRA_DETAILS.APIToken)
		fmt.Printf("[main.go] Project Key: %s \n", JIRA_DETAILS.ProjectKey)
		fmt.Printf("[main.go] Issue Type: %s \n", JIRA_DETAILS.IssueType)
		fmt.Printf("[main.go] Ticket For Problems: %v \n", JIRA_DETAILS.TicketForProblems)
		fmt.Printf("[main.go] Ticket For Problems: %v \n", JIRA_DETAILS.TicketForEvaluations)
		fmt.Println("[main.go] --- End Printing JIRA Input Details ---")

		fmt.Printf("[main.go] Dynatrace Tenant: %s \n", dynaTraceTenant)
		fmt.Printf("[main.go] Keptn Domain: %s \n", KEPTN_DETAILS.KeptnDomain)
		fmt.Printf("[main.go] Keptn Bridge URL: %s \n", KEPTN_DETAILS.KeptnBridgeURL)

	    // At this point, we have all mandatory input params. Proceed
		fmt.Println("[main.go] Got all input variables. Proceeding...")

		if JIRA_DETAILS.TicketForProblems {
		  fmt.Println("[main.go] Will create tickets for problems")
	    } else {
		  fmt.Println("[main.go] Will NOT create tickets for problems")
	    }

		if JIRA_DETAILS.TicketForEvaluations {
		  fmt.Println("[main.go] Will create tickets for evaluations")
	    } else {
		  fmt.Println("[main.go] Will NOT create tickets for evaluations")
	    }
    }
}
