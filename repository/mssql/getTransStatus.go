package mssql

import (
	"context"
	"database/sql"
	"malawi-getstatus/models"
	"time"
)

func (r *repository) GetTransStatus(transId int) (*models.TransEntity, error) {

	transStatus := new(models.TransEntity)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	const query = `SELECT TRANStransstid FROM africainv.dbo.TRANS WHERE TRANSID = @transId`
	err := r.db.QueryRowContext(ctx, query, sql.NamedArg{Name: "transId", Value: transId}).Scan(&transStatus.TransStatus)
	if err != nil {
		return nil, err
	}
	return transStatus, nil
}
