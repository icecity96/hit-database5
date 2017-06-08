package sqlparser

import (
	"testing"
	"log"
)

func TestParser(t *testing.T) {
	Parser("SELECT [ ENAME = ’Mary’ & DNAME = ’Research’ ] ( EMPLOYEE JOIN DEPARTMENT )")
	log.Println(LogicPlan(ASTtree.Query1Node.SelectOrProjection).String())
}

func TestParser2(t *testing.T) {
	Parser("PROJECTION [ ENAME ] ( SELECT [ ESALARY < 3000 ] ( EMPLOYEE JOIN SELECT [ PNO = ’P1’ ] ( WORKS_ON JOIN PROJECT ) ) )")
	log.Println(LogicPlan(ASTtree.Query1Node.SelectOrProjection).String())
}

func TestParser3(t *testing.T) {
	Parser("PROJECTION [ BDATE ] ( SELECT [ ENAME = ’John’ & DNAME = ’Research’ ] ( EMPLOYEE JOIN DEPARTMENT ) )")
	log.Println(LogicPlan(ASTtree.Query1Node.SelectOrProjection).String())
}