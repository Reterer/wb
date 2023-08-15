package main

import (
	"fmt"
	"sort"
	"strings"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func sortString(s string) string {
	r := []rune(s)
	sort.Slice(r, func(i, j int) bool {
		return r[i] < r[j]
	})
	return string(r)
}

func FindAnagramGroups(words []string) map[string][]string {
	// группируем аннаграммы
	g := make(map[string][]string)
	for _, word := range words {
		word = strings.ToLower(word)
		anagram := sortString(word)
		g[anagram] = append(g[anagram], word)
	}

	// формируем результат
	res := make(map[string][]string)
	for _, v := range g {
		if len(v) < 2 {
			continue
		}
		first := v[0]
		sort.Strings(v)
		var vRes []string
		vRes = append(vRes, v[0])
		for i := 1; i < len(v); i++ { // добавляем только уникальные слова
			if v[i-1] != v[i] {
				vRes = append(vRes, v[i])
			}
		}
		if len(vRes) < 2 { // У нас могло быть много уникальных слов
			continue
		}
		res[first] = vRes
	}
	return res
}

func main() {
	words := []string{"тяпка", "слиток", "столик", "пятак", "пятка", "листок"}
	fmt.Println(FindAnagramGroups(words))
}
