-- init.sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(150) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_moderator BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS materials (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    price_per_m2 DECIMAL(10,2) NOT NULL,
    lambda_value DECIMAL(5,3) NOT NULL,
    image_url VARCHAR(200),
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS calculation_requests (
    id SERIAL PRIMARY KEY,
    status VARCHAR(20) NOT NULL CHECK (status IN ('draft','deleted','formed','completed','rejected')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    creator_id INTEGER REFERENCES users(id) NOT NULL,
    formed_at TIMESTAMP,
    completed_at TIMESTAMP,
    moderator_id INTEGER REFERENCES users(id),
    indoor_temp DECIMAL(5,2) DEFAULT 22.0,
    outdoor_temp DECIMAL(5,2) DEFAULT -15.0,
    total_savings DECIMAL(10,2) DEFAULT 0
);

CREATE TABLE IF NOT EXISTS calculation_materials (
    id SERIAL PRIMARY KEY,
    request_id INTEGER REFERENCES calculation_requests(id) NOT NULL,
    material_id INTEGER REFERENCES materials(id) NOT NULL,
    area DECIMAL(10,2) NOT NULL,
    heat_loss DECIMAL(10,2) DEFAULT 0,
    savings DECIMAL(10,2) DEFAULT 0,
    UNIQUE(request_id, material_id)
);

-- Ограничение: один черновик на пользователя
CREATE UNIQUE INDEX IF NOT EXISTS unique_draft_per_user 
ON calculation_requests (creator_id) 
WHERE status = 'draft';

-- Наполнение тестовыми данными
INSERT INTO users (username, password_hash, is_moderator) VALUES 
('user1', 'hash1', false),
('moderator1', 'hash2', true)
ON CONFLICT (username) DO NOTHING;

INSERT INTO materials (name, description, price_per_m2, lambda_value, image_url) VALUES 
('Минеральная вата', 'Эффективный утеплитель с отличными звукоизоляционными свойствами', 650, 0.045, 'http://localhost:9000/images/mineral_wool.jpg'),
('Пенополистирол', 'Легкий и влагостойкий материал для утепления', 450, 0.038, 'http://localhost:9000/images/polystyrene.jpg'),
('PIR-плиты', 'Современный высокоэффективный утеплитель с низкой теплопроводностью', 1200, 0.028, 'http://localhost:9000/images/pir_plates.jpg'),
('Экструдированный пенополистирол', 'Прочный влагостойкий утеплитель для фундаментов и фасадов', 550, 0.033, 'http://localhost:9000/images/xps.jpg')
ON CONFLICT (id) DO NOTHING;