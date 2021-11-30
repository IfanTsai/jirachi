package parser

import (
	"github.com/IfanTsai/jirachi/common"
	"github.com/IfanTsai/jirachi/token"
)

type JNodeType int

const (
	Number JNodeType = iota
	BinOp
	UnaryOp
)

// JNode is general node structure of AST
type JNode struct {
	Type      JNodeType
	Token     *token.JToken
	LeftNode  *JNode // for BinOp
	RightNode *JNode // for BinOp
	Node      *JNode // for UnaryOp
	StartPos  *common.JPosition
	EndPos    *common.JPosition
}

func (n *JNode) String() string {
	switch n.Type {
	case BinOp:
		return "(" + n.LeftNode.String() + " " + n.Token.String() + " " + n.RightNode.String() + ")"
	case UnaryOp:
		return "(" + n.Token.String() + " " + n.Node.String() + ")"
	default:
		return n.Token.String()
	}
}
