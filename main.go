package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
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
	type Post struct {
		Msisdn        string `json:"msisdn"`
		Amount        int    `json:"amount"`
		Description   string `json:"description"`
		InvoiceNumber string `json:"invoice_number"`
	}

	pass := Auth{
		Wallet:   "500957",
		Password: "Test_Test_42",
	}
	post := Post{
		Msisdn:        "265882997445",
		Amount:        280,
		InvoiceNumber: "1252002",
		Description:   "Test1123",
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

	body, _ := json.Marshal(post)

	//Pass new buffer for request with URL to post.
	//This will make a post request and will share the JSON data
	//resp, err := http.Post("https://reqres.in/api/users", "application/json", bytes.NewBuffer(body))
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	token := "992|laravel_sanctum_WJZXQACwJah8W2HA3AuyHadq8Bx10GLFWO9Ma9zK43900d2c"
	bearer := "Bearer " + token
	//http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	//r, err := http.NewRequest("POST", "https://dev.payouts.tnmmpamba.co.mw/api/invoices/", bytes.NewBuffer(body))
	//if err != nil {
	//	panic(err)
	//}

	resp, err := http.Post("https://dev.payouts.tnmmpamba.co.mw/api/invoices/", "application/json",
		bytes.NewBuffer(body))

	// add authorization header to the req
	//resp.Header.Add("Authorization", bearer)
	resp.Header.Set("Authorization", bearer)
	if err != nil {
		log.Fatalf("impossible to build request: %s", err.Error())
	}

	defer resp.Body.Close()

	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)

	fmt.Println("json********************************")
	fmt.Println(res["json"])
	fmt.Println("json********************************")

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	bodyRes, _ := io.ReadAll(resp.Body)
	fmt.Println("response Body:", string(bodyRes))

	//req, err := http.NewRequest(http.MethodPost, postURL, bytes.NewBuffer(reqBody))
	//if err != nil {
	//	log.Fatalf("impossible to build request: %s", err)
	//}
	// add headers
	//r.Header.Add("Content-Type", "application/json")

	if err != nil {
		log.Error("failed to create a new request", err.Error())
		return err
	}

	//client := &http.Client{}
	//res, err := client.Do(r)
	//if err != nil {
	//	panic(err)
	//}
	//log.Printf("status Code: %d", res.StatusCode)
	//defer res.Body.Close()
	//
	//if err != nil {
	//	log.Error("Failed to read response body", err.Error())
	//	return nil, err
	//}

	//var responseBody = new(models.ChargeResponse)
	//err = json.Unmarshal(body, &responseBody)
	//if err != nil {
	//	log.Error("Failed to unmarshal response: ", err.Error())
	//	return nil, err
	//}
	//
	//body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))

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
