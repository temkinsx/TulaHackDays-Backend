-- Down: remove seeded achievements
DELETE
FROM achievements
WHERE name IN (
               'Первый отзыв',
               'Мастер отзывов',
               'Создатель мест',
               'Комментатор',
               'Фотограф',
               'Защитник здоровья'
    );
