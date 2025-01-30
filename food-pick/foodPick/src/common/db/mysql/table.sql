-- email, 비밀번호, 생년월일, 성별 
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    email VARCHAR(255),
    name VARCHAR(255),
    password VARCHAR(255),
    birth DATE,
    sex varchar(50),
    provider VARCHAR(50),
    push BOOLEAN DEFAULT TRUE,
    image varchar(1000) DEFAULT 'profile_default.png',
    role varchar(200)
);

INSERT INTO users (email, password, name,birth, sex, provider) VALUES ('test@test.com', 'asdasd123', '푸드픽맨','1990-01-01', 'male', 'test');
INSERT INTO users (email, password, name,birth, sex, provider) VALUES ('test01@test.com', 'asd123', '푸드픽맨','1990-01-01', 'male', 'test');
INSERT INTO users (email, password, name,birth, sex, provider) VALUES ('test02@test.com', 'asd123', '푸드픽맨','1990-01-01', 'male', 'test');
INSERT INTO users (email, password, name,birth, sex, provider) VALUES ('test03@test.com', 'asd123', '푸드픽맨','1990-01-01', 'male', 'test');
INSERT INTO users (email, password, name,birth, sex, provider) VALUES ('ryan@gmail.com', 'asdasd123', '푸드픽맨','1990-01-01', 'male', 'test');
INSERT INTO users (email, password, name,birth, sex, provider) VALUES ('joker@gmail.com', 'asdasd123', '푸드픽맨','1990-01-01', 'male', 'test');


CREATE TABLE tokens (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    refresh_expired_at INT,
    user_id INT,
    access_token VARCHAR(255),
    refresh_token VARCHAR(255),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE user_auths (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    auth_code VARCHAR(255),
    email VARCHAR(255),
    type VARCHAR(100)
);


-- 메타 데이터 테이블
  CREATE TABLE meta_tables (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    table_name VARCHAR(255) NOT NULL UNIQUE,
    table_description VARCHAR(255)
);


-- 아침, 점심, 저녁, 브런치, 간식, 야식
CREATE TABLE times (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    name VARCHAR(50) NOT NULL UNIQUE,
    image varchar(255) default 'category_default.png',
    description VARCHAR(255)
);

-- 한식, 중식, 일식, 양식, 분식,베트남 음식, 인도 음식, 패스트 푸드, 디저트, 퓨전 요리
CREATE TABLE types (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    name VARCHAR(50) NOT NULL UNIQUE,
    image varchar(255) default 'category_default.png',
    description VARCHAR(255)
);

-- 연인, 혼반, 가족, 다이어트, 회식, 친구
CREATE TABLE scenarios (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    name VARCHAR(50) NOT NULL UNIQUE,
    image varchar(255) default 'category_default.png',
    description VARCHAR(255)
);


-- 스트레스 해소, 피로 회복, 기분 전환, 제철 음식, 영양식, 특별한 날
CREATE TABLE themes(
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    name VARCHAR(50) NOT NULL UNIQUE,
    image varchar(255) default 'category_default.png',
    description VARCHAR(255)    
);




-- 신고 테이블
CREATE TABLE reports (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    user_id INT,
    reason varchar(1000),
    FOREIGN KEY (user_id) REFERENCES users(id)
);


-- meta table에 types, scenarios, times, themes 테이블을 저장하는 sql 문 만들어줘
INSERT INTO meta_tables (table_name, table_description) VALUES ('scenarios', '상황별'),('times', '시간별'),('types', '종류별'),('themes', '기분/테마별');

-- times 테이블에 아침, 점심, 저녁, 브런치, 간식, 야식 순으로 저장하는 sql 문 만들어줘
INSERT INTO times (name, description,image) VALUES ('아침', '아침','times/breakfast.png'), ('점심', '점심','times/lunch.png'), ('저녁', '저녁','times/dinner.png'), ('간식', '간식','times/snack.png'), ('야식', '야식','times/late night snack.png');

-- types 테이블에 한식, 중식, 일식, 양식, 분식,베트남 음식, 인도 음식, 패스트 푸드, 디저트, 퓨전 요리 순으로 저장하는 sql 문 만들어줘
INSERT INTO types (name, description,image) VALUES ('한식', '한식','types/korean food.png'), ('중식', '중식','types/chinese food.png'), ('일식', '일식','types/japanese food.png'), ('양식', '양식','types/western food.png'), ('분식', '분식','types/korean street food.png'),  ('패스트 푸드', '패스트 푸드','types/fast food.png'),('베트남 음식', '베트남 음식','types/vietnamese food.png'), ('인도 음식', '인도 음식','types/indian food.png'), ('디저트', '디저트','types/dessert.png'), ('퓨전 요리', '퓨전 요리','types/fusion cuisine.png');

-- scenarios 테이블에 연인, 혼반, 가족, 다이어트, 회식, 친구 순으로 저장하는 sql 문 만들어줘
INSERT INTO scenarios (name, description,image) VALUES ('연인', '연인','scenarios/couple.png'), ('혼밥', '혼밥','scenarios/eating alone.png'), ('가족', '가족','scenarios/family.png'), ('회식', '회식','scenarios/company dinner.png'), ('친구', '친구','scenarios/friend.png');


-- themes 테이블에 스트레스 해소, 피로 회복, 기분 전환, 제철 음식, 영양식, 특별한 날 순으로 저장하는 sql 문 만들어줘
INSERT INTO themes (name, description,image) VALUES ('스트레스 해소', '스트레스 해소','themes/stress.png'), ('해장', '해장','themes/hangover.png'),('피로 회복', '피로 회복','themes/fatigue recovery.png'), ('다이어트', '다이어트','themes/diet.png'), ('제철 음식', '제철 음식','themes/seasonal food.png');

CREATE TABLE user_tokens (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    user_id INT,
    token varchar(1000),
    FOREIGN KEY (user_id) REFERENCES users(id)
);




  

-- 영양소 테이블 용량, 칼로리, 탄수화물, 단백질, 지방 
create table nutrients (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    food_name varchar(255) UNIQUE,
    amount varchar(255),
    kcal DECIMAL(10, 2),
    carbohydrate DECIMAL(10, 2),
    protein DECIMAL(10, 2),
    fat DECIMAL(10, 2)
);

CREATE TABLE food_images (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,  -- 음식 이름으로 이미지 참조
    image VARCHAR(255) DEFAULT 'food_default.png',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

-- 음식 선택 개선

CREATE TABLE foods (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    name VARCHAR(255) NOT NULL,
    image_id INT,
    FOREIGN KEY (image_id) REFERENCES food_images(id)
);
-- INSERT INTO foods (name, food_image_id) VALUES ('김치찌개', NULL);

CREATE TABLE category_types (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    name VARCHAR(50) NOT NULL UNIQUE COMMENT '카테고리 유형'
);
INSERT INTO category_types (name) VALUES ('time'), ('type'), ('scenario'), ('theme');


CREATE TABLE categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    name VARCHAR(255) NOT NULL UNIQUE,
    type_id INT NOT NULL,
    FOREIGN KEY (type_id) REFERENCES category_types(id)
);
-- INSERT INTO categories (name, type_id) VALUES ('아침', 1), ('한식', 2);
INSERT INTO categories (name, type_id) VALUES 
    ('연인', (SELECT id FROM category_types WHERE name = 'scenario')),
    ('혼밥', (SELECT id FROM category_types WHERE name = 'scenario')),
    ('가족', (SELECT id FROM category_types WHERE name = 'scenario')),
    ('회식', (SELECT id FROM category_types WHERE name = 'scenario')),
    ('친구', (SELECT id FROM category_types WHERE name = 'scenario'));
    
    
INSERT INTO categories (name, type_id) VALUES 
    ('아침', (SELECT id FROM category_types WHERE name = 'time')),
    ('점심', (SELECT id FROM category_types WHERE name = 'time')),
    ('저녁', (SELECT id FROM category_types WHERE name = 'time')),
    ('간식', (SELECT id FROM category_types WHERE name = 'time')),
    ('야식', (SELECT id FROM category_types WHERE name = 'time'));


INSERT INTO categories (name, type_id) VALUES 
    ('한식', (SELECT id FROM category_types WHERE name = 'type')),
    ('중식', (SELECT id FROM category_types WHERE name = 'type')),
    ('일식', (SELECT id FROM category_types WHERE name = 'type')),
    ('양식', (SELECT id FROM category_types WHERE name = 'type')),
    ('분식', (SELECT id FROM category_types WHERE name = 'type')),
    ('베트남 음식', (SELECT id FROM category_types WHERE name = 'type')),
    ('인도 음식', (SELECT id FROM category_types WHERE name = 'type')),
    ('패스트 푸드', (SELECT id FROM category_types WHERE name = 'type')),
    ('디저트', (SELECT id FROM category_types WHERE name = 'type')),
    ('퓨전 요리', (SELECT id FROM category_types WHERE name = 'type'));


INSERT INTO categories (name, type_id) VALUES 
    ('스트레스 해소', (SELECT id FROM category_types WHERE name = 'theme')),
    ('해장', (SELECT id FROM category_types WHERE name = 'theme')),
    ('피로 회복', (SELECT id FROM category_types WHERE name = 'theme')),
    ('다이어트', (SELECT id FROM category_types WHERE name = 'theme')),
    ('제철 음식', (SELECT id FROM category_types WHERE name = 'theme'));



CREATE TABLE food_categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    food_id INT NOT NULL,
    category_id INT NOT NULL,
    FOREIGN KEY (food_id) REFERENCES foods(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);
-- INSERT INTO food_categories (food_id, category_id) VALUES (1, 1), -- 김치찌개는 아침 (1, 2), -- 김치찌개는 한식 (1, 3); -- 김치찌개는 가족

CREATE TABLE food_histories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    user_id INT NOT NULL, -- 선택한 유저 ID
    food_id INT NOT NULL, -- 선택된 음식 ID
    FOREIGN KEY (food_id) REFERENCES foods(id) ON DELETE CASCADE
);