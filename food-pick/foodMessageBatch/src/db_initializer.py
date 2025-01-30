import redis
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
from aws_initializer import AwsInitializer

# ✅ AWS 초기화 객체 생성
aws_init = AwsInitializer()

def load_ssm_parameters():
    """ AWS SSM에서 DB 및 Redis 관련 환경 변수 로드 """
    aws_init.init_aws()

    params = aws_init.get_params([
        "dev_food_mysql_user",
        "dev_food_mysql_password",
        "dev_common_mysql_host",
        "dev_common_mysql_port",
        "dev_food_mysql_db",
        "dev_common_openai_api_key",
        "dev_common_redis_host",
        "dev_common_redis_port",
        "dev_food_redis_user",
        "dev_food_redis_password"
    ])

    db_uri = f"mysql+pymysql://{params[9]}:{params[8]}@{params[0]}:{params[1]}/{params[5]}"
    return {
        "DB_URI": db_uri,
        "OPENAI_API_KEY": params[2],
        "REDIS_HOST": params[3],
        "REDIS_PORT": params[4],
        "REDIS_USERNAME": params[7],
        "REDIS_PASSWORD": params[6]
    }

# ✅ 환경 변수 로드
ssm_params = load_ssm_parameters()

# ✅ DB 연결 설정
engine = create_engine(ssm_params["DB_URI"], connect_args={"init_command": "SET time_zone = '+00:00'"})
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

openai_api_key = ssm_params["OPENAI_API_KEY"]

# ✅ Redis 연결 설정
redis_client = redis.Redis(
    host=ssm_params["REDIS_HOST"],
    port=ssm_params["REDIS_PORT"],
    username=ssm_params["REDIS_USERNAME"] if ssm_params["REDIS_USERNAME"] else None,
    password=ssm_params["REDIS_PASSWORD"] if ssm_params["REDIS_PASSWORD"] else None,
    decode_responses=True
)

# ✅ Redis 연결 테스트
try:
    redis_client.ping()
    print("✅ Redis 연결 성공!")
except redis.AuthenticationError:
    print("❌ Redis 인증 실패! 유저명과 비밀번호를 확인하세요.")
except Exception as e:
    print(f"❌ Redis 연결 오류: {e}")
