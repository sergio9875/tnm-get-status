package process

//
//// GetMessage get message from queue
//func (c *p.Controller) GetMessage(message string, messageData interface{}) error {
//
//	log.Info(*c.requestId, "trying to retrieve message body from message: ", message)
//	if err := json.Unmarshal([]byte(message), &messageData); err != nil {
//		log.Error(*c.requestId, "unable to retrieve message body: ", err.Error())
//		return err
//	}
//
//	log.Info(*c.requestId, "Successfully retrieved message: ", messageData)
//	return nil
//}
//
//func (c *p.Controller) SendRetryMessage(ctx context.Context, messageBody *models.IncomingRequest) error {
//	log.Info(*c.requestId, "start sendRetryMessage", messageBody)
//	log.Info(*c.requestId, enums.SuccessfullyPushed)
//	return c.SendMessage(ctx, messageBody, messageBody.Ttl)
//}
//
//func (c *p.Controller) SendMessage(ctx context.Context, messageBody *models.IncomingRequest, ttl string) error {
//	sqsB, err := json.Marshal(messageBody)
//	log.Info(*c.requestId, "messageBody", messageBody)
//	if err != nil {
//		log.Error(*c.requestId, "Error Create Message Body For SQS: ", err.Error())
//		return err
//
//	}
//	delaySeconds, err := strconv.ParseInt(ttl, 0, 32)
//	if err != nil {
//		log.Error(*c.requestId, "Error Cant Conver String Into Int32 ", err.Error())
//		return err
//	}
//	sqsMessage := &sqs.SendMessageInput{
//		DelaySeconds: int32(delaySeconds),
//		MessageBody:  aws.String(string(sqsB)),
//	}
//	c.sendSumoMessages(ctx, "sqsMessage: "+enums.SuccessfullyPushed, sqsMessage.MessageBody)
//	log.Info(*c.requestId, "Waiting delay seconds: ", delaySeconds)
//	_, err = c.sqsProducer.SendMsg(ctx, sqsMessage)
//	if err != nil {
//		log.Error(*c.requestId, "Error while pushing to sqs producer: ", err.Error())
//		return err
//	}
//	c.sendSumoMessages(ctx, enums.SuccessfullyPushed, messageBody)
//	return nil
//}
