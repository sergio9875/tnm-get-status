package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	log "malawi-getstatus/logger"
	"malawi-getstatus/process"
	"malawi-getstatus/utils"
	"os"
)

var invokeCount = 0
var controller *process.Controller

func Init() {
	controller = process.NewController(os.Getenv("SECRET_NAME"))
	invokeCount = 0
}

func init() {
	// used to init anything special
}

func LambdaHandler(ctx context.Context, sqsEvent events.SQSEvent) error {
	log.Debug("ROOT", "version: <GIT_HASH>")
	//stdout and stderr are sent to AWS CloudWatch Logs
	if invokeCount == 0 {
		Init()
	}

	invokeCount = invokeCount + 1
	if invokeCount > utils.SafeAtoi(utils.Getenv("MAX_INVOKE", "15"), 15) {
		// reset global variables to nil
		controller.ShutDown()
		Init()
	}

	for _, record := range sqsEvent.Records {
		controller.PreProcess(utils.StringPtr(uuid.New().String()))
		err := controller.Process(ctx, record)
		if err != nil {
			log.Fatalf("Lambda process failed %s", err.Error())
			return err
		}
		controller.PostProcess()
	}
	return nil
}

func main() {
	lambda.Start(LambdaHandler)
}
