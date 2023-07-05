package process

import (
	"context"
	"strconv"

	"malawi-getstatus/enums"
	log "malawi-getstatus/logger"
	"malawi-getstatus/models"
)

func (c *Controller) updateStatusRefund(ctx context.Context, responseBody *models.ResponseBody, messageBody *models.IncomingRequest) error {

	transid, err := strconv.Atoi(messageBody.TransId)
	if err != nil {
		log.Error(*c.requestId, "Cant Convert TransId to int  ", err.Error())
		return err
	}
	transrid, err := strconv.Atoi(messageBody.TransrId)
	if err != nil {
		log.Error(*c.requestId, "Cant Convert TransId to int  ", err.Error())
		return err
	}
	amount, err := strconv.ParseFloat(messageBody.Amount, 64)
	if err != nil {
		log.Error(*c.requestId, "Cant Convert Amount to int  ", err.Error())
		return err
	}

	if err := (*c.repository).UpdateTransRefund(transid, amount, GetPaymentCodeForRefundStatus(responseBody), transrid); err != nil {
		log.Error(*c.requestId, "Cant update Trans :  ", err.Error())
		return err
	}
	log.Info(*c.requestId, "Update Trans ", messageBody.TransId)
	log.Info(*c.requestId, "Update Transr ", messageBody.TransrId)
	c.sendSumoMessages(ctx, enums.SuccessfullyUpdate, messageBody.TransId)
	return nil
}
