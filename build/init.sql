-- build/init.sql
-- Создание базы данных
CREATE DATABASE "MaterialBD";

\c "MaterialBD";

-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(150) UNIQUE NOT NULL,
    password VARCHAR(128) NOT NULL,
    is_moderator BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица материалов
CREATE TABLE IF NOT EXISTS materials (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'действует' CHECK (status IN ('действует', 'удалён')),
    image_url VARCHAR(500),
    price_per_m2 FLOAT NOT NULL DEFAULT 0,
    lambda FLOAT NOT NULL DEFAULT 0,
    thickness FLOAT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица заявок
CREATE TABLE IF NOT EXISTS applications (
    id SERIAL PRIMARY KEY,
    status VARCHAR(50) DEFAULT 'черновик' CHECK (status IN ('черновик', 'удалён', 'сформирован', 'завершён', 'отклонён')),
    creator_id INTEGER NOT NULL REFERENCES users(id),
    total_area FLOAT NOT NULL DEFAULT 0,
    indoor_temp FLOAT NOT NULL DEFAULT 22,
    outdoor_temp FLOAT NOT NULL DEFAULT -15,
    total_savings FLOAT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    submitted_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    moderator_id INTEGER NULL REFERENCES users(id)
);

-- Таблица связи заявок и материалов (многие-ко-многим)
CREATE TABLE IF NOT EXISTS application_materials (
    application_id INTEGER NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    material_id INTEGER NOT NULL REFERENCES materials(id) ON DELETE CASCADE,
    area FLOAT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (application_id, material_id)
);

-- Создание индексов для улучшения производительности
CREATE INDEX IF NOT EXISTS idx_applications_creator_status ON materials_applications(creator_id, status);
CREATE INDEX IF NOT EXISTS idx_applications_status ON materials_applications(status);
CREATE INDEX IF NOT EXISTS idx_materials_status ON materials(status);
CREATE INDEX IF NOT EXISTS idx_application_materials_app_id ON application_materials(application_id);
CREATE INDEX IF NOT EXISTS idx_application_materials_mat_id ON application_materials(material_id);