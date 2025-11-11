# Image Storage

Сервис для хранения метаданных медиафайлов Telegram (изображения, видео, гифки)

###  Запуск 

Из директории deployments: ```docker compose up --build```  
Сайт будет запущен на http://localhost/

Удаление базы данных:   
```
docker compose down
docker volume rm deployments_pg_data
```

### API Endpoints

#### Загрузка метаданных медиафайла

Сохраняет tg_file_id медиафайла из Telegram с указанным типом и тегами.

```bash
curl -X POST "http://localhost/upload?user_id=123456789&tg_file_id=AgACAgIAAxkBAAIC...&file_type=photo" \
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
- `file_type` (string) - Тип файла: `photo`, `video` или `gif`

**Body (JSON):**
- `tags` ([]string) - Массив тегов (минимум 1 тег)

#### Получение медиафайлов по тегам

Возвращает информацию о медиафайлах (tg_file_id и file_type), которые содержат **все** указанные теги.  
Медиафайлы фильтруются по user_id (возвращаются файлы пользователя или публичные с user_id=0).  
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
  "images": [
    {
      "tg_file_id": "AgACAgIAAxkBAAIC...",
      "file_type": "photo"
    },
    {
      "tg_file_id": "BAACAgIAAxkBAAID...",
      "file_type": "video"
    }
  ]
}
```

**Query параметры:**
- `user_id` (int64) - ID пользователя Telegram

**Body (JSON):**
- `tags` ([]string) - Массив тегов для поиска (медиафайл должен иметь **все** эти теги)

**Структура ответа:**
- `exact_match` - медиафайлы, имеющие ВСЕ указанные теги
- `partial_match` - медиафайлы, имеющие хотя бы один из тегов (но не все), отсортированные по количеству совпадающих тегов

**Логика фильтрации:**
- `exact_match`: медиафайлы с точным совпадением всех тегов
- `partial_match`: медиафайлы с подмножеством тегов, отсортированные по количеству совпадений (больше совпадений = выше в списке)
- Возвращаются медиафайлы пользователя (`user_id`) или публичные (`user_id = 0`)
- Результаты отсортированы сначала по количеству совпадений, затем по `created_at DESC`

#### Получение всех медиафайлов пользователя

Возвращает все медиафайлы конкретного пользователя без фильтрации по тегам.

```bash
curl -X GET "http://localhost/user/images?user_id=123456789"
```

Ответ:
```json
{
  "images": [
    {
      "tg_file_id": "AgACAgIAAxkBAAIC...",
      "file_type": "photo",
      "tags": ["meme", "funny"]
    },
    {
      "tg_file_id": "BAACAgIAAxkBAAID...",
      "file_type": "video",
      "tags": ["cat", "cute"]
    }
  ]
}
```

**Query параметры:**
- `user_id` (int64) - ID пользователя Telegram

**Особенности:**
- Возвращает только медиафайлы указанного пользователя (НЕ включает публичные файлы)
- Результаты отсортированы по `created_at DESC` (новые первыми)
- Для каждого медиафайла возвращаются все его теги

**Повторная загрузка:**
- Если пользователь загружает тот же `tg_file_id` повторно с новыми тегами, изображение НЕ дублируется
- Новые теги добавляются к существующей записи
- Старые теги сохраняются

### Типы файлов

Поддерживаются следующие типы медиафайлов:
- `photo` - фотография
- `video` - видео
- `gif` - анимация/гифка

### Структура базы данных

**Таблица images:**
- `tg_file_id` (VARCHAR(255), PRIMARY KEY) - Telegram File ID
- `user_id` (BIGINT) - ID пользователя
- `file_type` (VARCHAR(50)) - Тип файла (photo, video, gif)
- `created_at` (TIMESTAMP) - Дата создания

**Таблица tags:**
- `tg_file_id` (VARCHAR(255)) - Ссылка на медиафайл
- `name` (VARCHAR(255)) - Название тега
- PRIMARY KEY (`tg_file_id`, `name`)
