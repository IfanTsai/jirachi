package lexer

import (
	mapset "github.com/deckarep/golang-set"
)

type JNodeType int

const (
	NUMBER JNodeType = iota
	OP
)

type JNode struct {
	Type      JNodeType
	Token     *JToken
	LeftNode  *JNode
	RightNode *JNode
}

func (n *JNode) String() string {
	if n.Type == OP {
		return "(" + n.LeftNode.String() + " " + n.Token.String() + " " + n.RightNode.String() + ")"
	}

	return n.Token.String()
}

type JParser struct {
	TokenIndex   int
	Tokens       []*JToken
	CurrentToken *JToken
}

func NewJParser(tokens []*JToken, tokenIndex int) *JParser {
	return &JParser{
		Tokens:     tokens,
		TokenIndex: tokenIndex,
	}
}

func (p *JParser) Parse() *JNode {
	p.advance()

	return p.expr()
}

func (p *JParser) advance() {
	p.TokenIndex++
	if p.TokenIndex < len(p.Tokens) {
		p.CurrentToken = p.Tokens[p.TokenIndex]
	}
}

func (p *JParser) factor() *JNode {
	currentToken := p.CurrentToken
	if currentToken == nil {
		return nil
	}

	switch currentToken.Type {
	case INT, FLOAT:
		p.advance()
		return &JNode{
			Type:  NUMBER,
			Token: currentToken,
		}
	}

	return nil
}

func (p *JParser) term() *JNode {
	return p.binOp(p.factor, mapset.NewSet(MUL, DIV))
}

func (p *JParser) expr() *JNode {
	return p.binOp(p.term, mapset.NewSet(PLUS, MINUS))
}

func (p *JParser) binOp(fun func() *JNode, ops mapset.Set) *JNode {
	leftNode := fun()

	for p.CurrentToken != nil && ops.Contains(p.CurrentToken.Type) {
		opToken := p.CurrentToken
		p.advance()
		rightNode := fun()
		leftNode = &JNode{
			Type:      OP,
			LeftNode:  leftNode,
			RightNode: rightNode,
			Token:     opToken,
		}
	}

	return leftNode
}
