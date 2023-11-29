package proces

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	log "malawi-getstatus/logger"
	"os"
)

// SqsApi defines the interface for the GetQueueUrl and SendMessage functions.
// We use this interface to test the functions using a mocked service.
type SqsApi interface {
	GetQueueUrl(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

	SendMessage(ctx context.Context,
		params *sqs.SendMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

// SQSProducer defines the object in which the sqs reference is held
type SQSProducer struct {
	queueName *string
	queueUrl  *string
	client    SqsApi
}

// GetQueueURL gets the URL of an Amazon SQS queue.
// Inputs:
//
//	c is the context of the method call, which includes the AWS Region.
//	api is the interface that defines the method call.
//	input defines the input arguments to the service call.
//
// Output:
//
//	If success, a GetQueueUrlOutput object containing the result of the service call and nil.
//	Otherwise, nil and an error from the call to GetQueueUrl.
func (sqs *SQSProducer) GetQueueURL(c context.Context, input *sqs.GetQueueUrlInput, optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error) {
	return sqs.client.GetQueueUrl(c, input, optFns...)
}

// SendMsg sends a message to an Amazon SQS queue.
// Inputs:
//
//	c is the context of the method call, which includes the AWS Region.
//	api is the interface that defines the method call.
//	input defines the input arguments to the service call.
//
// Output:
//
//	If success, a SendMessageOutput object containing the result of the service call and nil.
//	Otherwise, nil and an error from the call to SendMessage.
func (sqs *SQSProducer) SendMsg(c context.Context, input *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	input.QueueUrl = sqs.queueUrl
	return sqs.client.SendMessage(c, input, optFns...)
}

func CreateSQSClient() *sqs.Client {
	cfg, _ := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
	)

	// Configure a client with debug logging enabled
	if log.ValidateAgainstConfiguredLogLevel(log.TRACE) {
		cfg.ClientLogMode = aws.LogRequestWithBody | aws.LogResponse
	}

	return sqs.NewFromConfig(cfg)
}

func NewSQSProducerFromName(ctx context.Context, api SqsApi, queueName *string) (*SQSProducer, error) {
	// Get URL of queue
	gQInput := &sqs.GetQueueUrlInput{
		QueueName: queueName,
	}

	result, err := api.GetQueueUrl(ctx, gQInput)
	if err != nil {
		return nil, err
	}

	return &SQSProducer{
		queueName: queueName,
		queueUrl:  result.QueueUrl,
		client:    api,
	}, nil
}

func NewSQSProducerFromUrl(_ context.Context, api SqsApi, queueUrl *string) (*SQSProducer, error) {
	return &SQSProducer{
		queueName: queueUrl,
		queueUrl:  queueUrl,
		client:    api,
	}, nil
}
