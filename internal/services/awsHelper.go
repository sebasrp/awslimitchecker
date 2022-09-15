package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/servicequotas"
)

var conf *Config = &Config{}

type Config struct {
	Session       *session.Session
	DynamoDb      DynamodbClientInterface
	Eks           EksClientInterface
	Kinesis       KinesisClientInterface
	Rds           RdsClientInterface
	S3            S3ClientInterface
	ServiceQuotas SvcQuotaClientInterface
}

var InitializeConfig = initializeConfig

func initializeConfig(awsprofile string, region string) (*Config, error) {
	sess, err := createAwsSession(awsprofile, region)
	if err != nil {
		return &Config{}, fmt.Errorf("unable to create a session to aws with error: %v", err)
	}

	conf = &Config{
		Session:       &sess,
		DynamoDb:      dynamodb.New(&sess),
		Eks:           eks.New(&sess),
		Kinesis:       kinesis.New(&sess),
		Rds:           rds.New(&sess),
		S3:            s3.New(&sess),
		ServiceQuotas: servicequotas.New(&sess),
	}

	return conf, nil
}

func createAwsSession(awsprofile string, region string) (session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewSharedCredentials("", awsprofile)},
	)
	if err != nil {
		fmt.Errorf("Unable to create AWS session, %v", err)
	}
	return *sess, err
}
