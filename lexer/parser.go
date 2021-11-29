package lexer

import (
	"github.com/IfanTsai/jirachi/pkg/set"
	"github.com/pkg/errors"
)

type JNodeType int

const (
	Number JNodeType = iota
	BinOp
	UnaryOp
)

type JNode struct {
	Type      JNodeType
	Token     *JToken
	LeftNode  *JNode // for BinOp
	RightNode *JNode // for BinOp
	Node      *JNode // for UnaryOp
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

func (p *JParser) Parse() (*JNode, error) {
	p.advance()

	ast, err := p.expr()
	if err != nil {
		return nil, err
	}

	if p.CurrentToken.Type != EOF {
		return nil, errors.Wrap(&JInvalidSyntaxError{
			JError: &JError{
				StartPos: p.CurrentToken.StartPos,
				EndPos:   p.CurrentToken.EndPos,
			},
			Details: "Expected '+', '-', '*' or '/'",
		}, "failed to parse expr")
	}

	return ast, nil
}

func (p *JParser) advance() {
	p.TokenIndex++
	if p.TokenIndex < len(p.Tokens) {
		p.CurrentToken = p.Tokens[p.TokenIndex]
	}
}

/**
 * expr: term ( (PLUS | MINUS) term )*
 * term : factor ( ( MUL | DIV ) factor )*
 * factor: INT | FLOAT
 *         ( PLUS | MINUS ) factor
 *         LPAREN expr RPAREN
 */
func (p *JParser) factor() (*JNode, error) {
	currentToken := p.CurrentToken

	switch currentToken.Type {
	case INT, FLOAT:
		p.advance()

		return &JNode{
			Type:  Number,
			Token: currentToken,
		}, nil
	case PLUS, MINUS:
		p.advance()
		factor, err := p.factor()
		if err != nil {
			return nil, err
		}

		return &JNode{
			Type:  UnaryOp,
			Token: currentToken,
			Node:  factor,
		}, nil
	case LPAREN:
		p.advance()
		expr, err := p.expr()
		if err != nil {
			return nil, err
		}

		if p.CurrentToken.Type != RPAREN {
			return nil, errors.Wrap(&JInvalidSyntaxError{
				JError: &JError{
					StartPos: currentToken.StartPos,
					EndPos:   currentToken.EndPos,
				},
				Details: "Expected ')'",
			}, "failed to parse LPAREN")
		}

		p.advance()

		return expr, nil
	default:
		return nil, errors.Wrap(&JInvalidSyntaxError{
			JError: &JError{
				StartPos: currentToken.StartPos,
				EndPos:   currentToken.EndPos,
			},
			Details: "Expected '+', '-', '*' or '/'",
		}, "failed to parse factor")
	}
}

func (p *JParser) term() (*JNode, error) {
	return p.binOp(p.factor, set.NewSet(MUL, DIV))
}

func (p *JParser) expr() (*JNode, error) {
	return p.binOp(p.term, set.NewSet(PLUS, MINUS))
}

func (p *JParser) binOp(getNodeFunc func() (*JNode, error), ops *set.Set) (*JNode, error) {
	leftNode, err := getNodeFunc()
	if err != nil {
		return nil, err
	}

	for ops.Contains(p.CurrentToken.Type) {
		opToken := p.CurrentToken
		p.advance()
		rightNode, err := getNodeFunc()
		if err != nil {
			return nil, err
		}

		leftNode = &JNode{
			Type:      BinOp,
			LeftNode:  leftNode,
			RightNode: rightNode,
			Token:     opToken,
		}
	}

	return leftNode, nil
}
