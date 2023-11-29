package models

type RedisMessage struct {
	RedisKey string `json:"redisKey"`
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
	CallbackUrl     string `json:"CallbackUrl" validate:"omitempty,url"`
	Action          string `json:"Action" validate:"required,actionTypes"`
	URLQuery        string `json:"URLQuery,omitempty" validate:"omitempty,url"`
	URLToken        string `json:"URLToken,omitempty" validate:"omitempty,url"`
	CellphoneNumber string `json:"CellphoneNumber" validate:"omitempty,cellnumber"`
	Amount          int    `json:"Amount"`
	Wallet          string `json:"Wallet" redact:"complete"`
	Password        string `json:"Password" redact:"complete"`
	TransId         string `json:"TransId"`
	MbtId           string `json:"MbtId"`
	Description     string `json:"Description,omitempty"`
	IsInvoice       bool   `json:"IsInvoice,omitempty"`
	IsRefund        bool   `json:"IsRefund,omitempty"`
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
