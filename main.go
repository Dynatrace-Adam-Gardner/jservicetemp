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
    
	case keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName):
		log.Printf("Processing Evaluation.Finished Event")

		eventData := &keptnv2.EvaluationFinishedEventData{}
		parseKeptnCloudEventPayload(event, eventData)

		// Handle evaluation.finished and return any errors
    // See eventhandlers.go
		return HandleEvaluationFinishedEvent(myKeptn, event, eventData)
    
}
