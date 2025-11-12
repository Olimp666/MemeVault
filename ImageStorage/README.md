# Image Storage

Сервис для хранения метаданных медиафайлов Telegram (изображения, видео, гифки)

## Содержание

### API Endpoints
- [Загрузка метаданных медиафайла](#загрузка-метаданных-медиафайла) - `POST /upload`
- [Получение медиафайлов по тегам](#получение-медиафайлов-по-тегам) - `POST /images`
- [Получение всех медиафайлов пользователя](#получение-всех-медиафайлов-пользователя) - `GET /user/images`
- [Удаление конкретной картинки](#удаление-конкретной-картинки) - `DELETE /image/delete`
- [Удаление всех картинок пользователя](#удаление-всех-картинок-пользователя) - `DELETE /user/images/delete`
- [Замена тегов картинки](#замена-тегов-картинки) - `PUT /image/tags`
- [Генерация описания к картинке](#генерация-описания-к-картинке) - `POST /image/generate-description`

### Дополнительная информация
- [Типы файлов](#типы-файлов)
- [Запуск сервиса](#запуск)

---

## Запуск 

Из директории deployments: ```docker compose up --build```  
Сайт будет запущен на http://localhost/

Удаление базы данных:   
```
docker compose down
docker volume rm deployments_pg_data
```

---

## API Endpoints

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

Возвращает информацию о медиафайлах (tg_file_id, file_type и tags), которые содержат указанные теги.  
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
  "exact_match": [
    {
      "tg_file_id": "AgACAgIAAxkBAAIC...",
      "file_type": "photo",
      "tags": ["meme", "funny", "cat"]
    }
  ],
  "partial_match": [
    {
      "tg_file_id": "BAACAgIAAxkBAAID...",
      "file_type": "video",
      "tags": ["meme", "dog"]
    }
  ]
}
```

**Query параметры:**
- `user_id` (int64) - ID пользователя Telegram

**Body (JSON):**
- `tags` ([]string) - Массив тегов для поиска

**Структура ответа:**
- `exact_match` - медиафайлы, имеющие ВСЕ указанные теги (с полным списком всех тегов файла)
- `partial_match` - медиафайлы, имеющие хотя бы один из тегов (но не все), отсортированные по количеству совпадающих тегов (с полным списком всех тегов файла)

**Логика фильтрации:**
- `exact_match`: медиафайлы с точным совпадением всех указанных тегов
- `partial_match`: медиафайлы с подмножеством тегов, отсортированные по количеству совпадений (больше совпадений = выше в списке)
- Возвращаются медиафайлы пользователя (`user_id`) или публичные (`user_id = 0`)
- Для каждого медиафайла возвращаются ВСЕ его теги (не только совпавшие)
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

#### Удаление конкретной картинки

Удаляет конкретную картинку пользователя по tg_file_id.

```bash
curl -X DELETE "http://localhost/image/delete?user_id=123456789&tg_file_id=AgACAgIAAxkBAAIC..."
```

Ответ (HTTP 200):
```
Image deleted successfully
```

**Query параметры:**
- `user_id` (int64) - ID пользователя Telegram
- `tg_file_id` (string) - File ID из Telegram API

#### Удаление всех картинок пользователя

Удаляет все картинки конкретного пользователя.

```bash
curl -X DELETE "http://localhost/user/images/delete?user_id=123456789"
```

Ответ (HTTP 200):
```
All user images deleted successfully
```

**Query параметры:**
- `user_id` (int64) - ID пользователя Telegram

#### Замена тегов картинки

Заменяет все теги конкретной картинки пользователя на новый набор тегов.

**Важно:** Нельзя заменять теги для публичных изображений (user_id = 0).

```bash
curl -X PUT "http://localhost/image/tags?user_id=123456789&tg_file_id=AgACAgIAAxkBAAIC..." \
  -H "Content-Type: application/json" \
  -d '{
    "tags": ["новый_тег1", "новый_тег2"]
  }'
```

Ответ (HTTP 200):
```
Tags replaced successfully
```

**Query параметры:**
- `user_id` (int64) - ID пользователя Telegram (не должен быть 0)
- `tg_file_id` (string) - File ID из Telegram API

**Body (JSON):**
- `tags` ([]string) - Новый массив тегов (минимум 1 тег)

**Особенности:**
- Все старые теги удаляются и заменяются новыми
- Нельзя использовать для публичных изображений (user_id = 0)

#### Генерация описания к картинке

Генерирует текстовое описание для загруженной картинки.

**Примечание:** Метод пока не реализован и возвращает заглушку.

```bash
curl -X POST "http://localhost/image/generate-description" \
  -F "image=@/path/to/image.jpg"
```

Ответ (HTTP 200):
```json
{
  "description": "метод пока не реализован"
}
```

**Form Data:**
- `image` (file) - Файл изображения (максимум 10MB)

### Типы файлов

Поддерживаются следующие типы медиафайлов:
- `photo` - фотография
- `video` - видео
- `gif` - анимация/гифка
