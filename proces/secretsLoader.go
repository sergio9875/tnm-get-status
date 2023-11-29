package proces

// If you need more information about configurations or implementing the sample code, visit the AWS docs:
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/setting-up.html

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/config"
	"malawi-getstatus/models"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	log "malawi-getstatus/logger"
)

type SecretManager interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(options *secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

// SecretIDHolder model
type SecretIDHolder struct {
	SecretID string
	Client   SecretManager
}

func (sh *SecretIDHolder) getSecret() (*string, error) {
	log.Println("SYSTEM", "loading secret: ", sh.SecretID)
	//Create a Secrets Manager client
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(sh.SecretID),
		VersionStage: aws.String("AWSCURRENT"),
	}

	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html

	result, err := sh.Client.GetSecretValue(context.TODO(), input)
	if err != nil {
		log.Error("SYSTEM", "error retrieving secret", err)
		return nil, err
	}

	// Decrypts secret using the associated KMS CMK.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	var secretString string
	if result.SecretString != nil {
		secretString = *result.SecretString
	} else {
		if result.SecretBinary != nil {
			decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
			len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
			if err != nil {
				log.Error("SYSTEM", "Base64 Decode Error:", err)
				return nil, err
			}
			secretString = string(decodedBinarySecretBytes[:len])
		}
	}

	return &secretString, nil
}

// CreateSMService func
func CreateSMClient() *secretsmanager.Client {
	cfg, _ := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	return secretsmanager.NewFromConfig(cfg)
}

func (sh *SecretIDHolder) loadSecrets(secrets *models.SecretModel) *models.SecretModel {
	tempSH := &SecretIDHolder{
		SecretID: "",
		Client:   sh.Client,
	}
	for _, secretName := range secrets.Secrets {
		log.Infof("SYSTEM", "Loading inner secret: %s", secretName)
		tempSH.SecretID = secretName
		secretValue, err := tempSH.getSecret()
		if err != nil {
			log.Error("SYSTEM", "Inner secret error: "+err.Error())
			return nil
		}
		//log.Println("SYSTEM", "[" + *secretValue + "]")
		secrets = secrets.Merge(secretValue)
	}
	return secrets
}

// LoadSecret func
func (sh *SecretIDHolder) LoadSecret() *models.SecretModel {
	secretValue, error := sh.getSecret()
	if error != nil {
		log.Error("SYSTEM", error.Error())
		return nil
	}
	//var secret *models.SecretModel
	secret := &models.SecretModel{}
	err := json.Unmarshal([]byte(*secretValue), secret)

	if err != nil {
		return nil
	}

	// Need to implement multi secret loading here
	// secret is the parent of all inner secrets,
	// so now we need to load children then check
	// what part of the parent it is then merge
	if secret.Secrets != nil && len(secret.Secrets) > 0 {
		secret = sh.loadSecrets(secret)
	}
	return secret
}
