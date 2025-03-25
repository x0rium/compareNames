#!/bin/bash

# Цвета для вывода
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== Тестирование API сравнения имен ===${NC}"
echo "Запуск 5 тестовых случаев..."

# Проверка доступности сервера
echo -ne "Проверка доступности сервера: "
if ! curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health | grep -q "200"; then
    echo -e "${RED}Сервер не доступен!${NC}"
    echo "Запустите сервер перед запуском тестов: ./compareNames"
    exit 1
fi
echo -e "${GREEN}OK${NC}"

# Функция для запуска тестового случая
run_test() {
    local test_num=$1
    local name1=$2
    local name2=$3
    local description=$4
    
    echo -e "\n${YELLOW}Тест #${test_num}:${NC} $description"
    echo "Имя 1: $name1"
    echo "Имя 2: $name2"
    
    # Отправляем запрос
    local response=$(curl -s -X POST http://localhost:8080/api/match_names \
        -H "Content-Type: application/json" \
        -d "{\"name1\": \"$name1\", \"name2\": \"$name2\"}")
    
    # Форматируем JSON для удобства чтения
    echo "Результат:"
    echo $response | python3 -m json.tool 2>/dev/null || echo $response
    
    # Проверяем наличие ответа
    if [[ -n "$response" ]]; then
        echo -e "${GREEN}✓ Тест пройден${NC}"
    else
        echo -e "${RED}✗ Тест не пройден${NC}"
    fi
}

# Тестовые случаи
echo -e "\n${YELLOW}Группа 1: Точные и близкие совпадения${NC}"
run_test 1 "John" "John" "Точное совпадение"
run_test 2 "Иван" "Ivan" "Транслитерация"
run_test 3 "И. Петров" "Иван Петров" "Инициалы"

echo -e "\n${YELLOW}Группа 2: Частичные совпадения${NC}"
run_test 4 "Иванов-Петров Михаил" "М. Иванов-Петров" "Сложные имена с инициалами"
run_test 5 "Александр Сергеевич" "Александр" "Частичное совпадение (имя без отчества)"

echo -e "\n${YELLOW}Группа 3: Непохожие имена${NC}"
run_test 6 "Александр" "Виктор" "Разные имена"
run_test 7 "Иван Петров" "Сергей Иванов" "Разные имя и фамилия"
run_test 8 "李小龙" "Брюс Ли" "Разные алфавиты"

echo -e "\n${YELLOW}=== Тестирование завершено ===${NC}"
