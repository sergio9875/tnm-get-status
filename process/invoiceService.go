package process

import (
	"context"
	"errors"
	"fmt"
	"malawi-getstatus/enums"
	log "malawi-getstatus/logger"
	"malawi-getstatus/models"
	service "malawi-getstatus/services"
	"time"
)

func (c *Controller) InvoiceProcess(ctx context.Context, messageBody *models.IncomingRequest, redisBody *models.RedisMessage) error {
	fmt.Println("Start_Invoice_Process")
	mbtStatus, err := c.GetMbtTransStatus(messageBody.MbtId)
	if err != nil {
		log.Info(*c.requestId, "transactionStatus", err.Error())
		c.sendSumoMessages(ctx, err.Error(), nil)
		return err
	}

	log.Info("transactionStatus", mbtStatus)
	if mbtStatus != enums.Pending {
		log.Infof(*c.requestId, "stopping function, callback already received, transaction completed")
		return nil
	}
	
	token, err := service.GetToken(messageBody.URLToken, messageBody.Wallet, messageBody.Password)
	if err != nil {
		log.Infof("ERR_MSG: %s\n", err.Error())
		return err
	}

	fmt.Println("token_	Invoice_@@@@@@@@@@@@_Process....")

	responseBody := new(models.TnmBodyResponse)
	responseBody, err = service.SendGetRequest(messageBody.MbtId, token.Data.Token, messageBody.URLQuery)

	if err != nil {
		c.sendSumoMessages(ctx, err.Error(), nil)
		log.Infof(*c.requestId, "ERROR_INFO: %s", err.Error())
		return err
	}
	fmt.Println("isPaid____##########____....", responseBody.Paid)
	c.sendSumoMessages(ctx, enums.MalawiResponse+"Invoice", responseBody)

	if responseBody.Paid == true {
		return c.SendCallBackRequest(ctx, messageBody, responseBody.ReceiptNumber)
	}
	return c.SendRetryMessage(ctx, messageBody, redisBody)

}

// SendCallBackRequest send transaction to callback End point
func (c *Controller) SendCallBackRequest(ctx context.Context, body *models.IncomingRequest, receiptNumber string) error {
	log.Infof(*c.requestId, "Start Tnm Malawi Callback Process", body)
	cbResponse := new(models.PaymentGatewayResponse)

	payload := c.mapPaymentGatewayRequest(body, receiptNumber)
	c.sendSumoMessages(ctx, "payment callback request", payload)
	log.Infof(*c.requestId, "trying to send request to Tnm Malawi - Callback endpoint",
		payload, "to:", body.CallbackUrl)
	if err := (*c.httpClient).PostWithJsonResponse(body.CallbackUrl,
		make(map[string]string, 0), payload, cbResponse); err != nil {
		log.Info("___ERROR___ : Cant Send Post Request ", err.Error())
		log.Error(*c.requestId, "999;Invalid Deatils;", err.Error())
		return err
	}
	log.Infof(*c.requestId, "callback response", cbResponse)
	c.sendSumoMessages(ctx, "payment gateway response", cbResponse)
	if cbResponse.Code != enums.Success {
		_ = errors.New(fmt.Sprintf("statusCode: %v and errMsg: %v", payload.ReceiptCode, payload.ResultDescription))
	}
	return nil
}

func (c *Controller) mapPaymentGatewayRequest(msgBody *models.IncomingRequest, receiptNumber string) *models.CallbackRequest {
	return &models.CallbackRequest{
		ReceiptNumber:     receiptNumber,
		ReceiptCode:       "0",
		ResultDescription: enums.SuccessfullyDesc,
		ResultTime:        time.Now().Format("2006-01-02 15:04:05"),
		TransactionId:     msgBody.MbtId,
		Success:           true,
	}
}
