package compare

import (
	"time"

	"github.com/x0rium/compareNames/matcher/similarity"
)

// SameAlphabetResult содержит результат сравнения имен на одном алфавите
type SameAlphabetResult struct {
	ExactMatch                bool    // Точное совпадение
	Score                     int     // Оценка совпадения (0-100)
	MatchType                 string  // Тип совпадения ("match", "possible_match", "no_match")
	BestMatch1                string  // Лучшее совпадение для первого имени
	BestMatch2                string  // Лучшее совпадение для второго имени
	LevenshteinScore          float64 // Оценка по алгоритму Левенштейна
	JaroWinklerScore          float64 // Оценка по алгоритму Джаро-Винклера
	PhoneticScore             float64 // Фонетическая оценка
	DoubleMetaphoneScore      float64 // Оценка по алгоритму Double Metaphone
	CosineScore               float64 // Оценка по косинусному сходству
	AdditionalAttributesScore float64 // Оценка по дополнительным атрибутам
	ProcessingTime            int64   // Время обработки в миллисекундах
}

// ConfigProvider интерфейс для получения конфигурационных параметров
type ConfigProvider interface {
	GetLevenshteinWeight() float64
	GetJaroWinklerWeight() float64
	GetPhoneticWeight() float64
	GetDoubleMetaphoneWeight() float64
	GetCosineWeight() float64
	GetAdditionalAttrsWeight() float64
	GetNGramSize() int
}

// DefaultConfigProvider реализация ConfigProvider с значениями по умолчанию
type DefaultConfigProvider struct{}

func (d DefaultConfigProvider) GetLevenshteinWeight() float64     { return 0.30 }
func (d DefaultConfigProvider) GetJaroWinklerWeight() float64     { return 0.20 }
func (d DefaultConfigProvider) GetPhoneticWeight() float64        { return 0.15 }
func (d DefaultConfigProvider) GetDoubleMetaphoneWeight() float64 { return 0.20 }
func (d DefaultConfigProvider) GetCosineWeight() float64          { return 0.05 }
func (d DefaultConfigProvider) GetAdditionalAttrsWeight() float64 { return 0.10 }
func (d DefaultConfigProvider) GetNGramSize() int                 { return 3 }

// CompareSameAlphabet сравнивает имена на одном алфавите
func CompareSameAlphabet(name1, name2 string, attrs MatchAttributes, config interface{}, startTime time.Time) SameAlphabetResult {
	// Получаем конфигурацию
	var cfg ConfigProvider
	if c, ok := config.(ConfigProvider); ok {
		cfg = c
	} else {
		cfg = DefaultConfigProvider{}
	}

	// Получаем вариации имен
	name1Variations := getNameVariations(name1)
	name2Variations := getNameVariations(name2)

	// Ограничиваем количество вариаций для улучшения производительности
	maxVariations := 5
	if len(name1Variations) > maxVariations {
		name1Variations = name1Variations[:maxVariations]
	}
	if len(name2Variations) > maxVariations {
		name2Variations = name2Variations[:maxVariations]
	}

	maxScore := 0.0
	bestMatch1 := ""
	bestMatch2 := ""

	for _, var1 := range name1Variations {
		for _, var2 := range name2Variations {
			parts1 := normalizeNameParts(var1)
			parts2 := normalizeNameParts(var2)

			score := CompareNameParts(parts1, parts2, cfg)

			if score > maxScore {
				maxScore = score
				bestMatch1 = var1
				bestMatch2 = var2
			}
		}
	}

	// Если не нашли хороших совпадений
	if maxScore == 0.0 {
		return SameAlphabetResult{
			ExactMatch:     false,
			Score:          0,
			MatchType:      "no_match",
			ProcessingTime: time.Since(startTime).Milliseconds(),
		}
	}

	// Вычисление отдельных оценок для лучшего совпадения
	levScore := similarity.LevenshteinSimilarity(bestMatch1, bestMatch2)
	jaroScore := similarity.JaroWinklerSimilarity(bestMatch1, bestMatch2)
	phoneticScore := similarity.PhoneticSimilarity(bestMatch1, bestMatch2)
	doubleMetaphoneScore := similarity.DoubleMetaphoneSimilarity(bestMatch1, bestMatch2)
	cosineScore := similarity.CosineSimilarity(bestMatch1, bestMatch2, cfg.GetNGramSize())

	// Учет дополнительных атрибутов
	attrsScore := 0.0
	attrBonus := 0.0

	if attrs != nil {
		for attrName, attrValue := range attrs {
			if attrName == "birth_date" && attrValue.Match {
				attrsScore += 1.0
				attrBonus += 10.0 // Увеличиваем бонус к итоговой оценке
			} else if attrName == "country" && attrValue.Match {
				attrsScore += 0.5
				attrBonus += 5.0 // Увеличиваем бонус к итоговой оценке
			} else if attrName == "citizenship" && attrValue.Match {
				attrsScore += 0.7 // Новый атрибут для гражданства
				attrBonus += 7.0
			}
		}

		// Нормализация оценки дополнительных атрибутов
		if attrsScore > 1.0 {
			attrsScore = 1.0
		}
	}

	// Расчет итогового балла с учетом всех алгоритмов
	finalScore := (levScore*cfg.GetLevenshteinWeight()+
		jaroScore*cfg.GetJaroWinklerWeight()+
		phoneticScore*cfg.GetPhoneticWeight()+
		doubleMetaphoneScore*cfg.GetDoubleMetaphoneWeight()+
		cosineScore*cfg.GetCosineWeight()+
		attrsScore*cfg.GetAdditionalAttrsWeight())*100 + attrBonus

	// Округление до целого числа
	roundedScore := int(finalScore + 0.5)
	if roundedScore > 100 {
		roundedScore = 100
	}

	// Определение типа совпадения
	matchType := "no_match"
	if roundedScore >= 90 { // MinExactMatchScore
		matchType = "match"
	} else if roundedScore >= 70 { // MinPossibleMatchScore
		matchType = "possible_match"
	}

	// Формирование результата
	return SameAlphabetResult{
		ExactMatch:                false,
		Score:                     roundedScore,
		MatchType:                 matchType,
		BestMatch1:                bestMatch1,
		BestMatch2:                bestMatch2,
		LevenshteinScore:          levScore,
		JaroWinklerScore:          jaroScore,
		PhoneticScore:             phoneticScore,
		DoubleMetaphoneScore:      doubleMetaphoneScore,
		CosineScore:               cosineScore,
		AdditionalAttributesScore: attrsScore,
		ProcessingTime:            time.Since(startTime).Milliseconds(),
	}
}

// CompareExactMatch проверяет точное совпадение имен после предобработки
func CompareExactMatch(name1, name2 string) bool {
	// Предобработка имен
	processedName1 := preprocessName(name1)
	processedName2 := preprocessName(name2)

	// Проверка на пустые имена
	if processedName1 == "" || processedName2 == "" {
		return false
	}

	// Проверка на точное совпадение после предобработки
	return processedName1 == processedName2
}

// preprocessName предобработка имени: приведение к нижнему регистру, удаление лишних символов
func preprocessName(name string) string {
	// Приведение к нижнему регистру
	// Удаление лишних символов
	// Нормализация пробелов
	// ...
	// Эта функция должна быть реализована в соответствии с требованиями проекта
	// Здесь представлена упрощенная версия
	return name
}
