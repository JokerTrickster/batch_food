import json
from aws_initializer import AwsInitializer

# 전역 AWS 초기화 객체
aws_init = AwsInitializer()
aws_initialized = False


def lambda_handler(event, context):
    global aws_initialized

    # AWS 초기화 - 처음 한 번만 수행
    if not aws_initialized:
        aws_init.init_aws()
        aws_initialized = True

    param_value = aws_init.get_param('dev_common_mysql_host')

    if param_value is None:
        return {
            "statusCode": 500,
            "body": json.dumps({
                "message": "SSM 파라미터를 가져오는 데 실패했습니다."
            })
        }

    return {
        "statusCode": 200,
        "body": json.dumps({
            "message": "hello world11111",
            "parameter_value": param_value
        }),
    }

if __name__ == "__main__":

    aws_init.init_aws()
    response = lambda_handler({}, None)
    print(response['body'])
