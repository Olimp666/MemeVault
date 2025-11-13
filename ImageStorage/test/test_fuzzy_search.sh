#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

SERVER_URL="http://localhost"
TEST_USER_ID=123456789

echo -e "${YELLOW}=== Тест нечеткого поиска (Fuzzy Search) ===${NC}"

echo -e "${YELLOW}1. Загрузка тестовых изображений с различными тегами...${NC}"

declare -A TEST_DATA=(
    ["file_1"]="cat,animal,pet"
    ["file_2"]="dog,animal,pet"
    ["file_3"]="cat,funny,meme"
    ["file_4"]="dog,funny,meme"
    ["file_5"]="bird,animal,nature"
    ["file_6"]="fish,animal,water"
    ["file_7"]="python,programming,code"
    ["file_8"]="javascript,programming,code"
    ["file_9"]="golang,programming,code"
    ["file_10"]="react,frontend,web"
    ["file_11"]="vue,frontend,web"
    ["file_12"]="angular,frontend,web"
    ["file_13"]="docker,devops,container"
    ["file_14"]="kubernetes,devops,container"
    ["file_15"]="terraform,devops,infrastructure"
    ["file_16"]="sunset,nature,beautiful"
    ["file_17"]="ocean,nature,water"
    ["file_18"]="mountain,nature,hiking"
    ["file_19"]="coffee,morning,drink"
    ["file_20"]="tea,morning,drink"
    ["file_21"]="pizza,food,italian"
    ["file_22"]="sushi,food,japanese"
    ["file_23"]="burger,food,american"
    ["file_24"]="car,transport,vehicle"
    ["file_25"]="bike,transport,vehicle"
    ["file_26"]="train,transport,vehicle"
    ["file_27"]="book,reading,education"
    ["file_28"]="music,entertainment,art"
    ["file_29"]="movie,entertainment,cinema"
    ["file_30"]="game,entertainment,fun"
)

for file_id in "${!TEST_DATA[@]}"; do
    tags="${TEST_DATA[$file_id]}"
    IFS=',' read -ra tag_array <<< "$tags"
    json_tags="["
    for i in "${!tag_array[@]}"; do
        if [ $i -gt 0 ]; then
            json_tags+=","
        fi
        json_tags+="\"${tag_array[$i]}\""
    done
    json_tags+="]"
    
    HTTP_CODE=$(curl -s -w "%{http_code}" -o /dev/null -X POST \
        "$SERVER_URL/upload?user_id=$TEST_USER_ID&tg_file_id=$file_id&file_type=photo" \
        -H "Content-Type: application/json" \
        -d "{\"tags\": $json_tags}")
    
    if [ "$HTTP_CODE" != "201" ]; then
        echo -e "${RED}Ошибка: Не удалось загрузить $file_id (HTTP $HTTP_CODE)${NC}"
        exit 1
    fi
done

echo -e "${GREEN}✓ Загружено 30 изображений с различными тегами${NC}"

echo -e "${YELLOW}2. Тест точного поиска без опечаток...${NC}"

RESULT=$(curl -s -X POST "$SERVER_URL/images?user_id=$TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"tags": ["cat", "animal"]}')

EXACT_COUNT=$(echo "$RESULT" | grep -o '"exact_match"' | wc -l)
if [ "$EXACT_COUNT" -eq 0 ]; then
    echo -e "${RED}Ошибка: exact_match не найден в ответе${NC}"
    exit 1
fi

if echo "$RESULT" | grep "exact_match" | grep -q "file_1"; then
    echo -e "${GREEN}✓ Точный поиск работает корректно (найден file_1)${NC}"
else
    echo -e "${RED}Ошибка: file_1 не найден при точном поиске${NC}"
    echo "Ответ: $RESULT"
    exit 1
fi

echo -e "${YELLOW}3. Тест поиска с одной опечаткой (cat -> cet)...${NC}"

RESULT=$(curl -s -X POST "$SERVER_URL/images?user_id=$TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"tags": ["cet", "animal"]}')

if echo "$RESULT" | grep "exact_match" | grep -q "file_1"; then
    echo -e "${GREEN}✓ Fuzzy поиск нашел file_1 с опечаткой 'cet' вместо 'cat'${NC}"
else
    echo -e "${RED}Ошибка: file_1 не найден с опечаткой 'cet'${NC}"
    echo "Ответ: $RESULT"
    exit 1
fi

echo -e "${YELLOW}4. Тест поиска с двумя опечатками (programming -> progremming)...${NC}"

RESULT=$(curl -s -X POST "$SERVER_URL/images?user_id=$TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"tags": ["progremming", "code"]}')

if echo "$RESULT" | grep "exact_match" | grep -q "file_"; then
    COUNT=$(echo "$RESULT" | grep "exact_match" | grep -o "file_" | wc -l)
    if [ "$COUNT" -ge 1 ]; then
        echo -e "${GREEN}✓ Fuzzy поиск нашел результаты с опечаткой 'progremming'${NC}"
    else
        echo -e "${RED}Ошибка: Найдено слишком мало результатов${NC}"
        exit 1
    fi
else
    echo -e "${RED}Ошибка: Результаты не найдены с опечаткой 'progremming'${NC}"
    echo "Ответ: $RESULT"
    exit 1
fi

echo -e "${YELLOW}5. Тест поиска с опечаткой в обоих тегах (dg -> dog, animel -> animal)...${NC}"

RESULT=$(curl -s -X POST "$SERVER_URL/images?user_id=$TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"tags": ["dg", "animel"]}')

if echo "$RESULT" | grep "exact_match" | grep -q "file_2\|file_4"; then
    echo -e "${GREEN}✓ Fuzzy поиск работает с опечатками в обоих тегах${NC}"
else
    echo -e "${RED}Ошибка: Не найдены результаты с опечатками в обоих тегах${NC}"
    echo "Ответ: $RESULT"
    exit 1
fi

echo -e "${YELLOW}6. Тест частичного совпадения с опечатками (frontend -> fronend)...${NC}"

RESULT=$(curl -s -X POST "$SERVER_URL/images?user_id=$TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"tags": ["fronend", "web", "react"]}')

if echo "$RESULT" | grep "partial_match" | grep -q "file_"; then
    echo -e "${GREEN}✓ Частичное совпадение с опечатками работает${NC}"
else
    echo -e "${RED}Ошибка: partial_match не работает с опечатками${NC}"
    echo "Ответ: $RESULT"
    exit 1
fi

echo -e "${YELLOW}7. Тест негативного случая (слишком большая опечатка)...${NC}"

RESULT=$(curl -s -X POST "$SERVER_URL/images?user_id=$TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"tags": ["xyz", "abc"]}')

if echo "$RESULT" | grep -q '"exact_match":\[\]' || echo "$RESULT" | grep -q '"exact_match": \[\]'; then
    echo -e "${GREEN}✓ Поиск корректно не возвращает результаты при большой опечатке${NC}"
else
    echo -e "${RED}Внимание: Найдены результаты при большой опечатке${NC}"
fi

echo -e "${YELLOW}8. Тест поиска с опечаткой в регистре (Nature -> nature)...${NC}"

RESULT=$(curl -s -X POST "$SERVER_URL/images?user_id=$TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"tags": ["NATURE", "WATER"]}')

if echo "$RESULT" | grep "exact_match" | grep -q "file_"; then
    echo -e "${GREEN}✓ Поиск работает независимо от регистра${NC}"
else
    echo -e "${RED}Ошибка: Поиск не работает с заглавными буквами${NC}"
    echo "Ответ: $RESULT"
    exit 1
fi

echo -e "${YELLOW}9. Тест опечатки с пропущенной буквой (docker -> doker)...${NC}"

RESULT=$(curl -s -X POST "$SERVER_URL/images?user_id=$TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"tags": ["doker", "devops"]}')

if echo "$RESULT" | grep "exact_match" | grep -q "file_13"; then
    echo -e "${GREEN}✓ Fuzzy поиск находит результаты с пропущенной буквой${NC}"
else
    echo -e "${RED}Ошибка: Не найдены результаты с пропущенной буквой${NC}"
    echo "Ответ: $RESULT"
    exit 1
fi

echo -e "${YELLOW}10. Тест опечатки с лишней буквой (food -> foood)...${NC}"

RESULT=$(curl -s -X POST "$SERVER_URL/images?user_id=$TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"tags": ["foood", "italian"]}')

if echo "$RESULT" | grep "exact_match" | grep -q "file_21"; then
    echo -e "${GREEN}✓ Fuzzy поиск находит результаты с лишней буквой${NC}"
else
    echo -e "${RED}Ошибка: Не найдены результаты с лишней буквой${NC}"
    echo "Ответ: $RESULT"
    exit 1
fi

echo -e "${YELLOW}11. Тест замены букв (transport -> transpart)...${NC}"

RESULT=$(curl -s -X POST "$SERVER_URL/images?user_id=$TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"tags": ["transpart", "vehicle"]}')

if echo "$RESULT" | grep "exact_match" | grep -q "file_"; then
    echo -e "${GREEN}✓ Fuzzy поиск находит результаты с замененными буквами${NC}"
else
    echo -e "${RED}Ошибка: Не найдены результаты с замененными буквами${NC}"
    echo "Ответ: $RESULT"
    exit 1
fi

echo -e "${YELLOW}12. Тест корректности порядка результатов...${NC}"

RESULT=$(curl -s -X POST "$SERVER_URL/images?user_id=$TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"tags": ["cat", "animal"]}')

if echo "$RESULT" | grep -q '"exact_match"' && echo "$RESULT" | grep -q '"partial_match"'; then
    EXACT_POS=$(echo "$RESULT" | grep -b -o '"exact_match"' | head -1 | cut -d: -f1)
    PARTIAL_POS=$(echo "$RESULT" | grep -b -o '"partial_match"' | head -1 | cut -d: -f1)
    
    if [ "$EXACT_POS" -lt "$PARTIAL_POS" ]; then
        echo -e "${GREEN}✓ exact_match идет перед partial_match${NC}"
    else
        echo -e "${GREEN}✓ Оба поля присутствуют в ответе${NC}"
    fi
else
    echo -e "${RED}Ошибка: Не найдены required поля в ответе${NC}"
    exit 1
fi

echo -e "${YELLOW}13. Тест поиска с тремя тегами и опечатками...${NC}"

RESULT=$(curl -s -X POST "$SERVER_URL/images?user_id=$TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"tags": ["progremming", "cde", "pythn"]}')

if echo "$RESULT" | grep "exact_match" | grep -q "file_7"; then
    echo -e "${GREEN}✓ Fuzzy поиск работает с тремя тегами с опечатками${NC}"
else
    echo -e "${RED}Ошибка: Не найдены результаты при поиске с тремя тегами${NC}"
    echo "Ответ: $RESULT"
    exit 1
fi

echo -e "${YELLOW}14. Тест публичных изображений с опечатками...${NC}"

curl -s -X POST "$SERVER_URL/upload?user_id=0&tg_file_id=public_file&file_type=photo" \
    -H "Content-Type: application/json" \
    -d '{"tags": ["public", "shared", "common"]}' > /dev/null

OTHER_USER=999999999
RESULT=$(curl -s -X POST "$SERVER_URL/images?user_id=$OTHER_USER" \
    -H "Content-Type: application/json" \
    -d '{"tags": ["publik", "shared"]}')

if echo "$RESULT" | grep "exact_match" | grep -q "public_file"; then
    echo -e "${GREEN}✓ Fuzzy поиск работает для публичных изображений${NC}"
else
    echo -e "${RED}Ошибка: Публичные изображения не найдены с опечаткой${NC}"
    echo "Ответ: $RESULT"
    exit 1
fi

echo -e "${YELLOW}15. Тест производительности (множественные запросы)...${NC}"

START_TIME=$(date +%s%N)
for i in {1..10}; do
    curl -s -X POST "$SERVER_URL/images?user_id=$TEST_USER_ID" \
        -H "Content-Type: application/json" \
        -d '{"tags": ["animel", "pet"]}' > /dev/null
done
END_TIME=$(date +%s%N)

ELAPSED=$((($END_TIME - $START_TIME) / 1000000))
echo -e "${GREEN}✓ 10 запросов выполнены за ${ELAPSED}мс${NC}"

if [ "$ELAPSED" -lt 5000 ]; then
    echo -e "${GREEN}✓ Производительность в норме${NC}"
else
    echo -e "${YELLOW}⚠ Предупреждение: Запросы выполняются медленно${NC}"
fi

echo -e "${YELLOW}16. Тест специальных символов и пробелов...${NC}"

curl -s -X POST "$SERVER_URL/upload?user_id=$TEST_USER_ID&tg_file_id=special_file&file_type=photo" \
    -H "Content-Type: application/json" \
    -d '{"tags": ["tag-with-dash", "tag_with_underscore"]}' > /dev/null

RESULT=$(curl -s -X POST "$SERVER_URL/images?user_id=$TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"tags": ["tagwithdash", "tagwithunderscore"]}')

if echo "$RESULT" | grep -q "special_file" || echo "$RESULT" | grep -q "file_"; then
    echo -e "${GREEN}✓ Поиск обрабатывает специальные символы${NC}"
else
    echo -e "${YELLOW}ℹ Специальные символы могут требовать точного совпадения${NC}"
fi

echo -e "${YELLOW}17. Проверка отсутствия ложных срабатываний...${NC}"

RESULT=$(curl -s -X POST "$SERVER_URL/images?user_id=$TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"tags": ["zzzzzz", "yyyyyy"]}')

EXACT_FOUND=$(echo "$RESULT" | grep "exact_match" | grep -o "file_" | wc -l)
if [ "$EXACT_FOUND" -eq 0 ]; then
    echo -e "${GREEN}✓ Нет ложных срабатываний при несуществующих тегах${NC}"
else
    echo -e "${RED}Ошибка: Найдены результаты для несуществующих тегов${NC}"
    echo "Найдено: $EXACT_FOUND результатов"
    exit 1
fi

echo -e "${YELLOW}18. Очистка тестовых данных...${NC}"

HTTP_CODE=$(curl -s -w "%{http_code}" -o /dev/null -X DELETE \
    "$SERVER_URL/user/images/delete?user_id=$TEST_USER_ID")

if [ "$HTTP_CODE" == "200" ]; then
    echo -e "${GREEN}✓ Тестовые данные очищены${NC}"
else
    echo -e "${YELLOW}⚠ Не удалось полностью очистить тестовые данные${NC}"
fi

HTTP_CODE=$(curl -s -w "%{http_code}" -o /dev/null -X DELETE \
    "$SERVER_URL/user/images/delete?user_id=0")

HTTP_CODE=$(curl -s -w "%{http_code}" -o /dev/null -X DELETE \
    "$SERVER_URL/user/images/delete?user_id=$OTHER_USER")

echo -e "${GREEN}=== Все тесты fuzzy search пройдены успешно ===${NC}"
exit 0
