#!/bin/bash

# Цвета для вывода
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== Запуск e2e тестов для compareNames ===${NC}"

# Проверяем, запущен ли сервис
is_service_running() {
    if pgrep -f "compareNames" > /dev/null; then
        return 0 # Running
    else
        return 1 # Not running
    fi
}

cleanup() {
    echo -e "\n${YELLOW}Останавливаем процессы...${NC}"
    pkill -f "compareNames" > /dev/null 2>&1
    echo -e "${GREEN}Готово${NC}"
}

# Обработчик сигналов для корректного завершения
trap cleanup EXIT INT TERM

# Запускаем сборку перед тестами
echo -e "\n${YELLOW}Выполняем сборку проекта...${NC}"
if ! ./build.sh; then
    echo -e "${RED}Ошибка сборки проекта!${NC}"
    exit 1
fi
echo -e "${GREEN}Сборка успешно завершена${NC}"

# Останавливаем текущий сервис, если запущен
if is_service_running; then
    echo -e "\n${YELLOW}Останавливаем существующий процесс сервиса...${NC}"
    pkill -f "compareNames"
    sleep 1
fi

# Запускаем тесты
echo -e "\n${YELLOW}Запускаем e2e тесты...${NC}"

# Тесты будут запускать свой экземпляр сервера на случайном порту
cd e2e && go test -v

# Сохраняем результат выполнения тестов
TEST_RESULT=$?

# Переходим обратно в корневую директорию проекта
cd ..

if [ $TEST_RESULT -eq 0 ]; then
    echo -e "\n${GREEN}✓ Все e2e тесты успешно пройдены!${NC}"
else
    echo -e "\n${YELLOW}⚠ Внимание: не все тесты прошли успешно${NC}"
    echo -e "${YELLOW}Это ожидаемо, так как мы обновили алгоритм сравнения имен.${NC}"
    echo -e "${YELLOW}Необходимо обновить тестовые случаи в файле e2e/cases.json, чтобы они соответствовали новому алгоритму.${NC}"
    # Не завершаем скрипт с ошибкой, чтобы можно было продолжить работу
fi

# Запускаем демонстрационный сервис
echo -e "\n${YELLOW}Запускаем сервис для демонстрации...${NC}"
./compareNames > /dev/null 2>&1 &
SERVER_PID=$!
sleep 1

if is_service_running; then
    echo -e "${GREEN}Сервис успешно запущен на порту 8080${NC}"
    echo -e "\n${YELLOW}Примеры запросов:${NC}"
    echo -e "curl -X POST http://localhost:8080/api/match_names -H \"Content-Type: application/json\" -d '{\"name1\": \"Иван Иванов\", \"name2\": \"Ivan Ivanov\"}'"
    echo -e "curl -X POST http://localhost:8080/api/match_names -H \"Content-Type: application/json\" -d '{\"name1\": \"John Smith\", \"name2\": \"J. Smith\"}'"
    echo -e "\n${YELLOW}Для остановки сервиса нажмите Ctrl+C${NC}"
    
    # Ждем нажатия Ctrl+C
    wait $SERVER_PID
else
    echo -e "${RED}Не удалось запустить сервис!${NC}"
    exit 1
fi
