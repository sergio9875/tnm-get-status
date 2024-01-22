package process

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"malawi-getstatus/cache"
	"malawi-getstatus/models"
	repo "malawi-getstatus/repository"
	"malawi-getstatus/request"
	"malawi-getstatus/utils"
)

// Controller container
type Controller struct {
	secretHolder *SecretIDHolder
	sqsProducer  *SQSProducer
	sumoProducer *SQSProducer
	config       *models.SecretModel
	repository   *repo.Repository
	cacheClient  *cache.Cache
	httpClient   *request.IRequest
	requestId    *string
}

func NewController(secret string) *Controller {
	controller := Controller{
		requestId: utils.StringPtr("ROOT"),
	}
	controller.initSecret(secret)
	controller.initSumoProducer()
	controller.initRepository()
	controller.initCacheService()
	controller.initClient()
	return &controller
}

func (c *Controller) ShutDown() {
	c.config = nil
	c.sqsProducer = nil
	c.secretHolder = nil
	c.sumoProducer = nil
	c.cacheClient = nil
	c.httpClient = nil
	c.repository = nil
}

func (c *Controller) PreProcess(pid *string) {
	c.requestId = pid
}

func (c *Controller) PostProcess() {
	c.requestId = utils.StringPtr("ROOT")
}

func (c *Controller) Process(ctx context.Context, message events.SQSMessage) error {

	c.sendSumoMessages(ctx, "start Get-Status TNM Malawi", message)
	var messageBody = new(models.IncomingRequest)
	var redisMessage = new(models.RedisMessage)
	var err error

	if err = c.GetMessage(message.Body, messageBody); err != nil {
		c.sendSumoMessages(ctx, err.Error(), nil)
		return err
	}

	//if err = c.GetMessage(message.Body, redisMessage); err != nil {
	//	c.sendSumoMessages(ctx, err.Error(), nil)
	//	return err
	//}
	//fmt.Println("MESSAGE KEY FROM BODY", redisMessage)
	//if err = c.getCache(ctx, redisMessage.RedisKey, messageBody); err != nil {
	//	c.sendSumoMessages(ctx, err.Error(), nil)
	//	return err
	//}
	//
	//fmt.Println("MESSAGE BODY FROM REDIS", messageBody)
	//
	//fmt.Println("counter", messageBody.Counter)
	//fmt.Println("maxRetry", messageBody.MaxRetry)
	//// update counter on redis
	//if err = c.updateRedisCounter(ctx, redisMessage.RedisKey, messageBody.Counter); err != nil {
	//	c.sendSumoMessages(ctx, err.Error(), nil)
	//	return err
	//}

	// initiate the queue with value provided on terminal settings
	//c.initSqsProducer(messageBody.QueueName)

	//os.Exit(2)

	//log.Infof(*c.requestId, "message body")
	//check if counter is bigger that max retry

	//if utils.SafeAtoi(messageBody.Counter, 0) >= utils.SafeAtoi(messageBody.MaxRetry, 0) {
	//	log.Info(*c.requestId, "break sendRetryMessage Counter Over limit ", messageBody.Counter)
	//	return nil
	//}

	if messageBody.IsRefund == "true" {
		return c.RefundProcess(ctx, messageBody, redisMessage)
	}

	if messageBody.IsInvoice == "true" {
		return c.InvoiceProcess(ctx, messageBody, redisMessage)
	}

	return nil
}
