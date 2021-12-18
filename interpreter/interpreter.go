package interpreter

import (
	"fmt"
	"reflect"

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
	resValue, err := NewJInterpreter(context).Interpreter(ast)
	if err != nil {
		return nil, err
	}

	return resValue, nil
}

type JInterpreter struct {
	Context *common.JContext
}

func NewJInterpreter(context *common.JContext) *JInterpreter {
	return &JInterpreter{
		Context: context,
	}
}

func (i *JInterpreter) Interpreter(ast parser.JNode) (JValue, error) {
	return i.visit(ast)
}

func (i *JInterpreter) visit(node parser.JNode) (JValue, error) {
	switch node.Type() {
	case parser.Number:
		return i.visitNumberNode(node.(*parser.JNumberNode))
	case parser.String:
		return i.visitStringNode(node.(*parser.JStringNode))
	case parser.List:
		return i.visitListNode(node.(*parser.JListNode))
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
	case parser.ForExpr:
		return i.visitForExprNode(node.(*parser.JForExprNode))
	case parser.WhileExpr:
		return i.visitWhileExprNode(node.(*parser.JWhileExprNode))
	case parser.FuncDefExpr:
		return i.visitFunDefNode(node.(*parser.JFuncDefNode))
	case parser.CallExpr:
		return i.visitCallExprNode(node.(*parser.JCallExprNode))
	case parser.IndexExpr:
		return i.visitIndexExprNode(node.(*parser.JIndexExprNode))
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

func (i *JInterpreter) visitNumberNode(node *parser.JNumberNode) (JValue, error) {
	return NewJNumber(node.Token.Value).SetJPos(node.StartPos, node.EndPos).SetJContext(i.Context), nil
}

func (i *JInterpreter) visitStringNode(node *parser.JStringNode) (JValue, error) {
	return NewJString(node.Token.Value).SetJPos(node.StartPos, node.EndPos).SetJContext(i.Context), nil
}

func (i *JInterpreter) visitListNode(node *parser.JListNode) (JValue, error) {
	elementValues := make([]JValue, len(node.ElementNodes))
	for index := range node.ElementNodes {
		value, err := i.visit(node.ElementNodes[index])
		if err != nil {
			return nil, err
		}
		elementValues[index] = value
	}

	return NewJList(elementValues).SetJPos(node.StartPos, node.EndPos).SetJContext(i.Context), nil
}

func (i *JInterpreter) visitVarAssignNode(node *parser.JVarAssignNode) (JValue, error) {
	varName := node.Token.Value

	varValue, err := i.visit(node.Node)
	if err != nil {
		return nil, err
	}

	i.Context.SymbolTable.Set(varName, varValue)

	return varValue, nil
}

func (i *JInterpreter) visitVarAccessNode(node *parser.JVarAccessNode) (JValue, error) {
	varName := node.Token.Value
	varValue := i.Context.SymbolTable.Get(varName)
	if varValue == nil {
		return nil, errors.Wrap(&common.JRunTimeError{
			JError: &common.JError{
				StartPos: node.StartPos,
				EndPos:   node.EndPos,
			},
			Context: i.Context,
			Details: fmt.Sprintf("'%v' is not defined", varName),
		}, "failed to access variable")
	}

	return varValue.(JValue).Copy().SetJPos(node.StartPos, node.EndPos).SetJContext(i.Context), nil
}

func (i *JInterpreter) visitBinOpNode(node *parser.JBinOpNode) (JValue, error) {
	leftValue, err := i.visit(node.LeftNode)
	if err != nil {
		return nil, err
	}

	rightValue, err := i.visit(node.RightNode)
	if err != nil {
		return nil, err
	}

	var resValue JValue

	switch node.Token.Type {
	case token.PLUS:
		resValue, err = leftValue.AddTo(rightValue)
	case token.MINUS:
		resValue, err = leftValue.SubBy(rightValue)
	case token.MUL:
		resValue, err = leftValue.MulBy(rightValue)
	case token.DIV:
		resValue, err = leftValue.DivBy(rightValue)
	case token.POW:
		resValue, err = leftValue.PowBy(rightValue)
	case token.EE:
		resValue, err = leftValue.EqualTo(rightValue)
	case token.NE:
		resValue, err = leftValue.NotEqualTo(rightValue)
	case token.LT:
		resValue, err = leftValue.LessThan(rightValue)
	case token.LTE:
		resValue, err = leftValue.LessThanOrEqualTo(rightValue)
	case token.GT:
		resValue, err = leftValue.GreaterThan(rightValue)
	case token.GTE:
		resValue, err = leftValue.GreaterThanOrEqualTo(rightValue)
	case token.KEYWORD:
		switch node.Token.Value {
		case token.AND:
			resValue, err = leftValue.AndBy(rightValue)
		case token.OR:
			resValue, err = leftValue.OrBy(rightValue)
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
		return nil, errors.WithMessage(err, "failed to visit bin op node")
	}

	return resValue.SetJPos(node.StartPos, node.EndPos), nil
}

func (i *JInterpreter) visitUnaryOpNode(node *parser.JUnaryOpNode) (JValue, error) {
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
		return nil, errors.WithMessage(err, "failed to visit unary op node")
	}

	return number.SetJPos(node.StartPos, node.EndPos), nil
}

func (i *JInterpreter) visitIfExprNode(node *parser.JIfExprNode) (JValue, error) {
	for index := range node.CaseNodes {
		condition := node.CaseNodes[index][0]
		expr := node.CaseNodes[index][1]

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

	if node.ElseCaseNode != nil {
		elseValue, err := i.visit(node.ElseCaseNode)
		if err != nil {
			return nil, err
		}

		return elseValue.SetJContext(i.Context), nil
	}

	// eg. IF false THEN 123
	return NewJNumber(nil), nil
}

func (i *JInterpreter) visitWhileExprNode(node *parser.JWhileExprNode) (JValue, error) {
	var res JValue
	var resElementValues []JValue

	for {
		condition, err := i.visit(node.ConditionNode)
		if err != nil {
			return nil, err
		}

		if !condition.IsTrue() {
			break
		}

		res, err = i.visit(node.BodyNode)
		if err != nil {
			return nil, err
		}

		resElementValues = append(resElementValues, res)
	}

	if res == nil {
		return NewJNumber(nil), nil
	}

	return NewJList(resElementValues).SetJPos(node.StartPos, node.EndPos).SetJContext(i.Context), nil
}

func (i *JInterpreter) visitForExprNode(node *parser.JForExprNode) (JValue, error) {
	startNumber, err := i.visit(node.StartValueNode)
	if err != nil {
		return nil, err
	}

	endNumber, err := i.visit(node.EndValueNode)
	if err != nil {
		return nil, err
	}

	var stepNumber JValue
	if node.StepValueNode != nil {
		stepNumber, err = i.visit(node.StepValueNode)
		if err != nil {
			return nil, err
		}
	} else {
		stepNumber = NewJNumber(1)
	}

	isFloat := false
	var resElementValues []JValue

	if _, ok := startNumber.GetValue().(float64); ok {
		isFloat = true
	} else if _, ok := endNumber.GetValue().(float64); ok {
		isFloat = true
	} else if _, ok := stepNumber.GetValue().(float64); ok {
		isFloat = true
	}

	var res JValue

	if isFloat {
		var start, end, step float64
		if reflect.TypeOf(startNumber.GetValue()).Kind() == reflect.Int {
			start = float64(reflect.ValueOf(startNumber.GetValue()).Int())
		} else {
			start = reflect.ValueOf(startNumber.GetValue()).Float()
		}

		if reflect.TypeOf(endNumber.GetValue()).Kind() == reflect.Int {
			end = float64(reflect.ValueOf(endNumber.GetValue()).Int())
		} else {
			end = reflect.ValueOf(endNumber.GetValue()).Float()
		}

		if reflect.TypeOf(stepNumber.GetValue()).Kind() == reflect.Int {
			step = float64(reflect.ValueOf(stepNumber.GetValue()).Int())
		} else {
			step = reflect.ValueOf(stepNumber.GetValue()).Float()
		}

		for j := start; ; j += step {
			if (step > 0 && j >= end) || (step < 0 && j <= end) {
				break
			}

			i.Context.SymbolTable.Set(node.Token.Value, NewJNumber(j))

			res, err = i.visit(node.BodyNode)
			if err != nil {
				return nil, err
			}

			resElementValues = append(resElementValues, res)
		}
	} else {
		start := int(reflect.ValueOf(startNumber.GetValue()).Int())
		end := int(reflect.ValueOf(endNumber.GetValue()).Int())
		step := int(reflect.ValueOf(stepNumber.GetValue()).Int())

		for j := start; ; j += step {
			if (step > 0 && j >= end) || (step < 0 && j <= end) {
				break
			}

			i.Context.SymbolTable.Set(node.Token.Value, NewJNumber(j))

			res, err = i.visit(node.BodyNode)
			if err != nil {
				return nil, err
			}

			resElementValues = append(resElementValues, res)
		}
	}

	if res == nil {
		return NewJNumber(nil), nil
	}

	return NewJList(resElementValues).SetJPos(node.StartPos, node.EndPos).SetJContext(i.Context), nil
}

func (i *JInterpreter) visitFunDefNode(node *parser.JFuncDefNode) (JValue, error) {
	argNames := make([]string, len(node.ArgTokens))
	var ok bool
	for index := range node.ArgTokens {
		if argNames[index], ok = node.ArgTokens[index].Value.(string); !ok {
			return nil, errors.Wrap(&common.JRunTimeError{
				JError: &common.JError{
					StartPos: node.StartPos,
					EndPos:   node.EndPos,
				},
				Context: i.Context,
				Details: "arg token value is not string",
			}, "failed to visit function definition node")
		}
	}

	var funcName interface{}
	if node.Token != nil {
		funcName = node.Token.Value
	}

	functionValue := NewJFunction(funcName, argNames, node.BodyNode).
		SetJPos(node.StartPos, node.EndPos).
		SetJContext(i.Context)

	if funcName != nil {
		i.Context.SymbolTable.Set(funcName, functionValue)
	}

	return functionValue, nil
}

func (i *JInterpreter) visitCallExprNode(node *parser.JCallExprNode) (JValue, error) {
	callValue, err := i.visit(node.CallNode)
	if err != nil {
		return nil, err
	}

	argValues := make([]JValue, len(node.ArgNodes))
	for index := range node.ArgNodes {
		argValue, err := i.visit(node.ArgNodes[index])
		if err != nil {
			return nil, err
		}
		argValues[index] = argValue
	}

	returnValue, err := callValue.Execute(argValues)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to visit call expression node")
	}

	return returnValue, nil
}

func (i *JInterpreter) visitIndexExprNode(node *parser.JIndexExprNode) (JValue, error) {
	indexNodeValue, err := i.visit(node.IndexNode)
	if err != nil {
		return nil, err
	}

	indexExprValue, err := i.visit(node.IndexExpr)
	if err != nil {
		return nil, err
	}

	resValue, err := indexNodeValue.Index(indexExprValue)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to visit index expression node")
	}

	return resValue, nil
}
