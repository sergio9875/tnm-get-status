package services

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	log "malawi-getstatus/logger"
	"malawi-getstatus/models"
	"net/http"
	"strconv"
	"time"
)

func GetToken(PostURL string, wallet string, password string) (*models.TokenResponse, error) {
	log.Info("trying to retrieve access token")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{
		Transport: tr,
		Timeout:   40 * time.Second,
	}

	body := models.Auth{
		Wallet:   wallet,
		Password: password,
	}
	marshalled, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("impossible to marshall token config: %s", err.Error())
	}

	//url3 := "http://localhost:8888/token"
	req, err := http.NewRequest(http.MethodPost, PostURL, bytes.NewReader(marshalled))
	if err != nil {
		log.Fatalf("impossible to build request: %s", err.Error())
	}

	// Add Headers
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

	var tokenResponse *models.TokenResponse
	err = json.Unmarshal(resBody, &tokenResponse)
	if err != nil {
		log.Error("Failed to unmarshal Response_Token: ", err.Error())

	}
	if res.StatusCode != http.StatusOK {
		return nil,
			errors.New(fmt.Sprintf("statusCode: %v and errMsg: %v", res.Status, tokenResponse.Errors))
	}
	return tokenResponse, nil
}
