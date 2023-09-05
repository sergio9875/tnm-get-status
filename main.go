package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/uuid"
	"io"
	log "malawi-getstatus/logger"
	"malawi-getstatus/process"
	"malawi-getstatus/utils"
	"net/http"
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

// LambdaHandler - Listen to S3 events and start processing
func LambdaHandler(ctx context.Context, sqsEvent events.SQSEvent) error {

	log.Info("START PROCESS")

	// JSON body
	body1 := []byte(`{
    "wallet": "500957",
    "password" : "Test_Test_42"
	}`)

	url := "https://dev.payouts.tnmmpamba.co.mw/api/authenticate"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body1))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := io.ReadAll(resp.Body)
	log.Println("response Body:", body)
	fmt.Println("response Body:", string(body))

	log.Info("END PROCESS")

	log.Debug("ROOT", "version: <GIT_HASH>")

	if invokeCount == 0 {
		Init()
	}

	invokeCount = invokeCount + 1
	if invokeCount > *utils.SafeAtoi(utils.Getenv("MAX_INVOKE", "15"), aws.Int(15)) {
		// reset global variables to nil
		controller.ShutDown()
		Init()
	}

	for _, record := range sqsEvent.Records {
		controller.PreProcess(utils.StringPtr(uuid.New().String()))
		if err := controller.Process(ctx, record); err != nil {
			return err
		}
		controller.PostProcess()
	}

	return nil
}

func main() {
	lambda.Start(LambdaHandler)
}
