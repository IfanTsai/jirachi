package parser

import (
	"fmt"

	"github.com/IfanTsai/jirachi/common"

	"github.com/IfanTsai/jirachi/token"
)

type JNodeType int

const (
	Base JNodeType = iota
	Number
	VarAssign
	VarAccess
	BinOp
	UnaryOp
	IfExpr
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
	return "(" + n.Token.String() + " " + n.Node.String() + ")"
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
	Cases    [][2]JNode
	ElseCase JNode
}

func (n *JIfExprNode) Type() JNodeType {
	return IfExpr
}

func (n *JIfExprNode) String() string {
	return n.Token.String()
}
