package process

import (
	"context"
	"fmt"
	"malawi-getstatus/enums"
	log "malawi-getstatus/logger"
	"malawi-getstatus/models"
	service "malawi-getstatus/services"
	"os"
	"strconv"
	"strings"
)

func (c *Controller) RefundProcess(ctx context.Context, messageBody *models.IncomingRequest, redisBody *models.RedisMessage) error {

	token, err := service.GetToken(messageBody.URLToken, messageBody.Wallet, messageBody.Password)
	if err != nil {
		log.Infof("ERR_MSG: %s\n", err.Error())
		return err
	}

	fmt.Println("token@@@@@@@@@@@@", token)

	responseBody := new(models.TnmBodyResponse)
	responseBody, err = service.SendGetRequest(messageBody.TransId, token.Data.Token, messageBody.URLQuery)

	if err != nil {
		c.sendSumoMessages(ctx, err.Error(), nil)
		log.Infof(*c.requestId, "ERROR_INFO: %s", err.Error())
		return err
	}

	c.sendSumoMessages(ctx, enums.MalawiResponse+"Refund", responseBody)
	//if responseBody, err = service.SendGetRequest(messageBody.TransId, token.Data.Token, messageBody.URLQuery); err != nil {
	//	c.sendSumoMessages(ctx, err.Error(), nil)
	//	log.Infof(*c.requestId, "The error is "+err.Error(), nil)
	//	return err
	//}
	//log.Infof(*c.requestId, enums.MalawiResponse+"Refund", responseBody)
	//c.sendSumoMessages(ctx, enums.MalawiResponse+"Refund", responseBody)

	fmt.Println("data-reversed", responseBody.ReversedAt)
	os.Exit(2)
	if responseBody.Reversed == true {
		return c.UpdateRefund(ctx, responseBody, messageBody)
	}
	return c.SendRetryMessage(ctx, messageBody, redisBody)

}

func (c *Controller) getRefundStatus(transrId int) (int, error) {
	transrSettings, err := (*c.repository).GetRefundStatus(transrId)
	if err != nil {
		return -1, err
	}
	return transrSettings.TRANSRtransrstatusid, nil
}

func (c *Controller) UpdateRefund(ctx context.Context, responseBody *models.TnmBodyResponse, messageBody *models.IncomingRequest) error {

	transId, err := strconv.Atoi(messageBody.TransId)
	if err != nil {
		log.Error(*c.requestId, "Cant Convert TransId to int  ", err.Error())
		return err
	}
	transRid, err := strconv.Atoi(messageBody.TransrId)
	if err != nil {
		log.Error(*c.requestId, "Cant Convert TransId to int  ", err.Error())
		return err
	}

	amount, err := strconv.ParseFloat(strings.TrimSpace(messageBody.Amount), 64)
	if err != nil {
		log.Error(*c.requestId, "Cant Convert Amount to int  ", err.Error())
		return err
	}

	fmt.Println("trans-ID", transId)
	fmt.Println("transR-ID", transRid)
	fmt.Println("amount", amount)
	//os.Exit(2)
	if err := (*c.repository).UpdateTransRefund(transId, amount, GetPaymentCodeForRefundStatus(responseBody), transRid); err != nil {
		log.Error(*c.requestId, "Cant update Trans :  ", err.Error())
		return err
	}
	log.Info(*c.requestId, "Update Trans ", messageBody.TransId)
	log.Info(*c.requestId, "Update Transr ", messageBody.TransrId)
	c.sendSumoMessages(ctx, enums.SuccessfullyUpdate, messageBody.TransId)
	return nil
}
