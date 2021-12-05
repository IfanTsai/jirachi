package interpreter

import (
	"fmt"

	"github.com/IfanTsai/jirachi/common"
	"github.com/IfanTsai/jirachi/lexer"
	"github.com/IfanTsai/jirachi/parser"
	"github.com/IfanTsai/jirachi/token"
	"github.com/pkg/errors"
)

var GlobalSymbolTable = common.NewJSymbolTable(nil).
	Set("NULL", NewJNumber(0)).
	Set("TRUE", NewJNumber(1)).
	Set("FALSE", NewJNumber(0))

func Run(filename, text string) (interface{}, error) {
	// generate tokens
	tokens, err := lexer.NewJLexer(filename, text).MakeTokens()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to make tokens")
	}

	// generate AST
	ast, err := parser.NewJParser(tokens, -1).Parse()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse tokens")
	}

	// run program
	context := common.NewJContext("<program>", GlobalSymbolTable, nil, nil)
	number, err := NewJInterpreter(context).Interpreter(ast)
	if err != nil {
		return nil, err
	}

	return number.Value, nil
}

type JInterpreter struct {
	Context *common.JContext
}

func NewJInterpreter(context *common.JContext) *JInterpreter {
	return &JInterpreter{
		Context: context,
	}
}

func (i *JInterpreter) Interpreter(ast parser.JNode) (*JNumber, error) {
	return i.visit(ast)
}

func (i *JInterpreter) visit(node parser.JNode) (*JNumber, error) {
	switch node.Type() {
	case parser.Number:
		return i.visitNumberNode(node.(*parser.JNumberNode))
	case parser.BinOp:
		return i.visitBinOpNode(node.(*parser.JBinOpNode))
	case parser.UnaryOp:
		return i.visitUnaryOpNode(node.(*parser.JUnaryOpNode))
	case parser.VarAssign:
		return i.visitVarAssignNode(node.(*parser.JVarAssignNode))
	case parser.VarAccess:
		return i.visitVarAccessNode(node.(*parser.JVarAccessNode))
	case parser.IfExpr:
		return i.visitIfExprNode(node.(*parser.JIfExprNode))
	default:
		return nil, errors.Wrap(&common.JInvalidSyntaxError{
			JError: &common.JError{
				StartPos: node.GetStartPos(),
				EndPos:   node.GetEndPos(),
			},
			Details: "Expected number, bin op or unary op",
		}, "failed to visit node")
	}
}

func (i *JInterpreter) visitNumberNode(node *parser.JNumberNode) (*JNumber, error) {
	return NewJNumber(node.Token.Value).SetJPos(node.StartPos, node.EndPos).SetJContext(i.Context), nil
}

func (i *JInterpreter) visitVarAssignNode(node *parser.JVarAssignNode) (*JNumber, error) {
	varName := node.Token.Value

	value, err := i.visit(node.Node)
	if err != nil {
		return nil, err
	}

	i.Context.SymbolTable.Set(varName, value)

	return value, nil
}

func (i *JInterpreter) visitIfExprNode(node *parser.JIfExprNode) (*JNumber, error) {
	for index := range node.Cases {
		condition := node.Cases[index][0]
		expr := node.Cases[index][1]

		conditionValue, err := i.visit(condition)
		if err != nil {
			return nil, err
		}

		if conditionValue.IsTrue() {
			exprValue, err := i.visit(expr)
			if err != nil {
				return nil, err
			}

			return exprValue.SetJContext(i.Context), nil
		}
	}

	if node.ElseCase != nil {
		elseValue, err := i.visit(node.ElseCase)
		if err != nil {
			return nil, err
		}

		return elseValue.SetJContext(i.Context), nil
	}

	// eg. IF false THEN 123
	return NewJNumber(nil), nil
}

func (i *JInterpreter) visitVarAccessNode(node *parser.JVarAccessNode) (*JNumber, error) {
	varName := node.Token.Value
	value := i.Context.SymbolTable.Get(varName)
	if value == nil {
		return nil, errors.Wrap(&common.JRunTimeError{
			JError: &common.JError{
				StartPos: node.StartPos,
				EndPos:   node.EndPos,
			},
			Context: i.Context,
			Details: fmt.Sprintf("'%v' is not defined", varName),
		}, "failed to access variable")
	}

	return NewJNumber(value.(*JNumber).Value).SetJPos(node.StartPos, node.EndPos).SetJContext(i.Context), nil
}

func (i *JInterpreter) visitBinOpNode(node *parser.JBinOpNode) (*JNumber, error) {
	leftNumber, err := i.visit(node.LeftNode)
	if err != nil {
		return nil, err
	}

	rightNUmber, err := i.visit(node.RightNode)
	if err != nil {
		return nil, err
	}

	var resNumber *JNumber

	switch node.Token.Type {
	case token.PLUS:
		resNumber, err = leftNumber.AddTo(rightNUmber)
	case token.MINUS:
		resNumber, err = leftNumber.SubBy(rightNUmber)
	case token.MUL:
		resNumber, err = leftNumber.MulBy(rightNUmber)
	case token.DIV:
		resNumber, err = leftNumber.DivBy(rightNUmber)
	case token.POW:
		resNumber, err = leftNumber.PowBy(rightNUmber)
	case token.EE:
		resNumber, err = leftNumber.EqualTo(rightNUmber)
	case token.NE:
		resNumber, err = leftNumber.NotEqualTo(rightNUmber)
	case token.LT:
		resNumber, err = leftNumber.LessThen(rightNUmber)
	case token.LTE:
		resNumber, err = leftNumber.LessThenOrEqualTo(rightNUmber)
	case token.GT:
		resNumber, err = leftNumber.GreaterThen(rightNUmber)
	case token.GTE:
		resNumber, err = leftNumber.GreaterThenOrEqualTo(rightNUmber)
	case token.KEYWORD:
		switch node.Token.Value {
		case token.AND:
			resNumber, err = leftNumber.AndBy(rightNUmber)
		case token.OR:
			resNumber, err = leftNumber.OrBy(rightNUmber)
		}
	default:
		return nil, errors.Wrap(&common.JInvalidSyntaxError{
			JError: &common.JError{
				StartPos: node.StartPos,
				EndPos:   node.EndPos,
			},
			Details: "Expected '+', '-', '*' or '/'",
		}, "failed to visit bin op node")
	}

	if err != nil {
		return nil, err
	}

	return resNumber.SetJPos(node.StartPos, node.EndPos), nil
}

func (i *JInterpreter) visitUnaryOpNode(node *parser.JUnaryOpNode) (*JNumber, error) {
	number, err := i.visit(node.Node)
	if err != nil {
		return nil, err
	}
	switch {
	case node.Token.Type == token.MINUS:
		number, err = number.MulBy(NewJNumber(-1))
	case node.Token.Match(token.KEYWORD, token.NOT):
		number, err = number.Not()
	}

	if err != nil {
		return nil, err
	}

	return number.SetJPos(node.StartPos, node.EndPos), nil
}
