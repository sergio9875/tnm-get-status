package services

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	log "malawi-getstatus/logger"
	"malawi-getstatus/models"
	"net/http"
	"net/url"
	"strconv"
)

func SendGetRequest(transactionId string, token string, urlGetStatus string) (*models.TnmBodyResponse, error) {

	log.Info("trying to get query-status, transID)", transactionId)
	var bearer = "Bearer " + token
	base, err := url.Parse(urlGetStatus)
	if err != nil {
		log.Fatalf(": %s", err.Error())
	}
	//transactionId = "1350868"
	// Path params
	base.Path += transactionId

	// Query params
	params := url.Values{}
	base.RawQuery = params.Encode()

	fmt.Printf("Encoded URL is %q\n", base.String())

	//url3 := "http://localhost:8888/chargeError"
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{}
	req, err := http.NewRequest("GET", base.String(), nil)
	if err != nil {
		// handle error
	}
	req.Header.Set("Content-Type", "application/json")
	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("impossible to send request: %s", err.Error())
	}
	log.Printf("status Code: %d", strconv.Itoa(res.StatusCode))

	rBody, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("response body from Get-Status #########")
	log.Infof("%s", string(rBody))

	var responseBody = new(models.TnmResponse)
	err = json.Unmarshal(rBody, &responseBody)
	if err != nil {
		log.Error("Failed to unmarshal response: ", err.Error())
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil,
			errors.New(fmt.Sprintf("statusCode: %v and errMsg: %v", res.Status, responseBody.Errors))
	}
	log.Info("RSP: Lambda <--- TNM MALAWI: ", responseBody.Data)

	return &responseBody.Data, nil
}
