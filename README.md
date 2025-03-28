# CompareNames

Библиотека и API для нечёткого сравнения имён (ФИО) на русском и английском языках с учётом транслитерации. Позволяет определить, относятся ли два разных написания имени к одному и тому же человеку.

![GitHub](https://img.shields.io/github/license/x0rium/compareNames)
![Go Version](https://img.shields.io/badge/go-%3E%3D1.18-blue)

## 🚀 Возможности

- **Мультиязычность**: Сравнение имён на кириллице и латинице с автоматической транслитерацией
- **Устойчивость к вариациям**:
  - Поддержка различных стандартов транслитерации (ISO 9, ГОСТ 7.79-2000, BGN/PCGN, UNGEGN)
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

## 🔍 Алгоритм работы

Программа выполняет следующие шаги при сравнении имён:

1. **Предобработка имён**:
   - Приведение к нижнему регистру
   - Удаление лишних символов
   - Преобразование дефисов в пробелы для корректной обработки двойных фамилий
   - Удаление лишних пробелов

2. **Проверка на инициалы**:
   - Если одно из имён содержит инициалы (например, "И.И." или одиночные буквы), применяется специальная логика сопоставления
   - Инициалы сравниваются с первыми буквами полного имени
   - Учитываются возможные транслитерации инициалов

3. **Определение алфавита**:
   - Определяется, написаны ли имена кириллицей или латиницей
   - Если алфавиты отличаются, используется специальная логика сравнения с транслитерацией

4. **Сравнение имён на разных алфавитах** (если применимо):
   - Генерируются возможные транслитерации согласно различным стандартам
   - Сравниваются все возможные варианты транслитерации
   - Выбирается наилучшее соответствие

5. **Сравнение имён на одном алфавите**:
   - Применяются алгоритмы расстояния Левенштейна и Джаро-Винклера
   - Применяются фонетические алгоритмы (Soundex, Double Metaphone)
   - Вычисляется оценка косинусного сходства

6. **Вычисление итоговой оценки**:
   - Базовая оценка: взвешенное среднее всех метрик (Левенштейн, Джаро-Винклер, фонетические, косинусная)
   - Бонусы за транслитерацию между алфавитами (до 12%)
   - Бонусы за перестановки частей ФИО (до 12%)
   - Бонусы за обработку инициалов (до 15%)
   - Бонусы за обработку дефисных имен (до 8%)
   - Бонусы за уменьшительные/альтернативные формы имён (до 12%)

7. **Определение типа совпадения**:
   - На основе итоговой оценки определяется `matchType` (см. раздел "Интерпретация результатов")

8. **Логирование и возврат результата**:
   - Сомнительные совпадения (possible_match) логируются для дальнейшего анализа
   - Результат возвращается с детальной информацией о метриках

## ⚙️ Настройка конфигурации

Библиотека CompareNames предоставляет гибкую систему настройки, которая позволяет адаптировать алгоритм сравнения под конкретные сценарии использования. Вы можете изменить веса различных алгоритмов, пороговые значения для определения типа совпадения, включить или отключить определённые функции, и многое другое.

### Изменение конфигурации

#### В коде (Go)

```go
package main

import (
	"github.com/x0rium/compareNames/matcher"
)

func main() {
	// Создаём конфигурацию с параметрами по умолчанию
	config := matcher.DefaultConfig()
	
	// Изменяем параметры конфигурации
	config.LevenshteinWeight = 0.3            // Увеличиваем вес Левенштейна
	config.PhoneticWeight = 0.2               // Уменьшаем вес фонетического сравнения
	config.JaroWinklerThreshold = 0.8         // Снижаем порог для Джаро-Винклера
	config.NGramSize = 2                      // Уменьшаем размер n-грамм
	config.EnableNamePartPermutation = false  // Отключаем учёт перестановок
	
	// Используем изменённую конфигурацию при сравнении имён
	result := matcher.MatchNames("Иванов Иван", "Ivan Ivanov", nil, &config)
}
```

#### Через API (JSON)

```json
{
  "name1": "Иванов Иван Иванович",
  "name2": "Ivanov Ivan",
  "config": {
    "levenshtein_weight": 0.3,
    "jaro_winkler_weight": 0.3,
    "phonetic_weight": 0.2,
    "double_metaphone_weight": 0.2,
    "jaro_winkler_threshold": 0.8,
    "enable_transliteration": true,
    "transliteration_standards": ["gost", "iso9"],
    "enable_name_part_permutation": true,
    "ngram_size": 3,
    "enable_caching": true,
    "enable_logging": true
  }
}
```

### Параметры конфигурации

#### Веса алгоритмов

Определяют влияние каждого алгоритма на итоговый результат. Сумма всех весов должна быть равна 1.0.

| Параметр | Тип | Значение по умолчанию | Описание |
|----------|-----|------------------------|----------|
| `LevenshteinWeight` | float64 | 0.2 | Вес алгоритма Левенштейна. Увеличьте для большей чувствительности к опечаткам и заменам символов. |
| `JaroWinklerWeight` | float64 | 0.3 | Вес алгоритма Джаро-Винклера. Увеличьте для лучшей обработки имён с общим префиксом. |
| `PhoneticWeight` | float64 | 0.3 | Вес фонетических алгоритмов (Soundex). Увеличьте для лучшей обработки фонетических вариаций. |
| `DoubleMetaphoneWeight` | float64 | 0.2 | Вес алгоритма Double Metaphone. Увеличьте для лучшей обработки многоязычных имён. |
| `CosineWeight` | float64 | 0.0 | Вес косинусного сходства. Установите значение > 0 для включения этого алгоритма. |
| `AdditionalAttrsWeight` | float64 | 0.0 | Вес дополнительных атрибутов. Установите значение > 0 при использовании дополнительных атрибутов. |

#### Пороговые значения

Определяют границы для классификации результатов сравнения.

| Параметр | Тип | Значение по умолчанию | Описание |
|----------|-----|------------------------|----------|
| `JaroWinklerThreshold` | float64 | 0.85 | Пороговое значение для алгоритма Джаро-Винклера. Значения ниже порога считаются несовпадением. |
| `LevenshteinPrefixScale` | float64 | 0.1 | Коэффициент масштабирования для префикса в алгоритме Левенштейна. |
| `ExactMatchThreshold` | int | 100 | Пороговое значение для классификации точного совпадения (exact_match). |
| `MatchThreshold` | int | 90 | Пороговое значение для классификации совпадения (match). |
| `PossibleMatchThreshold` | int | 70 | Пороговое значение для классификации возможного совпадения (possible_match). |

#### Параметры транслитерации

Настройки для обработки имён в разных алфавитах.

| Параметр | Тип | Значение по умолчанию | Описание |
|----------|-----|------------------------|----------|
| `EnableTransliteration` | bool | true | Включает/отключает транслитерацию. Отключите, если сравниваете имена только в одном алфавите. |
| `TransliterationStandards` | []string | ["gost", "iso9", "bgnpcgn", "ungegn"] | Список используемых стандартов транслитерации. Чем больше стандартов, тем более гибкое, но медленное сравнение. |

#### Параметры перестановки

Настройки для обработки перестановок частей имени.

| Параметр | Тип | Значение по умолчанию | Описание |
|----------|-----|------------------------|----------|
| `EnableNamePartPermutation` | bool | true | Включает/отключает учёт перестановок частей имени. Отключите для ускорения, если порядок частей имени фиксирован. |

#### Другие параметры

| Параметр | Тип | Значение по умолчанию | Описание |
|----------|-----|------------------------|----------|
| `NGramSize` | int | 3 | Размер n-грамм для косинусного сходства. Меньшие значения увеличивают чувствительность к небольшим изменениям. |
| `EnableCaching` | bool | true | Включает/отключает кэширование результатов. Отключите при ограниченной памяти или для экономии ресурсов. |
| `MaxCacheSize` | int | 1000 | Максимальный размер кэша. Увеличьте для больших наборов данных, уменьшите при ограниченной памяти. |
| `EnableLogging` | bool | true | Включает/отключает логирование. Отключите в продакшене для повышения производительности. |

### Рекомендации по настройке

- **Для самой высокой точности**: Увеличьте `LevenshteinWeight` и `JaroWinklerWeight`, установите более высокие пороговые значения.
- **Для большей толерантности к опечаткам**: Увеличьте `LevenshteinWeight`, уменьшите `JaroWinklerThreshold`.
- **Для сравнения имён на разных языках**: Увеличьте `DoubleMetaphoneWeight` и `PhoneticWeight`.
- **Для оптимальной производительности**: Используйте только необходимые стандарты транслитерации, отключите логирование и кэширование.
- **Для сравнения с инициалами**: Сохраните высокий вес `JaroWinklerWeight`, включите `EnableNamePartPermutation`.

## 📂 Структура проекта

```
compareNames/
├── matcher/           # Основная библиотека сравнения имён
│   ├── similarity/    # Реализация алгоритмов сравнения строк
│   ├── translit/      # Функции для транслитерации
│   ├── compare/       # Различные стратегии сравнения имён
│   └── utils/         # Вспомогательные функции
├── api/               # REST API для сравнения имён
└── e2e/               # End-to-end тесты
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
    "birth_date": true,
    "country": false
  },
  "config": {
    "levenshtein_weight": 0.3,
    "jaro_winkler_weight": 0.25,
    "phonetic_weight": 0.2,
    "double_metaphone_weight": 0.15,
    "cosine_weight": 0.1,
    "ngram_size": 3,
    "enable_caching": true
  }
}
```

> Примечание: Параметры `attributes` и `config` являются необязательными.

**Примеры ответов**:

1. Точное совпадение:
```json
{
  "exact_match": true,
  "score": 100,
  "match_type": "exact_match",
  "processing_time_ms": 1
}
```

2. Совпадение (например, транслитерация):
```json
{
  "exact_match": false,
  "score": 99,
  "match_type": "match",
  "best_match1": "иванов иван",
  "best_match2": "ivanov ivan",
  "levenshtein_score": 1.0,
  "jaro_winkler_score": 1.0,
  "phonetic_score": 1.0,
  "double_metaphone_score": 1.0,
  "processing_time_ms": 2
}
```

3. Возможное совпадение (требуется проверка):
```json
{
  "exact_match": false,
  "score": 81,
  "match_type": "possible_match",
  "best_match1": "иванов иван петрович",
  "best_match2": "иванов иван иванович",
  "levenshtein_score": 0.89,
  "jaro_winkler_score": 0.91,
  "phonetic_score": 0.67,
  "processing_time_ms": 3,
  "from_cache": false
}
```

4. Несовпадение:
```json
{
  "exact_match": false,
  "score": 29,
  "match_type": "no_match",
  "levenshtein_score": 0.62,
  "jaro_winkler_score": 0.45,
  "phonetic_score": 0.0,
  "processing_time_ms": 2
}
```

## 🧪 Тестирование

### End-to-end тесты

Проект содержит end-to-end (e2e) тесты, которые проверяют правильность работы API и логики сравнения имён в различных сценариях. Тесты используют реальные примеры имён и проверяют соответствие результатов ожидаемым значениям.

Для запуска e2e тестов используйте скрипт:

```bash
./run-e2e.sh
```

Тесты проверяют различные сценарии:
- Точные совпадения
- Транслитерация (кириллица <-> латиница)
- Обработка инициалов
- Разные формы написания имён
- Перестановки частей ФИО
- Имена с дефисами
- Опечатки и неточности

## 🛠️ Установка и сборка

### Установка библиотеки

```bash
go get github.com/x0rium/compareNames
```

### Сборка приложений

```bash
# Сборка API сервера
go build -o compareNames main.go

# Запуск API сервера
./compareNames
```

## 📊 Интерпретация результатов

| Тип совпадения | Оценка | Описание |
|----------------|--------|----------|
| `exact_match`  | = 100  | Имена полностью идентичны (после предобработки) |
| `match`        | > 90   | Имена с высокой вероятностью относятся к одному человеку |
| `possible_match` | 70-90 | Имена могут относиться к одному человеку, требуется дополнительная проверка |
| `no_match`     | < 70   | Имена, вероятно, относятся к разным людям |

## 📄 Лицензия

MIT
