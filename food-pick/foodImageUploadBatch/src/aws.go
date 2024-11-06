package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	AwsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

var AwsClientSsm *ssm.Client

func InitAws() error {

	awsConfig, err := AwsConfig.LoadDefaultConfig(context.TODO(),
		AwsConfig.WithRegion("ap-northeast-2"))
	if err != nil {
		return err
	}
	AwsClientSsm = ssm.NewFromConfig(awsConfig)

	fmt.Println("aws ssm 초기화 완료")
	return nil
}

func AwsGetParams(paths []string) ([]string, error) {
	ctx := context.TODO()
	// get ssm param
	params, err := AwsClientSsm.GetParameters(ctx, &ssm.GetParametersInput{
		Names:          paths,
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}
	result := make([]string, len(paths))
	for i, path := range paths {
		val := ""
		for _, parameter := range params.Parameters {
			if strings.Contains(aws.ToString(parameter.ARN), path) {
				val = aws.ToString(parameter.Value)
				break
			}
		}
		result[i] = val
	}
	return result, nil
}

func AwsGetParam(path string) (string, error) {
	ctx := context.TODO()
	// get ssm param
	param, err := AwsClientSsm.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(path),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", err
	}
	return aws.ToString(param.Parameter.Value), nil
}
