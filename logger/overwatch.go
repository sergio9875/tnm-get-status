package logger

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
	"time"
)

// SumoConfig for SumoLogic
type SumoConfig struct {
	Host           string `json:"host"`
	LogLevel       string `json:"log_level"`
	CategoryPrefix string `json:"category_prefix"`
	PushToSumo     bool   `json:"push_to_sumo"`
}

var batchSize = 10
var config *SumoConfig
var logMessageList []types.SendMessageBatchRequestEntry

var hierarchy = []string{TRACE, DEBUG, INFO, WARN, ERROR}
var hierarchyIndex = 0

var sumoCategoryKey string

func findIndex(key string) int {
	for i := 0; i < len(hierarchy); i++ {
		if key == hierarchy[i] {
			return i
		}
	}
	return 4
}

func InitOverwatch(sumoCategory string, sumoConfig *SumoConfig) {
	config = sumoConfig
	hierarchyIndex = findIndex(config.LogLevel)
	logMessageList = make([]types.SendMessageBatchRequestEntry, 0)
	sumoCategoryKey = sumoCategory
}

func CloseOverwatch() {
	logMessageList = make([]types.SendMessageBatchRequestEntry, 0)
}

func AddToSumo(level string, requestId string, format string, source string, a ...interface{}) {
	if config == nil {
		return
	}
	currentLogIndex := findIndex(level)

	if currentLogIndex >= hierarchyIndex {
		var redacted = make([]interface{}, len(a))
		for idx, item := range a {
			if item != nil {
				redacted[idx] = Sanitizer(item)
			}
		}
		//msg := fmt.Sprint(redacted...)
		AppendToList(level, requestId, format, source, redacted...)
	}
}

type sumoPayload struct {
	Level     string      `json:"level"`
	RequestId string      `json:"request_id"`
	Time      time.Time   `json:"time"`
	Source    string      `json:"source"`
	Data      interface{} `json:"data"`
}

func AppendToList(level string, requestId string, format string, source string, val ...interface{}) {
	id := uuid.New().String()

	var detail interface{}
	if format != "" {
		detail = fmt.Sprintf(format, val...)
	} else if len(val) > 1 {
		detail = val
	} else {
		detail = val[0]
	}
	body := map[string]interface{}{
		"sumoPayload": sumoPayload{
			Level:     level,
			RequestId: requestId,
			Time:      time.Now().UTC(),
			Source:    source,
			Data:      detail,
		},
		"host":     config.Host,
		"category": fmt.Sprintf("%s/%s", config.CategoryPrefix, sumoCategoryKey),
	}

	out, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	stringifyBody := string(out)

	var message = types.SendMessageBatchRequestEntry{
		Id:          &id,
		MessageBody: &stringifyBody,
		//MessageGroupId:         aws.String(requestId),
		//MessageDeduplicationId: aws.String(uuid.New().String()),
	}
	logMessageList = append(logMessageList, message)
}

func GetLogs() *[]sqs.SendMessageBatchInput {
	messageBatches := make([]sqs.SendMessageBatchInput, 0)
	if len(logMessageList) > 0 {

		for i := 0; i < len(logMessageList); i += batchSize {
			var listEntries []types.SendMessageBatchRequestEntry

			if len(logMessageList[i:]) < batchSize {
				listEntries = logMessageList[i:]
			} else {
				outerLimit := i + batchSize
				listEntries = logMessageList[i:outerLimit]
			}
			messageBatches = append(messageBatches, sqs.SendMessageBatchInput{
				Entries: listEntries,
			})
		}
	}
	return &messageBatches
}
