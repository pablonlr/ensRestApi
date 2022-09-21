package main

import (
	"io/ioutil"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

//remove accents from the word
func removeAccents(s string) (string, error) {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, err := transform.String(t, s)
	if err != nil {
		return "", err
	}
	return output, nil
}

type SpanishRegister struct {
	WordsInDictionary map[string]bool
}

func (register *SpanishRegister) wordInSpanish(word string) bool {
	if len(word) < 2 {
		return false
	}
	noAccentWord, err := removeAccents(strings.ToLower(word))
	if err != nil {
		noAccentWord = word
	}

	_, ok := register.WordsInDictionary[noAccentWord]
	if ok {
		return true
	}
	if noAccentWord[len(noAccentWord)-1] == 's' {
		noAccentWordNoPlural := noAccentWord[:len(noAccentWord)-1]
		_, okPlural := register.WordsInDictionary[noAccentWordNoPlural]
		if okPlural {
			return true
		}
		if len(noAccentWordNoPlural) > 2 && noAccentWordNoPlural[len(noAccentWordNoPlural)-1] == 'e' {
			noAccentWordNoPluralES := noAccentWordNoPlural[:len(noAccentWordNoPlural)-1]
			_, okPluralES := register.WordsInDictionary[noAccentWordNoPluralES]
			if okPluralES {
				return true
			}
		}
	}
	return false

}

//Load all the spanish words from an external TXT and return a new spanish register
func NewRegisterFromTXTDictionary(txtPath string) (*SpanishRegister, error) {
	resp, err := ioutil.ReadFile(txtPath)
	if err != nil {
		return nil, err
	}
	//get the words splited by empty space
	words := strings.Fields(string(resp))
	wordsReg := make(map[string]bool)
	for _, x := range words {
		noAccentWord, err := removeAccents(x)
		if err != nil {
			return nil, err
		}
		wordsReg[noAccentWord] = true
	}
	return &SpanishRegister{
		WordsInDictionary: wordsReg,
	}, nil
}

//wordsInSpanishFilter recieve an array of words and return a subarray of the spanish words registered in the RAE dictionary
func (register *SpanishRegister) wordsInSpanishFilter(words []string) (spanishWords []string) {
	for _, x := range words {
		if register.wordInSpanish(x) {
			spanishWords = append(spanishWords, x)
		}
	}
	return
}
