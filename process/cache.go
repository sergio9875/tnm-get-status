package process

import (
	"context"
	"github.com/mitchellh/mapstructure"
	log "malawi-getstatus/logger"
	"malawi-getstatus/models"
	"malawi-getstatus/utils"
	"strconv"
)

func (c *Controller) getCache(ctx context.Context, redisKey string, cache *models.IncomingRequest) error {
	var err error
	var result map[string]string

	log.Info(*c.requestId, "trying to retrieve cache from redis: ", redisKey)

	if result, err = (*c.cacheClient).HGetAll(ctx, redisKey); err != nil {
		log.Error(*c.requestId, "unable to retrieve cache from redis: ", err.Error())
		return err
	}

	if err = mapstructure.Decode(result, &cache); err != nil {
		return err
	}

	log.Info(*c.requestId, "Successfully retrieved cache from redis: ")
	return nil
}

func (c *Controller) updateRedisCounter(ctx context.Context, redisKey string, counter string) error {
	log.Info(*c.requestId, "trying to update cache from redis for counter: ", redisKey, counter)

	newCounter := utils.SafeAtoi(counter, 0) + 1

	if err := (*c.cacheClient).HSet(ctx, redisKey, "counter", strconv.Itoa(newCounter)); err != nil {
		log.Error(*c.requestId, "error while trying to update redis counter", err.Error())
		return err
	}
	return nil
}
