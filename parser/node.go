package parser

import (
	"fmt"
	"strings"

	"github.com/IfanTsai/jirachi/common"

	"github.com/IfanTsai/jirachi/token"
)

type JNodeType int

const (
	Base JNodeType = iota
	Number
	String
	List
	Map
	VarAssign
	VarIndexAssign
	VarAccess
	BinOp
	UnaryOp
	IfExpr
	ForExpr
	WhileExpr
	FuncDefExpr
	CallExpr
	IndexExpr
	ReturnExpr
	ContinueExpr
	BreakExpr
)

// JNode is general node interface of AST
type JNode interface {
	fmt.Stringer
	Type() JNodeType
	GetToken() *token.JToken
	GetStartPos() *common.JPosition
	GetEndPos() *common.JPosition
}

// JBaseNode is general node structure of AST
type JBaseNode struct {
	Token    *token.JToken
	StartPos *common.JPosition
	EndPos   *common.JPosition
}

func (n *JBaseNode) Type() JNodeType {
	return Base
}

func (n *JBaseNode) String() string {
	return n.Token.String()
}

func (n *JBaseNode) GetToken() *token.JToken {
	return n.Token
}

func (n *JBaseNode) GetStartPos() *common.JPosition {
	return n.StartPos
}

func (n *JBaseNode) GetEndPos() *common.JPosition {
	return n.EndPos
}

// JNumberNode is number node structure of AST
type JNumberNode struct {
	*JBaseNode
}

func (n *JNumberNode) Type() JNodeType {
	return Number
}

func (n *JNumberNode) String() string {
	return n.Token.String()
}

// JStringNode is string node structure of AST
type JStringNode struct {
	*JBaseNode
}

func (s *JStringNode) Type() JNodeType {
	return String
}

// JListNode is list node structure of AST
type JListNode struct {
	*JBaseNode
	ElementNodes      []JNode
	IsBlockStatements bool
}

func (l *JListNode) Type() JNodeType {
	return List
}

func (l *JListNode) String() string {
	strBuilder := strings.Builder{}
	strBuilder.WriteByte('[')
	for index, element := range l.ElementNodes {
		if index != 0 {
			strBuilder.WriteString(", ")
		}

		strBuilder.WriteString(element.String())
	}
	strBuilder.WriteByte(']')

	return strBuilder.String()
}

// JMapNode is map node structure of AST
type JMapNode struct {
	*JBaseNode
	ElementMap map[JNode]JNode
}

func (m *JMapNode) Type() JNodeType {
	return Map
}

func (m *JMapNode) String() string {
	strBuilder := strings.Builder{}
	strBuilder.WriteByte('{')
	firstKey := true
	for key, value := range m.ElementMap {
		if !firstKey {
			strBuilder.WriteString(", ")
		}

		firstKey = false
		strBuilder.WriteString(key.String())
		strBuilder.WriteString(": ")
		strBuilder.WriteString(value.String())
	}
	strBuilder.WriteByte('}')

	return strBuilder.String()
}

// JIndexExprNode is index expression node structure of AST
type JIndexExprNode struct {
	*JBaseNode
	IndexNode JNode
	IndexExpr JNode
}

func (i *JIndexExprNode) Type() JNodeType {
	return IndexExpr
}

func (i *JIndexExprNode) String() string {
	return i.IndexNode.String() + "[" + i.IndexExpr.String() + "]"
}

// JBinOpNode is binary operation node structure of AST
type JBinOpNode struct {
	*JBaseNode
	LeftNode  JNode
	RightNode JNode
}

func (n *JBinOpNode) Type() JNodeType {
	return BinOp
}

func (n *JBinOpNode) String() string {
	return "(" + n.LeftNode.String() + " " + n.Token.String() + " " + n.RightNode.String() + ")"
}

// JUnaryOpNode is unary operation node structure of AST
type JUnaryOpNode struct {
	*JBaseNode
	Node JNode
}

func (n *JUnaryOpNode) Type() JNodeType {
	return UnaryOp
}

func (n *JUnaryOpNode) String() string {
	return "(" + n.Token.String() + " " + n.Node.String() + ")"
}

// JVarAssignNode is variable assign node structure of AST
type JVarAssignNode struct {
	*JBaseNode
	Node JNode
}

func (n *JVarAssignNode) Type() JNodeType {
	return VarAssign
}

func (n *JVarAssignNode) String() string {
	return "(" + n.Token.String() + " = " + n.Node.String() + ")"
}

// JVarIndexAssignNode is variable index assign node structure of AST
type JVarIndexAssignNode struct {
	*JVarAssignNode
	IndexExprNode *JIndexExprNode
}

func (n *JVarIndexAssignNode) Type() JNodeType {
	return VarIndexAssign
}

func (n *JVarIndexAssignNode) String() string {
	return "(" + n.IndexExprNode.String() + " = " + n.Node.String() + ")"
}

// JVarAccessNode is variable access node structure of AST
type JVarAccessNode struct {
	*JBaseNode
}

func (n *JVarAccessNode) Type() JNodeType {
	return VarAccess
}

func (n *JVarAccessNode) String() string {
	return n.Token.String()
}

// JIfExprNode is if expression node structure of AST
type JIfExprNode struct {
	*JBaseNode
	CaseNodes    [][2]JNode
	ElseCaseNode JNode
}

func (n *JIfExprNode) Type() JNodeType {
	return IfExpr
}

func (n *JIfExprNode) String() string {
	strBuilder := strings.Builder{}
	strBuilder.WriteString("if ")
	for index, caseNode := range n.CaseNodes {
		if index != 0 {
			strBuilder.WriteString(" else if ")
		}

		strBuilder.WriteString("(")
		strBuilder.WriteString(caseNode[0].String())
		strBuilder.WriteString(") {")
		strBuilder.WriteString(caseNode[1].String())
		strBuilder.WriteString("}")
	}

	if n.ElseCaseNode != nil {
		strBuilder.WriteString(" else ")
		strBuilder.WriteString(n.ElseCaseNode.String())
	}

	return strBuilder.String()
}

// JForExprNode is for expression node structure of AST
type JForExprNode struct {
	*JBaseNode        // JBaseNode.Token is variable name token
	StartValueNode    JNode
	EndValueNode      JNode
	StepValueNode     JNode
	BodyNode          JNode
	IsBlockStatements bool
}

func (n *JForExprNode) Type() JNodeType {
	return ForExpr
}

// JWhileExprNode is while expression node structure of AST
type JWhileExprNode struct {
	*JBaseNode
	ConditionNode     JNode
	BodyNode          JNode
	IsBlockStatements bool
}

func (n *JWhileExprNode) Type() JNodeType {
	return WhileExpr
}

func (n *JWhileExprNode) String() string {
	strBuilder := strings.Builder{}
	strBuilder.WriteString("while (")
	strBuilder.WriteString(n.ConditionNode.String())
	strBuilder.WriteString(") {")
	strBuilder.WriteString(n.BodyNode.String())
	strBuilder.WriteString("}")

	return strBuilder.String()
}

// JFuncDefNode is function definition node structure of AST
type JFuncDefNode struct {
	*JBaseNode // JBaseNode.Token is function name token
	ArgTokens  []*token.JToken
	BodyNode   JNode
}

func (n *JFuncDefNode) Type() JNodeType {
	return FuncDefExpr
}

func (n *JFuncDefNode) String() string {
	strBuilder := strings.Builder{}
	strBuilder.WriteString("(<FUNCTION> ")
	if n.Token != nil {
		strBuilder.WriteString(n.Token.String())
	}

	strBuilder.WriteString(" <args>(")
	for index, argToken := range n.ArgTokens {
		if index != 0 {
			strBuilder.WriteByte(' ')
		}
		strBuilder.WriteString(argToken.String())
	}
	strBuilder.WriteString(") ")

	strBuilder.WriteString("<body>" + n.BodyNode.String())

	strBuilder.WriteByte(')')

	return strBuilder.String()
}

// JCallExprNode is call expression node structure of AST
type JCallExprNode struct {
	*JBaseNode // JBaseNode fields are not use
	CallNode   JNode
	ArgNodes   []JNode
}

func (n *JCallExprNode) Type() JNodeType {
	return CallExpr
}

func (n *JCallExprNode) String() string {
	strBuilder := strings.Builder{}

	strBuilder.WriteString("(<FUNCTION> ")
	strBuilder.WriteString(n.CallNode.String())

	strBuilder.WriteString(" <args>(")
	for index, argNode := range n.ArgNodes {
		if index != 0 {
			strBuilder.WriteByte(' ')
		}
		strBuilder.WriteString(argNode.String())
	}
	strBuilder.WriteByte(')')

	strBuilder.WriteByte(')')

	return strBuilder.String()
}

type JReturnNode struct {
	*JBaseNode
	ReturnNode JNode
}

func (n *JReturnNode) Type() JNodeType {
	return ReturnExpr
}

func (n *JReturnNode) String() string {
	return n.Token.String() + " " + n.ReturnNode.String()
}

type JContinueNode struct {
	*JBaseNode
}

func (n *JContinueNode) Type() JNodeType {
	return ContinueExpr
}

func (n *JContinueNode) String() string {
	return n.Token.String()
}

type JBreakNode struct {
	*JBaseNode
}

func (n *JBreakNode) Type() JNodeType {
	return BreakExpr
}

func (n *JBreakNode) String() string {
	return n.Token.String()
}
