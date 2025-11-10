# Image Storage

Сервис для хранения изображений

###  Запуск 

Из директории deployments: ```docker compose up --build```  
Сайт будет запущен на http://localhost/

Удаление базы данных:   
```
docker compose down
docker volume rm deployments_pg_data
```

### API Endpoints

#### Загрузка изображения

```bash
curl -X POST http://localhost/upload \
  -F "image=@test_image.jpg" \
  -F "user_id=1" \
  -F "tags=[\"meme\",\"funny\"]"
```

Ответ:
```json
{
  "image_id": 1
}
```

#### Получение изображений по тегу

```bash
curl -X GET "http://localhost/images?tag=meme"
```

Ответ (изображения кодируются в base64):
```json
{
  "images": [
    "/9j/4AAQSkZJRgABAQAA...",
    "iVBORw0KGgoAAAANSUhE..."
  ]
}
```

Для декодирования и сохранения изображения:
```bash
curl -s "http://localhost/images?tag=meme" | jq -r '.images[0]' | base64 -d > image.jpg
```

### Тестирование

Для проверки загрузки и скачивания изображений запустите скрипт из директории `test`:

```bash
cd test
./test_upload_download.sh
```

Скрипт загружает тестовое изображение, скачивает его обратно и проверяет, что файлы совпадают.
