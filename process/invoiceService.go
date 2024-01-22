package process

import (
	"context"
	"malawi-getstatus/enums"
	log "malawi-getstatus/logger"
	"malawi-getstatus/models"
	service "malawi-getstatus/services"
	"time"
)

func (c *Controller) InvoiceProcess(ctx context.Context, messageBody *models.IncomingRequest, redisBody *models.RedisMessage) error {

	transactionStatus, err := c.GetTransactionStatus(messageBody.TransId)
	if err != nil {
		log.Info(*c.requestId, "transactionStatus", err.Error())
		c.sendSumoMessages(ctx, err.Error(), nil)
		return err
	}
	log.Info("transactionStatus", transactionStatus)
	if transactionStatus == enums.TransactionCompleted {
		log.Infof(*c.requestId, "stopping function, callback already received, transaction completed")
		return err
	}

	if transactionStatus != enums.Pending {
		log.Info(*c.requestId, "Status not Pending ", transactionStatus)
		return err
	}

	token, err := service.GetToken(messageBody.URLToken, messageBody.Wallet, messageBody.Password)
	if err != nil {
		log.Infof("ERR_MSG: %s\n", err.Error())
		return err
	}

	//fmt.Println("Process....")
	//fmt.Println("token@@@@@@@@@@@@", token.Data.Token)

	responseBody := new(models.TnmResponse)
	responseBody, err = c.SendGetRequest(messageBody.TransId, token.Data.Token, messageBody.URLQuery)
	//c.sendSumoMessages(ctx, err.Error(), nil)
	if err != nil {
		log.Infof(*c.requestId, "ERROR_INFO: %s", err.Error())
		return err
	}

	c.sendSumoMessages(ctx, enums.MalawiResponse+"Invoice", responseBody)
	//os.Exit(2)

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
