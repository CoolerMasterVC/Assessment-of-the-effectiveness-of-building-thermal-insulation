-- build/fill.sql
-- Заполнение пользователей
INSERT INTO users (login, password, is_moderator) VALUES
  ('user', 'user', false),
  ('admin', 'admin', true),
  ('user2', 'user2', false)
ON CONFLICT (login) DO NOTHING;

-- Заполнение теплоизоляционных материалов
INSERT INTO materials (name, description, status, image_url, price_per_m2, lambda, thickness) VALUES
  (
    'Минеральная вата',
    'Эффективный утеплитель с отличными звукоизоляционными свойствами',
    'действует',
    'http://localhost:9000/images/mineral_wool.jpg',
    650,
    0.045,
    0.1
  ),
  (
    'Пенополистирол',
    'Легкий и влагостойкий материал для утепления',
    'действует',
    'http://localhost:9000/images/polystyrene.jpg',
    450,
    0.038,
    0.05
  ),
  (
    'PIR-плиты',
    'Современный высокоэффективный утеплитель с низкой теплопроводностью',
    'действует',
    'http://localhost:9000/images/pir_plates.jpg',
    1200,
    0.028,
    0.08
  )
ON CONFLICT (name) DO NOTHING;

-- Заполнение заявок
INSERT INTO applications (status, creator_id, total_area, indoor_temp, outdoor_temp, created_at) VALUES
  ('черновик', 1, 37.5, 22.0, -15.0, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;

-- Заполнение связи заявок и материалов
INSERT INTO application_materials (application_id, material_id, area) VALUES
  (1, 1, 15.5),
  (1, 3, 22.0)
ON CONFLICT DO NOTHING;