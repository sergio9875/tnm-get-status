package process

import (
	"strconv"
)

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
