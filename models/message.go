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
	Url             string `json:"url,omitempty"`
	ApiKey          string `json:"apiKey" redact:"complete"`
	ApiSecret       string `json:"apiSecret" redact:"complete"`
	AcquireRoute    string `json:"acquireRoute,omitempty"`
	Action          string `json:"action,omitempty"`
	UrlQuery        string `json:"urlQuery,omitempty"`
	TranType        int    `json:"tranType,omitempty"`
	OriginalTransId string `json:"originalTransId,omitempty"`
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
	UrlQuery        string `json:"UrlQuery"`
	TranType        int    `json:"TranType"`
	OriginalTransId string `json:"OriginalTransId"`
}

type QueryStatus struct {
	ApiKey       string      `json:"apiKey" redact:"complete"`
	ApiSecret    string      `json:"apiSecret" redact:"complete"`
	AcquireRoute string      `json:"acquireRoute,omitempty"`
	RouteParams  RouteParams `json:"routeParams,omitempty"`
}

type ResponseBody struct {
	ConversationId    string `json:"ConversationId,omitempty"`
	ResponseTime      string `json:"ResponseTime,omitempty"`
	TransId           string `json:"TransId,omitempty"`
	OriginalTransId   string `json:"OriginalTransId,omitempty"`
	ResultCode        string `json:"ResultCode,omitempty"`
	ResultDescription string `json:"ResultDescription,omitempty"`
}
