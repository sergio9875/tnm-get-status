package process

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

//func (c *Controller) SendGetStatus(ctx context.Context, request *models.IncomingRequest) (*models.TnmResponse, error) {
//	// first need to retrieve access token
//	fmt.Println("send request to get-status", request)
//	token, err := c.SendTokenRequest(request.URLToken, request.Wallet, request.Password)
//	if err != nil {
//		log.Error("Failed to Get Token Request: ", error.Error)
//		return nil, nil
//	}
//	fmt.Println("BEFORE GET REQUEST", token)
//	//os.Exit(2)
//	//return c.SendGetRequest(request.MbtId, token, request.URLQuery)
//	return nil, nil
//}

//func (c *Controller) SendTokenRequest(PostURL string, wallet string, password string) (token string, responseError *models.ResponseError) {
//	log.Info("trying to retrieve access token", PostURL)
//	tr := &http.Transport{
//		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
//	}
//	client := http.Client{
//		Transport: tr,
//		Timeout:   40 * time.Second,
//	}
//	// Prepare request body
//	body := models.Auth{
//		Wallet:   wallet,
//		Password: password,
//	}
//	marshalled, err := json.Marshal(body)
//	if err != nil {
//		log.Fatalf("impossible to marshall token config: %s", err.Error())
//	}
//
//	url3 := "http://localhost:8888/tokenError"
//	req, err := http.NewRequest(http.MethodPost, url3, bytes.NewReader(marshalled))
//	if err != nil {
//		log.Fatalf("impossible to build request: %s", err.Error())
//	}
//	// add headers
//	req.Header.Set("Content-Type", "application/json")
//	res, err := client.Do(req)
//	if err != nil {
//		log.Fatalf("impossible to send request: %s", err.Error())
//	}
//	log.Printf("status Code: %d", strconv.Itoa(res.StatusCode))
//
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//
//		}
//	}(res.Body)
//	resBody, err := io.ReadAll(res.Body)
//	if err != nil {
//		log.Fatalf("impossible to read all body of response: %s", err.Error())
//	}
//	log.Printf("res body: %s", string(resBody))
//
//	if res.StatusCode != http.StatusOK {
//		fmt.Println("error with a token")
//		var errorResponse *models.ErrorMalawiMessage
//		err = json.Unmarshal(resBody, &errorResponse)
//		//
//		//	//return  "nil", nil
//		//
//		return "", &models.ResponseError{
//			StatusCode: res.StatusCode,
//			Err:        errors.New(errorResponse.Message),
//		}
//		//
//		//	//return &models.ApiResult{
//		//	//	Message:    errorResponse.Message,
//		//	//	StatusCode: res.StatusCode,
//		//	//	Response:   string(resBody),
//		//	//}, nil
//		//
//	}
//	fmt.Println("88**********************************************************************88")
//	//os.Exit(2)
//	var tokenResponse *models.TokenResponse
//	err = json.Unmarshal(resBody, &tokenResponse)
//	if err != nil {
//		log.Error("Failed to unmarshal Response_Token: ", err.Error())
//		return tokenResponse.Data.Token, nil
//	}
//	fmt.Printf("tokenResponse: %s\n", tokenResponse)
//	//return &models.ApiResult{
//	//	Message:    tokenResponse.Message,
//	//	StatusCode: res.StatusCode,
//	//	Response:   string(resBody),
//	//	Token:      tokenResponse.Data.Token,
//	//}, nil
//	return tokenResponse.Data.Token, nil
//}

func (c *Controller) SendGetRequest(transactionId string, token string, urlGetStatus string) (*models.TnmResponse, error) {

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

	url3 := "http://localhost:8888/chargeError"
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{}
	req, err := http.NewRequest("GET", url3, nil)
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
			errors.New(fmt.Sprintf("message: %v\n", responseBody.Message) +
				fmt.Sprintf("statusCode: %v\n", res.StatusCode) +
				fmt.Sprintf("errors: %v\n", responseBody.Errors))
		//errors.New(fmt.Sprintf("message: %v\n,statusCode: %v\n, errors: %v\n", responseBody.Message, res.StatusCode, responseBody.Errors))
	}
	log.Info("RSP: Lambda <--- TNM MALAWI: ", responseBody)

	return responseBody, nil
}
