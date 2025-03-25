package similarity

import (
	"math"
	"strings"
)

// StringVector представляет вектор n-грамм для строки
type StringVector map[string]int

// LevenshteinDistance вычисляет расстояние Левенштейна между двумя строками
func LevenshteinDistance(s1, s2 string) int {
	// Приведение к нижнему регистру для игнорирования регистра
	s1, s2 = strings.ToLower(s1), strings.ToLower(s2)

	// Если одна из строк пустая, расстояние равно длине другой строки
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	// Создаем матрицу для динамического программирования
	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i // Заполняем первый столбец
	}
	for j := range matrix[0] {
		matrix[0][j] = j // Заполняем первую строку
	}

	// Заполняем матрицу
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}
			matrix[i][j] = min(
				matrix[i-1][j]+1,      // удаление
				matrix[i][j-1]+1,      // вставка
				matrix[i-1][j-1]+cost, // замена
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

// LevenshteinSimilarity вычисляет оценку схожести на основе расстояния Левенштейна
// Возвращает число от 0 до 1, где 1 означает полное совпадение
func LevenshteinSimilarity(s1, s2 string) float64 {
	if len(s1) == 0 && len(s2) == 0 {
		return 1.0
	}
	if len(s1) == 0 || len(s2) == 0 {
		return 0.0
	}

	// Для коротких строк (<=3 символа) - специальная обработка
	if len(s1) <= 3 || len(s2) <= 3 {
		if strings.ToLower(s1) == strings.ToLower(s2) {
			return 1.0
		}
		return 0.0
	}

	// Префиксное совпадение (для бонуса)
	prefixLen := 0
	for i := 0; i < min(len(s1), len(s2)); i++ {
		if strings.ToLower(string(s1[i])) == strings.ToLower(string(s2[i])) {
			prefixLen++
		} else {
			break
		}
	}
	prefixBonus := float64(prefixLen) * 0.1 // 0.1 - это масштабный коэффициент для префикса

	// Основная оценка Левенштейна
	distance := LevenshteinDistance(s1, s2)
	maxLen := max(len(s1), len(s2))
	baseScore := 1.0 - float64(distance)/float64(maxLen)

	// Возвращаем финальную оценку с бонусом за префикс, не превышая 1.0
	return math.Min(1.0, baseScore+prefixBonus)
}

// JaroSimilarity вычисляет схожесть Джаро между двумя строками
func JaroSimilarity(s1, s2 string) float64 {
	// Приведение к нижнему регистру для игнорирования регистра
	s1, s2 = strings.ToLower(s1), strings.ToLower(s2)

	// Проверка на пустые строки
	if len(s1) == 0 && len(s2) == 0 {
		return 1.0
	}
	if len(s1) == 0 || len(s2) == 0 {
		return 0.0
	}

	// Если строки идентичны, возвращаем 1.0
	if s1 == s2 {
		return 1.0
	}

	// Для коротких строк - специальная обработка
	if len(s1) <= 3 || len(s2) <= 3 {
		if strings.ToLower(s1) == strings.ToLower(s2) {
			return 1.0
		}
		return 0.0
	}

	// Вычисляем расстояние для поиска совпадений
	matchDistance := max(len(s1), len(s2))/2 - 1
	if matchDistance < 0 {
		matchDistance = 0
	}

	// Отмечаем символы, которые совпадают
	matches1 := make([]bool, len(s1))
	matches2 := make([]bool, len(s2))
	matchCount := 0

	// Находим совпадающие символы
	for i := 0; i < len(s1); i++ {
		// Вычисляем нижнюю и верхнюю границы для поиска
		start := max(0, i-matchDistance)
		end := min(i+matchDistance+1, len(s2))

		for j := start; j < end; j++ {
			// Если символ уже совпал с другим или не совпадает - пропускаем
			if matches2[j] || s1[i] != s2[j] {
				continue
			}
			// Отмечаем совпадения
			matches1[i] = true
			matches2[j] = true
			matchCount++
			break
		}
	}

	// Если нет совпадений, возвращаем 0
	if matchCount == 0 {
		return 0.0
	}

	// Подсчитываем транспозиции (количество символов, которые не на своих местах)
	transpositions := 0
	j := 0
	for i := 0; i < len(s1); i++ {
		if !matches1[i] {
			continue
		}
		// Находим следующий совпадающий символ в s2
		for !matches2[j] {
			j++
		}
		// Если символы не совпадают, увеличиваем счетчик транспозиций
		if s1[i] != s2[j] {
			transpositions++
		}
		j++
	}

	// Окончательный подсчет транспозиций (деленный на 2, так как алгоритм считает транспозицию дважды)
	transpositions /= 2

	// Вычисляем метрику Джаро
	return (float64(matchCount)/float64(len(s1)) +
		float64(matchCount)/float64(len(s2)) +
		float64(matchCount-transpositions)/float64(matchCount)) / 3.0
}

// JaroWinklerSimilarity вычисляет схожесть Джаро-Винклера между двумя строками
func JaroWinklerSimilarity(s1, s2 string) float64 {
	// Вычисляем базовую оценку Джаро
	jaro := JaroSimilarity(s1, s2)

	// Если схожесть Джаро ниже порога, возвращаем ее без модификации
	if jaro < 0.7 {
		return jaro
	}

	// Приводим к нижнему регистру для корректного сравнения
	s1Lower, s2Lower := strings.ToLower(s1), strings.ToLower(s2)

	// Находим длину общего префикса (максимум 4 символа)
	prefixLen := 0
	for i := 0; i < min(len(s1Lower), len(s2Lower), 4); i++ {
		if s1Lower[i] == s2Lower[i] {
			prefixLen++
		} else {
			break
		}
	}

	// Константа масштабирования для Винклера (обычно 0.1)
	p := 0.1

	// Вычисляем и возвращаем схожесть Джаро-Винклера
	return jaro + float64(prefixLen)*p*(1-jaro)
}

// CosineSimilarity вычисляет косинусное сходство между двумя строками
// используя n-граммы символов (по умолчанию триграммы)
func CosineSimilarity(s1, s2 string, gramSize int) float64 {
	// Проверка на пустые строки
	if len(s1) == 0 && len(s2) == 0 {
		return 1.0
	}
	if len(s1) == 0 || len(s2) == 0 {
		return 0.0
	}

	// Если не указан размер n-граммы, используем по умолчанию триграммы
	if gramSize <= 0 {
		gramSize = 3
	}

	// Приводим к нижнему регистру
	s1Lower := strings.ToLower(s1)
	s2Lower := strings.ToLower(s2)

	// Создаем векторы n-грамм
	vec1 := createNGramVector(s1Lower, gramSize)
	vec2 := createNGramVector(s2Lower, gramSize)

	// Вычисляем косинусное сходство
	return computeCosineSimilarity(vec1, vec2)
}

// Создает вектор n-грамм для строки
func createNGramVector(s string, n int) StringVector {
	vector := make(StringVector)

	// Корректируем размер n-граммы, если строка короче n
	if len(s) < n {
		n = len(s)
		if n == 0 {
			return vector
		}
	}

	// Создаем n-граммы и подсчитываем их частоту
	for i := 0; i <= len(s)-n; i++ {
		ngram := s[i : i+n]
		vector[ngram]++
	}

	return vector
}

// Вычисляет косинусное сходство между двумя векторами n-грамм
func computeCosineSimilarity(vec1, vec2 StringVector) float64 {
	// Вычисляем скалярное произведение
	dotProduct := 0.0
	for ngram, count1 := range vec1 {
		if count2, ok := vec2[ngram]; ok {
			dotProduct += float64(count1 * count2)
		}
	}

	// Вычисляем длину (норму) каждого вектора
	norm1 := 0.0
	for _, count := range vec1 {
		norm1 += float64(count * count)
	}
	norm1 = math.Sqrt(norm1)

	norm2 := 0.0
	for _, count := range vec2 {
		norm2 += float64(count * count)
	}
	norm2 = math.Sqrt(norm2)

	// Проверка на деление на ноль
	if norm1 == 0.0 || norm2 == 0.0 {
		return 0.0
	}

	// Возвращаем косинусное сходство
	return dotProduct / (norm1 * norm2)
}

// PhoneticSimilarity вычисляет фонетическую схожесть двух строк
// с использованием комбинации фонетических алгоритмов
func PhoneticSimilarity(s1, s2 string) float64 {
	if len(s1) == 0 || len(s2) == 0 {
		return 0.0
	}

	// Разбиваем строки на слова
	words1 := strings.Fields(s1)
	words2 := strings.Fields(s2)

	// Если нет слов - возвращаем 0
	if len(words1) == 0 || len(words2) == 0 {
		return 0.0
	}

	// Получаем базовую оценку на основе фонетического сходства
	baseScore := 0.0

	// Если строки находятся в одном алфавите, используем прямое сравнение
	// иначе выполняем транслитерацию и сравнение
	if isSameScript(s1, s2) {
		// Прямое фонетическое сравнение
		baseScore = directPhoneticCompare(words1, words2)
	} else {
		// Транслитерация и сравнение
		baseScore = translitPhoneticCompare(words1, words2)
	}

	// Добавляем бонус за совпадение первых букв
	if len(words1) > 0 && len(words2) > 0 &&
		strings.ToLower(string(words1[0][0])) == strings.ToLower(string(words2[0][0])) {
		baseScore += 0.1
	}

	// Ограничиваем результат диапазоном [0, 1]
	if baseScore > 1.0 {
		baseScore = 1.0
	}

	return baseScore
}

// Проверяет, используют ли строки один и тот же алфавит
func isSameScript(s1, s2 string) bool {
	// Упрощенная проверка - сравниваем первые символы
	if len(s1) == 0 || len(s2) == 0 {
		return true // Условно считаем, что пустые строки в одном алфавите
	}

	// Проверяем, оба ли символа ASCII или оба не ASCII
	isASCII1 := s1[0] < 128
	isASCII2 := s2[0] < 128

	return isASCII1 == isASCII2
}

// Прямое фонетическое сравнение без транслитерации
func directPhoneticCompare(words1, words2 []string) float64 {
	// Считаем фонетически похожие слова
	matches := 0
	totalPairs := max(len(words1), len(words2))

	for _, word1 := range words1 {
		for _, word2 := range words2 {
			// Сравниваем слова напрямую
			if strings.EqualFold(word1, word2) || arePhoneticallyClose(word1, word2) {
				matches++
				break
			}
		}
	}

	return float64(matches) / float64(totalPairs)
}

// Фонетическое сравнение с транслитерацией
func translitPhoneticCompare(words1, words2 []string) float64 {
	// В этой функции можно использовать транслитерацию
	// Для простоты используем прямое сравнение
	return directPhoneticCompare(words1, words2)
}

// Проверяет, фонетически близки ли слова
func arePhoneticallyClose(word1, word2 string) bool {
	// Слова фонетически близки, если:
	// 1. Первые буквы совпадают
	// 2. Длины отличаются не более чем на 2 символа
	// 3. Расстояние Левенштейна <= 2

	// Проверка на пустые слова
	if len(word1) == 0 || len(word2) == 0 {
		return false
	}

	// Первые буквы должны совпадать
	if strings.ToLower(string(word1[0])) != strings.ToLower(string(word2[0])) {
		return false
	}

	// Длины не должны отличаться более чем на 2 символа
	if abs(len(word1)-len(word2)) > 2 {
		return false
	}

	// Расстояние Левенштейна должно быть <= 2
	if LevenshteinDistance(word1, word2) <= 2 {
		return true
	}

	return false
}

// Модуль числа
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// Заглушки для функций, которые определены в других файлах пакета
// Эти функции нужны только для компиляции
// Реальная реализация находится в соответствующих файлах

// Вспомогательные функции
func min(a, b int, rest ...int) int {
	result := a
	if b < result {
		result = b
	}
	for _, v := range rest {
		if v < result {
			result = v
		}
	}
	return result
}

func max(a, b int, rest ...int) int {
	result := a
	if b > result {
		result = b
	}
	for _, v := range rest {
		if v > result {
			result = v
		}
	}
	return result
}
