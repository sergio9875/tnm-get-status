package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	log "malawi-getstatus/logger"
	"malawi-getstatus/process"
	"malawi-getstatus/utils"
	"net/http"
	"os"
	"strconv"
	"time"
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
func callClient(jsonStr string) error {
	//start := time.Now()
	var err error
	var bearer = "Bearer " + jsonStr

	url := "https://dev.payouts.tnmmpamba.co.mw/api/invoices/1350604"

	request, err := http.NewRequest(http.MethodGet, url, nil)

	request.Header.Add("Authorization", bearer)
	request.Header.Set("Content-Type", "application/json") // => your content-type
	client := http.Client{Timeout: 5 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	log.Println("response from get:", string(body))
	log.Println(string(body))
	return nil
}

func callPost(token string) error {
	log.Println("TOKEN!!!!!!!!!!!!", token)
	// JSON body
	//	body := []byte(`{
	//		"msisdn": "265882997445",
	//		"amount": 100,
	//		"description":"some narration",
	//"invoice_number": "1252005"
	//	}`)

	url := "https://payouts.tnmmpamba.co.mw/api/invoices/refund/AJ950B60NF"
	// create request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		panic(err)
	}

	// set headers
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/json")

	client := http.Client{Timeout: 5 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	log.Println("response from get:", string(body))
	log.Println(string(body))
	return nil
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

	log.Info("trying to retrieve access token")
	marshalled, err := json.Marshal(pass)
	if err != nil {
		log.Fatalf("impossible to marshall token config: %s", err.Error())
	}
	client := &http.Client{}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest(http.MethodPost, "https://dev.payouts.tnmmpamba.co.mw/api/authenticate", bytes.NewReader(marshalled))
	if err != nil {
		log.Fatalf("impossible to build request: %s", err.Error())
	}
	// add headers
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("impossible to send request: %s", err.Error())
	}
	log.Printf("status Code: %d", strconv.Itoa(res.StatusCode))

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("impossible to read all body of response: %s", err.Error())
	}
	log.Printf("res body token: %s", string(resBody))

	var tokenResponse *TokenResponse
	err = json.Unmarshal(resBody, &tokenResponse)
	if err != nil {
		log.Error("Failed to unmarshal Response_Token: ", err.Error())
		return nil
	}
	log.Printf("TOKEN____RES***: %s", tokenResponse.Data.Token)

	//resultApi := callClient(tokenResponse.Data.Token)
	resultApi := callPost(tokenResponse.Data.Token)
	log.Println("reultapi", resultApi)

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
