package logger

import (
	"bou.ke/monkey"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/stretchr/testify/assert"
	"malawi-getstatus/utils"
	"testing"
	"time"
)

func Test_findIndex(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "good",
			args: args{
				key: DEBUG,
			},
			want: 1,
		},
		{
			name: "good unknown",
			args: args{
				key: "UNKNOWN",
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findIndex(tt.args.key); got != tt.want {
				t.Errorf("findIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitOverwatch(t *testing.T) {
	type args struct {
		sumoConfig *SumoConfig
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "good",
			args: args{
				sumoConfig: &SumoConfig{
					LogLevel: WARN,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitOverwatch("123", tt.args.sumoConfig)

			assert.NotNil(t, config)
			assert.Equal(t, hierarchyIndex, 3)
			assert.NotNil(t, logMessageList)
		})
	}
}

func TestCloseOverwatch(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "good"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logMessageList = append(logMessageList, types.SendMessageBatchRequestEntry{})
			assert.Equal(t, len(logMessageList), 1)
			CloseOverwatch()
			assert.Equal(t, len(logMessageList), 0)
		})
	}
}

func TestAddToSumo(t *testing.T) {
	wayback := time.Date(1974, time.May, 19, 1, 2, 3, 4, time.UTC)
	patch := monkey.Patch(time.Now, func() time.Time { return wayback })
	defer patch.Unpatch()
	type args struct {
		sumoLevel string
		level     string
		requestId string
		format    string
		source    string
		a         []interface{}
	}
	tests := []struct {
		name   string
		args   args
		expect int
		want   *string
	}{
		{
			name: "good no format",
			args: args{
				DEBUG,
				DEBUG,
				"123",
				"",
				"source1",
				[]interface{}{"qwerty"},
			},
			expect: 1,
		},
		{
			name: "good with format",
			args: args{
				DEBUG,
				DEBUG,
				"123",
				"%s",
				"source1",
				[]interface{}{"qwerty"},
			},
			expect: 1,
		},
		{
			name: "good no log",
			args: args{
				INFO,
				DEBUG,
				"123",
				"%s",
				"source1",
				[]interface{}{"qwerty"},
			},
			expect: 0,
		},
		{
			name: "good multiple",
			args: args{
				TRACE,
				DEBUG,
				"123",
				"",
				"source1",
				[]interface{}{"qwerty", 1, 1.2},
			},
			expect: 1,
			want:   utils.StringPtr("{\"category\":\"/connector/safaricom/callback\",\"host\":\"\",\"sumoPayload\":{\"level\":\"DEBUG\",\"request_id\":\"123\",\"time\":\"1974-05-19T01:02:03.000000004Z\",\"source\":\"source1\",\"data\":[\"qwerty\",1,1.2]}}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitOverwatch("connector/safaricom/callback", &SumoConfig{
				LogLevel: tt.args.sumoLevel,
			})
			AddToSumo(tt.args.level, tt.args.requestId, tt.args.format, tt.args.source, tt.args.a...)
			assert.Equal(t, len(logMessageList), tt.expect)
			if tt.want != nil {
				assert.Equal(t, *tt.want, *logMessageList[0].MessageBody)
			}
			CloseOverwatch()
		})
	}
}

func TestGetLogs(t *testing.T) {
	type args struct {
		sumoLevel string
		level     string
		requestId string
		format    string
		source    string
		a         []interface{}
	}
	tests := []struct {
		name   string
		args   args
		repeat int
		expect int
	}{
		{
			name: "good 11 entries equals 2",
			args: args{
				DEBUG,
				DEBUG,
				"123",
				"",
				"source1",
				[]interface{}{"qwerty"},
			},
			repeat: 11,
			expect: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitOverwatch("123", &SumoConfig{
				LogLevel: tt.args.sumoLevel,
			})
			for i := 0; i < tt.repeat; i++ {
				AddToSumo(tt.args.level, tt.args.requestId, tt.args.format, tt.args.source, tt.args.a...)
			}
			assert.Equal(t, len(logMessageList), tt.repeat)
			logs := GetLogs()
			assert.Equal(t, len(*logs), tt.expect)
			CloseOverwatch()
		})
	}
}
