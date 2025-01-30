package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func AwsSsmGetParam(path string) (string, error) {
	ctx := context.TODO()
	param, err := AwsClientSsm.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(path),
		WithDecryption: PointerTrue(),
	})
	if err != nil {
		return "", err
	}

	return aws.ToString(param.Parameter.Value), nil
}
func AwsSsmGetParams(paths []string) ([]string, error) {
	ctx := context.TODO()
	params, err := AwsClientSsm.GetParameters(ctx, &ssm.GetParametersInput{
		Names:          paths,
		WithDecryption: PointerTrue(),
	})
	if err != nil {
		return nil, err
	}
	var values []string
	for _, param := range params.Parameters {
		values = append(values, aws.ToString(param.Value))
	}
	return values, nil

}

func PointerTrue() *bool {
	t := true
	return &t
}
