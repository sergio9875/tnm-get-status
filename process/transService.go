package process

import (
	"strconv"
)

func (c *Controller) GetMbtTransStatus(mbtId string) (int, error) {
	MbtID, err := strconv.ParseInt(mbtId, 10, 64)
	if err != nil {
		return -1, err
	}
	transSettings, err := (*c.repository).GetMbtStatus(int(MbtID))
	if err != nil {
		return -1, err
	}
	return transSettings.MbtStatus, nil
}
