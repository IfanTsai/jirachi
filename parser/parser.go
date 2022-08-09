package parser

import (
	"fmt"

	"github.com/IfanTsai/go-lib/set"

	"github.com/IfanTsai/jirachi/common"
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

	ast, err := p.statements(false)
	if err != nil {
		return nil, err
	}

	if p.CurrentToken.Type != token.EOF {
		return nil, p.createInvalidSyntaxError(
			"number, identifier, '+', '-', '(', '[', 'if', 'for', 'while', '*', '/' or '^'",
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

func (p *JParser) backTo(tokenIndex int) {
	p.TokenIndex = tokenIndex
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
	startPos := p.CurrentToken.StartPos

	p.advance()

	var elementNodes []JNode
	if p.CurrentToken.Type != token.RSQUARE {
		isFirstElement := true
		for isFirstElement || p.CurrentToken.Type == token.COMMA {
			if !isFirstElement {
				p.advance()
			} else {
				isFirstElement = false
			}

			expr, err := p.expr()
			if err != nil {
				return nil, p.createInvalidSyntaxError(
					"Expected ']', 'if', 'for', 'while', 'fun', number, identifier, '+', '-', '(', '[', or 'not'",
					"list expression",
				)
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

func (p *JParser) mapExpr() (JNode, error) {
	startPos := p.CurrentToken.StartPos

	p.advance()

	elementMap := make(map[JNode]JNode)
	if p.CurrentToken.Type != token.RBRACE {
		isFirstElement := true
		for isFirstElement || p.CurrentToken.Type == token.COMMA {
			if !isFirstElement {
				p.advance()
			} else {
				isFirstElement = false
			}

			KeyExpr, err := p.expr()
			if err != nil {
				return nil, p.createInvalidSyntaxError(
					"Expected '}', 'if', 'for', 'while', 'fun', number, identifier, '+', '-', '(', '[', or 'not'",
					"map expression",
				)
			}

			if p.CurrentToken.Type != token.COLON {
				return nil, p.createInvalidSyntaxError("':'", "map expression")
			}

			p.advance()

			valueExpr, err := p.expr()
			if err != nil {
				return nil, err
			}

			elementMap[KeyExpr] = valueExpr
		}

		if p.CurrentToken.Type != token.RBRACE {
			return nil, p.createInvalidSyntaxError("',' or '}'", "map expression")
		}
	}

	p.advance()
	return &JMapNode{
		JBaseNode: &JBaseNode{
			StartPos: startPos,
			EndPos:   p.CurrentToken.EndPos.Copy().Back(nil),
		},
		ElementMap: elementMap,
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

	var body JNode
	isBlock := false

	if p.CurrentToken.Type == token.NEWLINE {
		isBlock = true
		p.advance()

		body, err = p.statements(true)
		if err != nil {
			return nil, err
		}

		if !p.CurrentToken.Match(token.KEYWORD, token.END) {
			return nil, p.createInvalidSyntaxError(fmt.Sprintf("'%s'", token.END), "while expression")
		}

		p.advance()
	} else {
		body, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return &JWhileExprNode{
		JBaseNode: &JBaseNode{
			StartPos: conditionExpr.GetStartPos(),
			EndPos:   body.GetEndPos(),
		},
		ConditionNode:     conditionExpr,
		BodyNode:          body,
		IsBlockStatements: isBlock,
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

	var body JNode
	isBlock := false

	if p.CurrentToken.Type == token.NEWLINE {
		isBlock = true
		p.advance()

		body, err = p.statements(true)
		if err != nil {
			return nil, err
		}

		if !p.CurrentToken.Match(token.KEYWORD, token.END) {
			return nil, p.createInvalidSyntaxError(fmt.Sprintf("'%s'", token.END), "for expression")
		}

		p.advance()
	} else {
		body, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return &JForExprNode{
		JBaseNode: &JBaseNode{
			Token:    varNameToken,
			StartPos: varNameToken.StartPos,
			EndPos:   body.GetEndPos(),
		},
		StartValueNode:    startEXpr,
		EndValueNode:      endExpr,
		StepValueNode:     stepExpr,
		BodyNode:          body,
		IsBlockStatements: isBlock,
	}, nil
}

func (p *JParser) ifExpr() (JNode, error) {
	cases, elseCase, err := p.parseIfExprCases(token.IF)
	if err != nil {
		return nil, err
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

func (p *JParser) elifExpr() ([][2]JNode, JNode, error) {
	return p.parseIfExprCases(token.ELIF)
}

func (p *JParser) elseExpr() (JNode, error) {
	var err error
	var elseCase JNode

	if !p.CurrentToken.Match(token.KEYWORD, token.ELSE) {
		return nil, err
	}

	p.advance()

	if p.CurrentToken.Type == token.NEWLINE {
		elseCase, err = p.statements(true)
		if err != nil {
			return nil, err
		}

		if !p.CurrentToken.Match(token.KEYWORD, token.END) {
			return nil, p.createInvalidSyntaxError(fmt.Sprintf("'%s'", token.END), "else if expression")
		}

		p.advance()
	} else {
		elseCase, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return elseCase, nil
}

func (p *JParser) parseElifOrElseExpr() ([][2]JNode, JNode, error) {
	if p.CurrentToken.Match(token.KEYWORD, token.ELIF) {
		return p.elifExpr()
	} else if p.CurrentToken.Match(token.KEYWORD, token.ELSE) {
		elseCase, err := p.elseExpr()

		return nil, elseCase, err
	}

	return nil, nil, nil
}

func (p *JParser) parseIfExprCases(caseKeyword string) ([][2]JNode, JNode, error) {
	if !p.CurrentToken.Match(token.KEYWORD, caseKeyword) {
		return nil, nil, p.createInvalidSyntaxError(fmt.Sprintf("'%s'", caseKeyword), "case expression")
	}

	p.advance()

	condition, err := p.expr()
	if err != nil {
		return nil, nil, err
	}

	if !p.CurrentToken.Match(token.KEYWORD, token.THEN) {
		return nil, nil, p.createInvalidSyntaxError(fmt.Sprintf("'%s'", token.THEN), "case expression")
	}

	p.advance()

	var (
		cases    [][2]JNode // { condition expression, body }
		elseCase JNode
		newCases [][2]JNode
	)

	if p.CurrentToken.Type == token.NEWLINE {
		p.advance()

		statementNodes, err := p.statements(true)
		if err != nil {
			return nil, nil, err
		}

		cases = append(cases, [2]JNode{condition, statementNodes})

		if p.CurrentToken.Match(token.KEYWORD, token.END) {
			p.advance()
		} else {
			newCases, elseCase, err = p.parseElifOrElseExpr()
			if err != nil {
				return nil, nil, err
			}

			cases = append(cases, newCases...)
		}
	} else {
		expr, err := p.statement()
		if err != nil {
			return nil, nil, err
		}

		cases = append(cases, [2]JNode{condition, expr})

		newCases, elseCase, err = p.parseElifOrElseExpr()
		if err != nil {
			return nil, nil, err
		}

		cases = append(cases, newCases...)
	}

	return cases, elseCase, nil
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

	var (
		body JNode
		err  error
	)

	if p.CurrentToken.Type == token.ARROW {
		p.advance()

		body, err = p.expr()
		if err != nil {
			return nil, err
		}
	} else if p.CurrentToken.Type == token.NEWLINE {
		p.advance()

		body, err = p.statements(true)
		if err != nil {
			return nil, err
		}

		if !p.CurrentToken.Match(token.KEYWORD, token.END) {
			return nil, p.createInvalidSyntaxError(fmt.Sprintf("'%s'", token.END), "function definition expression")
		}

		p.advance()
	} else {
		return nil, p.createInvalidSyntaxError("'->' or NEWLINE", "function definition expression")
	}

	return &JFuncDefNode{
		JBaseNode: &JBaseNode{
			Token: varNameToken,
		},
		ArgTokens: argTokens,
		BodyNode:  body,
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
	case token.LBRACE:
		return p.mapExpr()
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

	return nil, p.createInvalidSyntaxError("number, identifier, '(', '[' or 'NOT'", "factor")
}

func (p *JParser) call() (JNode, error) {
	atom, err := p.atom()
	if err != nil {
		return nil, err
	}

	if _, ok := atom.(*JVarAccessNode); ok {
		if val, ok := atom.GetToken().Value.(string); ok && val == "@" {
			if p.CurrentToken.Type == token.STRING {
				currentToken := p.CurrentToken
				p.advance()

				return &JCallExprNode{
					JBaseNode: &JBaseNode{
						StartPos: atom.GetStartPos(),
						EndPos:   p.CurrentToken.EndPos,
					},
					CallNode: atom,
					ArgNodes: []JNode{
						&JStringNode{
							JBaseNode: &JBaseNode{
								Token:    currentToken,
								StartPos: currentToken.StartPos,
								EndPos:   currentToken.EndPos,
							},
						},
					},
				}, nil
			}
		}
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
					"Expected ')', 'if', 'for', 'while', 'fun', number, identifier, '+', '-', '(', '[' or 'not'",
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

func (p *JParser) statement() (JNode, error) {
	currentToken := p.CurrentToken

	if currentToken.IsKeyWord() {
		switch p.CurrentToken.Value {
		case token.RETURN:
			p.advance()

			tokenIndex := p.TokenIndex
			expr, err := p.expr()
			if err != nil {
				p.backTo(tokenIndex)
			}

			var endPos *common.JPosition
			if expr == nil {
				endPos = currentToken.EndPos
			} else {
				endPos = expr.GetEndPos()
			}

			return &JReturnNode{
				JBaseNode: &JBaseNode{
					Token:    currentToken,
					StartPos: currentToken.StartPos,
					EndPos:   endPos,
				},
				ReturnNode: expr,
			}, nil
		case token.BREAK:
			p.advance()

			return &JBreakNode{
				JBaseNode: &JBaseNode{
					Token:    currentToken,
					StartPos: currentToken.StartPos,
					EndPos:   currentToken.EndPos,
				},
			}, nil
		case token.CONTINUE:
			p.advance()

			return &JContinueNode{
				JBaseNode: &JBaseNode{
					Token:    currentToken,
					StartPos: currentToken.StartPos,
					EndPos:   p.CurrentToken.EndPos,
				},
			}, nil
		}
	}

	return p.expr()
}

func (p *JParser) statements(isBlock bool) (JNode, error) {
	startPos := p.CurrentToken.StartPos

	for p.CurrentToken.Type == token.NEWLINE {
		p.advance()
	}

	statementNode, err := p.statement()
	if err != nil {
		return nil, err
	}

	statementNodes := []JNode{statementNode}

	moreStatements := false
	for {
		for p.CurrentToken.Type == token.NEWLINE {
			p.advance()
			moreStatements = true
		}

		if !moreStatements || p.CurrentToken.Type == token.EOF {
			break
		}

		tokenIndex := p.TokenIndex
		statementNode, err = p.statement()
		if err != nil {
			p.backTo(tokenIndex)

			break
		}

		statementNodes = append(statementNodes, statementNode)
	}

	if len(statementNodes) == 1 {
		return statementNodes[0], nil
	}

	return &JListNode{
		JBaseNode: &JBaseNode{
			StartPos: startPos,
			EndPos:   statementNodes[len(statementNodes)-1].GetEndPos(),
		},
		ElementNodes:      statementNodes,
		IsBlockStatements: isBlock,
	}, nil
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

func (p *JParser) createInvalidSyntaxError(expected, parseType string) error {
	return errors.Wrap(&common.JInvalidSyntaxError{
		JError: &common.JError{
			StartPos: p.CurrentToken.StartPos,
			EndPos:   p.CurrentToken.EndPos,
		},
		Details: "Expected " + expected,
	}, "failed to parse "+parseType)
}
