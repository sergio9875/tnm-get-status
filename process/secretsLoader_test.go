package process

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// Define a mock struct to be used in your unit tests
type mockSecretsManagerClient struct {
	CallCount          *int
	FakeGetSecretValue func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(options *secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

var smCallCount = 0

func (m *mockSecretsManagerClient) GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(options *secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	return m.FakeGetSecretValue(ctx, params, optFns...)
}

//func TestTest(t *testing.T) {
//  secretHolder := &SecretIDHolder{
//    SecretID: "configuration/treasury/config",
//    Client: CreateSMClient(),
//  }
//  _ = secretHolder.LoadSecret()
//}

func TestCreateSMService(t *testing.T) {
	_ = os.Setenv("AWS_REGION", "eu-west-1")
	svc := CreateSMClient()
	if svc == nil {
		t.Errorf("TestCreateSMService: expected none nil")
	}
}

func TestSecretLoad(t *testing.T) {
	// Setup Test
	// mockSvc := &mockSecretsManagerClient{}

	smCallCount = 0
	var sh = &SecretIDHolder{
		SecretID: "123wq",
		Client: &mockSecretsManagerClient{
			FakeGetSecretValue: func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(options *secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
				smCallCount++
				switch smCallCount {
				case 1: // Success
					if *params.VersionStage == "AWSCURRENT" {
						return &secretsmanager.GetSecretValueOutput{
							SecretString: aws.String(`
{
  "port": 8000,
  "routes": [
    {
      "version": "/v5",
      "legacy": {
        "endpoint": "https://secure.3gdirectpay.com/API/v5/",
        "requests": [
          "createToken",
          "verifyToken"
        ]
      },
      "reroutes": {
        "basePath": "http://localhost:8800",
        "paths": {
          "chargeTokenCreditCard": "/chargeTokenCreditCard"
        }
      }
    }]
}
                `),
						}, nil
					}
				case 2: // error getting secret
					return nil, errors.New("something strange is a foot")
				case 3: // broken json
					if *params.VersionStage == "AWSCURRENT" {
						return &secretsmanager.GetSecretValueOutput{
							SecretString: aws.String(`}`),
						}, nil
					}
				case 4: // no secret value
					if *params.VersionStage == "AWSCURRENT" {
						return &secretsmanager.GetSecretValueOutput{
							SecretString: nil,
							SecretBinary: nil,
						}, nil
					}
				case 5:
					if *params.VersionStage == "AWSCURRENT" {
						return &secretsmanager.GetSecretValueOutput{
							SecretString: nil,
							SecretBinary: []byte("e30="),
						}, nil
					}
				case 6:
					if *params.VersionStage == "AWSCURRENT" {
						return &secretsmanager.GetSecretValueOutput{
							SecretString: nil,
							SecretBinary: []byte("e300="),
						}, nil
					}
				}
				return nil, errors.New("no such version")
			},
		},
	}
	// 1 Success
	_ = sh.LoadSecret()
	// 2 Error getting secret
	sm := sh.LoadSecret()
	if sm != nil {
		t.Errorf("TestSecretLoad failed getting secret, was expecting a nil model")
	}
	// 3 Can not Unmarshall
	sm = sh.LoadSecret()
	if sm != nil {
		t.Errorf("TestSecretLoad fail unmarshally, was expecting a nil model")
	}
	// 4 no secret value
	sm = sh.LoadSecret()
	if sm != nil {
		t.Errorf("TestSecretLoad fail no secret, was expecting a nil model")
	}
	// 5 binary secret
	sm = sh.LoadSecret()
	if sm == nil {
		t.Errorf("TestSecretLoad fail binary secret success, was expecting a none nil model")
	}
	// 6 binary secret bad base64
	sm = sh.LoadSecret()
	if sm != nil {
		t.Errorf("TestSecretLoad fail binary secret bad, was expecting a nil model")
	}
}

func TestMultiSecretLoad(t *testing.T) {
	// Setup Test
	// mockSvc := &mockSecretsManagerClient{}

	smCallCount = 0
	var sh = &SecretIDHolder{
		SecretID: "123wq",
		Client: &mockSecretsManagerClient{
			FakeGetSecretValue: func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(options *secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
				smCallCount++
				switch smCallCount {
				case 1, 3, 5: // Success
					if *params.VersionStage == "AWSCURRENT" {
						return &secretsmanager.GetSecretValueOutput{
							SecretString: aws.String(`
{
  "secrets": [ "fake_secrets" ],
  "port": 8000
}
                `),
						}, nil
					}
				case 2:
					if *params.VersionStage == "AWSCURRENT" {
						return &secretsmanager.GetSecretValueOutput{
							SecretString: aws.String(`
{
  "db": {
    "treasury": {
      "dialect": "mysql"
    }
  }
}
                `),
						}, nil
					}
				case 4: // broken inner secret value
					return nil, errors.New("something strange is a foot")
				case 6: // broken inner secret value
					if *params.VersionStage == "AWSCURRENT" {
						return &secretsmanager.GetSecretValueOutput{
							SecretString: aws.String(`
{
  "db": {
    "treasury": {
      "dialect": "mysql,
    }
  }
}
                `),
						}, nil
					}
				}
				return nil, errors.New("no such version")
			},
		},
	}
	// 1 + 2 Success
	_ = sh.LoadSecret()
	// 3 + 4 No inner secret found
	sm := sh.LoadSecret()
	if sm != nil {
		t.Errorf("TestSecretLoad fail unmarshally, was expecting a nil model")
	}
	// 5 + 6 Can not Unmarshall inner
	sm = sh.LoadSecret()
	if sm != nil {
		t.Errorf("TestSecretLoad fail unmarshally, was expecting a nil model")
	}
}
