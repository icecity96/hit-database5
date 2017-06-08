package sqlparser

import (
	"myGo/parser"
	"log"
	"fmt"
	"hit-database5/sqlparser/utils"
	"os"
)

// 文法
var g = &parser.Grammar{ []*parser.Rule{
	{"Query", []string{"Query1"}},
	{"Query1", []string{"SelectStmt"}},
	{"Query1", []string{"ProjectionStmt"}},
	{"Table",[]string{"SelectStmt"}},
	{"Table",[]string{"ProjectionStmt"}},
	{"Table",[]string{"JoinStmt"}},
	{"Table",[]string{"id"}},
	{"SelectStmt",[]string{"select","[","Conditions","]","(","Table",")"}},
	{"ProjectionStmt",[]string{"projection","[","List","]","(","Table",")"}},
	{"JoinStmt",[]string{"Table","join","Table"}},
	{"Conditions",[]string{"Expr","Op1","Conditions"}},
	{"Conditions",[]string{"Expr"}},
	{"Op1",[]string{"|"}},
	{"Op1",[]string{"&"}},
	{"Expr",[]string{"id","Op2","id"}},
	{"Op2",[]string{">"}},
	{"Op2",[]string{">="}},
	{"Op2",[]string{"<"}},
	{"Op2",[]string{"<="}},
	{"Op2",[]string{"="}},
	{"Op2",[]string{"!="}},
	{"List",[]string{"id",",","List"}},
	{"List",[]string{"id"}},
},nil}

func Parser(query string) {
	g.CollectSymbols()
	ac := parser.ComputeActions(g)
	p := parser.NewParser(ac)
	scanner := NewScanner(query)
	for {
		tok := scanner.Scan()
		ok,_ := parse(p,tok,"Query")
		if ok {
			break
		}
	}
}

var ASTtree QueryNode

func parse(p *parser.Parser,token *Token,start string) (bool,error) {
	for {
		// get action
		action,ok := p.Actions[p.Stack[len(p.Stack)-1]][token.TokenType]

		if !ok {
			log.Printf("Unexpert token:(%v,%v)",token.TokenType,token.TokenValue)
		}

		switch action.(type) {
		// 移入
		case parser.Shift:
			nextState := action.(parser.Shift).State
			p.Data = append(p.Data,token)
			p.Stack = append(p.Stack,nextState)
			return false, nil
		case parser.Reduce:
			rule := action.(parser.Reduce).Rule

			// NOTE here print reduce, only used for debug
			// log.Printf("input %v => reduce %s -> %s\n", token.TokenType, rule.Pattern, rule.Symbol)
			makeNode(p,rule)

			popCount := len(rule.Pattern)
			p.Stack = p.Stack[0:len(p.Stack)-popCount]
			// acc
			if rule.Symbol == start {
				ASTtree = p.Data[0].(QueryNode)
				log.Println(ASTtree.String())
				return true, nil
			}
			state := p.Stack[len(p.Stack)-1]
			action, ok = p.Actions[state][rule.Symbol]
			if _, well := action.(parser.Reduce); !ok || well {
				panic(fmt.Errorf("parse error, bad next State"))
			}

			p.Stack = append(p.Stack, action.(parser.Shift).State)
		default:
			return false, fmt.Errorf("unkonw action!")
		}
	}
}

func makeNode(p *parser.Parser, rule *parser.Rule) Node {
	switch rule.Symbol {
	case "Query":
		res := QueryNode{p.Data[len(p.Data)-1].(Query1Node)}
		p.Data[len(p.Data)-1] = res
		return res
	case "Query1":
		res := Query1Node{p.Data[len(p.Data)-1].(Node)}
		p.Data[len(p.Data)-1] = res
		return res
	case "Table":
		// Table ->  Node
		var res TableNode
		if rule.Pattern[0] == "id" {
			res = TableNode{p.Data[len(p.Data)-1].(*Token).TokenValue,nil}
		} else if rule.Pattern[0] == "JoinStmt"{
			res = TableNode{fmt.Sprintf("j%d",jtemp),p.Data[len(p.Data)-1].(Node)}
			jtemp++
		} else if rule.Pattern[0] == "SelectStmt" {
			res = TableNode{fmt.Sprintf("s%d",stemp),p.Data[len(p.Data)-1].(Node)}
			stemp++
		} else if rule.Pattern[0] == "ProjectionStmt" {
			res = TableNode{fmt.Sprintf("p%d",ptemp),p.Data[len(p.Data)-1].(Node)}
			ptemp++
		}
		p.Data = p.Data[:len(p.Data)-len(rule.Pattern)]
		p.Data = append(p.Data,res)
		return res
	case "SelectStmt":
		res := SelectNode{"select",
						p.Data[len(p.Data)-5].(ConditionNode),
						p.Data[len(p.Data)-2].(TableNode)}
		p.Data = p.Data[:len(p.Data)-len(rule.Pattern)]
		p.Data = append(p.Data,res)
		return res
	case "ProjectionStmt":
		res := ProjectionNode{"projection",
						p.Data[len(p.Data)-5].(ListNode),
						p.Data[len(p.Data)-2].(TableNode)}
		p.Data = p.Data[:len(p.Data)-len(rule.Pattern)]
		p.Data = append(p.Data,res)
		return res
	case "JoinStmt":
		var res JoinNode
		rnode := p.Data[len(p.Data)-1].(TableNode)
		lnode := p.Data[len(p.Data)-3].(TableNode)
		res = JoinNode{&lnode,"join",&rnode}
		p.Data = p.Data[:len(p.Data)-len(rule.Pattern)]
		p.Data = append(p.Data,res)
		return res
	case "Conditions":
		var res ConditionNode
		// Conditions -> Expr
		if len(rule.Pattern) == 1 {
			res = ConditionNode{p.Data[len(p.Data)-1].(ExprsNode),"",nil}
		} else {
			node := p.Data[len(p.Data)-1].(ConditionNode)
			res = ConditionNode{
				p.Data[len(p.Data)-3].(ExprsNode),
				p.Data[len(p.Data)-2].(*Token).TokenValue,
				&node}
		}
		p.Data = p.Data[:len(p.Data)-len(rule.Pattern)]
		p.Data = append(p.Data,res)
		return res
	case "Op1":
		return nil
	case "Expr":
		var res ExprsNode
		res = ExprsNode{p.Data[len(p.Data)-3].(*Token).TokenValue,
						p.Data[len(p.Data)-2].(*Token).TokenValue,
						p.Data[len(p.Data)-1].(*Token).TokenValue}
		p.Data = p.Data[:len(p.Data)-len(rule.Pattern)]
		p.Data = append(p.Data,res)
		return res
	case "Op2":
		return nil
	case "List":
		var res ListNode
		// List -> id
		if len(rule.Pattern) == 1 {
			res = ListNode{[]string{p.Data[len(p.Data)-1].(*Token).TokenValue}}
		} else {
			newAttrubute := p.Data[len(p.Data)-3].(*Token).TokenValue
			res = ListNode{append(p.Data[len(p.Data)-1].(ListNode).Attribute,newAttrubute)}
		}
		p.Data = p.Data[:len(p.Data)-len(rule.Pattern)]
		p.Data = append(p.Data,res)
		return res
	default:
		return nil
	}
}

func LogicPlan(astTree Node) Node {
	realAstTree := astTree
	switch realAstTree.(type) {
	case SelectNode:
		root := realAstTree.(SelectNode)

		// 直接作用于原始表上了，不可再下推
		if root.Table.Table == nil {
			break
		}

		switch root.Table.Table.(type) {
		case ProjectionNode:
			// 直接对下层优化
			LogicPlan(root.Table.Table)
		case SelectNode:
			condition1 := realAstTree.(SelectNode).Conditions
			condition2 := root.Table.Table.(SelectNode).Conditions
			condition := mergeCondition(condition1,condition2)
			selectNode := root.Table.Table.(SelectNode)
			selectNode.Conditions = condition
			astTree = selectNode
			LogicPlan(astTree)
		case JoinNode:
			conditions := root.Conditions.Condition()
			childJoin := root.Table.Table.(JoinNode)
			// 条件分配
			var leftConditions []ExprsNode
			var rightConditions []ExprsNode

			leftRawTables := childJoin.Ltable.RawTables()
			rightRawTables := childJoin.Rtable.RawTables()

			leftConditions = conditionsDistribute(conditions,leftRawTables)
			rightConditions = conditionsDistribute(conditions,rightRawTables)

			//若一边没分配则改边不做处理,否则下推select
			if len(leftConditions) != 0 {
				var tempCondition ConditionNode
				temp := &tempCondition
				for i := 0; i < len(leftConditions); i++ {
					if i == len(leftConditions) - 1 {
						temp.Exprs = leftConditions[i]
						temp.Op1 = ""
						temp.Conditions = nil
						break
					}
					temp.Exprs = leftConditions[i]
					temp.Op1 = "&"
					temp.Conditions = new(ConditionNode)
					temp = temp.Conditions
				}
				var ltable = &TableNode{fmt.Sprintf("s%d",stemp),SelectNode{"select",tempCondition,*childJoin.Ltable}}
				stemp++
				childJoin.Ltable = ltable
			}
			if len(rightConditions) != 0 {
				var tempCondition ConditionNode
				temp := &tempCondition
				for i := 0; i < len(rightConditions); i++ {
					if i == len(rightConditions) - 1 {
						temp.Exprs = rightConditions[i]
						temp.Op1 = ""
						temp.Conditions = nil
						break
					}
					temp.Exprs = rightConditions[i]
					temp.Op1 = "&"
					temp.Conditions = new(ConditionNode)
					temp = temp.Conditions
				}
				var rtable = &TableNode{fmt.Sprintf("s%d",stemp),SelectNode{"select",tempCondition,*childJoin.Rtable}}
				stemp++
				childJoin.Rtable = rtable
			}
			childJoin.Ltable.Table = LogicPlan(childJoin.Ltable.Table)
			childJoin.Rtable.Table = LogicPlan(childJoin.Rtable.Table)
			astTree = childJoin
		}

	case ProjectionNode:
		root := realAstTree.(ProjectionNode)
		// 直接作用于原始表上了，不可再下推
		if root.Table.Table == nil {
			break
		}
		// 作用于中间表，可能可下推
		switch root.Table.Table.(type) {
		case ProjectionNode:
			childProjection := root.Table.Table.(ProjectionNode)
			intersect := utils.SliceIntersect(childProjection.List.Attribute,root.List.Attribute)
			diff := utils.SliceDiff(root.List.Attribute,childProjection.List.Attribute)
			// 如果两次投影的交集为空，或者上层投影中的属性在下层投影中未出现，那么结果必定为空，无需继续优化
			if len(intersect) < 1 || len(diff) >= 1 {
				log.Println("The result is Empty, Dont need to do more work!")
				os.Exit(1)
			} else {
				// 上层投影是下层投影的一个子集,合并两次投影成上层投影,第一层优化
				root = childProjection
				root.List.Attribute = intersect
				astTree = root
				// 递归优化
				astTree = LogicPlan(astTree)
			}
		case SelectNode:
			// 在若下层操作是select，继续下推projections可能出现表达式不等价，为避免这种
			// 可能性.索性遇到select就停止下推,然后尝试优化select
			root.Table.Table = LogicPlan(root.Table.Table)
			astTree = root
		case JoinNode:
			// 对于JoinNode 首先要获取其下所有原始的Table，以便分配projections里面的属性
			attribute := root.List.Attribute
			childJoin := root.Table.Table.(JoinNode)
			// 投影属性分配
			var leftAttribute []string
			var rightAttribute []string

			leftRawTables := childJoin.Ltable.RawTables()
			rightRawTables := childJoin.Rtable.RawTables()

			leftAttribute = attributeDistribute(attribute,leftRawTables)
			rightAttribute = attributeDistribute(attribute,rightRawTables)

			// 若两边分配属性的和小于原有投影属性个数，则必然为空
			if len(utils.SliceUnique(utils.SliceMerge(leftAttribute,rightAttribute))) < len(attribute) {
				log.Println("The result is Empty, Dont need to do more work!")
				os.Exit(1)
			}
			// 若有一边没有分配到属性,则该边为空
			if len(leftAttribute) < 1 {
				astTree = ProjectionNode{"projection",ListNode{rightAttribute},*childJoin.Rtable}
				LogicPlan(astTree)
			} else if len(rightAttribute) < 1 {
				astTree = ProjectionNode{"projection",ListNode{leftAttribute},*childJoin.Ltable}
				LogicPlan(astTree)
			} else {
				// TODO JOIN两边变形 两个由projections构成的table
				var ltable = &TableNode{fmt.Sprintf("p%d",ptemp),ProjectionNode{"projection",ListNode{leftAttribute},*childJoin.Ltable}}
				ptemp++
				var rtable = &TableNode{fmt.Sprintf("p%d",ptemp),ProjectionNode{"projection",ListNode{rightAttribute},*childJoin.Rtable}}
				ptemp++
				childJoin.Rtable = rtable
				childJoin.Ltable = ltable
				childJoin.Ltable.Table = LogicPlan(childJoin.Ltable.Table)
				childJoin.Rtable.Table = LogicPlan(childJoin.Rtable.Table)
				astTree = childJoin
			}
		}

	default:
		log.Println("ERROR,Unknown Type")
	}
	return astTree
}

func mergeCondition(condition1, condition2 ConditionNode) ConditionNode {
	var condition ConditionNode = condition1
	temp := &condition
	for temp.Conditions != nil {
		temp = temp.Conditions
	}
	var temp2 ConditionNode = condition2
	temp.Conditions = &temp2
	return condition
}

func conditionsDistribute(conditions []ExprsNode,table []string) (res []ExprsNode) {
	for _, v1 := range conditions {
		for _,v2 := range table {
			if v1.Lvalue[0] == v2[0] {
				res = append(res,v1)
				break
			}
		}
	}
	return
}

func attributeDistribute(attribute, table []string) (res []string) {
	for _,v1 := range attribute {
		for _,v2 := range table {
			if v1[0] == v2[0] {
				res = append(res,v1)
				break
			}
		}
	}
	return
}