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
	"reflect"
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

	type Auth struct {
		Wallet   string `json:"wallet"`
		Password string `json:"password"`
	}

	pass := Auth{
		Wallet:   "500957",
		Password: "Test_Test_42",
	}
	type TokenRes struct {
		Token     string `json:"token,omitempty"`
		ExpiresAt string `json:"expires_at,omitempty"`
	}
	type TokenResponse struct {
		Message string      `json:"message,omitempty"`
		Errors  interface{} `json:"errors,omitempty"`
		Trace   interface{} `json:"trace,omitempty"`
		Data    TokenRes    `json:"data,omitempty"`
	}

	//Convert User to byte using Json.Marshal
	//Ignoring error.
	fmt.Println(reflect.TypeOf(pass))

	//	body, _ := json.Marshal(pass)
	//body_auth, _ := json.Marshal(pass)

	//Pass new buffer for request with URL to post.
	//This will make a post request and will share the JSON data
	//resp, err := http.Post("https://reqres.in/api/users", "application/json", bytes.NewBuffer(body))
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	token := "986|laravel_sanctum_3Wj8d8RHaBCRySQjEEwyrIH1o9VPGNiRQUD0LQ3Nb1b1f97a"
	var bearer = "Bearer " + token
	var jsonStr = []byte(`{"msisdn":"265882997445", "amount:"100", "description":"narration", "invoice_number":"1242959"}`)
	req, err := http.NewRequest("POST", "https://dev.payouts.tnmmpamba.co.mw/api/invoices", bytes.NewBuffer(jsonStr))
	if err != nil {
		// handle error
	}

	req.Header.Set("Content-Type", "application/json")

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	fmt.Println("response Status:", res.Status)
	fmt.Println("response Headers:", res.Header)
	body, _ := io.ReadAll(res.Body)
	fmt.Println("response Body:", string(body))

	//resp, err := http.Post("https://dev.payouts.tnmmpamba.co.mw/api/authenticate", "application/json", bytes.NewBuffer(body))
	////req.Header.Add("Authorization", "Bearer ...")
	//// An error is returned if something goes wrong
	//if err != nil {
	//	panic(err)
	//}
	////Need to close the response stream, once response is read.
	////Hence defer close. It will automatically take care of it.
	//defer resp.Body.Close()
	//
	////Check response code, if New user is created then read response.
	//if resp.StatusCode == http.StatusOK {
	//	body, err := ioutil.ReadAll(resp.Body)
	//	if err != nil {
	//		//Failed to read response.
	//		panic(err)
	//	}
	//
	//	//Convert bytes to String and print
	//	jsonStr := string(body)
	//
	//	fmt.Println("Response: ", jsonStr)
	//
	//	//client := &http.Client{}
	//
	//	//var tokenResponse TokenResponse
	//	//
	//	//err = json.Unmarshal(body, &tokenResponse)
	//	//if err != nil {
	//	//	log.Error("Failed to unmarshal Response_Token: ", err.Error())
	//	//	return nil
	//	//}
	//
	//} else {
	//	//The status is not Created. print the error.
	//	fmt.Println("Get failed with error: ", resp.Status)
	//}

	log.Info("END PROCESS")
	os.Exit(2)

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
