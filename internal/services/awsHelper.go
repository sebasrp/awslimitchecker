package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/aws/aws-sdk-go/service/sns"
)

var conf *Config = &Config{}

type Config struct {
	Session        *session.Session
	Acm            AcmClientInterface
	Autoscaling    AutoscalingClientInterface
	Cloudformation CloudformationClientInterface
	DynamoDb       DynamodbClientInterface
	Ec2            Ec2ClientInterface
	Eks            EksClientInterface
	ElastiCache    ElastiCacheClientInterface
	Elb            ElbClientInterface   // for classic load balancers
	Elbv2          Elbv2ClientInterface // for ALB, NLB load balancers
	Iam            IamClientInterface
	Kinesis        KinesisClientInterface
	Rds            RdsClientInterface
	S3             S3ClientInterface
	ServiceQuotas  SvcQuotaClientInterface
	Sns            SnsClientInterface
}

var InitializeConfig = initializeConfig

func initializeConfig(region string) error {
	sess, err := createAwsSession(region)
	if err != nil {
		return fmt.Errorf("unable to create a session to aws with error: %v", err)
	}

	conf = &Config{
		Session:        &sess,
		Acm:            acm.New(&sess),
		Autoscaling:    autoscaling.New(&sess),
		Cloudformation: cloudformation.New(&sess),
		DynamoDb:       dynamodb.New(&sess),
		Ec2:            ec2.New(&sess),
		Eks:            eks.New(&sess),
		ElastiCache:    elasticache.New(&sess),
		Elb:            elb.New(&sess),   // for classic load balancers
		Elbv2:          elbv2.New(&sess), // for ALB and NLB load balancers
		Iam:            iam.New(&sess),
		Kinesis:        kinesis.New(&sess),
		Rds:            rds.New(&sess),
		S3:             s3.New(&sess),
		ServiceQuotas:  servicequotas.New(&sess),
		Sns:            sns.New(&sess),
	}

	return conf, nil
}

func createAwsSession(region string) (session.Session, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		fmt.Printf("Unable to create AWS session, %v", err)
	}
	return *sess, err
}
