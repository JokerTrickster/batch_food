import json
import datetime
import pytz  # ✅ 타임존 변환을 위한 라이브러리
from db_initializer import SessionLocal, openai_api_key,redis_client  # ✅ DB 및 Redis import
from langchain_openai import ChatOpenAI
from langchain.prompts import PromptTemplate
from langchain.schema.runnable import RunnableLambda
from sqlalchemy import text

# ✅ LangChain LLM 초기화
llm = ChatOpenAI(
    model_name="gpt-4",
    temperature=0.3,
    openai_api_key= openai_api_key  # ✅ OpenAI API 키는 AWS SSM에서 로드
)

# ✅ 성별+연령대별 많이 먹은 음식 조회 (최대 20개)
def get_popular_foods_by_age_gender(age, sex, meal_time):
    age_group = (age // 10) * 10  
    redis_key = f"foods:{age_group}:{sex}:{meal_time}"

    # ✅ Redis에서 캐시 확인
    cached_foods = redis_client.get(redis_key)
    if cached_foods:
        print(f"🔹 Redis 캐시 히트: {redis_key}")
        return json.loads(cached_foods)

    print(f"🔹 Redis 캐시 없음: {redis_key}. MySQL에서 조회 후 저장.")

    session = SessionLocal()
    try:
        query = text("""
        SELECT f.name
        FROM food_histories fh
        JOIN foods f ON fh.food_id = f.id
        JOIN users u ON fh.user_id = u.id
        WHERE u.sex = :sex 
        AND (YEAR(CURDATE()) - YEAR(u.birth)) BETWEEN :age_min AND :age_max
        GROUP BY f.name
        ORDER BY COUNT(fh.id) DESC
        LIMIT 20;
        """)

        result = session.execute(query, {
            "sex": sex, 
            "age_min": age_group, 
            "age_max": age_group + 9
        }).fetchall()

        popular_foods = [row[0] for row in result] if result else []
        
        # ✅ Redis에 저장 (6시간 유지)
        redis_client.setex(redis_key, 21600, json.dumps(popular_foods))
        return popular_foods
    finally:
        session.close()

# ✅ LLM을 사용해 음식 1개 추천
def recommend_food_by_llm(age_group, sex, meal_time, food_list):
    prompt_template = PromptTemplate(
        template="""
        {age_group}대 {sex} 사용자들이 최근 {meal_time}에 {foods}을(를) 많이 먹었습니다.
        취향을 분석하여 {meal_time}에 가장 어울리는 음식 하나를 추천해 주세요.

        - 반드시 음식 1개만 추천하세요.
        - 다른 추가 설명은 제공하지 마세요.
        - 응답 예시:
          "김치찌개"
        """,
        input_variables=["age_group", "sex", "meal_time", "foods"]
    )

    def extract_content(ai_message):
        return ai_message.content.strip()

    llm_chain = prompt_template | llm | RunnableLambda(extract_content)

    response = llm_chain.invoke({
        "age_group": age_group,
        "sex": sex,
        "meal_time": meal_time,
        "foods": ", ".join(food_list) if food_list else "추천할 음식이 없습니다."
    })

    # ✅ 추천 결과 Redis에 저장 (6시간 유지)
    redis_key = f"foods:{age_group}:{sex}:{meal_time}"
    redis_client.setex(redis_key, 21600, response)

    return response

# ✅ 모든 연령대(10대~60대) 실행하는 함수
def generate_food_recommendations():
    for age_group in range(0, 70, 10):  # 0대 ~ 60대
        for sex in ["male", "female"]:
            for meal_time in ["lunch", "dinner"]:
                redis_key = f"foods:{age_group}:{sex}:{meal_time}"
                
                # ✅ Redis 캐시 확인 후 스킵
                if redis_client.exists(redis_key):
                    print(f"✅ 이미 추천된 음식 있음: {redis_key}")
                    continue

                # ✅ 인기 음식 조회
                popular_foods = get_popular_foods_by_age_gender(age_group, sex, meal_time)

                # ✅ LLM을 사용해 음식 1개 추천
                if not popular_foods:
                    recommended_food = recommend_food_by_llm(age_group, sex, meal_time, [])
                else:
                    recommended_food = recommend_food_by_llm(age_group, sex, meal_time, popular_foods)

                print(f"🍽️ {age_group}대 {sex} {meal_time} 추천: {recommended_food}")
                
# ✅ 현재 시간을 기준으로 meal_time (lunch/dinner) 선택
def get_current_meal_time():
    """ 현재 시간(UTC)을 KST로 변환한 후 lunch 또는 dinner를 결정하는 함수 """
    kst = pytz.timezone('Asia/Seoul')
    current_time = datetime.datetime.now(pytz.utc).astimezone(kst)  # ✅ UTC → KST 변환
    current_hour = current_time.hour

    # ✅ 11:00~15:00 → lunch / 17:00~22:00 → dinner
    if 11 <= current_hour < 15:
        return "lunch"
    elif 17 <= current_hour < 22:
        return "dinner"
    else:
        return "lunch"  # ✅ 기본적으로 lunch를 선택

# ✅ Lambda 핸들러 (메인 실행 함수)
def lambda_handler(event, context):
    print("🔹 Lambda 실행 시작")

    # ✅ 음식 추천 배치 실행
    generate_food_recommendations()

    print("🔹 Lambda 실행 완료")
    return {
        "statusCode": 200,
        "body": json.dumps({
            "message": "음식 추천 배치 실행 완료"
        }, ensure_ascii=False)
    }


# ✅ 프로그램 실행
if __name__ == "__main__":
    response = lambda_handler({}, None)
    print(response)
