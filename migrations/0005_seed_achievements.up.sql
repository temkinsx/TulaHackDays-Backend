-- Seed default achievements
INSERT INTO achievements (name, description, icon, points, type, condition, is_active)
VALUES ('Первый отзыв', 'Оставьте свой первый отзыв', '⭐', 10, 'first_review', '{"reviews_count": 1}', true),
       ('Мастер отзывов', 'Оставьте 50 отзывов', '🏆', 100, 'review_master', '{"reviews_count": 50}', true),
       ('Создатель мест', 'Добавьте 10 новых объектов', '📍', 150, 'place_creator', '{"places_count": 10}', true),
       ('Комментатор', 'Оставьте 20 комментариев', '💬', 50, 'commenter', '{"comments_count": 20}', true),
       ('Фотограф', 'Добавьте 25 фотографий', '📸', 75, 'photographer', '{"photos_count": 25}', true),
       ('Защитник здоровья', 'Получите 500 очков', '❤️', 200, 'health_advocate', '{"points": 500}',
        true)
ON CONFLICT (name) DO NOTHING;
