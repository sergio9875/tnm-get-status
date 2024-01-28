package process

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"malawi-getstatus/enums"
	log "malawi-getstatus/logger"
	"malawi-getstatus/models"
	"strconv"
)

func (c *Controller) sendSumoMessages(ctx context.Context, message string, params interface{}) {

	if params != nil {
		params = fmt.Sprintf("%+v", params)
	}

	sumo := &models.SumoPusherMessage{
		Category: "mno",
		SumoPayload: models.SumoPayload{
			Stack:   *c.requestId,
			Message: "[mos-status-check] " + message,
			Params:  params,
		},
	}
	//log.Info(*c.requestId, "Sumo Params : ", sumo.SumoPayload.Params)
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

// GetMessage get message from queue
func (c *Controller) GetMessage(message string, messageData interface{}) error {
	if err := json.Unmarshal([]byte(message), &messageData); err != nil {
		log.Error(*c.requestId, "unable to retrieve message body: ", err.Error())
		return err
	}
	return nil
}

func (c *Controller) SendRetryMessage(ctx context.Context, messageBody *models.IncomingRequest, redisMessage *models.RedisMessage) error {
	log.Info(*c.requestId, "start sendRetryMessage", messageBody)
	log.Info(*c.requestId, enums.SuccessfullyPushed)
	return c.SendMessage(ctx, redisMessage, messageBody.Ttl)
}

func (c *Controller) SendMessage(ctx context.Context, messageBody *models.RedisMessage, ttl string) error {
	sqsB, err := json.Marshal(messageBody)
	log.Info(*c.requestId, "messageBody", messageBody)
	if err != nil {
		log.Error(*c.requestId, "Error Create Message Body For SQS: ", err.Error())
		return err

	}
	delaySeconds, err := strconv.ParseInt(ttl, 0, 32)
	if err != nil {
		log.Error(*c.requestId, "Error Cant Conver String Into Int32 ", err.Error())
		return err
	}
	sqsMessage := &sqs.SendMessageInput{
		DelaySeconds: int32(delaySeconds),
		MessageBody:  aws.String(string(sqsB)),
	}
	c.sendSumoMessages(ctx, "sqsMessage: "+enums.SuccessfullyPushed, sqsMessage.MessageBody)
	log.Info(*c.requestId, "Waiting delay seconds: ", delaySeconds)
	_, err = c.sqsProducer.SendMsg(ctx, sqsMessage)
	if err != nil {
		log.Error(*c.requestId, "Error while pushing to sqs producer: ", err.Error())
		return err
	}
	c.sendSumoMessages(ctx, enums.SuccessfullyPushed, messageBody)
	return nil
}
