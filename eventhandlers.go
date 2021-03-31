package main

import (
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2" // make sure to use v2 cloudevents here
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	jira "gopkg.in/andygrunwald/go-jira.v1"
	"io/ioutil"
	"log"
)

func HandleEvaluationTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.EvaluationTriggeredEventData) error {
	log.Println("Handling evaluation.triggered Event: %s", incomingEvent.Context.GetID())
	log.Println("[eventhandlers.go] HandleEvaluationTriggeredEvent Incoming Event Type: ", incomingEvent.Type())
	// Got sh.keptn.event.evaluation.triggered
	// Send a jira-service.started event

	//------------------------------------
	// 1. Send task started event
	//------------------------------------
	log.Println("Printing task started data object")
	log.Println(data)
	id, err := myKeptn.SendTaskStartedEvent(data, ServiceName)

	log.Println("[eventhandlers.go] Task Started ID: ", id)

	if err != nil {
		errMsg := fmt.Sprintf("Failed to send task started CloudEvent (%s), aborting...", err.Error())
		log.Println(errMsg)
		return err
	}

	return nil

}

func HandleEvaluationFinishedEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.EvaluationFinishedEventData) error {
	log.Println("[eventhandlers.go] Handling evaluation.finished Event: %s", incomingEvent.Context.GetID())

	log.Println("[eventhandlers.go] HandleEvaluationFinishedEvent Incoming Event Type: ", incomingEvent.Type())

	if incomingEvent.Source() == ServiceName {
		// skip evaluation.finished, it has been sent out by jira-service
		log.Println("[eventhandlers.go] Received an evaluation.finished event from jira-service so stop processing so as to not get into a recursion.")
		return nil
	}

	//------------------------------------
	// 2. Do work
	//------------------------------------

	//------------------------------------
	// 3. Send task finished event
	//------------------------------------
	outputData := &keptnv2.EventData{
		Status:  keptnv2.StatusSucceeded,
		Result:  keptnv2.ResultPass,
		Message: "jira-service finished....",
	}

	log.Println("Printing task finished data object")
	log.Println(outputData)
	log.Println("-----------------------------------")

	id, err := myKeptn.SendTaskFinishedEvent(outputData, ServiceName)

	log.Println("[eventhandlers.go] Task Finished ID: ", id)

	if err != nil {
		errMsg := fmt.Sprintf("Failed to send task finished CloudEvent (%s), aborting...", err.Error())
		log.Println(errMsg)
		return err
	}

	return nil
}

func createJIRATicket(summary string, description string) string {

	log.Println("[eventhandlers.go] Creating JIRA Ticket Now...")

	tp := jira.BasicAuthTransport{
		Username: JIRA_DETAILS.Username,
		Password: JIRA_DETAILS.APIToken,
	}

	jiraClient, err := jira.NewClient(tp.Client(), JIRA_DETAILS.BaseURL)
	if err != nil {
		panic(err)
	}

	i := jira.Issue{
		Fields: &jira.IssueFields{
			Assignee: &jira.User{
				AccountID: JIRA_DETAILS.AssigneeID,
			},
			Reporter: &jira.User{
				AccountID: JIRA_DETAILS.ReporterID,
			},
			Description: description,
			Type: jira.IssueType{
				Name: JIRA_DETAILS.IssueType,
			},
			Project: jira.Project{
				Key: JIRA_DETAILS.ProjectKey,
			},
			Summary: summary,
		},
	}

	// Create ticket
	issue, response, err := jiraClient.Issue.Create(&i)

	if err != nil {
		bodyBytes, err2 := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println(err2)
		}
		bodyString := string(bodyBytes)
		log.Println(bodyString)
		log.Println(err)
	} else {
		log.Println("[eventhandlers.go] Created ticket successfully: ", issue.Key)
	}

	return issue.Key

}
