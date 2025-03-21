# CompareNames

Библиотека и API для нечёткого сравнения имён (ФИО) на русском и английском языках с учётом транслитерации. Позволяет определить, относятся ли два разных написания имени к одному и тому же человеку.

![GitHub](https://img.shields.io/github/license/x0rium/compareNames)
![Go Version](https://img.shields.io/badge/go-%3E%3D1.18-blue)

## 🚀 Возможности

- **Мультиязычность**: Сравнение имён на кириллице и латинице с автоматической транслитерацией
- **Устойчивость к вариациям**:
  - Поддержка украинского алфавита и имён
  - Обработка двойных фамилий (через дефис)
  - Устойчивость к опечаткам и неточностям
  - Сравнение имён в разном порядке (ФИО, ИФО и т.д.)
  - Корректная обработка инициалов
- **Гибкость настройки**:
  - Настраиваемые веса алгоритмов для разных сценариев
  - Возможность использования дополнительных атрибутов для повышения точности
- **Производительность**:
  - Кэширование результатов и вариаций имён
  - Оптимизированные алгоритмы сравнения
- **Интеграция**:
  - REST API для использования в других системах
  - Демонстрационное CLI приложение

## 🧠 Алгоритмы сравнения

Библиотека использует комбинацию нескольких алгоритмов для достижения максимальной точности:

| Алгоритм | Описание | Преимущества |
|----------|----------|--------------|
| **Расстояние Левенштейна** | Измеряет количество операций редактирования для преобразования одной строки в другую | Хорошо обрабатывает опечатки и замены символов |
| **Расстояние Джаро-Винклера** | Даёт более высокую оценку строкам с общим префиксом | Эффективен для имён с одинаковым началом |
| **Фонетические алгоритмы** | Учитывают схожесть произношения слов | Устойчивость к фонетическим вариациям |
| **Double Metaphone** | Улучшенный фонетический алгоритм с поддержкой разных языков | Работает с многоязычными именами |
| **Косинусное сходство** | Сравнивает сходство наборов n-грамм | Эффективно при перестановке слов |

## 📂 Структура проекта

```
compareNames/
├── matcher/           # Основная библиотека сравнения имён
│   ├── similarity/    # Реализация алгоритмов сравнения строк
│   ├── translit/      # Функции для транслитерации
│   └── testdata/      # Тестовые данные
├── api/               # REST API для сравнения имён
└── cmd/demo/          # Демонстрационное CLI приложение
```

## 📋 Примеры использования

### Базовое сравнение имён

```go
package main

import (
	"fmt"

	"github.com/x0rium/compareNames/matcher"
)

func main() {
	name1 := "Иванов Иван Иванович"
	name2 := "Ivanov Ivan"

	// Сравнение с настройками по умолчанию
	result := matcher.MatchNames(name1, name2, nil, nil)

	fmt.Printf("Результат сравнения: %s (оценка: %d%%)\n", 
		result.MatchType, result.Score)
	fmt.Printf("Наилучшее соответствие: %s <-> %s\n", 
		result.BestMatch1, result.BestMatch2)
}
```

### Сравнение с дополнительными атрибутами

```go
package main

import (
	"fmt"

	"github.com/x0rium/compareNames/matcher"
)

func main() {
	name1 := "Иванов Иван Иванович"
	name2 := "Ivanov Ivan"
	
	// Создаем атрибуты для сравнения
	attrs := matcher.CreateAttributes()
	matcher.AddAttribute(attrs, "birth_date", true)  // Даты рождения совпадают
	matcher.AddAttribute(attrs, "country", false)    // Страны не совпадают
	
	// Сравнение с атрибутами
	result := matcher.MatchNamesWithAttributes(name1, name2, attrs)
	
	matcher.PrintMatchResult(result)
}
```

### Настройка параметров сравнения

```go
package main

import (
	"fmt"

	"github.com/x0rium/compareNames/matcher"
)

func main() {
	name1 := "Иванов Иван Иванович"
	name2 := "Ivanov Ivan"
	
	// Создаем пользовательскую конфигурацию
	config := matcher.DefaultConfig()
	config.LevenshteinWeight = 0.4       // Увеличиваем вес Левенштейна
	config.PhoneticWeight = 0.25         // Увеличиваем вес фонетического сравнения
	config.EnableCaching = true          // Включаем кэширование
	
	// Сравнение с пользовательской конфигурацией
	result := matcher.MatchNames(name1, name2, nil, &config)
	
	matcher.PrintMatchResult(result)
}
```

## 🌐 REST API

### Запуск API сервера

```bash
go run main.go -port 8080
```

### Использование API

**Endpoint**: `/api/match_names` (POST)

**Пример запроса**:

```json
{
  "name1": "Иванов Иван Иванович",
  "name2": "Ivanov Ivan",
  "attributes": {
    "birth_date": { "match": true },
    "country": { "match": false }
  },
  "config": {
    "levenshtein_weight": 0.3,
    "jaro_winkler_weight": 0.25,
    "phonetic_weight": 0.2,
    "cosine_weight": 0.15,
    "additional_attributes_weight": 0.1,
    "ngram_size": 3,
    "enable_caching": true
  },
  "disable_cache": false
}
```

**Пример ответа**:

```json
{
  "exact_match": false,
  "score": 85,
  "match_type": "match",
  "best_match1": "иванов иван иванович",
  "best_match2": "ivanov ivan",
  "levenshtein_score": 0.8,
  "jaro_winkler_score": 0.9,
  "phonetic_score": 0.95,
  "cosine_score": 0.87,
  "additional_attributes_score": 0.5,
  "processing_time_ms": 3,
  "from_cache": false
}
```

## 🖥️ Демонстрационное приложение

Для запуска демонстрационного приложения:

```bash
go run cmd/demo/demo.go
```

Демо-приложение предлагает два режима работы:
1. Запуск предустановленных примеров
2. Интерактивный режим для ручного сравнения имён

## 🛠️ Установка и сборка

### Установка библиотеки

```bash
go get github.com/x0rium/compareNames
```

### Сборка приложений

```bash
# Сборка API сервера
go build -o name-matching-api main.go

# Сборка демо-приложения
go build -o name-matching-demo cmd/demo/demo.go
```

## 📊 Интерпретация результатов

| Тип совпадения | Оценка | Описание |
|----------------|--------|----------|
| `exact_match`  | ≥ 90   | Имена считаются идентичными |
| `match`        | 70-89  | Имена, вероятно, относятся к одному человеку |
| `possible_match` | 50-69 | Имена могут относиться к одному человеку, требуется дополнительная проверка |
| `no_match`     | < 50   | Имена, вероятно, относятся к разным людям |

## 📄 Лицензия

MIT
