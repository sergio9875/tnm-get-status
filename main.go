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
	"io/ioutil"
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

	//url := "https://jsonplaceholder.typicode.com/users/1"
	//url := "https://dev.payouts.tnmmpamba.co.mw/api/payments/NT543Y5YT5R432"

	//If the struct variable names does not match with json attributes
	//then you can define the json attributes actual name after json:attname as shown below.
	type User struct {
		Name string `json:"name"`
		Job  string `json:"job"`
	}
	type Auth struct {
		Wallet   string `json:"wallet"`
		Password string `json:"password"`
	}

	//Create user struct which need to post.
	//user := User{
	//	Name: "Test User1",
	//	Job:  "Go lang Developer",
	//}
	//Create user struct which need to post.
	pass := Auth{
		Wallet:   "500957",
		Password: "Test_Test_42",
	}

	//Convert User to byte using Json.Marshal
	//Ignoring error.
	body, _ := json.Marshal(pass)
	//body_auth, _ := json.Marshal(pass)

	//Pass new buffer for request with URL to post.
	//This will make a post request and will share the JSON data
	//resp, err := http.Post("https://reqres.in/api/users", "application/json", bytes.NewBuffer(body))
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Post("https://dev.payouts.tnmmpamba.co.mw/api/authenticate", "application/json", bytes.NewBuffer(body))

	// An error is returned if something goes wrong
	if err != nil {
		panic(err)
	}
	//Need to close the response stream, once response is read.
	//Hence defer close. It will automatically take care of it.
	defer resp.Body.Close()

	//Check response code, if New user is created then read response.
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			//Failed to read response.
			panic(err)
		}

		//Convert bytes to String and print
		jsonStr := string(body)

		fmt.Println("Response: ", jsonStr)

	} else {
		//The status is not Created. print the error.
		fmt.Println("Get failed with error: ", resp.Status)
	}

	//body1 := []byte(`{
	//"wallet": "500957",
	//"password" : "Test_Test_42"
	//}`)
	//
	//url := "https://dev.payouts.tnmmpamba.co.mw/api/authenticate"
	//req, err := http.NewRequest("POST", url, bytes.NewBuffer(body1))
	//req.Header.Set("X-Custom-Header", "myvalue")
	//req.Header.Set("Content-Type", "application/json")
	//
	//tr := &http.Transport{
	//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//}
	//
	//client := &http.Client{Transport: tr}
	//resp, err := client.Do(req)
	//if err != nil {
	//	panic(err)
	//}
	//defer resp.Body.Close()
	//
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	//body, _ := io.ReadAll(resp.Body)
	//log.Println("response Body:", body)
	//fmt.Println("response Body:", string(body))

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
