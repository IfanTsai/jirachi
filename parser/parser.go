package parser

import (
	"github.com/IfanTsai/jirachi/common"
	"github.com/IfanTsai/jirachi/pkg/set"
	"github.com/IfanTsai/jirachi/token"
	"github.com/pkg/errors"
)

type getNodeFunc func() (*JNode, error)

type JParser struct {
	TokenIndex   int
	Tokens       []*token.JToken
	CurrentToken *token.JToken
}

func NewJParser(tokens []*token.JToken, tokenIndex int) *JParser {
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

	if p.CurrentToken.Type != token.EOF {
		return nil, errors.Wrap(&common.JInvalidSyntaxError{
			JError: &common.JError{
				StartPos: p.CurrentToken.StartPos,
				EndPos:   p.CurrentToken.EndPos,
			},
			Details: "Expected '+', '-', '*', '/' or '^'",
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

func (p *JParser) back() {
	p.TokenIndex--
	if p.TokenIndex >= 0 {
		p.CurrentToken = p.Tokens[p.TokenIndex]
	}
}

func (p *JParser) ifExpr() (*JNode, error) {
	if !p.CurrentToken.Match(token.KEYWORD, token.IF) {
		return nil, errors.Wrap(&common.JInvalidSyntaxError{
			JError: &common.JError{
				StartPos: p.CurrentToken.StartPos,
				EndPos:   p.CurrentToken.EndPos,
			},
			Details: "Expected 'IF'",
		}, "failed to parse if expression")
	}

	var elseCase *JNode
	cases := make([][2]*JNode, 0, 3)

	cases, err := p.parseThenExpr(cases)
	if err != nil {
		return nil, err
	}

	for p.CurrentToken.Match(token.KEYWORD, token.ELIF) {
		cases, err = p.parseThenExpr(cases)
		if err != nil {
			return nil, err
		}
	}

	if p.CurrentToken.Match(token.KEYWORD, token.ELSE) {
		p.advance()

		elseCase, err = p.expr()
		if err != nil {
			return nil, err
		}
	}

	var endPos *common.JPosition
	if elseCase != nil {
		endPos = elseCase.EndPos
	} else {
		endPos = cases[len(cases)-1][0].EndPos
	}

	return &JNode{
		Type:     IfExpr,
		Cases:    cases,
		ElseCase: elseCase,
		StartPos: cases[0][0].StartPos,
		EndPos:   endPos,
	}, nil

}

func (p *JParser) atom() (*JNode, error) {
	currentToken := p.CurrentToken

	switch currentToken.Type {
	case token.INT, token.FLOAT:
		p.advance()

		return &JNode{
			Type:     Number,
			Token:    currentToken,
			StartPos: currentToken.StartPos,
			EndPos:   currentToken.EndPos,
		}, nil
	case token.IDENTIFIER:
		p.advance()

		return &JNode{
			Type:     VarAccess,
			Token:    currentToken,
			StartPos: currentToken.StartPos,
			EndPos:   currentToken.EndPos,
		}, nil
	case token.LPAREN:
		p.advance()
		expr, err := p.expr()
		if err != nil {
			return nil, err
		}

		if p.CurrentToken.Type != token.RPAREN {
			return nil, errors.Wrap(&common.JInvalidSyntaxError{
				JError: &common.JError{
					StartPos: currentToken.StartPos,
					EndPos:   currentToken.EndPos,
				},
				Details: "Expected ')'",
			}, "failed to parse LPAREN")
		}

		p.advance()

		return expr, nil
	case token.KEYWORD:
		switch currentToken.Value {
		case token.IF:
			return p.ifExpr()
		}
	}

	return nil, errors.Wrap(&common.JInvalidSyntaxError{
		JError: &common.JError{
			StartPos: currentToken.StartPos,
			EndPos:   currentToken.EndPos,
		},
		Details: "Expected int, float, identifier, '+', '-', '(' or 'NOT'",
	}, "failed to parse factor")
}

func (p *JParser) power() (*JNode, error) {
	return p.binOp(p.atom, set.NewSet(token.POW), p.factor)
}

func (p *JParser) factor() (*JNode, error) {
	currentToken := p.CurrentToken

	switch currentToken.Type {
	case token.PLUS, token.MINUS:
		p.advance()
		factor, err := p.factor()
		if err != nil {
			return nil, err
		}

		return &JNode{
			Type:     UnaryOp,
			Token:    currentToken,
			Node:     factor,
			StartPos: currentToken.StartPos,
			EndPos:   factor.EndPos,
		}, nil

	default:
		return p.power()
	}
}

func (p *JParser) term() (*JNode, error) {
	return p.binOp(p.factor, set.NewSet(token.MUL, token.DIV), nil)
}

func (p *JParser) arithmeticExpr() (*JNode, error) {
	return p.binOp(p.term, set.NewSet(token.PLUS, token.MINUS), nil)
}

func (p *JParser) compareExpr() (*JNode, error) {
	currentToken := p.CurrentToken

	if currentToken.Match(token.KEYWORD, token.NOT) {
		p.advance()

		compExpr, err := p.compareExpr()
		if err != nil {
			return nil, err
		}

		return &JNode{
			Type:     UnaryOp,
			Token:    currentToken,
			Node:     compExpr,
			StartPos: currentToken.StartPos,
			EndPos:   compExpr.EndPos,
		}, nil
	}

	return p.binOp(p.arithmeticExpr, set.NewSet(token.EE, token.NE, token.LT, token.LTE, token.GT, token.GTE), nil)
}

func (p *JParser) expr() (*JNode, error) {
	// check if it is an assignment expression
	if p.CurrentToken.Type == token.IDENTIFIER {
		varToken := p.CurrentToken
		p.advance()

		if p.CurrentToken.Type == token.EQ {
			p.advance()

			expr, err := p.expr()
			if err != nil {
				return nil, err
			}

			return &JNode{
				Type:     VarAssign,
				Token:    varToken,
				Node:     expr,
				StartPos: varToken.StartPos,
				EndPos:   expr.EndPos,
			}, nil
		} else if p.CurrentToken.Type == token.IDENTIFIER {
			// no support consecutive identifiers
			return nil, errors.Wrap(&common.JInvalidSyntaxError{
				JError: &common.JError{
					StartPos: p.CurrentToken.StartPos,
					EndPos:   p.CurrentToken.EndPos,
				},
				Details: "Expected '+', '-', '*', '/' or '^'",
			}, "failed to parse expr")
		} else {
			// go back when it is not an assignment expression
			p.back()
		}
	}

	return p.binOp(p.compareExpr, set.NewSet(token.AND, token.OR), nil)
}

func (p *JParser) binOp(getNodeFuncA getNodeFunc, ops *set.Set, getNodeFuncB getNodeFunc) (*JNode, error) {
	if getNodeFuncB == nil {
		getNodeFuncB = getNodeFuncA
	}

	leftNode, err := getNodeFuncA()
	if err != nil {
		return nil, err
	}

	for (p.CurrentToken.Type == token.KEYWORD && ops.Contains(p.CurrentToken.Value)) ||
		ops.Contains(p.CurrentToken.Type) {

		opToken := p.CurrentToken
		p.advance()
		rightNode, err := getNodeFuncB()
		if err != nil {
			return nil, err
		}

		leftNode = &JNode{
			Type:      BinOp,
			LeftNode:  leftNode,
			RightNode: rightNode,
			Token:     opToken,
			StartPos:  leftNode.StartPos,
			EndPos:    rightNode.EndPos,
		}
	}

	return leftNode, nil
}

func (p *JParser) parseThenExpr(cases [][2]*JNode) ([][2]*JNode, error) {
	p.advance()

	condition, err := p.expr()
	if err != nil {
		return nil, err
	}

	if !p.CurrentToken.Match(token.KEYWORD, token.THEN) {
		return nil, errors.Wrap(&common.JInvalidSyntaxError{
			JError: &common.JError{
				StartPos: p.CurrentToken.StartPos,
				EndPos:   p.CurrentToken.EndPos,
			},
			Details: "Expected 'THEN'",
		}, "failed to parse if expression")
	}

	p.advance()

	expr, err := p.expr()
	if err != nil {
		return nil, err
	}

	cases = append(cases, [2]*JNode{condition, expr})

	return cases, nil
}
