package process

import (
	"context"
	"golang.org/x/exp/slices"
	"malawi-getstatus/enums"
	log "malawi-getstatus/logger"
	"malawi-getstatus/models"
	service "malawi-getstatus/services"
	"strconv"
)

func (c *Controller) RefundProcess(ctx context.Context, messageBody *models.IncomingRequest, redisBody *models.RedisMessage) error {
	transrid, err := strconv.Atoi(messageBody.TransId)

	if err != nil {
		c.sendSumoMessages(ctx, err.Error(), nil)
		log.Infof(*c.requestId, "The error is "+err.Error(), nil)
		return err
	}

	transrStatus, err := c.getRefundStatus(transrid)
	log.Infof(*c.requestId, "transrStatus", transrStatus)

	if err != nil {
		c.sendSumoMessages(ctx, err.Error(), nil)
		log.Infof(*c.requestId, "The error is "+err.Error(), nil)
		return err
	}
	transrStatusArray := []int{2, 5, 6, 9}
	if slices.Contains(transrStatusArray, transrStatus) {
		responseBody := new(models.TnmResponse)
		token, err := service.GetToken(messageBody.URLToken, messageBody.Wallet, messageBody.Password)
		if err != nil {
			log.Infof("ERR_MSG: %s\n", err.Error())
			return err
		}
		if responseBody, err = c.SendGetRequest(messageBody.TransId, token.Data.Token, messageBody.URLQuery); err != nil {
			c.sendSumoMessages(ctx, err.Error(), nil)
			log.Infof(*c.requestId, "The error is "+err.Error(), nil)
			return err
		}
		log.Infof(*c.requestId, enums.MalawiResponse+"Refund", responseBody)
		c.sendSumoMessages(ctx, enums.MalawiResponse+"Refund", responseBody)

		if responseBody.Data.Reversed == true {
			return c.UpdateRefund(ctx, &responseBody.Data, messageBody)
		}
		return c.SendRetryMessage(ctx, messageBody, redisBody)

	}
	return nil
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
	amount, err := strconv.ParseFloat(messageBody.Amount, 64)
	if err != nil {
		log.Error(*c.requestId, "Cant Convert Amount to int  ", err.Error())
		return err
	}

	if err := (*c.repository).UpdateTransRefund(transId, amount, GetPaymentCodeForRefundStatus(responseBody), transRid); err != nil {
		log.Error(*c.requestId, "Cant update Trans :  ", err.Error())
		return err
	}
	log.Info(*c.requestId, "Update Trans ", messageBody.TransId)
	log.Info(*c.requestId, "Update Transr ", messageBody.TransrId)
	c.sendSumoMessages(ctx, enums.SuccessfullyUpdate, messageBody.TransId)
	return nil
}
