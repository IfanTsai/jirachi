package parser

import (
	"fmt"

	"github.com/IfanTsai/jirachi/common"
	"github.com/IfanTsai/jirachi/pkg/set"
	"github.com/IfanTsai/jirachi/token"
	"github.com/pkg/errors"
)

type getNodeFunc func() (JNode, error)

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

func (p *JParser) Parse() (JNode, error) {
	p.advance()

	ast, err := p.expr()
	if err != nil {
		return nil, err
	}

	if p.CurrentToken.Type != token.EOF {
		return nil, p.createInvalidSyntaxError(
			"number, identifier, '+', '-', '(', '[', 'IF', 'FOR', 'WHILE', '*', '/' or '^'",
			"expression",
		)
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

func (p *JParser) indexExpr(atomNode JNode) (JNode, error) {
	startPos := p.CurrentToken.StartPos

	p.advance()

	expr, err := p.expr()
	if err != nil {
		return nil, err
	}

	if p.CurrentToken.Type != token.RSQUARE {
		return nil, p.createInvalidSyntaxError("']'", "index expression")
	}

	p.advance()

	indexExprNode := &JIndexExprNode{
		JBaseNode: &JBaseNode{
			StartPos: startPos,
			EndPos:   p.CurrentToken.EndPos.Copy().Back(nil),
		},
		IndexNode: atomNode,
		IndexExpr: expr,
	}

	if p.CurrentToken.Type == token.EQ {
		p.advance()

		expr, err := p.expr()
		if err != nil {
			return nil, err
		}

		return &JVarIndexAssignNode{
			JVarAssignNode: &JVarAssignNode{
				JBaseNode: &JBaseNode{
					Token:    atomNode.GetToken(),
					StartPos: startPos,
					EndPos:   expr.GetEndPos(),
				},
				Node: expr,
			},
			IndexExprNode: indexExprNode,
		}, nil
	}

	return indexExprNode, nil
}

func (p *JParser) listExpr() (JNode, error) {
	if p.CurrentToken.Type != token.LSQUARE {
		return nil, p.createInvalidSyntaxError("'['", "list expression")
	}

	startPos := p.CurrentToken.StartPos

	p.advance()

	var elementNodes []JNode
	if p.CurrentToken.Type != token.RSQUARE {
		expr, err := p.expr()
		if err != nil {
			return nil, p.createInvalidSyntaxError(
				"Expected ']', 'IF', 'FOR', 'WHILE', 'FUN', int, float, identifier, '+', '-', '(', '[', or 'NOT'",
				"list expression",
			)
		}

		elementNodes = append(elementNodes, expr)

		for p.CurrentToken.Type == token.COMMA {
			p.advance()

			expr, err := p.expr()
			if err != nil {
				return nil, err
			}
			elementNodes = append(elementNodes, expr)
		}

		if p.CurrentToken.Type != token.RSQUARE {
			return nil, p.createInvalidSyntaxError("',' or ']'", "list expression")
		}
	}

	p.advance()
	return &JListNode{
		JBaseNode: &JBaseNode{
			StartPos: startPos,
			EndPos:   p.CurrentToken.EndPos.Copy().Back(nil),
		},
		ElementNodes: elementNodes,
	}, nil
}

func (p *JParser) whileExpr() (JNode, error) {
	if !p.CurrentToken.Match(token.KEYWORD, token.WHILE) {
		return nil, p.createInvalidSyntaxError(fmt.Sprintf("'%s'", token.WHILE), "while expression")
	}

	p.advance()

	conditionExpr, err := p.expr()
	if err != nil {
		return nil, err
	}

	if !p.CurrentToken.Match(token.KEYWORD, token.THEN) {
		return nil, p.createInvalidSyntaxError(fmt.Sprintf("'%s'", token.THEN), "while expression")
	}

	p.advance()

	bodyExpr, err := p.expr()
	if err != nil {
		return nil, err
	}

	return &JWhileExprNode{
		JBaseNode: &JBaseNode{
			StartPos: conditionExpr.GetStartPos(),
			EndPos:   bodyExpr.GetEndPos(),
		},
		ConditionNode: conditionExpr,
		BodyNode:      bodyExpr,
	}, nil
}

func (p *JParser) forExpr() (JNode, error) {
	if !p.CurrentToken.Match(token.KEYWORD, token.FOR) {
		return nil, p.createInvalidSyntaxError(fmt.Sprintf("'%s'", token.FOR), "for expression")
	}

	p.advance()

	if p.CurrentToken.Type != token.IDENTIFIER {
		return nil, p.createInvalidSyntaxError("identifier", "for expression")
	}

	varNameToken := p.CurrentToken

	p.advance()

	if p.CurrentToken.Type != token.EQ {
		return nil, p.createInvalidSyntaxError("'='", "for expression")
	}

	p.advance()

	startEXpr, err := p.expr()
	if err != nil {
		return nil, err
	}

	if !p.CurrentToken.Match(token.KEYWORD, token.TO) {
		return nil, p.createInvalidSyntaxError(fmt.Sprintf("'%s'", token.TO), "for expression")
	}

	p.advance()

	endExpr, err := p.expr()
	if err != nil {
		return nil, err
	}

	var stepExpr JNode
	if p.CurrentToken.Match(token.KEYWORD, token.STEP) {
		p.advance()

		stepExpr, err = p.expr()
		if err != nil {
			return nil, err
		}
	}

	if !p.CurrentToken.Match(token.KEYWORD, token.THEN) {
		return nil, p.createInvalidSyntaxError(fmt.Sprintf("'%s'", token.THEN), "for expression")
	}

	p.advance()

	bodyExpr, err := p.expr()
	if err != nil {
		return nil, err
	}

	return &JForExprNode{
		JBaseNode: &JBaseNode{
			Token:    varNameToken,
			StartPos: varNameToken.StartPos,
			EndPos:   bodyExpr.GetEndPos(),
		},
		StartValueNode: startEXpr,
		EndValueNode:   endExpr,
		StepValueNode:  stepExpr,
		BodyNode:       bodyExpr,
	}, nil
}

func (p *JParser) ifExpr() (JNode, error) {
	if !p.CurrentToken.Match(token.KEYWORD, token.IF) {
		return nil, p.createInvalidSyntaxError(fmt.Sprintf("'%s'", token.IF), "if expression")
	}

	var elseCase JNode
	cases := make([][2]JNode, 0, 3)

	cases, err := p.parseElifThenExpr(cases)
	if err != nil {
		return nil, err
	}

	for p.CurrentToken.Match(token.KEYWORD, token.ELIF) {
		cases, err = p.parseElifThenExpr(cases)
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
		endPos = elseCase.GetEndPos()
	} else {
		endPos = cases[len(cases)-1][0].GetEndPos()
	}

	return &JIfExprNode{
		JBaseNode: &JBaseNode{
			StartPos: cases[0][0].GetStartPos(),
			EndPos:   endPos,
		},
		CaseNodes:    cases,
		ElseCaseNode: elseCase,
	}, nil
}

func (p *JParser) funcDef() (JNode, error) {
	if !p.CurrentToken.Match(token.KEYWORD, token.FUN) {
		return nil, p.createInvalidSyntaxError(fmt.Sprintf("'%s'", token.FUN), "function definition")
	}

	p.advance()

	var varNameToken *token.JToken
	if p.CurrentToken.Type == token.IDENTIFIER {
		varNameToken = p.CurrentToken

		p.advance()
	}

	if p.CurrentToken.Type != token.LPAREN {
		return nil, p.createInvalidSyntaxError("'('", "function definition")
	}

	p.advance()

	var argTokens []*token.JToken
	if p.CurrentToken.Type == token.IDENTIFIER {
		argTokens = append(argTokens, p.CurrentToken)

		p.advance()

		for p.CurrentToken.Type == token.COMMA {
			p.advance()

			if p.CurrentToken.Type != token.IDENTIFIER {
				return nil, p.createInvalidSyntaxError("identifier", "function definition")
			}

			argTokens = append(argTokens, p.CurrentToken)

			p.advance()
		}
	}

	if p.CurrentToken.Type != token.RPAREN {
		return nil, p.createInvalidSyntaxError("')'", "function definition")
	}

	p.advance()

	if p.CurrentToken.Type != token.ARROW {
		return nil, p.createInvalidSyntaxError("'->'", "function definition")
	}

	p.advance()

	bodyNode, err := p.expr()
	if err != nil {
		return nil, err
	}

	return &JFuncDefNode{
		JBaseNode: &JBaseNode{
			Token: varNameToken,
		},
		ArgTokens: argTokens,
		BodyNode:  bodyNode,
	}, nil
}

func (p *JParser) atom() (JNode, error) {
	currentToken := p.CurrentToken

	switch currentToken.Type {
	case token.INT, token.FLOAT:
		p.advance()

		return &JNumberNode{
			JBaseNode: &JBaseNode{
				Token:    currentToken,
				StartPos: currentToken.StartPos,
				EndPos:   currentToken.EndPos,
			},
		}, nil
	case token.STRING:
		p.advance()

		return &JStringNode{
			JBaseNode: &JBaseNode{
				Token:    currentToken,
				StartPos: currentToken.StartPos,
				EndPos:   currentToken.EndPos,
			},
		}, nil
	case token.IDENTIFIER:
		p.advance()

		atomNode := &JVarAccessNode{
			JBaseNode: &JBaseNode{
				Token:    currentToken,
				StartPos: currentToken.StartPos,
				EndPos:   currentToken.EndPos,
			},
		}

		if p.CurrentToken.Type == token.LSQUARE {
			return p.indexExpr(atomNode)
		} else {
			return atomNode, nil
		}
	case token.LPAREN:
		p.advance()
		expr, err := p.expr()
		if err != nil {
			return nil, err
		}

		if p.CurrentToken.Type != token.RPAREN {
			return nil, p.createInvalidSyntaxError("')'", "LPAREN")
		}

		p.advance()

		return expr, nil
	case token.LSQUARE:
		return p.listExpr()
	case token.KEYWORD:
		switch currentToken.Value {
		case token.IF:
			return p.ifExpr()
		case token.FOR:
			return p.forExpr()
		case token.WHILE:
			return p.whileExpr()
		case token.FUN:
			return p.funcDef()
		}
	}

	return nil, p.createInvalidSyntaxError("number, identifier, '+', '-', '(', '[' or 'NOT'", "factor")
}

func (p *JParser) call() (JNode, error) {
	atom, err := p.atom()
	if err != nil {
		return nil, err
	}

	if p.CurrentToken.Type == token.LPAREN {
		p.advance()

		var argNodes []JNode
		if p.CurrentToken.Type == token.RPAREN {
			p.advance()
		} else {
			expr, err := p.expr()
			if err != nil {
				return nil, p.createInvalidSyntaxError(
					"Expected ')', 'IF', 'FOR', 'WHILE', 'FUN', int, float, identifier, '+', '-', '(', '[' or 'NOT'",
					"call expression",
				)
			}

			argNodes = append(argNodes, expr)

			for p.CurrentToken.Type == token.COMMA {
				p.advance()

				expr, err := p.expr()
				if err != nil {
					return nil, err
				}
				argNodes = append(argNodes, expr)
			}

			if p.CurrentToken.Type != token.RPAREN {
				return nil, p.createInvalidSyntaxError("',' or ')'", "call expression")
			}

			p.advance()
		}

		var endPos *common.JPosition
		if len(argNodes) > 0 {
			endPos = argNodes[len(argNodes)-1].GetEndPos()
		} else {
			endPos = atom.GetEndPos()
		}

		return &JCallExprNode{
			JBaseNode: &JBaseNode{
				StartPos: atom.GetStartPos(),
				EndPos:   endPos,
			},
			CallNode: atom,
			ArgNodes: argNodes,
		}, nil
	}

	return atom, nil
}

func (p *JParser) power() (JNode, error) {
	return p.binOp(p.call, set.NewSet(token.POW), p.factor)
}

func (p *JParser) factor() (JNode, error) {
	currentToken := p.CurrentToken

	switch currentToken.Type {
	case token.PLUS, token.MINUS:
		p.advance()
		factor, err := p.factor()
		if err != nil {
			return nil, err
		}

		return &JUnaryOpNode{
			JBaseNode: &JBaseNode{
				Token:    currentToken,
				StartPos: currentToken.StartPos,
				EndPos:   factor.GetEndPos(),
			},
			Node: factor,
		}, nil

	default:
		return p.power()
	}
}

func (p *JParser) term() (JNode, error) {
	return p.binOp(p.factor, set.NewSet(token.MUL, token.DIV), nil)
}

func (p *JParser) arithmeticExpr() (JNode, error) {
	return p.binOp(p.term, set.NewSet(token.PLUS, token.MINUS), nil)
}

func (p *JParser) compareExpr() (JNode, error) {
	currentToken := p.CurrentToken

	if currentToken.Match(token.KEYWORD, token.NOT) {
		p.advance()

		compExpr, err := p.compareExpr()
		if err != nil {
			return nil, err
		}

		return &JUnaryOpNode{
			JBaseNode: &JBaseNode{
				Token:    currentToken,
				StartPos: currentToken.StartPos,
				EndPos:   compExpr.GetEndPos(),
			},
			Node: compExpr,
		}, nil
	}

	return p.binOp(p.arithmeticExpr, set.NewSet(token.EE, token.NE, token.LT, token.LTE, token.GT, token.GTE), nil)
}

func (p *JParser) expr() (JNode, error) {
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

			return &JVarAssignNode{
				JBaseNode: &JBaseNode{
					Token:    varToken,
					StartPos: varToken.StartPos,
					EndPos:   expr.GetEndPos(),
				},
				Node: expr,
			}, nil
		} else if p.CurrentToken.Type == token.IDENTIFIER {
			// no support consecutive identifiers
			return nil, p.createInvalidSyntaxError("'+', '-', '*', '/' or '^'", "expression")
		} else {
			// go back when it is not an assignment expression
			p.back()
		}
	}

	return p.binOp(p.compareExpr, set.NewSet(token.AND, token.OR), nil)
}

func (p *JParser) binOp(getNodeFuncA getNodeFunc, ops *set.Set, getNodeFuncB getNodeFunc) (JNode, error) {
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

		leftNode = &JBinOpNode{
			JBaseNode: &JBaseNode{
				Token:    opToken,
				StartPos: leftNode.GetStartPos(),
				EndPos:   rightNode.GetEndPos(),
			},
			LeftNode:  leftNode,
			RightNode: rightNode,
		}
	}

	return leftNode, nil
}

func (p *JParser) parseElifThenExpr(cases [][2]JNode) ([][2]JNode, error) {
	p.advance()

	condition, err := p.expr()
	if err != nil {
		return nil, err
	}

	if !p.CurrentToken.Match(token.KEYWORD, token.THEN) {
		return nil, p.createInvalidSyntaxError(fmt.Sprintf("'%s'", token.THEN), "if expression")
	}

	p.advance()

	expr, err := p.expr()
	if err != nil {
		return nil, err
	}

	cases = append(cases, [2]JNode{condition, expr})

	return cases, nil
}

func (p *JParser) createInvalidSyntaxError(expected, parseType string) error {
	return errors.Wrap(&common.JInvalidSyntaxError{
		JError: &common.JError{
			StartPos: p.CurrentToken.StartPos,
			EndPos:   p.CurrentToken.EndPos,
		},
		Details: "Expected " + expected,
	}, "failed to parse "+parseType)
}
