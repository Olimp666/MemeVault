#!/bin/bash

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Настройки
SERVER_URL="http://localhost"
USER_ID=123456789
TG_FILE_ID="AgACAgIAAxkBAAIC_test_file_id_12345"
TAGS='["meme", "funny", "test"]'

echo -e "${YELLOW}=== Тест сохранения и получения tg_file_id ===${NC}"

echo -e "${YELLOW}1. Загрузка метаданных изображения (tg_file_id) на сервер...${NC}"

# Загрузка tg_file_id с тегами
HTTP_CODE=$(curl -s -w "%{http_code}" -o /tmp/upload_response.txt -X POST \
  "$SERVER_URL/upload?user_id=$USER_ID&tg_file_id=$TG_FILE_ID" \
  -H "Content-Type: application/json" \
  -d "{\"tags\": [\"meme\", \"funny\", \"test\"]}")

if [ "$HTTP_CODE" != "201" ]; then
    echo -e "${RED}Ошибка: Неожиданный HTTP код: $HTTP_CODE${NC}"
    echo "Ответ сервера:"
    cat /tmp/upload_response.txt
    rm -f /tmp/upload_response.txt
    exit 1
fi

echo -e "${GREEN}✓ Метаданные успешно сохранены (HTTP 201)${NC}"
echo "TG File ID: $TG_FILE_ID"
echo "Tags: meme, funny, test"

echo -e "${YELLOW}2. Получение tg_file_id по всем тегам...${NC}"

# Получение изображений по тегам (должны быть все теги)
GET_RESPONSE=$(curl -s -X POST "$SERVER_URL/images?user_id=$USER_ID" \
  -H "Content-Type: application/json" \
  -d '{"tags": ["meme", "funny", "test"]}')

if [ $? -ne 0 ]; then
    echo -e "${RED}Ошибка: Не удалось получить изображения${NC}"
    exit 1
fi

echo "Ответ сервера: $GET_RESPONSE"

# Проверка, что получен массив tg_file_ids
if ! echo "$GET_RESPONSE" | grep -q '"tg_file_ids"'; then
    echo -e "${RED}Ошибка: В ответе отсутствует поле tg_file_ids${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Получен ответ от сервера${NC}"

# Извлечение tg_file_ids из ответа
RETURNED_FILE_ID=$(echo "$GET_RESPONSE" | grep -o "\"$TG_FILE_ID\"" | head -1 | tr -d '"')

if [ -z "$RETURNED_FILE_ID" ]; then
    echo -e "${RED}Ошибка: Загруженный tg_file_id не найден в результатах${NC}"
    echo "Ожидали: $TG_FILE_ID"
    echo "Получили: $GET_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✓ TG File ID найден в результатах${NC}"

echo -e "${YELLOW}3. Проверка фильтрации по тегам (частичное совпадение)...${NC}"

# Попытка получить с неполным набором тегов (должна вернуть результат)
GET_PARTIAL=$(curl -s -X POST "$SERVER_URL/images?user_id=$USER_ID" \
  -H "Content-Type: application/json" \
  -d '{"tags": ["meme", "funny"]}')

if echo "$GET_PARTIAL" | grep -q "\"$TG_FILE_ID\""; then
    echo -e "${GREEN}✓ Частичный набор тегов работает корректно${NC}"
else
    echo -e "${RED}Ошибка: Изображение не найдено с частичным набором тегов${NC}"
    exit 1
fi

echo -e "${YELLOW}4. Проверка фильтрации с несуществующим тегом...${NC}"

# Попытка получить с тегом, которого нет (не должна вернуть результат)
GET_WRONG=$(curl -s -X POST "$SERVER_URL/images?user_id=$USER_ID" \
  -H "Content-Type: application/json" \
  -d '{"tags": ["nonexistent"]}')

if echo "$GET_WRONG" | grep -q "\"$TG_FILE_ID\""; then
    echo -e "${RED}Ошибка: Изображение найдено с несуществующим тегом${NC}"
    exit 1
else
    echo -e "${GREEN}✓ Фильтрация по тегам работает корректно${NC}"
fi

echo -e "${YELLOW}5. Тест публичных изображений (user_id=0)...${NC}"

# Загрузка публичного изображения
PUBLIC_FILE_ID="AgACAgIAAxkBAAIC_public_file_id_67890"
HTTP_CODE=$(curl -s -w "%{http_code}" -o /dev/null -X POST \
  "$SERVER_URL/upload?user_id=0&tg_file_id=$PUBLIC_FILE_ID" \
  -H "Content-Type: application/json" \
  -d '{"tags": ["public", "meme"]}')

if [ "$HTTP_CODE" != "201" ]; then
    echo -e "${RED}Ошибка: Не удалось загрузить публичное изображение${NC}"
    exit 1
fi

# Другой пользователь должен видеть публичное изображение
OTHER_USER_ID=999999999
GET_PUBLIC=$(curl -s -X POST "$SERVER_URL/images?user_id=$OTHER_USER_ID" \
  -H "Content-Type: application/json" \
  -d '{"tags": ["public"]}')

if echo "$GET_PUBLIC" | grep -q "\"$PUBLIC_FILE_ID\""; then
    echo -e "${GREEN}✓ Публичные изображения доступны всем пользователям${NC}"
else
    echo -e "${RED}Ошибка: Публичное изображение не доступно другим пользователям${NC}"
    exit 1
fi

echo -e "${YELLOW}6. Сравнение загруженного и полученного tg_file_id...${NC}"

if [ "$RETURNED_FILE_ID" == "$TG_FILE_ID" ]; then
    echo -e "${GREEN}✓ УСПЕХ: TG File ID совпадают!${NC}"
    echo "Загружено: $TG_FILE_ID"
    echo "Получено:  $RETURNED_FILE_ID"
    echo -e "${GREEN}=== Все тесты пройдены успешно ===${NC}"
    
    # Очистка
    rm -f /tmp/upload_response.txt
    exit 0
else
    echo -e "${RED}✗ ОШИБКА: TG File ID различаются!${NC}"
    echo "Загружено: $TG_FILE_ID"
    echo "Получено:  $RETURNED_FILE_ID"
    exit 1
fi
