package process

import (
	"golang.org/x/net/context"
	"malawi-getstatus/cache/redis"
	log "malawi-getstatus/logger"
	"malawi-getstatus/repository/mssql"
	"malawi-getstatus/request"
)

func (c *Controller) initCacheService() {
	cacheClient, err := redis.NewCache(c.config.Cache)
	if err != nil {
		log.Fatalf(*c.requestId, "Lambda init failed on cache service: %v", err)
	}
	c.cacheClient = &cacheClient
}

func (c *Controller) initRepository() {
	localRepo, err := mssql.NewRepository(c.config.Database.Africainv)
	if err != nil {
		log.Fatalf(*c.requestId, "Lambda init failed on repository: %s", err)
	}
	c.repository = &localRepo
}

func (c *Controller) initSecret(secret string) {
	c.secretHolder = &SecretIDHolder{
		SecretID: secret,
		Client:   CreateSMClient(),
	}
	c.config = c.secretHolder.LoadSecret()
}

func (c *Controller) initSqsProducer(queueName string) {
	var err error
	c.sqsProducer, err = NewSQSProducerFromUrl(context.TODO(), CreateSQSClient(),
		&queueName)
	if err != nil {
		log.Fatalf(*c.requestId, "Lambda init failed on sqs producer: %v", err)
	}
}

func (c *Controller) initSumoProducer() {
	var err error
	c.sumoProducer, err = NewSQSProducerFromUrl(context.TODO(), CreateSQSClient(),
		&c.config.Services.SumoPusherUrl)
	if err != nil {
		log.Fatalf(*c.requestId, "Lambda init failed on sqs producer: %v", err)
	}
}

func (c *Controller) initClient() {
	httpClient, err := request.NewClient()

	if err != nil {
		log.Fatalf(*c.requestId, "Lambda init failed on http client service: %v", err)
	}

	c.httpClient = &httpClient
}
