package Tokenizer

import (
	"strings"
	"unicode"

	"github.com/ibcode15/BetterFs/internal/stemmer"
)

type _Tokenizer struct {
	Runes []rune
	index int
	size  int
}

func (tokenizer _Tokenizer) Init(str string) _Tokenizer {
	tokenizer.Runes = []rune(str)
	for i, r := range tokenizer.Runes {
		tokenizer.Runes[i] = unicode.ToLower(r)
	}
	tokenizer.index = 0
	tokenizer.size = len(tokenizer.Runes)
	return tokenizer
}

func trimLeftSpace2(s string) string {
	return strings.TrimLeftFunc(s, unicode.IsSpace)
}

func trimLeftSpace(r []rune) []rune {
	index := 0
	size := len(r)
	for unicode.IsSpace(r[index]) {
		index += 1
		if index >= size {
			break
		}
	}

	return r[index:]
}

type Token struct {
	HasToken bool
	value    string
}

func (T Token) noToken() Token {
	T.HasToken = false
	return T
}

func (T Token) foundToken(token []rune) Token {
	T.HasToken = true
	T.value = string(token)
	return T

}

func (T Token) foundWordToken(token []rune) Token {
	T.HasToken = true
	T.value = stemmer.StemString(string(token))
	return T

}

type StringIter = func(func(string) bool)

func (tokenizer *_Tokenizer) next_token() Token {
	if tokenizer.index >= tokenizer.size {
		return Token{}.noToken()
	}

	strippedString := trimLeftSpace(tokenizer.Runes[tokenizer.index:])

	tokenizer.index = tokenizer.size - len(strippedString)

	if tokenizer.index >= tokenizer.size {
		return Token{}.noToken()
	}

	start := tokenizer.index
	current_char := tokenizer.Runes[tokenizer.index]

	if unicode.IsLetter(current_char) {
		for unicode.IsLetter(current_char) || unicode.IsNumber(current_char) {

			tokenizer.index += 1
			if tokenizer.index >= tokenizer.size {
				break
			}

			current_char = tokenizer.Runes[tokenizer.index]
		}

		return Token{}.foundWordToken(tokenizer.Runes[start:tokenizer.index])

	} else if unicode.IsNumber(current_char) {
		for unicode.IsNumber(current_char) {
			tokenizer.index += 1
			if tokenizer.index >= tokenizer.size {
				break
			}
			current_char = tokenizer.Runes[tokenizer.index]
		}
		return Token{}.foundToken(tokenizer.Runes[start:tokenizer.index])
	} else {
		tokenizer.index += 1
		return Token{}.foundToken([]rune{tokenizer.Runes[tokenizer.index-1]})
	}

}

func (tokenizer *_Tokenizer) IterateTokens() StringIter {
	return func(yield func(string) bool) {
		next_token := tokenizer.next_token()
		for next_token.HasToken {
			if !yield(next_token.value) {
				return
			}
			next_token = tokenizer.next_token()
		}
	}

}

func (tokenizer *_Tokenizer) ToArray() []string {
	output := []string{}
	for i := range tokenizer.IterateTokens() {
		output = append(output, i)
	}
	return output
}

func CreateTokenizer(str string) _Tokenizer {
	return _Tokenizer{}.Init(str)
}
