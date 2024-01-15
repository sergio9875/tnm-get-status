package models

type RedisMessage struct {
	RedisKey string `json:"redisKey"`
}

type Message struct {
	ReferenceId      string `json:"referenceId"`
	ServiceName      string `json:"serviceName"`
	PaymentReference string `json:"paymentReference"`
	QueueName        string `json:"queueName"`
	Ttl              string `json:"ttl"`
	MaxRetry         string `json:"maxRetry"`
	Counter          string `json:"counter"`
	ConsumerKey      string `json:"consumerKey" redact:"complete"`
	ConsumerSecret   string `json:"consumerSecret" redact:"complete"`
	AcquireRoute     string `json:"acquireRoute" redact:"complete"`
	UrlQuery         string `json:"urlQuery" redact:"complete"`
	UrlToken         string `json:"urlToken" redact:"complete"`
	Action           string `json:"action"`
	IsRefund         string `json:"isRefund,omitempty"`
	TransrId         string `json:"transrId,omitempty"`
	TransId          string `json:"transId"`
}

type SumoPusherMessage struct {
	Category    string      `json:"category"`
	Fields      string      `json:"fields,omitempty"`
	SumoPayload SumoPayload `json:"sumoPayload,omitempty"`
}

type SumoPayload struct {
	Stack   string      `json:"stack"`
	Message string      `json:"message"`
	Params  interface{} `json:"params"`
}

type IncomingRequest struct {
	Action          string `json:"action" validate:"required,actionTypes"`
	URLQuery        string `json:"URLQuery,omitempty" validate:"omitempty,url"`
	URLToken        string `json:"URLToken,omitempty" validate:"omitempty,url"`
	CellphoneNumber string `json:"CellphoneNumber" validate:"omitempty,cellnumber"`
	Amount          string `json:"Amount"`
	Ttl             string `json:"ttl"`
	MaxRetry        string `json:"MaxRetry"`
	Wallet          string `json:"Wallet" redact:"complete"`
	Password        string `json:"Password" redact:"complete"`
	TransId         string `json:"TransId"`
	TransrId        string `json:"TransrId"`
	MbtId           string `json:"MbtId"`
	Description     string `json:"Description,omitempty"`
	IsInvoice       string `json:"IsInvoice,omitempty"`
	IsRefund        string `json:"IsRefund,omitempty"`
	Counter         string `json:"Counter"`
	QueueName       string `json:"queueName"`
	CallbackUrl     string `json:"callbackUrl"`
}

type CallbackRequest struct {
	ReceiptNumber     string `json:"receipt_number,omitempty"`
	ReceiptCode       string `json:"receipt_code,omitempty"`
	ResultDescription string `json:"result_description,omitempty"`
	ResultTime        string `json:"result_time,omitempty"`
	TransactionId     string `json:"transaction_id,"`
	Success           bool   `json:"success,omitempty"`
}

type PaymentGatewayResponse struct {
	Code         string `json:"code"`
	Explanation  string `json:"explanation"`
	RedirectURL  string `json:"redirectURL"`
	Instructions string `json:"instructions"`
	Details      struct {
		ResultCode int    `json:"ResultCode"`
		StatusCode string `json:"StatusCode"`
	} `json:"details"`
}

type TnmBodyResponse struct {
	InvoiceNumber         string `json:"invoice_number,omitempty"`
	Amount                string `json:"amount,omitempty"`
	Msisdn                string `json:"msisdn,omitempty"`
	ReceiptNumber         string `json:"receipt_number,omitempty"`
	SettledAt             string `json:"settled_at,omitempty"`
	Paid                  bool   `json:"paid"`
	ReversalTranscationId string `json:"reversal_transcation_id"`
	Reversed              bool   `json:"reversed"`
	ReversedAt            string `json:"reversed_at"`
	StatusCode            int    `json:"statusCode"`
	Message               string `json:"message"`
}

type TnmResponse struct {
	Message string          `json:"message"`
	Errors  interface{}     `json:"errors,omitempty"`
	Trace   interface{}     `json:"trace,omitempty"`
	Data    TnmBodyResponse `json:"data,omitempty"`
}

type ApiResult struct {
	Response   string `json:"response"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Token      string `json:"token,omitempty"`
}

type Response struct {
	DpoStatus     string `json:"dpoStatus"`
	DpoStatusDesc string `json:"dpoStatusDesc"`
	ResponseBody  string `json:"response"`
}
type RouteParams struct {
	Action          string `json:"Action"`
	UrlQuery        string `json:"UrlQuery" redact:"complete"`
	TranType        int    `json:"TranType"`
	OriginalTransId string `json:"OriginalTransId" redact:"complete"`
}

type TokenRes struct {
	Token     string `json:"token,omitempty"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

type TokenResponse struct {
	Message string      `json:"message,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Trace   interface{} `json:"trace,omitempty"`
	Data    TokenRes    `json:"data,omitempty"`
}

//
//type APIRequest struct {
//	Msisdn              string `json:"msisdn,omitempty"`
//	Amount              int    `json:"amount,omitempty"`
//	Receiver_type       int64  ` json:"receiver_type,omitempty"`
//	Receiver_identifier string `json:"receiver_identifier,omitempty" redact:"last4"`
//	Tran_id             string `json:"tran_id,omitempty"`
//	Receiver_msisdn     string `json:"receiver_msisdn,omitempty"`
//	Narration           string `json:"narration,omitempty"`
//	Description         string `json:"description,omitempty"`
//	Org_tran_id         string `json:"org_tran_id,omitempty"`
//	Invoice_number      string `json:"invoice_number"`
//}

type Auth struct {
	Wallet   string `json:"wallet" required:"wallet"`
	Password string `json:"password" required:"password"`
}

type QueryStatus struct {
	ApiKey       string      `json:"apiKey" redact:"complete"`
	ApiSecret    string      `json:"apiSecret" redact:"complete"`
	AcquireRoute string      `json:"acquireRoute" redact:"complete"`
	RouteParams  RouteParams `json:"routeParams"`
}

type ErrorMalawiMessage struct {
	Message string `json:"message"`
}

type ResponseBody struct {
	ConversationId    string `json:"ConversationId"`
	ResponseTime      string `json:"ResponseTime"`
	TransId           string `json:"TransId"`
	OriginalTransId   string `json:"OriginalTransId"`
	ResultCode        string `json:"ResultCode"`
	ResultDescription string `json:"ResultDescription"`
}
