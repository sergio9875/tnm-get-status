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
	Url             string `json:"url"`
	ApiKey          string `json:"apiKey" redact:"complete"`
	ApiSecret       string `json:"apiSecret" redact:"complete"`
	AcquireRoute    string `json:"acquireRoute"`
	Action          string `json:"action" validate:"required"`
	UrlQuery        string `json:"urlQuery" redact:"complete"`
	TranType        int    `json:"tranType"`
	OriginalTransId string `json:"originalTransId" redact:"complete"`
	TransId         string `json:"transId,omitempty"`
	TransrId        string `json:"transrId,omitempty"`
	Amount          string `json:"amount,omitempty"`
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

type QueryStatus struct {
	ApiKey       string      `json:"apiKey" redact:"complete"`
	ApiSecret    string      `json:"apiSecret" redact:"complete"`
	AcquireRoute string      `json:"acquireRoute" redact:"complete"`
	RouteParams  RouteParams `json:"routeParams"`
}

type ResponseBody struct {
	ConversationId    string `json:"ConversationId"`
	ResponseTime      string `json:"ResponseTime"`
	TransId           string `json:"TransId"`
	OriginalTransId   string `json:"OriginalTransId"`
	ResultCode        string `json:"ResultCode"`
	ResultDescription string `json:"ResultDescription"`
}
