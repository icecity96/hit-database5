package sqlparser

import (
	"strings"
	"regexp"
)

type Scanner struct {
	SqlWords []string
	currentIndex int
}

func NewScanner(s string) *Scanner {
	return &Scanner{SqlWords:strings.Split(s," "),currentIndex:0}
}

var keywords = map[string]bool {
	"SELECT":true,
	"JOIN":true,
	"PROJECTION":true,
}

type Token struct {
	TokenType	string
	TokenValue 	string
}

func (s *Scanner) Scan() *Token {
	if s.currentIndex == len(s.SqlWords) {
		return &Token{"EOF","EOF"}
	}
	str := s.SqlWords[s.currentIndex]
	s.currentIndex++
	// if keyword
	if keywords[str] == true {
		return &Token{TokenType:strings.ToLower(str),TokenValue:str}
	} else if ok,_:=regexp.MatchString("[^\\W]+",str); !ok {
		return &Token{str,str}
	} else {
		return &Token{"id",str}
	}
}
