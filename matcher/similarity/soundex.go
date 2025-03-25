package similarity

import (
	"strings"
	"unicode"

	"github.com/x0rium/compareNames/matcher/translit"
)

// Soundex реализует алгоритм Soundex для английских слов
func Soundex(word string) string {
	if len(word) == 0 {
		return "0000"
	}

	// Приводим к верхнему регистру и оставляем только буквы
	word = strings.ToUpper(word)
	letters := []rune{'0', '0', '0', '0'}

	// Сохраняем первую букву
	firstLetter := unicode.ToUpper(rune(word[0]))
	if unicode.IsLetter(firstLetter) {
		letters[0] = firstLetter
	} else {
		return "0000" // Если первый символ не буква, возвращаем "0000"
	}

	// Карта замены букв на цифры
	replacements := map[rune]rune{
		'B': '1', 'F': '1', 'P': '1', 'V': '1',
		'C': '2', 'G': '2', 'J': '2', 'K': '2', 'Q': '2', 'S': '2', 'X': '2', 'Z': '2',
		'D': '3', 'T': '3',
		'L': '4',
		'M': '5', 'N': '5',
		'R': '6',
	}

	// Индекс для заполнения массива letters
	j := 1
	prevCode := '0'

	// Обрабатываем остальные буквы
	for i := 1; i < len(word) && j < 4; i++ {
		c := unicode.ToUpper(rune(word[i]))

		// Пропускаем H и W
		if c == 'H' || c == 'W' {
			continue
		}

		// Получаем код для текущей буквы
		code, ok := replacements[c]
		if !ok {
			continue // Пропускаем символы, которые не в карте
		}

		// Добавляем код, если он отличается от предыдущего
		if code != prevCode && code != '0' {
			letters[j] = code
			j++
		}

		prevCode = code
	}

	return string(letters)
}

// RussianSoundex адаптированный алгоритм Soundex для русских имен
func RussianSoundex(word string) string {
	// Если слово на кириллице, транслитерируем его
	if translit.IsCyrillic(word) {
		word = translit.TranslitGOST(word)
	}

	// Применяем стандартный Soundex к транслитерированному слову
	return Soundex(word)
}

// SoundexSimilarity вычисляет схожесть на основе алгоритма Soundex
func SoundexSimilarity(s1, s2 string) float64 {
	// Проверка на пустые строки
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

	// Проверяем, есть ли кириллические символы
	isCyrillic1 := translit.IsCyrillic(s1)
	isCyrillic2 := translit.IsCyrillic(s2)

	// Выбираем подходящую функцию Soundex
	var soundexFunc func(string) string
	if isCyrillic1 || isCyrillic2 {
		soundexFunc = RussianSoundex
	} else {
		soundexFunc = Soundex
	}

	// Считаем совпадения Soundex кодов
	matchCount := 0
	totalWords := max(len(words1), len(words2))

	// Для каждого слова в первой строке
	for _, word1 := range words1 {
		soundex1 := soundexFunc(word1)

		// Ищем совпадение во второй строке
		for _, word2 := range words2 {
			soundex2 := soundexFunc(word2)
			if soundex1 == soundex2 {
				matchCount++
				break
			}
		}
	}

	// Возвращаем нормализованную оценку
	return float64(matchCount) / float64(totalWords)
}

// RussianSoundexSimilarity вычисляет схожесть на основе адаптированного алгоритма Soundex для русских имен
func RussianSoundexSimilarity(s1, s2 string) float64 {
	// Проверка на пустые строки
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

	// Считаем совпадения Soundex кодов
	matchCount := 0
	totalWords := max(len(words1), len(words2))

	// Для каждого слова в первой строке
	for _, word1 := range words1 {
		soundex1 := RussianSoundex(word1)

		// Ищем совпадение во второй строке
		for _, word2 := range words2 {
			soundex2 := RussianSoundex(word2)
			if soundex1 == soundex2 {
				matchCount++
				break
			}
		}
	}

	// Возвращаем нормализованную оценку
	return float64(matchCount) / float64(totalWords)
}

// SoundexMatch проверяет, совпадают ли Soundex коды двух слов
func SoundexMatch(word1, word2 string) bool {
	return Soundex(word1) == Soundex(word2)
}

// RussianSoundexMatch проверяет, совпадают ли адаптированные Soundex коды двух слов
func RussianSoundexMatch(word1, word2 string) bool {
	return RussianSoundex(word1) == RussianSoundex(word2)
}
