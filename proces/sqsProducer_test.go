package proces

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"malawi-getstatus/utils"
	"os"
	"testing"
)

var sqsCallCount = 0

type mockSqsManagerClient struct {
	FakeGetQueueUrl func(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)
	FakeSendMessage func(ctx context.Context,
		params *sqs.SendMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

func (m *mockSqsManagerClient) GetQueueUrl(ctx context.Context, params *sqs.GetQueueUrlInput, optFns ...func(options *sqs.Options)) (*sqs.GetQueueUrlOutput, error) {
	return m.FakeGetQueueUrl(ctx, params, optFns...)
}

func (m *mockSqsManagerClient) SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(options *sqs.Options)) (*sqs.SendMessageOutput, error) {
	return m.FakeSendMessage(ctx, params, optFns...)
}

func TestCreateSQSClient(t *testing.T) {
	_ = os.Setenv("AWS_REGION", "eu-west-1")
	svc := CreateSQSClient()
	if svc == nil {
		t.Errorf("TestCreateSQSClient: expected none nil")
	}
}

func TestNewSQSProducerFromQueueUrl(t *testing.T) {
	fakeClient := &mockSqsManagerClient{
		FakeGetQueueUrl: func(_ context.Context, params *sqs.GetQueueUrlInput, optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error) {
			return &sqs.GetQueueUrlOutput{
				QueueUrl: utils.StringPtr("fakeQueueUrl"),
			}, nil
		},
	}
	sqsProducer, err := NewSQSProducerFromUrl(context.TODO(), fakeClient, utils.StringPtr("fakeQueueUrl"))
	if sqsProducer == nil {
		t.Errorf("Was not expecting a nil producer on step %d", sqsCallCount)
	}
	if err != nil {
		t.Errorf("Was not expecting an error on step %d", sqsCallCount)
	}
	if sqsProducer.queueName == nil {
		t.Errorf("Was not expecting a nil producer.queueName on step %d", sqsCallCount)
	}
	if *sqsProducer.queueName != "fakeQueueUrl" {
		t.Errorf("producer.queueName expecting[%s], got[%s] on step %d", "fakeQueueUrl",
			*sqsProducer.queueName, sqsCallCount)
	}
	if sqsProducer.queueUrl == nil {
		t.Errorf("Was not expecting a nil producer.queueURL on step %d", sqsCallCount)
	}
	if *sqsProducer.queueUrl != "fakeQueueUrl" {
		t.Errorf("producer.queueName expecting[%s], got[%s] on step %d", "fakeQueueUrl",
			*sqsProducer.queueUrl, sqsCallCount)
	}
}

func TestNewSQSProducerFromName(t *testing.T) {
	sqsCallCount = 1
	fakeClient := &mockSqsManagerClient{
		FakeGetQueueUrl: func(_ context.Context, params *sqs.GetQueueUrlInput, optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error) {
			switch sqsCallCount {
			case 1: // all good
				return &sqs.GetQueueUrlOutput{
					QueueUrl: utils.StringPtr("fakeQueueUrl"),
				}, nil
			case 2: // error
				return nil, errors.New("something strange is a foot")
			}
			return nil, errors.New("invalid step in test")
		},
	}
	sqsProducer, err := NewSQSProducerFromName(context.TODO(), fakeClient, utils.StringPtr("fakeQueueName"))
	if sqsProducer == nil {
		t.Errorf("Was not expecting a nil producer on step %d", sqsCallCount)
	}
	if err != nil {
		t.Errorf("Was not expecting an error on step %d", sqsCallCount)
	}
	if sqsProducer.queueName == nil {
		t.Errorf("Was not expecting a nil producer.queueName on step %d", sqsCallCount)
	}
	if *sqsProducer.queueName != "fakeQueueName" {
		t.Errorf("producer.queueName expecting[%s], got[%s] on step %d", "fakeQueueName",
			*sqsProducer.queueName, sqsCallCount)
	}
	if sqsProducer.queueUrl == nil {
		t.Errorf("Was not expecting a nil producer.queueURL on step %d", sqsCallCount)
	}
	if *sqsProducer.queueUrl != "fakeQueueUrl" {
		t.Errorf("producer.queueName expecting[%s], got[%s] on step %d", "fakeQueueUrl",
			*sqsProducer.queueUrl, sqsCallCount)
	}

	sqsCallCount++
	sqsProducer, err = NewSQSProducerFromName(context.TODO(), fakeClient, utils.StringPtr("fakeQueueName"))
	if sqsProducer != nil {
		t.Errorf("Was not expecting a producer on step %d", sqsCallCount)
	}
	if err == nil {
		t.Errorf("Was not expecting a nil error on step %d", sqsCallCount)
	}
	if err != nil && err.Error() != "something strange is a foot" {
		t.Errorf("Incorrect error message[%s] on step %d", err.Error(), sqsCallCount)
	}

}

func TestSQSProducer_GetQueueURL(t *testing.T) {
	fakeClient := &mockSqsManagerClient{
		FakeGetQueueUrl: func(_ context.Context, params *sqs.GetQueueUrlInput, optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error) {
			return &sqs.GetQueueUrlOutput{
				QueueUrl: utils.StringPtr("fakeQueueUrl"),
			}, nil
		},
	}
	producer := &SQSProducer{
		client:    fakeClient,
		queueName: utils.StringPtr("mock"),
		queueUrl:  utils.StringPtr("mockUrl"),
	}
	result, _ := (*producer).GetQueueURL(context.TODO(), &sqs.GetQueueUrlInput{QueueName: utils.StringPtr("mock")})
	if *result.QueueUrl != "fakeQueueUrl" {
		t.Errorf("GetQueueUrl test failed, expected[%s], got[%s]", "fakeQueueUrl", *result.QueueUrl)
	}
}

func TestSQSProducer_GetQueueURLError(t *testing.T) {
	fakeClient := &mockSqsManagerClient{
		FakeGetQueueUrl: func(_ context.Context, params *sqs.GetQueueUrlInput, optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error) {
			return nil, errors.New("something strange is a foot")
		},
	}
	producer := &SQSProducer{
		client:    fakeClient,
		queueName: utils.StringPtr("mock"),
		queueUrl:  utils.StringPtr("mockUrl"),
	}
	_, err := (*producer).GetQueueURL(context.TODO(), &sqs.GetQueueUrlInput{QueueName: utils.StringPtr("mock")})
	if err == nil {
		t.Error("GetQueueUrlError test failed, expected[error], got[nil]")
	}
}

func TestSQSProducer_SendMsg(t *testing.T) {
	fakeClient := &mockSqsManagerClient{
		FakeSendMessage: func(_ context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
			return &sqs.SendMessageOutput{
				MessageId: utils.StringPtr("123"),
			}, nil
		},
	}
	producer := &SQSProducer{
		client:    fakeClient,
		queueName: utils.StringPtr("mock"),
		queueUrl:  utils.StringPtr("mockUrl"),
	}
	result, _ := (*producer).SendMsg(context.TODO(), &sqs.SendMessageInput{MessageBody: utils.StringPtr("{\"msg\":\"hello!\"}")})
	if *result.MessageId != "123" {
		t.Errorf("SendMsg test failed, expected[%s], got[%s]", "123", *result.MessageId)
	}
}

func TestSQSProducer_SendMsgError(t *testing.T) {
	fakeClient := &mockSqsManagerClient{
		FakeSendMessage: func(_ context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
			return nil, errors.New("something strange is a foot")
		},
	}
	producer := &SQSProducer{
		client:    fakeClient,
		queueName: utils.StringPtr("mock"),
		queueUrl:  utils.StringPtr("mockUrl"),
	}
	_, err := (*producer).SendMsg(context.TODO(), &sqs.SendMessageInput{MessageBody: utils.StringPtr("{\"msg\":\"hello!\"}")})
	if err == nil {
		t.Error("SendMsg test failed, expected[error], got[nil]")
	}
}
