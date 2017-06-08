package sqlparser

import (
	"testing"
	"log"
	"strings"
)

func TestScanner_Scan(t *testing.T) {
	scanner := NewScanner("SELECT [ ENAME = ’Mary’ & DNAME = ’Research’ ] ( EMPLOYEE JOIN DEPARTMENT )")
	for tok := scanner.Scan(); tok != nil; {
		log.Println("TokenType:",strings.ToLower(tok.TokenType),"TokenValue:",tok.TokenValue)
		tok = scanner.Scan()
	}
}
