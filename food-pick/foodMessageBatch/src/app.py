import json
import pytz  # ✅ 타임존 변환을 위한 라이브러리
from db_initializer import SessionLocal, openai_api_key, redis_client  # ✅ DB 및 Redis import
from langchain_openai import ChatOpenAI
from datetime import datetime, timedelta
from langchain.prompts import PromptTemplate
from langchain.schema.runnable import RunnableLambda
from sqlalchemy import text
import requests

# ✅ LangChain LLM 초기화
llm = ChatOpenAI(
    model_name="gpt-4",
    temperature=0.3,
    openai_api_key= openai_api_key  # ✅ OpenAI API 키는 AWS SSM에서 로드
)

       
# ✅ 사용자의 나이와 성별 조회 함수
def get_user_info(user_id):
    session = SessionLocal()
    try:
        query = text("""
        SELECT YEAR(CURDATE()) - YEAR(u.birth) AS age, u.sex
        FROM users u
        WHERE u.id = :user_id;
        """)
        result = session.execute(query, {"user_id": user_id}).fetchone()
        return {"age": result[0], "sex": result[1]} if result else None
    finally:
        session.close()
# ✅ 특정 시간대(점심/저녁)에 가장 많이 먹은 음식 TOP 10 조회 함수
def get_user_history(user_id, meal_time):
    session = SessionLocal()
    try:
        # ✅ 점심 & 저녁 시간 필터 설정
        if meal_time == "lunch":
            start_hour = 11
            end_hour = 15
        else:  # dinner
            start_hour = 16
            end_hour = 23

        query = text("""
        SELECT f.name, COUNT(*) AS eat_count
        FROM food_histories fh
        JOIN foods f ON fh.food_id = f.id
        WHERE fh.user_id = :user_id
        AND HOUR(CONVERT_TZ(fh.created_at, '+00:00', '+09:00')) BETWEEN :start_hour AND :end_hour
        GROUP BY f.name
        ORDER BY eat_count DESC 
        LIMIT 10;
        """)

        result = session.execute(query, {
            "user_id": user_id,
            "start_hour": start_hour,
            "end_hour": end_hour
        }).fetchall()

        return [row[0] for row in result] if result else []
    finally:
        session.close()

# ✅ 현재 시간을 기준으로 meal_time (lunch/dinner) 선택
def get_current_meal_time():
    kst = pytz.timezone('Asia/Seoul')
    current_time = datetime.now(pytz.utc).astimezone(kst)
    current_hour = current_time.hour

    if 11 <= current_hour < 15:
        return "lunch"
    elif 16 <= current_hour < 23:
        return "dinner"
    else:
        return "lunch"

# ✅ LLM을 활용한 음식 추천 함수 (최근 먹은 음식 리스트 기반)
def generate_food_from_history(user_name, meal_time, food_list):
    prompt_template = PromptTemplate(
        template="""
        {user_name} 님이 최근 {meal_time}에 {foods}을(를) 드셨습니다.
        취향을 분석하여 {meal_time}에 어울리는 새로운 음식 하나를 추천해 주세요.

        - 반드시 음식 1개만 추천하세요.
        - 다른 추가 설명은 제공하지 마세요.
        - 응답 예시:
          "김치찌개"
        """,
        input_variables=["user_name", "meal_time", "foods"]
    )

    llm_chain = prompt_template | llm | RunnableLambda(lambda x: x.content.strip())

    response = llm_chain.invoke({
        "user_name": user_name,
        "meal_time": meal_time,
        "foods": ", ".join(food_list)
    })

    return response.strip()


# ✅ LangChain을 활용한 음식 추천 함수
def recommend_food_for_user(user_id,user_name, meal_time):
    """ 사용자의 최근 음식 데이터를 기반으로 LangChain LLM을 활용한 음식 추천 """
    user_foods = get_user_history(user_id, meal_time)
    # ✅ 1️⃣ 사용자의 음식 기록이 있는 경우 → LLM을 사용하여 추천
    if user_foods:
        recommended_food = generate_food_from_history(user_name, meal_time, user_foods)
        if meal_time == "lunch":
            return f"{user_name} 님, 푸드픽에서 점심 메뉴로 {recommended_food}을(를) 추천드립니다!"
        else:
            return f"{user_name} 님, 푸드픽에서 저녁 메뉴로 {recommended_food}을(를) 추천드립니다!"

    # ✅ 2️⃣ 사용자의 음식 기록이 없는 경우 → Redis에서 연령대 + 성별 + 시간대별 추천 조회
    user_info = get_user_info(user_id)
    if user_info:
        age_group = (user_info["age"] // 10) * 10  # 10대, 20대, 30대 등으로 그룹화
        sex = user_info["sex"]

        # ✅ Redis 키: 연령대 + 성별 + 시간대별 캐싱된 추천
        age_sex_cache_key = f"foods:{age_group}:{sex}:{meal_time}"
        cached_food = redis_client.get(age_sex_cache_key)
        if meal_time == "lunch":
            return f"{user_name} 님, 푸드픽에서 점심 메뉴로 {cached_food}을(를) 추천드립니다!"
        else:
            return f"{user_name} 님, 푸드픽에서 저녁 메뉴로 {cached_food}을(를) 추천드립니다!"

    return f"{user_name} 님, {meal_time} 추천할 데이터가 없습니다."


def get_enable_push_user():
    session = SessionLocal()
    try:
        # 한 달 전의 날짜 계산
        one_month_ago = (datetime.utcnow() - timedelta(days=30)).strftime('%Y-%m-%d %H:%M:%S')

        query = text("""
        SELECT u.id, u.name
        FROM users u
        JOIN user_tokens ut ON u.id = ut.user_id
        WHERE u.push = 1 
        AND u.deleted_at IS NULL
        AND ut.deleted_at IS NULL
        AND ut.updated_at >= :one_month_ago
        GROUP BY u.id;
        """)

        result = session.execute(query, {"one_month_ago": one_month_ago}).fetchall()
        return [{"id": row[0], "name": row[1]} for row in result] if result else []
    finally:
        session.close()
        
# ✅ API 호출 함수 (유저 메시지 전송)
def send_user_message(user_id, meal_time, recommended_food):
    """ 음식 추천 결과를 사용자에게 메시지로 보내는 함수 """
    url = "https://food-pick-api.jokertrickster.com/v0.1/users/message"
    headers = {
        "Content-Type": "application/json"
    }
    if meal_time == "lunch":
        meal_time = "점심"
    else:
        meal_time = "저녁"
        
    # ✅ 요청 데이터 생성
    payload = {
        "role": "foodadmin",
        "userId": user_id,
        "title": f"{meal_time.capitalize()} 음식 추천",  # ✅ Title에 meal_time 반영
        "message": recommended_food
    }

    try:
        response = requests.post(url, headers=headers, json=payload)
        response.raise_for_status()  # ✅ HTTP 오류 발생 시 예외 처리
    except requests.exceptions.RequestException as e:
        print(f"❌ 메시지 전송 실패: {e}")

# ✅ Lambda 핸들러 (메인 실행 함수)
def lambda_handler(event, context):
    # 푸시가 활성화되어 있는 유저들을 가져온다. 
    userList = get_enable_push_user()
    
    for user in userList: 
        user_id = user['id']
        user_name = user['name']
        if user_id == 1:
            continue
        print(user_id, user_name)
        meal_time = get_current_meal_time()
        recommended_food = recommend_food_for_user(user_id,user_name, meal_time)
        send_user_message(user_id, meal_time, recommended_food)

    return {
        "statusCode": 200,
        "body": json.dumps({
            "message": "음식 메시지 전송 성공",
            "count": len(userList)
        }, ensure_ascii=False)
    }

# ✅ 프로그램 실행
if __name__ == "__main__":
    response = lambda_handler({}, None)
    print(response)
