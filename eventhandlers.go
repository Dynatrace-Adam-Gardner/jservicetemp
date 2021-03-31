package main

import (
	cloudevents "github.com/cloudevents/sdk-go/v2" // make sure to use v2 cloudevents here
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"log"
	"fmt"
)

func HandleEvaluationFinishedEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.EvaluationFinishedEventData) error {
	log.Printf("Handling evaluation.finished Event: %s", incomingEvent.Context.GetID())

	//------------------------------------
	// 1. Send task started event
	//------------------------------------
	_, err := myKeptn.SendTaskStartedEvent(data, ServiceName)

	if err != nil {
		errMsg := fmt.Sprintf("Failed to send task started CloudEvent (%s), aborting...", err.Error())
		log.Println(errMsg)
		return err
	}

	//------------------------------------
	// 2. Do work
	//------------------------------------

	//------------------------------------
	// 3. Send task finished event
	//------------------------------------
	_, err = myKeptn.SendTaskFinishedEvent(data, ServiceName)

	if err != nil {
		errMsg := fmt.Sprintf("Failed to send task finished CloudEvent (%s), aborting...", err.Error())
		log.Println(errMsg)
		return err
	}

	return nil
}

func createJIRATicket(summary string, description string) string {

	log.Println("[eventhandlers.go] Creating JIRA Ticket Now...")

	return ""

}
