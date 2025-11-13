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
FILE_TYPE="photo"
TAGS='["meme", "funny", "test"]'

echo -e "${YELLOW}=== Тест сохранения и получения tg_file_id ===${NC}"

echo -e "${YELLOW}1. Загрузка метаданных изображения (tg_file_id) на сервер...${NC}"

# Загрузка tg_file_id с тегами
HTTP_CODE=$(curl -s -w "%{http_code}" -o /tmp/upload_response.txt -X POST \
  "$SERVER_URL/upload?user_id=$USER_ID&tg_file_id=$TG_FILE_ID&file_type=$FILE_TYPE" \
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
echo "File Type: $FILE_TYPE"
echo "Tags: meme, funny, test"

echo -e "${YELLOW}2. Получение tg_file_id по всем тегам...${NC}"

GET_RESPONSE=$(curl -s -X POST "$SERVER_URL/images?user_id=$USER_ID" \
  -H "Content-Type: application/json" \
  -d '{"tags": ["meme", "funny", "test"]}')

if [ $? -ne 0 ]; then
    echo -e "${RED}Ошибка: Не удалось получить изображения${NC}"
    exit 1
fi

echo "Ответ сервера: $GET_RESPONSE"

if ! echo "$GET_RESPONSE" | grep -q '"exact_match"'; then
    echo -e "${RED}Ошибка: В ответе отсутствует поле exact_match${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Получен ответ от сервера${NC}"

RETURNED_FILE_ID=$(echo "$GET_RESPONSE" | grep -oP "\"exact_match\".*?\"tg_file_id\":\s*\"$TG_FILE_ID\"" | grep -o "$TG_FILE_ID")
RETURNED_FILE_TYPE=$(echo "$GET_RESPONSE" | grep -oP "\"exact_match\".*?\"tg_file_id\":\s*\"$TG_FILE_ID\".*?\"file_type\":\s*\"\K[^\"]+")

if [ -z "$RETURNED_FILE_ID" ]; then
    echo -e "${RED}Ошибка: Загруженный tg_file_id не найден в exact_match${NC}"
    echo "Ожидали: $TG_FILE_ID"
    echo "Получили: $GET_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✓ TG File ID найден в exact_match${NC}"

if [ "$RETURNED_FILE_TYPE" != "$FILE_TYPE" ]; then
    echo -e "${RED}Ошибка: File type не совпадает${NC}"
    echo "Ожидали: $FILE_TYPE"
    echo "Получили: $RETURNED_FILE_TYPE"
    exit 1
fi

echo -e "${GREEN}✓ File Type совпадает${NC}"

# Проверка наличия тегов в ответе exact_match
if ! echo "$GET_RESPONSE" | grep -q '"tags"'; then
    echo -e "${RED}Ошибка: В ответе exact_match отсутствуют теги${NC}"
    exit 1
fi

if echo "$GET_RESPONSE" | grep -q '"meme"' && echo "$GET_RESPONSE" | grep -q '"funny"' && echo "$GET_RESPONSE" | grep -q '"test"'; then
    echo -e "${GREEN}✓ Все теги присутствуют в exact_match${NC}"
else
    echo -e "${RED}Ошибка: Не все теги найдены в exact_match${NC}"
    echo "Получили: $GET_RESPONSE"
    exit 1
fi

echo -e "${YELLOW}3. Проверка фильтрации по тегам (частичное совпадение)...${NC}"

GET_PARTIAL=$(curl -s -X POST "$SERVER_URL/images?user_id=$USER_ID" \
  -H "Content-Type: application/json" \
  -d '{"tags": ["meme", "funny"]}')

if echo "$GET_PARTIAL" | grep -q "\"partial_match\""; then
    if echo "$GET_PARTIAL" | grep "\"partial_match\"" | grep -q "\"$TG_FILE_ID\""; then
        echo -e "${GREEN}✓ Изображение найдено в partial_match${NC}"
    else
        echo -e "${RED}Ошибка: Изображение не найдено в partial_match${NC}"
        echo "Получили: $GET_PARTIAL"
        exit 1
    fi
else
    echo -e "${RED}Ошибка: Отсутствует поле partial_match в ответе${NC}"
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
  "$SERVER_URL/upload?user_id=0&tg_file_id=$PUBLIC_FILE_ID&file_type=video" \
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
else
    echo -e "${RED}✗ ОШИБКА: TG File ID различаются!${NC}"
    echo "Загружено: $TG_FILE_ID"
    echo "Получено:  $RETURNED_FILE_ID"
    exit 1
fi

echo -e "${YELLOW}7. Тест получения всех изображений пользователя...${NC}"

GET_USER_IMAGES=$(curl -s -X GET "$SERVER_URL/user/images?user_id=$USER_ID")

if [ $? -ne 0 ]; then
    echo -e "${RED}Ошибка: Не удалось получить изображения пользователя${NC}"
    exit 1
fi

echo "Ответ сервера: $GET_USER_IMAGES"

if ! echo "$GET_USER_IMAGES" | grep -q '"images"'; then
    echo -e "${RED}Ошибка: В ответе отсутствует поле images${NC}"
    exit 1
fi

if ! echo "$GET_USER_IMAGES" | grep -q "\"$TG_FILE_ID\""; then
    echo -e "${RED}Ошибка: Изображение пользователя не найдено${NC}"
    exit 1
fi

if ! echo "$GET_USER_IMAGES" | grep -q '"tags"'; then
    echo -e "${RED}Ошибка: В ответе отсутствуют теги${NC}"
    exit 1
fi

if echo "$GET_USER_IMAGES" | grep -q '"meme"' && echo "$GET_USER_IMAGES" | grep -q '"funny"' && echo "$GET_USER_IMAGES" | grep -q '"test"'; then
    echo -e "${GREEN}✓ Все теги присутствуют в ответе${NC}"
else
    echo -e "${RED}Ошибка: Не все теги найдены в ответе${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Endpoint /user/images работает корректно${NC}"

echo -e "${YELLOW}8. Тест замены тегов...${NC}"

# Загружаем изображение с тегами
TEST_REPLACE_FILE_ID="test_replace_file_id_99999"
HTTP_CODE=$(curl -s -w "%{http_code}" -o /dev/null -X POST \
  "$SERVER_URL/upload?user_id=$USER_ID&tg_file_id=$TEST_REPLACE_FILE_ID&file_type=photo" \
  -H "Content-Type: application/json" \
  -d '{"tags": ["old_tag1", "old_tag2"]}')

if [ "$HTTP_CODE" != "201" ]; then
    echo -e "${RED}Ошибка: Не удалось загрузить изображение для теста замены тегов${NC}"
    exit 1
fi

# Заменяем теги
HTTP_CODE=$(curl -s -w "%{http_code}" -o /tmp/replace_response.txt -X PUT \
  "$SERVER_URL/image/tags?user_id=$USER_ID&tg_file_id=$TEST_REPLACE_FILE_ID" \
  -H "Content-Type: application/json" \
  -d '{"tags": ["new_tag1", "new_tag2", "new_tag3"]}')

if [ "$HTTP_CODE" != "200" ]; then
    echo -e "${RED}Ошибка: Не удалось заменить теги${NC}"
    cat /tmp/replace_response.txt
    exit 1
fi

# Проверяем, что теги заменились
GET_REPLACED=$(curl -s -X GET "$SERVER_URL/user/images?user_id=$USER_ID")

if echo "$GET_REPLACED" | grep -q "\"$TEST_REPLACE_FILE_ID\""; then
    if echo "$GET_REPLACED" | grep -q '"new_tag1"' && echo "$GET_REPLACED" | grep -q '"new_tag2"' && echo "$GET_REPLACED" | grep -q '"new_tag3"'; then
        if ! echo "$GET_REPLACED" | grep -q '"old_tag1"' && ! echo "$GET_REPLACED" | grep -q '"old_tag2"'; then
            echo -e "${GREEN}✓ Теги успешно заменены${NC}"
        else
            echo -e "${RED}Ошибка: Старые теги не были удалены${NC}"
            exit 1
        fi
    else
        echo -e "${RED}Ошибка: Новые теги не найдены${NC}"
        exit 1
    fi
else
    echo -e "${RED}Ошибка: Изображение не найдено после замены тегов${NC}"
    exit 1
fi

echo -e "${YELLOW}9. Тест защиты публичных изображений от замены тегов...${NC}"

# Пытаемся заменить теги у публичного изображения (должно быть запрещено)
REPLACE_PUBLIC=$(curl -s -X PUT "$SERVER_URL/image/tags?user_id=0&tg_file_id=$PUBLIC_FILE_ID" \
  -H "Content-Type: application/json" \
  -d '{"tags": ["should_not_work"]}')

if echo "$REPLACE_PUBLIC" | grep -q "cannot replace tags for default user"; then
    echo -e "${GREEN}✓ Защита публичных изображений работает${NC}"
else
    echo -e "${RED}Ошибка: Удалось заменить теги публичного изображения${NC}"
    echo "Ответ: $REPLACE_PUBLIC"
    exit 1
fi

echo -e "${YELLOW}10. Тест удаления конкретной картинки...${NC}"

# Загружаем изображение для удаления
DELETE_FILE_ID="test_delete_file_id_88888"
HTTP_CODE=$(curl -s -w "%{http_code}" -o /dev/null -X POST \
  "$SERVER_URL/upload?user_id=$USER_ID&tg_file_id=$DELETE_FILE_ID&file_type=video" \
  -H "Content-Type: application/json" \
  -d '{"tags": ["to_delete"]}')

if [ "$HTTP_CODE" != "201" ]; then
    echo -e "${RED}Ошибка: Не удалось загрузить изображение для удаления${NC}"
    exit 1
fi

# Проверяем что изображение есть
GET_BEFORE_DELETE=$(curl -s -X GET "$SERVER_URL/user/images?user_id=$USER_ID")
if ! echo "$GET_BEFORE_DELETE" | grep -q "\"$DELETE_FILE_ID\""; then
    echo -e "${RED}Ошибка: Изображение не найдено перед удалением${NC}"
    exit 1
fi

# Удаляем изображение
HTTP_CODE=$(curl -s -w "%{http_code}" -o /tmp/delete_response.txt -X DELETE \
  "$SERVER_URL/image/delete?user_id=$USER_ID&tg_file_id=$DELETE_FILE_ID")

if [ "$HTTP_CODE" != "200" ]; then
    echo -e "${RED}Ошибка: Не удалось удалить изображение${NC}"
    cat /tmp/delete_response.txt
    exit 1
fi

# Проверяем что изображение удалено
GET_AFTER_DELETE=$(curl -s -X GET "$SERVER_URL/user/images?user_id=$USER_ID")
if echo "$GET_AFTER_DELETE" | grep -q "\"$DELETE_FILE_ID\""; then
    echo -e "${RED}Ошибка: Изображение все еще существует после удаления${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Удаление конкретной картинки работает корректно${NC}"

echo -e "${YELLOW}11. Тест удаления всех картинок пользователя...${NC}"

# Создаем тестового пользователя с несколькими изображениями
TEST_DELETE_USER=777777777
curl -s -X POST "$SERVER_URL/upload?user_id=$TEST_DELETE_USER&tg_file_id=file_1&file_type=photo" \
  -H "Content-Type: application/json" -d '{"tags": ["tag1"]}' > /dev/null

curl -s -X POST "$SERVER_URL/upload?user_id=$TEST_DELETE_USER&tg_file_id=file_2&file_type=photo" \
  -H "Content-Type: application/json" -d '{"tags": ["tag2"]}' > /dev/null

curl -s -X POST "$SERVER_URL/upload?user_id=$TEST_DELETE_USER&tg_file_id=file_3&file_type=photo" \
  -H "Content-Type: application/json" -d '{"tags": ["tag3"]}' > /dev/null

# Проверяем что у пользователя есть изображения
GET_BEFORE_DELETE_ALL=$(curl -s -X GET "$SERVER_URL/user/images?user_id=$TEST_DELETE_USER")
IMAGE_COUNT=$(echo "$GET_BEFORE_DELETE_ALL" | grep -o "tg_file_id" | wc -l)

if [ "$IMAGE_COUNT" -lt 3 ]; then
    echo -e "${RED}Ошибка: Недостаточно изображений перед удалением (найдено: $IMAGE_COUNT)${NC}"
    exit 1
fi

# Удаляем все изображения пользователя
HTTP_CODE=$(curl -s -w "%{http_code}" -o /tmp/delete_all_response.txt -X DELETE \
  "$SERVER_URL/user/images/delete?user_id=$TEST_DELETE_USER")

if [ "$HTTP_CODE" != "200" ]; then
    echo -e "${RED}Ошибка: Не удалось удалить все изображения пользователя${NC}"
    cat /tmp/delete_all_response.txt
    exit 1
fi

# Проверяем что у пользователя нет изображений
GET_AFTER_DELETE_ALL=$(curl -s -X GET "$SERVER_URL/user/images?user_id=$TEST_DELETE_USER")
if echo "$GET_AFTER_DELETE_ALL" | grep -q '"images":\[\]' || echo "$GET_AFTER_DELETE_ALL" | grep -q '"images": \[\]'; then
    echo -e "${GREEN}✓ Удаление всех картинок пользователя работает корректно${NC}"
else
    echo -e "${RED}Ошибка: У пользователя остались изображения после удаления${NC}"
    echo "Ответ: $GET_AFTER_DELETE_ALL"
    exit 1
fi

echo -e "${YELLOW}12. Тест генерации описания...${NC}"

# Проверяем генерацию описания
GEN_DESC=$(curl -s -X POST "$SERVER_URL/image/generate-description" \
  -F "image=@test/test_image.jpg")

if echo "$GEN_DESC" | grep -q "метод пока не реализован"; then
    echo -e "${GREEN}✓ Генерация описания возвращает корректный ответ${NC}"
else
    echo -e "${RED}Ошибка: Генерация описания вернула неожиданный ответ${NC}"
    echo "Ответ: $GEN_DESC"
    exit 1
fi

echo -e "${YELLOW}13. Тест счетчика использования...${NC}"

# Загружаем тестовое изображение для проверки счетчика
USAGE_TEST_FILE_ID="test_usage_counter_file_id"
HTTP_CODE=$(curl -s -w "%{http_code}" -o /dev/null -X POST \
  "$SERVER_URL/upload?user_id=$USER_ID&tg_file_id=$USAGE_TEST_FILE_ID&file_type=photo" \
  -H "Content-Type: application/json" \
  -d '{"tags": ["usage_test"]}')

if [ "$HTTP_CODE" != "201" ]; then
    echo -e "${RED}Ошибка: Не удалось загрузить изображение для теста счетчика${NC}"
    exit 1
fi

# Увеличиваем счетчик первый раз
HTTP_CODE=$(curl -s -w "%{http_code}" -o /tmp/usage_response.txt -X POST \
  "$SERVER_URL/image/increment-usage?user_id=$USER_ID&tg_file_id=$USAGE_TEST_FILE_ID")

if [ "$HTTP_CODE" != "200" ]; then
    echo -e "${RED}Ошибка: Не удалось увеличить счетчик использования${NC}"
    cat /tmp/usage_response.txt
    exit 1
fi

RESPONSE=$(cat /tmp/usage_response.txt)
if echo "$RESPONSE" | grep -q "Usage count incremented successfully"; then
    echo -e "${GREEN}✓ Счетчик успешно увеличен (1-й раз)${NC}"
else
    echo -e "${RED}Ошибка: Неожиданный ответ при увеличении счетчика${NC}"
    exit 1
fi

# Увеличиваем счетчик второй раз
HTTP_CODE=$(curl -s -w "%{http_code}" -o /dev/null -X POST \
  "$SERVER_URL/image/increment-usage?user_id=$USER_ID&tg_file_id=$USAGE_TEST_FILE_ID")

if [ "$HTTP_CODE" != "200" ]; then
    echo -e "${RED}Ошибка: Не удалось увеличить счетчик второй раз${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Счетчик успешно увеличен (2-й раз)${NC}"

# Увеличиваем счетчик третий раз
HTTP_CODE=$(curl -s -w "%{http_code}" -o /dev/null -X POST \
  "$SERVER_URL/image/increment-usage?user_id=$USER_ID&tg_file_id=$USAGE_TEST_FILE_ID")

if [ "$HTTP_CODE" != "200" ]; then
    echo -e "${RED}Ошибка: Не удалось увеличить счетчик третий раз${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Счетчик успешно увеличен (3-й раз)${NC}"

# Пытаемся увеличить счетчик для несуществующего изображения (должна быть ошибка)
HTTP_CODE=$(curl -s -w "%{http_code}" -o /tmp/usage_error.txt -X POST \
  "$SERVER_URL/image/increment-usage?user_id=$USER_ID&tg_file_id=nonexistent_file_id")

if [ "$HTTP_CODE" == "500" ]; then
    if grep -q "image not found" /tmp/usage_error.txt; then
        echo -e "${GREEN}✓ Корректная обработка ошибки для несуществующего изображения${NC}"
    else
        echo -e "${RED}Ошибка: Неожиданное сообщение об ошибке${NC}"
        cat /tmp/usage_error.txt
        exit 1
    fi
else
    echo -e "${RED}Ошибка: Ожидался HTTP 500 для несуществующего изображения, получен $HTTP_CODE${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Счетчик использования работает корректно${NC}"

echo -e "${GREEN}=== Все тесты пройдены успешно ===${NC}"

# Очистка
rm -f /tmp/upload_response.txt /tmp/replace_response.txt /tmp/delete_response.txt /tmp/delete_all_response.txt
exit 0
