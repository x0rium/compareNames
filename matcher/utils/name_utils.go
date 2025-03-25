package utils

import (
	"strings"
)

// Nicknames карта уменьшительных имен
var Nicknames = map[string][]string{
	// Русские имена
	"александр":  {"саша", "шура", "саня", "алекс"},
	"алексей":    {"леша", "лёша", "алеша", "алёша", "лёха", "леха"},
	"анатолий":   {"толя", "толик"},
	"андрей":     {"андрюша", "дрюня"},
	"антон":      {"антоша", "тоша", "тоха"},
	"артем":      {"тема", "артемка", "тёма"},
	"борис":      {"боря", "борька"},
	"вадим":      {"вадик", "вадя"},
	"валентин":   {"валя", "валик"},
	"валерий":    {"валера", "валерка"},
	"василий":    {"вася", "васька", "васек", "васёк"},
	"виктор":     {"витя", "витька", "витек", "витёк"},
	"виталий":    {"виталик", "виталя"},
	"владимир":   {"вова", "володя", "вовка", "вовочка", "владик"},
	"владислав":  {"влад", "владик", "слава"},
	"вячеслав":   {"слава", "славик"},
	"геннадий":   {"гена", "генка", "геша"},
	"георгий":    {"гоша", "жора", "гера"},
	"григорий":   {"гриша", "гришка", "гриня"},
	"даниил":     {"даня", "данька", "данила"},
	"денис":      {"дениска", "деня"},
	"дмитрий":    {"дима", "димка", "митя"},
	"евгений":    {"женя", "женька", "жека"},
	"егор":       {"егорка", "гоша"},
	"иван":       {"ваня", "ванька", "ванечка"},
	"игорь":      {"игорек", "игорёк", "гарик"},
	"илья":       {"ильюша", "илюша"},
	"кирилл":     {"кирюша", "кир"},
	"константин": {"костя", "костик", "кост"},
	"леонид":     {"лёня", "леня", "лёнчик", "ленчик"},
	"максим":     {"макс", "максик", "максимка"},
	"михаил":     {"миша", "мишка", "миха", "мишаня"},
	"никита":     {"никитка", "ник", "никитос"},
	"николай":    {"коля", "колька", "николка", "ник"},
	"олег":       {"олежка", "олежек", "олежик"},
	"павел":      {"паша", "пашка", "павлик"},
	"петр":       {"петя", "петька", "петруха"},
	"роман":      {"рома", "ромка", "ромчик"},
	"сергей":     {"серега", "серёга", "сережа", "серёжа"},
	"станислав":  {"стас", "славик", "слава"},
	"степан":     {"стёпа", "степа", "стёпка", "степка"},
	"тимофей":    {"тима", "тимоха", "тимоша"},
	"федор":      {"федя", "федька", "федюня"},
	"юрий":       {"юра", "юрка", "юрчик"},
	"ярослав":    {"яра", "ярик", "слава"},

	// Женские имена
	"александра": {"саша", "шура", "саня", "алекс"},
	"алена":      {"аленка", "аленушка", "алёна", "алёнка", "алёнушка"},
	"алина":      {"алинка", "аля"},
	"анастасия":  {"настя", "настенька", "ася"},
	"анна":       {"аня", "анечка", "анька", "анюта"},
	"валентина":  {"валя", "валюша", "тина"},
	"валерия":    {"лера", "лерочка", "валя"},
	"вера":       {"верочка", "верка"},
	"виктория":   {"вика", "викуля", "викуся"},
	"галина":     {"галя", "галочка", "галка"},
	"дарья":      {"даша", "дашенька", "дашка"},
	"евгения":    {"женя", "женечка"},
	"екатерина":  {"катя", "катенька", "катюша", "катерина"},
	"елена":      {"лена", "леночка", "ленка", "еленка"},
	"елизавета":  {"лиза", "лизочка", "лизка", "лизавета"},
	"ирина":      {"ира", "ирочка", "иришка", "иринка"},
	"кристина":   {"кристи", "крис", "кристинка"},
	"лариса":     {"лара", "ларочка", "лариска"},
	"любовь":     {"люба", "любочка", "любаша"},
	"людмила":    {"люда", "людочка", "мила", "люся"},
	"маргарита":  {"рита", "риточка", "маргоша"},
	"марина":     {"мариша", "маришка", "мариночка"},
	"мария":      {"маша", "машенька", "машка", "маня"},
	"надежда":    {"надя", "наденька", "надюша"},
	"наталья":    {"наташа", "наташенька", "наталия", "ната"},
	"нина":       {"ниночка", "нинуля", "нинуша"},
	"оксана":     {"ксюша", "оксаночка", "ксана"},
	"ольга":      {"оля", "оленька", "олечка", "ольчик"},
	"полина":     {"поля", "полинка", "полюшка"},
	"светлана":   {"света", "светочка", "светик", "светланка"},
	"софья":      {"соня", "сонечка", "софа"},
	"татьяна":    {"таня", "танечка", "танюша"},
	"юлия":       {"юля", "юленька", "юлька"},
	"яна":        {"яночка", "янка"},

	// Английские имена
	"alexander":   {"alex", "al", "alec", "sandy", "sasha"},
	"anthony":     {"tony", "ant", "toni"},
	"benjamin":    {"ben", "benji", "benny"},
	"charles":     {"charlie", "chuck", "chaz"},
	"christopher": {"chris", "topher", "kit"},
	"daniel":      {"dan", "danny", "dani"},
	"david":       {"dave", "davey", "davy"},
	"edward":      {"ed", "eddie", "ted", "teddy"},
	"elizabeth":   {"liz", "lizzy", "beth", "betty", "eliza"},
	"james":       {"jim", "jimmy", "jamie"},
	"jennifer":    {"jen", "jenny", "jenn"},
	"john":        {"johnny", "jack", "jock"},
	"joseph":      {"joe", "joey", "jo"},
	"katherine":   {"kate", "katie", "kathy", "kat"},
	"margaret":    {"maggie", "meg", "peggy"},
	"matthew":     {"matt", "matty"},
	"michael":     {"mike", "mikey", "mick"},
	"nicholas":    {"nick", "nicky", "nico"},
	"patrick":     {"pat", "patty", "paddy"},
	"peter":       {"pete", "petey"},
	"richard":     {"rick", "ricky", "dick", "richie"},
	"robert":      {"rob", "robbie", "bob", "bobby"},
	"samuel":      {"sam", "sammy"},
	"steven":      {"steve", "stevie"},
	"thomas":      {"tom", "tommy"},
	"william":     {"will", "bill", "billy", "willy"},
}

// GetNameVariations генерирует различные вариации имени, включая перестановки и уменьшительные формы
func GetNameVariations(name string) []string {
	parts := NormalizeNameParts(name)
	if len(parts) == 0 {
		return []string{}
	}

	variationsMap := make(map[string]bool) // Используем map для удаления дубликатов

	// Добавляем исходное имя
	original := strings.Join(parts, " ")
	variationsMap[original] = true

	// Проверяем на наличие инициалов
	hasInitials := HasInitials(name)

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

	// Добавляем уменьшительные формы имен
	for i, part := range parts {
		// Проверяем, есть ли для этой части уменьшительные формы
		if nicknames, ok := Nicknames[strings.ToLower(part)]; ok {
			for _, nickname := range nicknames {
				// Создаем копию частей
				newParts := make([]string, len(parts))
				copy(newParts, parts)

				// Заменяем часть на уменьшительную форму
				newParts[i] = nickname

				// Добавляем вариацию
				variationsMap[strings.Join(newParts, " ")] = true

				// Если есть хотя бы 2 части, добавляем перестановки с уменьшительной формой
				if len(newParts) >= 2 {
					if len(newParts) == 2 {
						// ИФ или ФИ
						variationsMap[strings.Join([]string{newParts[1], newParts[0]}, " ")] = true
					} else if len(newParts) == 3 {
						// Добавляем основные перестановки
						variationsMap[strings.Join([]string{newParts[1], newParts[0], newParts[2]}, " ")] = true
						variationsMap[strings.Join([]string{newParts[0], newParts[2], newParts[1]}, " ")] = true
					}
				}
			}
		}
	}

	// Преобразуем map в слайс
	variations := make([]string, 0, len(variationsMap))
	for v := range variationsMap {
		variations = append(variations, v)
	}

	return variations
}

// IsNickname проверяет, является ли имя уменьшительной формой другого имени
// Возвращает полное имя и флаг, является ли имя уменьшительной формой
func IsNickname(name string) (string, bool) {
	nameLower := strings.ToLower(name)

	// Проверяем все полные имена
	for fullName, nicknames := range Nicknames {
		for _, nickname := range nicknames {
			if nickname == nameLower {
				return fullName, true
			}
		}
	}

	return "", false
}

// GetFullName возвращает полное имя для уменьшительной формы
// Если имя не является уменьшительной формой, возвращает исходное имя
func GetFullName(name string) string {
	if fullName, isNickname := IsNickname(name); isNickname {
		return fullName
	}
	return name
}

// GetAllNicknames возвращает все уменьшительные формы для имени
func GetAllNicknames(name string) []string {
	nameLower := strings.ToLower(name)

	// Проверяем, является ли имя полным
	if nicknames, ok := Nicknames[nameLower]; ok {
		return nicknames
	}

	// Проверяем, является ли имя уменьшительной формой
	if fullName, isNickname := IsNickname(nameLower); isNickname {
		return Nicknames[fullName]
	}

	return []string{}
}
