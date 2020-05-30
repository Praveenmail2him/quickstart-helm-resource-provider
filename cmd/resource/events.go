package resource

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/aws"
)

const callbackDelaySeconds = 30

func errorEvent(model *Model, err error) handler.ProgressEvent {
	log.Printf("Returning ERROR...")
	return handler.ProgressEvent{
		OperationStatus: handler.Failed,
		Message:         err.Error(),
		ResourceModel:   model,
	}
}

func successEvent(model *Model) handler.ProgressEvent {
	log.Printf("Returning SUCCESS...")
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		ResourceModel:   model,
	}
}

func inProgressEvent(model *Model, stage Stage) handler.ProgressEvent {
	log.Printf("Returning IN_PROGRESS next stage %v...\n", stage)
	return handler.ProgressEvent{
		OperationStatus: handler.InProgress,
		ResourceModel:   model,
		Message:         fmt.Sprintf("%v in progress\n", stage),
		CallbackContext: map[string]interface{}{
			"Stage":     stage,
			"StartTime": os.Getenv("StartTime"),
			"Name":      aws.StringValue(model.Name),
		},
		CallbackDelaySeconds: callbackDelaySeconds,
	}
}

func makeEvent(model *Model, nextStage Stage, err error) handler.ProgressEvent {
	timeout := checkTimeOut(os.Getenv("StartTime"), model.TimeOut)
	if timeout {
		return errorEvent(nil, errors.New("resource creation timed out"))
	}
	if err != nil {
		return errorEvent(model, err)
	}
	if nextStage == CompleteStage {
		return successEvent(model)
	}
	return inProgressEvent(model, nextStage)
}
