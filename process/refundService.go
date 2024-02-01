package process

import (
	"context"
	"fmt"
	"malawi-getstatus/enums"
	log "malawi-getstatus/logger"
	"malawi-getstatus/models"
	service "malawi-getstatus/services"
	"strconv"
	"strings"
)

func (c *Controller) RefundProcess(ctx context.Context, messageBody *models.IncomingRequest, redisBody *models.RedisMessage) error {
	fmt.Println("Start_Refund_Process")
	token, err := service.GetToken(messageBody.URLToken, messageBody.Wallet, messageBody.Password)
	if err != nil {
		log.Infof("ERR_MSG: %s\n", err.Error())
		return err
	}

	fmt.Println("token@@@@@@@@@@@@_Refund")
	responseBody := new(models.TnmBodyResponse)
	responseBody, err = service.SendGetRequest(messageBody.MbtId, token.Data.Token, messageBody.URLQuery)

	if err != nil {
		c.sendSumoMessages(ctx, err.Error(), nil)
		log.Infof(*c.requestId, "ERROR_INFO: %s", err.Error())
		return err
	}
	c.sendSumoMessages(ctx, enums.MalawiResponse+"Refund", responseBody)

	fmt.Println("isDataReversed", responseBody.Reversed)
	if responseBody.Reversed == enums.IsRefund {
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

	if err := (*c.repository).UpdateTransRefund(transId, amount, GetPaymentCodeForRefundStatus(responseBody), transRid); err != nil {
		log.Error(*c.requestId, "Cant update Trans :  ", err.Error())
		return err
	}
	log.Info(*c.requestId, "Update Trans ", messageBody.TransId)
	c.sendSumoMessages(ctx, enums.SuccessfullyUpdate, messageBody.TransId)
	return nil
}
