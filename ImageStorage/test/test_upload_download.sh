#!/bin/bash

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Настройки
SERVER_URL="http://localhost"
TEST_IMAGE="test_image.jpg"
DOWNLOADED_IMAGE="downloaded_image.jpg"
USER_ID=1
TAG="test"

echo -e "${YELLOW}=== Тест загрузки и скачивания изображения ===${NC}"

# Проверка наличия тестового изображения
if [ ! -f "$TEST_IMAGE" ]; then
    echo -e "${RED}Ошибка: Файл $TEST_IMAGE не найден!${NC}"
    echo "Пожалуйста, поместите тестовое изображение в директорию test/"
    exit 1
fi

echo -e "${YELLOW}1. Загрузка изображения на сервер...${NC}"

# Загрузка изображения
UPLOAD_RESPONSE=$(curl -s -X POST "$SERVER_URL/upload" \
  -F "image=@$TEST_IMAGE" \
  -F "user_id=$USER_ID" \
  -F "tags=[\"$TAG\"]")

if [ $? -ne 0 ]; then
    echo -e "${RED}Ошибка: Не удалось отправить запрос на загрузку${NC}"
    exit 1
fi

echo "Ответ сервера: $UPLOAD_RESPONSE"

# Извлечение image_id из ответа
IMAGE_ID=$(echo "$UPLOAD_RESPONSE" | grep -o '"image_id":[0-9]*' | grep -o '[0-9]*')

if [ -z "$IMAGE_ID" ]; then
    echo -e "${RED}Ошибка: Не удалось получить image_id из ответа${NC}"
    echo "Ответ сервера: $UPLOAD_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✓ Изображение успешно загружено с ID: $IMAGE_ID${NC}"

echo -e "${YELLOW}2. Получение изображения по тегу '$TAG'...${NC}"

# Получение изображений по тегу
GET_RESPONSE=$(curl -s -X GET "$SERVER_URL/images?tag=$TAG")

if [ $? -ne 0 ]; then
    echo -e "${RED}Ошибка: Не удалось получить изображения${NC}"
    exit 1
fi

# Проверка, что получен массив изображений
IMAGES_COUNT=$(echo "$GET_RESPONSE" | grep -o '"images":\[' | wc -l)

if [ "$IMAGES_COUNT" -eq 0 ]; then
    echo -e "${RED}Ошибка: В ответе отсутствует массив images${NC}"
    echo "Ответ сервера: $GET_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✓ Получен ответ от сервера${NC}"

# Извлечение первого изображения из массива (base64)
IMAGE_BASE64=$(echo "$GET_RESPONSE" | grep -o '"images":\["[^"]*"' | sed 's/"images":\["//' | sed 's/"$//')

if [ -z "$IMAGE_BASE64" ]; then
    echo -e "${RED}Ошибка: Не удалось извлечь данные изображения${NC}"
    exit 1
fi

echo -e "${YELLOW}3. Декодирование и сохранение изображения...${NC}"

# Декодирование base64 и сохранение в файл
echo "$IMAGE_BASE64" | base64 -d > "$DOWNLOADED_IMAGE"

if [ $? -ne 0 ]; then
    echo -e "${RED}Ошибка: Не удалось декодировать изображение${NC}"
    exit 1
fi

if [ ! -f "$DOWNLOADED_IMAGE" ]; then
    echo -e "${RED}Ошибка: Файл $DOWNLOADED_IMAGE не был создан${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Изображение сохранено в $DOWNLOADED_IMAGE${NC}"

echo -e "${YELLOW}4. Сравнение исходного и скачанного файлов...${NC}"

# Получение хешей файлов
ORIGINAL_HASH=$(md5sum "$TEST_IMAGE" | awk '{print $1}')
DOWNLOADED_HASH=$(md5sum "$DOWNLOADED_IMAGE" | awk '{print $1}')

echo "MD5 исходного файла:    $ORIGINAL_HASH"
echo "MD5 скачанного файла:   $DOWNLOADED_HASH"

# Сравнение файлов
if [ "$ORIGINAL_HASH" == "$DOWNLOADED_HASH" ]; then
    echo -e "${GREEN}✓ УСПЕХ: Файлы идентичны!${NC}"
    echo -e "${GREEN}=== Тест пройден успешно ===${NC}"
    
    # Удаление скачанного файла
    rm -f "$DOWNLOADED_IMAGE"
    exit 0
else
    echo -e "${RED}✗ ОШИБКА: Файлы различаются!${NC}"
    echo "Скачанный файл сохранен как $DOWNLOADED_IMAGE для анализа"
    exit 1
fi
