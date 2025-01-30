import boto3
import botocore

class AwsInitializer:
    def __init__(self, region_name='ap-northeast-2'):
        self.region_name = region_name
        self.ssm_client = None

    def init_aws(self):
        try:
            self.ssm_client = boto3.client('ssm', region_name=self.region_name)
            print("AWS SSM 초기화 완료")
        except botocore.exceptions.BotoCoreError as e:
            print(f"AWS 초기화 실패: {str(e)}")
            raise e

    def get_params(self, paths):
        try:
            response = self.ssm_client.get_parameters(Names=paths, WithDecryption=True)
            result = [param.get('Value', '') for param in response.get('Parameters', [])]
            return result
        except botocore.exceptions.BotoCoreError as e:
            print(f"SSM 파라미터 가져오기 실패: {str(e)}")
            return []

    def get_param(self, path):
        try:
            response = self.ssm_client.get_parameter(Name=path, WithDecryption=True)
            return response.get('Parameter', {}).get('Value', '')
        except botocore.exceptions.BotoCoreError as e:
            print(f"SSM 파라미터 가져오기 실패: {str(e)}")
            return None