package main

import (
	"fmt"
	"log"
	"strings"
)

var sqlKey []string = []string{"SELECT", "PROJECTION", "JOIN"}
var tableList [][]string = [][]string{
	{"ESSN", "ADDRESS", "SALARY", "SUPERSSN", "ENAME", "DNO"}, //EMPLOYEE
	{"DNO", "DNAME", "DNEMBER", "MGRSSN", "MGRSTARDATE"},      //DEPARTMENT
	{"PNAME", "PNO", "PLOCATION", "DNO"},                      //PROJECT
	{"HOURS", "ESSN", "PNO"},                                  // WORKSON
}

type node struct {
	Statment  string
	LnextNode *node
	RnextNode *node
	Content   string
}

func NewTree() *node {
	return &node{}
}

func main() {

}

func printOriginTree(sql string) *node {
	tree := NewTree()
	var res = tree
	log.Println("初始执行树")
	sqlWords := strings.Split(sql, " ")
	for i, v := range sqlWords {
		if v == "SELECT" || v == "PROJECTION" {
			tree.Statment = v
			var s string
			j := i + 1
			if sqlWords[j] == "[" {
				j++
				for j < len(sqlWords) && sqlWords[j] != "]" {
					s += sqlWords[j]
					j++
				}
				log.Println(v, s)
				i = j
			}
		} else if v == "[" || v == "]" || v == ")" {
			continue
		} else if v == "(" {
			tree.LnextNode = NewTree()
			tree = tree.LnextNode
		} else if v == "JOIN" {
			x := fmt.Sprintf("%s\t%s", sqlWords[i-1], sqlWords[i+1])
			log.Println(v, x)
			tree.Statment = v
			tree.Content = nil
			tree.LnextNode = NewTree()
			tree.LnextNode.Content = sqlWords[i-1]
			tree.RnextNode = NewTree()
			tree.RnextNode.Content = sqlWords[i-1]
			i++
		} else {
			if sqlWords[i+1] == "JOIN" {
				continue
			}
			tree.Content = tree.Content + v
		}
	}
	return res
}
