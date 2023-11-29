package proces

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/mitchellh/mapstructure"
	"malawi-getstatus/cache"
	"malawi-getstatus/cache/redis"
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
	cacheClient  *cache.Cache
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
	controller.initCacheService()
	return &controller
}

func (c *Controller) initSecret(secret string) {
	c.secretHolder = &SecretIDHolder{
		SecretID: secret,
		Client:   CreateSMClient(),
	}
	c.config = c.secretHolder.LoadSecret()

}

func (c *Controller) initCacheService() {
	cacheClient, err := redis.NewCache(c.config.Cache)
	if err != nil {
		log.Fatalf(*c.requestId, "Lambda init failed on cache service: %v", err)
	}
	c.cacheClient = &cacheClient
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
	fmt.Print("localRepo", localRepo)
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
	//c.sendSumoMessages(ctx, "start Get-Status TNM Malawi", message)
	var err error
	redisMessage := new(models.RedisMessage)
	msgBody := new(models.IncomingRequest)
	tnmResponseBody := new(models.ApiResult)

	if err = c.getMessage(message.Body, &msgBody); err != nil {
		c.sendSumoMessages(ctx, err.Error(), nil)
		return nil
	}

	if err = c.GetCache(ctx, redisMessage.RedisKey, redisMessage); err != nil {
		c.sendSumoMessages(ctx, err.Error(), nil)
		return err
	}

	log.Infof(*c.requestId, "message body", msgBody)
	// check if counter is bigger that max retry..
	//if utils.SafeAtoi(msgBody., 10) >= utils.SafeAtoi(messageBody.MaxRetry, 3) {
	//	log.Info(*c.requestId, "break sendRetryMessage Counter Over limit ", messageBody.Counter)
	//	return nil
	//}
	//
	log.Infof(*c.requestId, "check status by transId", msgBody.TransId)
	log.Info("URL", msgBody.URLQuery)
	// todo GET CURRENT STATUS OF TRANSACTION BY TRANS_ID
	//if messageBody.IsRefund == "true" {
	//	return c.RefundProcess(ctx, messageBody, redisMessage)
	//}
	transactionStatus, err := c.GetTransactionStatus(msgBody.TransId)
	if err != nil {
		log.Info(*c.requestId, "", err.Error())
		c.sendSumoMessages(ctx, err.Error(), nil)
		return err
	}
	fmt.Println("transactionStatus****", transactionStatus)
	if transactionStatus != enums.Pending {
		log.Info(*c.requestId, "Status not Pending ", transactionStatus)
		return nil
	}
	if tnmResponseBody, err = c.SendGetStatus(ctx, msgBody); err != nil {
		c.sendSumoMessages(ctx, err.Error(), nil)
		return err
	}
	log.Info("RSP: Lambda <--- TNM MALAWI: ", tnmResponseBody)

	//malawiRequest := c.mapTnmMalawiRequest(msgBody)

	//return response, nil
	//
	//log.Infof(*c.requestId, "trying to send request", msgBody)
	//url2 := "https://dev.payouts.tnmmpamba.co.mw/api/invoices/ST443Y5YT56532"
	//
	//// send Query getStatus to the Malawi.
	//
	//mtnResponseBody.ResultCode = "200"
	////log.Info("RSP: Lambda <--- TNM MALAWI: ", responseBody)
	//log.Infof(*c.requestId, "RESPONSE_FROM_MALAWI %v", tnmResponseBody)
	//
	//if mtnResponseBody.ResultCode == enums.StatusCode {
	//
	//	log.Infof(*c.requestId, "updateStatusRefund")
	//	return c.updateStatusRefund(ctx, mtnResponseBody, msgBody)
	//} else {
	//	log.Infof(*c.requestId, "Status is pending %v", mtnResponse)
	//
	//}

	return nil
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

	if err := json.Unmarshal([]byte(message), &messageData); err != nil {
		log.Error(*c.requestId, "unable to retrieve message body: ", err.Error())
		return err
	}

	log.Info(*c.requestId, "Successfully retrieved message: ", messageData)
	return nil
}

//func (c *Controller) mapTnmMalawiRequest(msgBody *models.IncomingRequest) *models.QueryStatus {
//
//	paymentDetails := models.RouteParams{
//		Action:          msgBody.Action,
//		UrlQuery:        msgBody.UrlQuery,
//		TranType:        msgBody.TranType,
//		OriginalTransId: msgBody.TransId,
//	}
//	return &models.QueryStatus{
//		ApiKey:       msgBody.ApiKey,
//		AcquireRoute: msgBody.AcquireRoute,
//		ApiSecret:    msgBody.ApiSecret,
//		RouteParams:  paymentDetails,
//	}
//
//}

// get cache
func (c *Controller) GetCache(ctx context.Context, redisKey string, cache *models.RedisMessage) error {
	var err error
	var result map[string]string

	log.Info(*c.requestId, "trying to retrieve cache from redis: ", redisKey)

	if result, err = (*c.cacheClient).HGetAll(ctx, redisKey); err != nil {
		log.Error(*c.requestId, "unable to retrieve cache from redis: ", err.Error())
		return err
	}

	if err = mapstructure.Decode(result, &cache); err != nil {
		return err
	}

	log.Info(*c.requestId, "Successfully retrieved cache from redis: ", cache)
	return nil
}
