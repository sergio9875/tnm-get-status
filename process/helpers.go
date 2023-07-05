package process

import (
	"encoding/json"
	log "malawi-getstatus/logger"
)

func (c *Controller) handleGetMessage(message string, messageData interface{}) error {
	if err := json.Unmarshal([]byte(message), &messageData); err != nil {
		log.Error(*c.requestId, "unable to retrieve message body: ", err.Error())
		return err
	}

	//log.Info(*c.requestId, "Successfully retrieved message: ", messageData)
	return nil
}
