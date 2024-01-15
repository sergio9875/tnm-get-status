package process

import (
	"malawi-getstatus/enums"
	"malawi-getstatus/models"
	"strconv"
	"time"
)

func GetPaymentGatewayCode(responseBody *models.ResponseBody) int {
	if responseBody.ResultCode == enums.Success {
		return 3
	}
	return 7
}

func GetTimeStamp() string {
	t := time.Now().UnixNano() / 1000000
	return strconv.Itoa(int(t))
}

func GetPaymentCodeForRefundStatus(responseBody *models.TnmBodyResponse) int {
	if responseBody.StatusCode == enums.StatusCode {
		return 3
	}
	return 4
}

func stringToFloat(str string) float64 {
	if floatnumber, err := strconv.ParseFloat(str, 64); err == nil {
		return floatnumber
	}
	return 0.00
}
