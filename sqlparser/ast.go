package sqlparser

import (
	"strings"
	"fmt"
	"hit-database5/sqlparser/utils"
)

type Node interface {
	String() string
}

type QueryNode struct {
	Query1Node Query1Node
}

func (node QueryNode) String() string {
	return node.Query1Node.String()
}

type Query1Node struct {
	SelectOrProjection Node
}

func (node Query1Node) String() string {
	return node.SelectOrProjection.String()
}

type SelectNode struct {
	KeyWord string
	Conditions ConditionNode
	Table TableNode
}

func (node SelectNode) String() string {
	res := "Node Type: SELECT {\n"
	res += "	CONDITIONS:" + node.Conditions.String()
	res += "\n	Table:" + node.Table.String()
	res += "\n}"
	return res
}

type ListNode struct {
	Attribute []string
}

func (node ListNode) String() string {
	return strings.Join(node.Attribute,",")
}

type ConditionNode struct {
	Exprs ExprsNode
	Op1   string	// | &
	Conditions *ConditionNode
}

func (node ConditionNode) Condition() []ExprsNode {
	var res []ExprsNode
	if node.Exprs.Lvalue == "" {
		return nil
	}
	res = append(res,node.Exprs)
	if node.Op1 != "" {
		res = append(res,node.Conditions.Condition()...)
	}
	return res
}

func (node ConditionNode) String() string {
	res := fmt.Sprintf("%s",node.Exprs.String())
	if node.Op1 == "" {
		return res
	} else {
		res += " " + node.Op1 + " " + node.Conditions.String()
	}
	return res
}

type ExprsNode struct {
	Lvalue string
	Op2	string // = | > | < | >= | <= | !=
	Rvalue string
}

func (node ExprsNode) String() string {
	return fmt.Sprintf("%s %s %s",node.Lvalue,node.Op2,node.Rvalue)
}

type TableNode struct {
	TableName string
	Table Node
}

var ptemp int
var stemp int
var jtemp int

func (node TableNode) RawTables() []string {
	var res []string
	if  node.Table == nil {
		res = append(res,node.TableName)
	} else {
		switch node.Table.(type) {
		case SelectNode:
			res = append(res,node.Table.(SelectNode).Table.RawTables()...)
		case ProjectionNode:
			res = append(res,node.Table.(ProjectionNode).Table.RawTables()...)
		case JoinNode:
			res = append(res,node.Table.(JoinNode).Rtable.RawTables()...)
			res = append(res,node.Table.(JoinNode).Ltable.RawTables()...)
		default:
			break
		}
	}
	return utils.SliceUnique(res)
}

func (node TableNode) Tables() []string {
	return []string{node.TableName}
}

func (node TableNode) String() string {
	res := fmt.Sprintf("Node Type: Table %s {\n",node.TableName)
	if node.Table != nil {
		res += node.Table.String()
	}
	res += "\n}"
	return res
}

type JoinNode struct {
	Ltable *TableNode
	KeyWord string
	Rtable *TableNode
}

func (node JoinNode) Tables() []string {
	var res []string
	res = append(res,node.Ltable.TableName)
	if node.KeyWord != "" {
		res = append(res,node.Rtable.Tables()...)
	}
	return res
}

func (node JoinNode) String() string {
	res := "Node Type: Join {\n"
	res += "Tables: " + strings.Join(node.Tables(),",") + "\n"
	// ltable explain
	switch node.Ltable.Table.(type) {
	case SelectNode:
		res += fmt.Sprintf("Table %s:\n",node.Ltable.TableName)
		res += node.Ltable.Table.(SelectNode).String() + "\n"
		res +="}"
	case ProjectionNode:
		res += fmt.Sprintf("Table %s {:\n",node.Ltable.TableName)
		res += node.Ltable.Table.(ProjectionNode).String() + "\n"
		res += "}"
	case JoinNode:
		res += fmt.Sprintf("Table %s:\n",node.Ltable.TableName)
		res += node.Ltable.Table.(JoinNode).String() + "\n"
		res += "}"
	default:
		break
	}
	switch node.Rtable.Table.(type) {
	case SelectNode:
		res += fmt.Sprintf("Table %s:\n",node.Rtable.TableName)
		res += node.Rtable.Table.(SelectNode).String() + "\n"
		res +="}"
	case ProjectionNode:
		res += fmt.Sprintf("Table %s {:\n",node.Rtable.TableName)
		res += node.Rtable.Table.(ProjectionNode).String() + "\n"
		res += "}"
	case JoinNode:
		res += fmt.Sprintf("Table %s:\n",node.Rtable.TableName)
		res += node.Rtable.Table.(JoinNode).String() + "\n"
		res += "}"
	default:
		break
	}
	res += "\n}"
	return res
}

type ProjectionNode struct {
	KeyWord string
	List ListNode
	Table TableNode
}

func (node ProjectionNode) String() string {
	res := "Node Type: Projection {\n"
	res += "	Attribute: " + node.List.String() + "\n"
	res += "	Table: " + node.Table.String()
	res += "\n}"
	return res
}
