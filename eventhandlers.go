package main

...

func HandleEvaluationFinishedEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.EvaluationFinishedEventData) error {
	log.Printf("Handling evaluation.finished Event: %s", incomingEvent.Context.GetID())

	//------------------------------------
	// 1. Send task started event
	//------------------------------------
	myKeptn.SendTaskStartedEvent(data, ServiceName)
  
  //------------------------------------
  // 2. Do work (this part works)
  //------------------------------------
  
  //------------------------------------
	// 3. Send task finished event
	//------------------------------------
	myKeptn.SendTaskFinishedEvent(data, ServiceName)
}
