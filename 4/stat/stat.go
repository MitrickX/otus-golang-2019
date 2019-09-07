package stat

import (
	"sort"
	"strings"
)

// is symbol a words separator
func isWordsSeparator(symbol rune) bool {
	return symbol == ' ' ||
		symbol == ',' ||
		symbol == '!' ||
		symbol == '?' ||
		symbol == '-' ||
		symbol == ';' ||
		symbol == '.' ||
		symbol == '\n' ||
		symbol == '\t'
}

func getTop(stat map[string]int, n int) []string {

	// total count of items
	count := len(stat)

	// boundary case
	if count == 0 {
		return nil
	}

	// prevent go out of bounds
	if n > count {
		n = count
	}

	// inner struct to represent word + count pair, so we can sort it
	type item struct {
		word  string
		count int
	}

	// list of items (word+count)
	var items []item

	// collect items from stat
	for word, count := range stat {
		items = append(items, item{word, count})
	}

	// sort by count with breaking tie by comparing strings
	sort.Slice(items, func(i, j int) bool {
		if items[i].count > items[j].count {
			return true
		} else if items[i].count < items[j].count {
			return false
		} else {
			return strings.Compare(items[i].word, items[j].word) == -1
		}
	})

	// here our result
	var words = make([]string, n)
	for i := 0; i < n; i++ {
		words[i] = items[i].word
	}

	return words
}

// Get from text top n most encountered words from text
func Top(text string, n int) []string {

	// boundary case
	if len(text) == 0 {
		return nil
	}

	// keep counf of encountered of words
	var stat = make(map[string]int)

	// current word - slice of runes
	var word []rune

	// input stream of runes with extra termination symbol, so no need dublicate loop body one more time after loop
	stream := []rune(text)
	stream = append(stream, '\n')

	for _, symbol := range stream {
		if !isWordsSeparator(symbol) {
			word = append(word, symbol)
		} else {
			// ignore preposition, aricles and etc
			if len(word) > 2 {
				stat[strings.ToLower(string(word))]++
			}
			word = word[0:0]
		}
	}

	// get top words from stat map
	return getTop(stat, n)
}

// just a sugar
func Top10(text string) []string {
	return Top(text, 10)
}
