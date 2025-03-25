package matcher

import (
	"math"
	"strings"
)

// min находит минимальное из трех чисел
func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// levenshteinDistance вычисляет расстояние Левенштейна между двумя строками
func levenshteinDistance(s, t string) int {
	// Длины строк
	m := len([]rune(s))
	n := len([]rune(t))

	// Создаем двумерный массив для динамического программирования
	d := make([][]int, m+1)
	for i := range d {
		d[i] = make([]int, n+1)
	}

	// Инициализация первой строки и столбца
	for i := 0; i <= m; i++ {
		d[i][0] = i
	}
	for j := 0; j <= n; j++ {
		d[0][j] = j
	}

	// Перебираем все символы
	sRunes := []rune(s)
	tRunes := []rune(t)

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			cost := 1
			if sRunes[i-1] == tRunes[j-1] {
				cost = 0
			}

			// Выбираем минимальное значение из трех возможных операций
			d[i][j] = min(
				d[i-1][j]+1,      // удаление
				d[i][j-1]+1,      // вставка
				d[i-1][j-1]+cost, // замена
			)
		}
	}

	return d[m][n]
}

// jaroWinklerSimilarity вычисляет сходство Джаро-Винклера между двумя строками
// Значение от 0 (нет сходства) до 1 (идентичные строки)
func jaroWinklerSimilarity(s1, s2 string) float64 {
	// Приводим к нижнему регистру и преобразуем в руны
	s1 = strings.ToLower(s1)
	s2 = strings.ToLower(s2)

	s1Runes := []rune(s1)
	s2Runes := []rune(s2)

	// Если одна из строк пуста, возвращаем 0
	if len(s1Runes) == 0 || len(s2Runes) == 0 {
		return 0.0
	}

	// Если строки идентичны, возвращаем 1
	if s1 == s2 {
		return 1.0
	}

	// Рассчитываем расстояние для поиска совпадений
	matchDistance := int(math.Max(float64(len(s1Runes)), float64(len(s2Runes)))/2.0) - 1
	if matchDistance < 0 {
		matchDistance = 0
	}

	// Инициализируем массивы для отслеживания совпадений
	s1Matches := make([]bool, len(s1Runes))
	s2Matches := make([]bool, len(s2Runes))

	// Считаем количество совпадений
	matchCount := 0
	for i := range s1Runes {
		start := int(math.Max(0, float64(i-matchDistance)))
		end := int(math.Min(float64(len(s2Runes)-1), float64(i+matchDistance)))

		for j := start; j <= end; j++ {
			if !s2Matches[j] && s1Runes[i] == s2Runes[j] {
				s1Matches[i] = true
				s2Matches[j] = true
				matchCount++
				break
			}
		}
	}

	// Если нет совпадений, возвращаем 0
	if matchCount == 0 {
		return 0.0
	}

	// Считаем количество транспозиций
	transpositions := 0
	j := 0
	for i := 0; i < len(s1Runes); i++ {
		if s1Matches[i] {
			for j < len(s2Runes) && !s2Matches[j] {
				j++
			}
			if j < len(s2Runes) && s1Runes[i] != s2Runes[j] {
				transpositions++
			}
			j++
		}
	}

	// Рассчитываем расстояние Джаро
	transpositions = transpositions / 2
	jaroSimilarity := (float64(matchCount)/float64(len(s1Runes)) +
		float64(matchCount)/float64(len(s2Runes)) +
		float64(matchCount-transpositions)/float64(matchCount)) / 3.0

	// Рассчитываем расстояние Джаро-Винклера
	// Находим длину общего префикса (максимум 4)
	prefixLength := 0
	maxPrefixLength := int(math.Min(4, math.Min(float64(len(s1Runes)), float64(len(s2Runes)))))

	for i := 0; i < maxPrefixLength; i++ {
		if s1Runes[i] == s2Runes[i] {
			prefixLength++
		} else {
			break
		}
	}

	// Scaling factor for how much the score is adjusted upwards for having common prefixes
	p := 0.1 // стандартный коэффициент

	// Финальный расчет
	jaroWinklerSimilarity := jaroSimilarity + float64(prefixLength)*p*(1-jaroSimilarity)

	return jaroWinklerSimilarity
}
