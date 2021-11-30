package interpreter

import (
	"github.com/IfanTsai/jirachi/common"
	"github.com/IfanTsai/jirachi/lexer"
	"github.com/IfanTsai/jirachi/parser"
	"github.com/IfanTsai/jirachi/token"
	"github.com/pkg/errors"
)

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
	number, err := NewJInterpreter().Visit(ast)
	if err != nil {
		return nil, err
	}

	return number.Value, nil
}

type JInterpreter struct {
}

func NewJInterpreter() *JInterpreter {
	return &JInterpreter{}
}

func (i *JInterpreter) Visit(node *parser.JNode) (*JNumber, error) {
	switch node.Type {
	case parser.Number:
		return i.visitNumberNode(node)
	case parser.BinOp:
		return i.visitBinOpNode(node)
	case parser.UnaryOp:
		return i.visitUnaryOpNode(node)
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
	return NewJNumber(node.Token.Value).SetPos(node.StartPos, node.EndPos), nil
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

	return resNumber.SetPos(node.StartPos, node.EndPos), nil
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

	return number.SetPos(node.StartPos, node.EndPos), nil
}
