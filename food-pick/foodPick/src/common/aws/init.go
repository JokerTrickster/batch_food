package aws

import (
	"context"
	"time"

	AwsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/ssm"

	"github.com/aws/aws-sdk-go-v2/aws"
)

var AwsClientSsm *ssm.Client
var awsClientSes *sesv2.Client
var awsClientS3 *s3.Client
var awsClientS3Uploader *manager.Uploader
var awsClientS3Downloader *manager.Downloader
var awsS3Signer *s3.PresignClient

type ImgType uint8

const (
	ImgTypeFood     = ImgType(0)
	ImgTypeCategory = ImgType(1)
	ImgTypeProfile  = ImgType(2)
)

type imgMetaStruct struct {
	bucket     func() string
	domain     func() string
	path       string
	width      int
	height     int
	expireTime time.Duration
}

func InitAws() error {
	var awsConfig aws.Config
	var err error

	awsConfig, err = AwsConfig.LoadDefaultConfig(context.TODO(), AwsConfig.WithRegion("ap-northeast-2"))
	if err != nil {
		return err
	}
	AwsClientSsm = ssm.NewFromConfig(awsConfig)
	awsClientSes = sesv2.NewFromConfig(awsConfig)
	awsClientS3 = s3.NewFromConfig(awsConfig)
	awsClientS3Uploader = manager.NewUploader(awsClientS3)
	awsClientS3Downloader = manager.NewDownloader(awsClientS3)
	awsClientSes = sesv2.NewFromConfig(awsConfig)
	awsS3Signer = s3.NewPresignClient(awsClientS3)
	err = InitAwsSes()
	if err != nil {
		return err
	}
	return nil
}
