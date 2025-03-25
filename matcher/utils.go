package matcher

import (
	"regexp"
	"sort"
	"strings"
)

// PreprocessName предобработка имени: приведение к нижнему регистру, удаление лишних символов
// Экспортированная версия метода для использования в демонстрационных приложениях
func (m *NameMatcher) PreprocessName(name string) string {
	if name == "" {
		return ""
	}

	// Приведение к нижнему регистру
	name = strings.ToLower(name)

	// Преобразование дефисов в пробелы для корректной обработки двойных фамилий
	name = strings.ReplaceAll(name, "-", " ")

	// Обработка апострофов - сохраняем их для имен типа О'Нил
	name = strings.ReplaceAll(name, "'", "")

	// Удаление лишних пробелов
	re := regexp.MustCompile(`\s+`)
	name = re.ReplaceAllString(name, " ")

	return strings.TrimSpace(name)
}

// preprocessName внутренний метод предобработки имени, использующий экспортированную версию
func (m *NameMatcher) preprocessName(name string) string {
	return m.PreprocessName(name)
}

// Проверяет, содержит ли имя инициалы
func (m *NameMatcher) hasInitials(name string) bool {
	// Проверяем наличие точек в имени
	if strings.Contains(name, ".") {
		return true
	}

	// Ищем одиночные буквы (инициалы) среди частей имени
	parts := strings.Fields(name)
	for _, part := range parts {
		if len(part) == 1 {
			return true
		}
	}
	return false
}

// Разбивает имя на части и возвращает их в нормализованном виде
func (m *NameMatcher) normalizeNameParts(name string) []string {
	// Проверяем, содержит ли имя инициалы
	hasInitials := strings.Contains(name, ".")

	// Если имя содержит инициалы, обрабатываем специальным образом
	if hasInitials {
		// Разбиваем имя на части
		parts := strings.Fields(name)
		result := make([]string, 0, len(parts))

		// Обрабатываем каждую часть
		for _, part := range parts {
			if part != "" {
				// Если это инициал с точкой, сохраняем только букву (без точки)
				if strings.HasSuffix(part, ".") && len(part) == 2 {
					// Добавляем только первую букву инициала
					result = append(result, string(part[0]))
				} else {
					// Иначе добавляем часть как обычно
					result = append(result, part)
				}
			}
		}

		// Если в результате только одна часть, это может быть фамилия с инициалами
		// Например, "Петров И. С." -> ["петров", "и", "с"]
		if len(result) == 1 && len(parts) > 1 {
			// Разбиваем имя на части, заменяя точки на пробелы
			name = strings.ReplaceAll(name, ".", " ")
			parts = strings.Fields(name)
			result = make([]string, 0, len(parts))

			// Обрабатываем каждую часть
			for _, part := range parts {
				if part != "" {
					result = append(result, part)
				}
			}
		}

		return result
	}

	// Если имя не содержит инициалов, обрабатываем как обычно
	parts := strings.Fields(name)
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		if part != "" {
			result = append(result, part)
		}
	}

	return result
}

// Генерирует различные вариации имени, включая перестановки
func (m *NameMatcher) getNameVariations(name string) []string {
	// Проверяем кэш
	if variations, ok := m.nameVariantions[name]; ok {
		return variations
	}

	parts := m.normalizeNameParts(name)
	if len(parts) == 0 {
		return []string{}
	}

	variationsMap := make(map[string]bool) // Используем map для удаления дубликатов

	// Добавляем исходное имя
	original := strings.Join(parts, " ")
	variationsMap[original] = true

	// Проверяем на наличие инициалов
	hasInitials := m.hasInitials(name)

	// Если есть хотя бы 2 части, добавляем перестановки
	if len(parts) >= 2 {
		// Добавляем все возможные перестановки для 2 и 3 частей
		if len(parts) == 2 {
			// ИФ или ФИ
			variationsMap[strings.Join([]string{parts[1], parts[0]}, " ")] = true
		} else if len(parts) == 3 {
			// ФИО -> ИФО
			variationsMap[strings.Join([]string{parts[1], parts[0], parts[2]}, " ")] = true
			// ФИО -> ФОИ
			variationsMap[strings.Join([]string{parts[0], parts[2], parts[1]}, " ")] = true
			// ФИО -> ОФИ
			variationsMap[strings.Join([]string{parts[2], parts[0], parts[1]}, " ")] = true
			// ФИО -> ОИФ
			variationsMap[strings.Join([]string{parts[2], parts[1], parts[0]}, " ")] = true
			// ФИО -> ИОФ
			variationsMap[strings.Join([]string{parts[1], parts[2], parts[0]}, " ")] = true

			// Также добавляем вариации с пропущенными частями
			// Без отчества
			variationsMap[strings.Join([]string{parts[0], parts[1]}, " ")] = true
			variationsMap[strings.Join([]string{parts[1], parts[0]}, " ")] = true
			// Без имени
			variationsMap[strings.Join([]string{parts[0], parts[2]}, " ")] = true
			variationsMap[strings.Join([]string{parts[2], parts[0]}, " ")] = true
			// Без фамилии
			variationsMap[strings.Join([]string{parts[1], parts[2]}, " ")] = true
			variationsMap[strings.Join([]string{parts[2], parts[1]}, " ")] = true
		}
	}

	// Если есть инициалы, добавляем специальную обработку
	if hasInitials {
		// Определяем, какие части являются инициалами
		var fullParts, initialParts []string
		for _, p := range parts {
			if strings.HasSuffix(p, ".") || len(p) == 1 {
				initialParts = append(initialParts, p)
			} else {
				fullParts = append(fullParts, p)
			}
		}

		// Комбинируем их в разных порядках
		if len(fullParts) > 0 && len(initialParts) > 0 {
			// Сначала полные части, потом инициалы
			variationsMap[strings.Join(append(fullParts, initialParts...), " ")] = true
			// Сначала инициалы, потом полные части
			variationsMap[strings.Join(append(initialParts, fullParts...), " ")] = true

			// Для каждой полной части добавляем вариацию с инициалом
			for _, fullPart := range fullParts {
				for _, initial := range initialParts {
					variationsMap[fullPart+" "+initial] = true
					variationsMap[initial+" "+fullPart] = true
				}
			}
		}
	}

	// Преобразуем map в слайс
	variations := make([]string, 0, len(variationsMap))
	for v := range variationsMap {
		variations = append(variations, v)
	}

	// Сортируем для стабильности результатов
	sort.Strings(variations)

	// Сохраняем в кэш
	m.nameVariantions[name] = variations

	return variations
}

// Функции для работы с кэшем перенесены в cache.go
