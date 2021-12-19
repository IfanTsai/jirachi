package builtin

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"

	"github.com/IfanTsai/jirachi/common"

	"github.com/IfanTsai/jirachi/interpreter/object"
)

var (
	NULL        = object.NewJNumber(0)
	TRUE        = object.NewJNumber(1)
	FALSE       = object.NewJNumber(0)
	Len         = object.NewJBuiltInFunction("len", []string{"value"}, executeLen)
	Type        = object.NewJBuiltInFunction("type", []string{"value"}, executeType)
	Print       = object.NewJBuiltInFunction("print", []string{"value"}, executePrint)
	Input       = object.NewJBuiltInFunction("input", []string{}, executeInput)
	InputNumber = object.NewJBuiltInFunction("input_number", []string{}, executeInputNumber)
	IsNumber    = object.NewJBuiltInFunction("is_number", []string{"value"}, executeIsNumber)
	IsString    = object.NewJBuiltInFunction("is_string", []string{"value"}, executeIsString)
	IsList      = object.NewJBuiltInFunction("is_list", []string{"value"}, executeIsList)
	IsFunction  = object.NewJBuiltInFunction("is_function", []string{"value"}, executeIsFunction)
)

func executeLen(function *object.JBuiltInFunction, args []object.JValue) (object.JValue, error) {
	arg := args[0]
	switch argValue := arg.(type) {
	case *object.JList:
		return object.NewJNumber(len(argValue.ElementValues)), nil
	case *object.JString:
		return object.NewJNumber(len(argValue.Value.(string))), nil
	}

	return nil, errors.Wrap(&common.JRunTimeError{
		JError: &common.JError{
			StartPos: function.StartPos,
			EndPos:   function.EndPos,
		},
		Context: function.GetContext(),
		Details: "First argument must be list or string",
	}, "failed to call len")
}

func executeType(function *object.JBuiltInFunction, args []object.JValue) (object.JValue, error) {
	var argType string
	arg := args[0]

	switch arg.(type) {
	case *object.JNumber:
		argType = object.Number
	case *object.JString:
		argType = object.String
	case *object.JList:
		argType = object.List
	case *object.JFunction:
		argType = object.Function
	case *object.JBuiltInFunction:
		argType = object.BuiltInFunction
	}

	return object.NewJString(argType), nil
}

func executeIsNumber(function *object.JBuiltInFunction, args []object.JValue) (object.JValue, error) {
	if _, ok := args[0].(*object.JNumber); !ok {
		return FALSE, nil
	}

	return TRUE, nil
}

func executeIsString(function *object.JBuiltInFunction, args []object.JValue) (object.JValue, error) {
	if _, ok := args[0].(*object.JString); !ok {
		return FALSE, nil
	}

	return TRUE, nil
}

func executeIsList(function *object.JBuiltInFunction, args []object.JValue) (object.JValue, error) {
	if _, ok := args[0].(*object.JList); !ok {
		return FALSE, nil
	}

	return TRUE, nil
}

func executeIsFunction(function *object.JBuiltInFunction, args []object.JValue) (object.JValue, error) {
	if _, ok := args[0].(*object.JFunction); !ok {
		return FALSE, nil
	}

	return TRUE, nil
}

func executePrint(function *object.JBuiltInFunction, args []object.JValue) (object.JValue, error) {
	fmt.Println(args[0])

	return nil, nil
}

func executeInput(function *object.JBuiltInFunction, args []object.JValue) (object.JValue, error) {
	textBytes, _, _ := bufio.NewReader(os.Stdin).ReadLine()

	return object.NewJString(string(textBytes)), nil
}

func executeInputNumber(function *object.JBuiltInFunction, args []object.JValue) (object.JValue, error) {
	textBytes, _, _ := bufio.NewReader(os.Stdin).ReadLine()
	text := string(textBytes)

	var number interface{}
	var err error

	number, err = strconv.Atoi(text)
	fmt.Println("=========> " + text)
	if err != nil {
		number, err = strconv.ParseFloat(text, 64)
		if err != nil {
			return nil, errors.Wrap(&common.JRunTimeError{
				JError: &common.JError{
					StartPos: function.StartPos,
					EndPos:   function.EndPos,
				},
				Context: function.GetContext(),
				Details: text + " must be an number",
			}, "failed to call input_number")
		}
	}

	return object.NewJNumber(number), nil
}
