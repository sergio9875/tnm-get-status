package process

import (
	"context"
	"fmt"
	"malawi-getstatus/enums"
	log "malawi-getstatus/logger"
	"malawi-getstatus/models"
	"time"
)

func (c *Controller) InvoiceProcess(ctx context.Context, messageBody *models.IncomingRequest, redisBody *models.RedisMessage) error {

	transactionStatus, err := c.GetTransactionStatus(messageBody.TransId)
	if err != nil {
		log.Info(*c.requestId, "", err.Error())
		c.sendSumoMessages(ctx, err.Error(), nil)
		return err
	}
	log.Info("transactionStatus", transactionStatus)
	if transactionStatus == enums.TransactionCompleted {
		log.Infof(*c.requestId, "stopping function, callback already received, transaction completed")
		return err
	}

	fmt.Println("transactionStatus****", transactionStatus)
	if transactionStatus != enums.Pending {
		log.Info(*c.requestId, "Status not Pending ", transactionStatus)
		return err
	}

	responseBody := new(models.TnmResponse)
	if responseBody, err = c.SendGetStatus(ctx, messageBody); err != nil {
		c.sendSumoMessages(ctx, err.Error(), nil)
		log.Infof(*c.requestId, "The error is "+err.Error(), nil)
		return err
	}
	log.Infof(*c.requestId, enums.MalawiResponse+"Invoice", responseBody)
	c.sendSumoMessages(ctx, enums.MalawiResponse+"Invoice", responseBody)

	fmt.Println("responseBody.Data.Paid", responseBody.Data.Paid)

	if responseBody.Data.Paid == true {
		return c.SendCallBackRequest(ctx, messageBody, responseBody.Data.ReceiptNumber)
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
	//if cbResponse.Code != enums.Success {
	//	//cbResponse.StatusCode = payload.RequestStatusCode
	//	//cbResponse.Explanation = body.StatusDescription
	//}
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
