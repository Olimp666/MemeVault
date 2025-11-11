# Image Storage

Сервис для хранения метаданных изображений

###  Запуск 

Из директории deployments: ```docker compose up --build```  
Сайт будет запущен на http://localhost/

Удаление базы данных:   
```
docker compose down
docker volume rm deployments_pg_data
```

### API Endpoints

#### Загрузка метаданных изображения

Сохраняет tg_file_id изображения из Telegram с указанными тегами.

```bash
curl -X POST "http://localhost/upload?user_id=123456789&tg_file_id=AgACAgIAAxkBAAIC..." \
  -H "Content-Type: application/json" \
  -d '{
    "tags": ["meme", "funny"]
  }'
```

Ответ (HTTP 201):
```
Created
```

**Query параметры:**
- `user_id` (int64) - ID пользователя Telegram
- `tg_file_id` (string) - File ID из Telegram API

**Body (JSON):**
- `tags` ([]string) - Массив тегов (минимум 1 тег)

#### Получение изображений по тегам

Возвращает tg_file_id изображений, которые содержат **все** указанные теги.  
Изображения фильтруются по user_id (возвращаются изображения пользователя или публичные с user_id=0).  
Результаты отсортированы по дате создания (новые первыми).

```bash
curl -X POST "http://localhost/images?user_id=123456789" \
  -H "Content-Type: application/json" \
  -d '{
    "tags": ["meme", "funny"]
  }'
```

Ответ:
```json
{
  "tg_file_ids": [
    "AgACAgIAAxkBAAIC...",
    "AgACAgIAAxkBAAID..."
  ]
}
```

**Query параметры:**
- `user_id` (int64) - ID пользователя Telegram

**Body (JSON):**
- `tags` ([]string) - Массив тегов для поиска (изображение должно иметь **все** эти теги)

**Логика фильтрации:**
- Возвращаются только изображения, у которых есть все указанные теги
- Возвращаются изображения, у которых `user_id` совпадает с переданным ИЛИ `user_id = 0` (публичные)
- Результаты отсортированы по `created_at DESC`

### Структура базы данных

**Таблица images:**
- `tg_file_id` (VARCHAR(255), PRIMARY KEY) - Telegram File ID
- `user_id` (BIGINT) - ID пользователя
- `created_at` (TIMESTAMP) - Дата создания

**Таблица tags:**
- `tg_file_id` (VARCHAR(255)) - Ссылка на изображение
- `name` (VARCHAR(255)) - Название тега
- PRIMARY KEY (`tg_file_id`, `name`)
