package mssql

import (
	"context"
	"database/sql"
	"time"
)

func (r *repository) UpdateTransRefund(transId int, amount float64, activeStatus int, transrId int) error {
	transactive := 0
	if activeStatus == 3 {
		transactive = 1
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	now := time.Now()
	currentTime := now.Format("2006-01-02 01:01:00")
	const query = `UPDATE africainv.dbo.TRANSR SET 
					TRANSRACTIVE = @transactive,
					TRANSRtransrstatusid = @activeStatus,
					TRANSRTERMINALRefundFeeAmount = @amount,
					TRANSRrefunddate = @currentTime
				  WHERE TRANSRtransid = @transId and transrid = @transrId`
	_, err := r.db.ExecContext(ctx, query, sql.NamedArg{Name: "transactive", Value: transactive}, sql.NamedArg{Name: "activeStatus", Value: activeStatus}, sql.NamedArg{Name: "transId", Value: transId}, sql.NamedArg{Name: "amount", Value: amount}, sql.NamedArg{Name: "currentTime", Value: currentTime}, sql.NamedArg{Name: "transrId", Value: transrId})
	if err != nil {
		return err
	}
	return nil
}
