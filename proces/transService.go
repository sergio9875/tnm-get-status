package proces

import "strconv"

// import (
//
//	"context"
//	"strconv"
//
//	log "malawi-getstatus/logger"
//	"malawi-getstatus/models"
//
// )
//
// func (c *Controller) updateStatusRefund(ctx context.Context, responseBody *models.ResponseBody, messageBody *models.IncomingRequest) error {
//
//		transId, err := strconv.Atoi(messageBody.TransId)
//		if err != nil {
//			log.Error(*c.requestId, "Cant Convert TransId to int  ", err.Error())
//			return err
//		}
//		transRid, err := strconv.Atoi(messageBody.TransrId)
//		if err != nil {
//			log.Error(*c.requestId, "Cant Convert TransId to int  ", err.Error())
//			return err
//		}
//		amount, err := strconv.ParseFloat(strconv.Itoa(messageBody.Amount), 64)
//		if err != nil {
//			log.Error(*c.requestId, "Cant Convert Amount to int  ", err.Error())
//			return err
//		}
//
//		if err := (*c.repository).UpdateTransRefund(transId, amount, GetPaymentCodeForRefundStatus(responseBody), transRid); err != nil {
//			log.Error(*c.requestId, "Cant update Trans :  ", err.Error())
//			return err
//		}
//		log.Info(*c.requestId, "Update Trans ", messageBody.TransId)
//		log.Info(*c.requestId, "Update Transr ", messageBody.TransrId)
//		//c.sendSumoMessages(ctx, enums.SuccessfullyUpdate, messageBody.TransId)
//		return nil
//	}
//func (c *Controller) GetTransactionStatus(transId string) (int, error) {
//	TRANSID, err := strconv.ParseInt(transId, 10, 64)
//	if err != nil {
//		return -1, err
//	}
//	transSettings, err := (*c.repository).GetTransStatus(int(TRANSID))
//	if err != nil {
//		return -1, err
//	}
//	return transSettings.TransStatus, nil
//}

func (c *Controller) GetTransactionStatus(transId string) (int, error) {
	TRANSID, err := strconv.ParseInt(transId, 10, 64)
	if err != nil {
		return -1, err
	}
	transSettings, err := (*c.repository).GetTransStatus(int(TRANSID))
	if err != nil {
		return -1, err
	}
	return transSettings.TransStatus, nil
}
