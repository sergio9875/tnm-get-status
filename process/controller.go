package process

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"malawi-getstatus/enums"
	log "malawi-getstatus/logger"
	"malawi-getstatus/models"
	repo "malawi-getstatus/repository"
	"malawi-getstatus/repository/mssql"
	"malawi-getstatus/request"
	"malawi-getstatus/utils"
)

// Controller container
type Controller struct {
	secretHolder *SecretIDHolder
	sqsProducer  *SQSProducer
	sumoProducer *SQSProducer
	repository   *repo.Repository
	config       *models.SecretModel
	requestId    *string
	httpClient   *request.IRequest
}

func (c *Controller) initClient() {
	httpClient, err := request.NewClient()

	if err != nil {
		log.Fatalf(*c.requestId, "Lambda init failed on http client service: %v", err)
	}

	c.httpClient = &httpClient
}

func NewController(secret string) *Controller {
	controller := Controller{
		requestId: utils.StringPtr("ROOT"),
	}
	controller.initSecret(secret)
	controller.initRepository()
	controller.initClient()
	controller.initSumoProducer()
	return &controller
}

func (c *Controller) initSecret(secret string) {
	c.secretHolder = &SecretIDHolder{
		SecretID: secret,
		Client:   CreateSMClient(),
	}
	c.config = c.secretHolder.LoadSecret()

}

func (c *Controller) initSumoProducer() {
	var err error
	c.sumoProducer, err = NewSQSProducerFromUrl(context.TODO(), CreateSQSClient(),
		&c.config.Services.SumoPusherUrl)
	if err != nil {
		log.Fatalf(*c.requestId, "Lambda init failed on sqs producer: %v", err)
	}
}

func (c *Controller) initSqsProducer(queueName string) {
	var err error
	c.sqsProducer, err = NewSQSProducerFromUrl(context.TODO(), CreateSQSClient(),
		&queueName)
	if err != nil {
		log.Fatalf(*c.requestId, "Lambda init failed on sqs producer: %v", err)
	}
}

func (c *Controller) initRepository() {
	localRepo, err := mssql.NewRepository(c.config.Database.Africainv)
	if err != nil {
		log.Fatalf(*c.requestId, "Lambda init failed on repository: %v", err)
	}
	c.repository = &localRepo
}

func (c *Controller) ShutDown() {
	c.config = nil
	c.sqsProducer = nil
	c.secretHolder = nil
	if c.repository != nil {
		err := (*c.repository).Close()
		if err != nil {
			return
		}
	}
	c.repository = nil
}

func (c *Controller) PreProcess(pid *string) {
	c.requestId = pid
}

func (c *Controller) PostProcess() {
	c.requestId = utils.StringPtr("ROOT")
}

func (c *Controller) Process(ctx context.Context, message events.SQSMessage) error {
	c.sendSumoMessages(ctx, "start tnm-malawi get callback process", message)

	var err error

	msgBody := new(models.IncomingRequest)
	mtnResponseBody := new(models.ResponseBody)
	mtnResponse := new(models.Response)

	if err = c.getMessage(message.Body, &msgBody); err != nil {
		c.sendSumoMessages(ctx, err.Error(), nil)
		return err
	}

	MalawiRequest := c.mapTnmMalawiRequest(msgBody)
	headers := make(map[string]string, 0)
	url := "https://mpgs.pgcoza.biz"
	fmt.Println(url)

	log.Infof(*c.requestId, "trying to send request to payment gateway",
		MalawiRequest, "to:", url)

	// send Query getStatus to the Malawi.
	if err := (*c.httpClient).PostWithJsonResponse(url, headers, MalawiRequest, mtnResponse); err != nil {
		return err
	}

	log.Infof(*c.requestId, "successfully retrieved payment gateway response %v", mtnResponse)

	if err = json.Unmarshal([]byte(mtnResponse.ResponseBody), mtnResponseBody); err != nil {
		log.Error(*c.requestId, "___ERROR___ : Can't Read Response From Malawi ", err.Error())
		return err
	}
	log.Infof(*c.requestId, "RES FROM MALAWI %v", mtnResponseBody)

	if mtnResponseBody.ResultCode == enums.StatusCode {

		log.Infof(*c.requestId, "updateStatusRefund")
		return c.updateStatusRefund(ctx, mtnResponseBody, msgBody)
	} else {
		log.Infof(*c.requestId, "Status is pending %v", mtnResponse.ResponseBody)

	}

	return err
}

func (c *Controller) sendSumoMessages(ctx context.Context, message string, params interface{}) {

	if params != nil {
		params = fmt.Sprintf("%+v", params)
	}

	sumo := &models.SumoPusherMessage{
		Category: "malawi",
		SumoPayload: models.SumoPayload{
			Stack:   *c.requestId,
			Message: "[tnm-malawi-status-check] " + message,
			Params:  params,
		},
	}
	log.Info(*c.requestId, "Sumo Params : ", sumo.SumoPayload.Params)
	messageBody, err := json.Marshal(sumo)
	if err != nil {
		log.Error(*c.requestId, "Error Create Message Body For SQS: ", err.Error())
		return
	}

	sqsMessage := &sqs.SendMessageInput{
		MessageBody: aws.String(string(messageBody)),
	}

	_, err = c.sumoProducer.SendMsg(ctx, sqsMessage)

	if err != nil {
		log.Error(*c.requestId, "Error while pushing to sqs producer: ", err.Error())
		return
	}

	log.Info(*c.requestId, enums.SuccessfullyPushed)
}

func (c *Controller) getMessage(message string, messageData interface{}) error {

	//log.Info(*c.requestId, "trying to retrieve message body from message: ", message)
	if err := json.Unmarshal([]byte(message), &messageData); err != nil {
		log.Error(*c.requestId, "unable to retrieve message body: ", err.Error())
		return err
	}

	//log.Info(*c.requestId, "Successfully retrieved message: ", messageData)
	return nil
}

func (c *Controller) mapTnmMalawiRequest(msgBody *models.IncomingRequest) *models.QueryStatus {

	paymentDetails := models.RouteParams{
		Action:          msgBody.Action,
		UrlQuery:        msgBody.UrlQuery,
		TranType:        msgBody.TranType,
		OriginalTransId: msgBody.OriginalTransId,
	}
	return &models.QueryStatus{
		ApiKey:       msgBody.ApiKey,
		AcquireRoute: msgBody.AcquireRoute,
		ApiSecret:    msgBody.ApiSecret,
		RouteParams:  paymentDetails,
	}

}
