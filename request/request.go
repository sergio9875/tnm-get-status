package request

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type IRequest interface {
	Get(url string, headers map[string]string) ([]byte, error)
	GetWithJsonResponse(url string, headers map[string]string, responseData interface{}) error
	Post(url string, headers map[string]string, requestData interface{}) ([]byte, error)
	PostWithJsonResponse(url string, headers map[string]string, requestData interface{}, responseData interface{}) error
}

type Request struct {
	Client *http.Client
}

func NewClient() (IRequest, error) {
	httpClient := &http.Client{}

	return &Request{Client: httpClient}, nil
}

// Get return []byte from the get request
func (r *Request) Get(url string, headers map[string]string) ([]byte, error) {
	var body []byte
	var err error

	request, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	if len(headers) > 0 {
		for k, v := range headers {
			request.Header.Set(k, v)
		}
	}

	resp, errR := r.Client.Do(request)
	if errR != nil {
		return nil, errR
	}
	defer resp.Body.Close()

	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	}

	return body, nil
}

// Post send a post request with a timeout of 10 sec per default
func (r *Request) Post(url string, headers map[string]string, requestData interface{}) ([]byte, error) {

	if r.Client.Timeout == 0 {
		r.Client.Timeout = 20 * time.Second
	}

	var body []byte
	var err error

	if body, err = json.Marshal(requestData); err != nil {
		return nil, err
	}

	request, errRequest := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))

	if errRequest != nil {
		return nil, errRequest
	}

	if len(headers) > 0 {
		for k, v := range headers {
			request.Header.Set(k, v)
		}
	}

	resp, errR := r.Client.Do(request)
	if errR != nil {
		return nil, errR
	}
	defer resp.Body.Close()

	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	}

	return body, nil
}

func (r *Request) PostWithJsonResponse(url string,
	headers map[string]string,
	requestData interface{},
	responseData interface{}) error {

	var body []byte
	var err error

	if body, err = r.Post(url, headers, requestData); err != nil {
		return err
	}

	if err = json.Unmarshal(body, &responseData); err != nil {
		return err
	}

	return nil
}

func (r *Request) GetWithJsonResponse(url string,
	headers map[string]string,
	responseData interface{}) error {
	var body []byte
	var err error

	if body, err = r.Get(url, headers); err != nil {
		return err
	}

	if err = json.Unmarshal(body, &responseData); err != nil {
		return err
	}

	return nil
}
