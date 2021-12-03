package interpreter

import (
	"fmt"

	"github.com/IfanTsai/jirachi/common"
	"github.com/IfanTsai/jirachi/lexer"
	"github.com/IfanTsai/jirachi/parser"
	"github.com/IfanTsai/jirachi/token"
	"github.com/pkg/errors"
)

var GlobalSymbolTable = common.NewJSymbolTable(nil).Set("null", NewJNumber(0))

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
	number, err := NewJInterpreter(context).Visit(ast)
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

func (i *JInterpreter) Visit(node *parser.JNode) (*JNumber, error) {
	switch node.Type {
	case parser.Number:
		return i.visitNumberNode(node)
	case parser.BinOp:
		return i.visitBinOpNode(node)
	case parser.UnaryOp:
		return i.visitUnaryOpNode(node)
	case parser.VarAssign:
		return i.visitVarAssignNode(node)
	case parser.VarAccess:
		return i.visitVarAccessNode(node)
	default:
		return nil, errors.Wrap(&common.JInvalidSyntaxError{
			JError: &common.JError{
				StartPos: node.StartPos,
				EndPos:   node.EndPos,
			},
			Details: "Expected number, bin op or unary op",
		}, "failed to visit node")
	}
}

func (i *JInterpreter) visitNumberNode(node *parser.JNode) (*JNumber, error) {
	return NewJNumber(node.Token.Value).SetJPos(node.StartPos, node.EndPos).SetJContext(i.Context), nil
}

func (i *JInterpreter) visitVarAssignNode(node *parser.JNode) (*JNumber, error) {
	varName := node.Token.Value

	value, err := i.Visit(node.Node)
	if err != nil {
		return nil, err
	}

	i.Context.SymbolTable.Set(varName, value)

	return value, nil
}

func (i *JInterpreter) visitVarAccessNode(node *parser.JNode) (*JNumber, error) {
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

func (i *JInterpreter) visitBinOpNode(node *parser.JNode) (*JNumber, error) {
	leftNumber, err := i.Visit(node.LeftNode)
	if err != nil {
		return nil, err
	}

	rightNUmber, err := i.Visit(node.RightNode)
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

func (i *JInterpreter) visitUnaryOpNode(node *parser.JNode) (*JNumber, error) {
	number, err := i.Visit(node.Node)
	if err != nil {
		return nil, err
	}
	if node.Token.Type == token.MINUS {
		number, err = number.MulBy(NewJNumber(-1))
		if err != nil {
			return nil, err
		}
	}

	return number.SetJPos(node.StartPos, node.EndPos), nil
}
