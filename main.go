package main

/* For backwards compatibility, we need to cater for the following legacy events:
*  - "sh.keptn.events.evaluation-done"
*  - "sh.keptn.event.problem.open"
*  - sh.keptn.events.problem"
*
*  And the following new event types:
* - "sh.keptn.event.evaluation.finished"
* - "sh.keptn.event.jira-service.triggered"
 */

import (
	"context"
	"errors"
	"log"
	"os"

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

	log.Printf("[main.go] 18:59 build")
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
