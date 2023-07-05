package redis

import (
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"

	"context"
	"log"
	"malawi-getstatus/models"
	"malawi-getstatus/utils"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	redisHost string
	redisPort string
)

func TestMain(m *testing.M) {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	parts := strings.Split(mr.Addr(), ":")
	redisHost = parts[0]
	redisPort = parts[1]

	code := m.Run()
	mr.Close()
	os.Exit(code)
}

func TestNewCache(t *testing.T) {
	port, _ := strconv.Atoi(redisPort)
	db := 0
	dbP := &db
	type args struct {
		config *models.Cache
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Good test",
			args{
				&models.Cache{
					Type:     utils.StringPtr("redis"),
					Host:     utils.StringPtr(redisHost),
					Port:     &port,
					Password: utils.StringPtr(""),
					Database: dbP,
				},
			},
			false,
		},
		{
			"Type Error",
			args{
				&models.Cache{
					Type:     utils.StringPtr("redisV3"),
					Host:     utils.StringPtr(redisHost),
					Port:     &port,
					Password: utils.StringPtr(""),
					Database: dbP,
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCache(tt.args.config)
			if (err != nil) != tt.wantErr {
				log.Printf("%v \n", tt.args)
				t.Errorf("NewCache() 1 error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("NewCache() 2 got = %v, wantErr %v", got, tt.wantErr)
			}
		})
	}
}

func Test_holder_HSet(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := &Holder{db}

	args := make(map[string]string, 0)
	args["field"] = "0"

	mock.ExpectHSet("key", args).SetVal(0)

	result := repo.HSet(context.TODO(), "key", "field", "0")
	assert.Nil(t, result)
}

func Test_holder_Expire(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := &Holder{db}

	args := make(map[string]string, 0)
	args["field"] = "0"

	timeNow := time.Second

	mock.ExpectExpire("key", timeNow).SetVal(true)

	result := repo.Expire(context.TODO(), "key", timeNow)
	assert.Nil(t, result)
}
