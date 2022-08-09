package interpreter

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"

	"github.com/pkg/errors"

	"github.com/IfanTsai/jirachi/common"

	"github.com/IfanTsai/jirachi/interpreter/object"
)

var (
	NULL        = object.NewJNull()
	TRUE        = object.NewJNumber(1)
	FALSE       = object.NewJNumber(0)
	Len         = object.NewJBuiltInFunction("len", []string{"value"}, ExecuteLen)
	Type        = object.NewJBuiltInFunction("type", []string{"value"}, ExecuteType)
	Print       = object.NewJBuiltInFunction("print", []string{"value"}, ExecutePrint)
	Println     = object.NewJBuiltInFunction("println", []string{"value"}, ExecutePrintln)
	Input       = object.NewJBuiltInFunction("input", []string{}, ExecuteInput)
	InputNumber = object.NewJBuiltInFunction("input_number", []string{}, ExecuteInputNumber)
	IsNumber    = object.NewJBuiltInFunction("is_number", []string{"value"}, ExecuteIsNumber)
	IsString    = object.NewJBuiltInFunction("is_string", []string{"value"}, ExecuteIsString)
	IsList      = object.NewJBuiltInFunction("is_list", []string{"value"}, ExecuteIsList)
	IsFunction  = object.NewJBuiltInFunction("is_function", []string{"value"}, ExecuteIsFunction)
	RunShell    = object.NewJBuiltInFunction("run_shell", []string{"text"}, ExecuteRunShell)
	RunScript   = object.NewJBuiltInFunction("run", []string{"filename"}, ExecuteRun)
)

func ExecuteLen(function *object.JBuiltInFunction, args []object.JValue) (object.JValue, error) {
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

func ExecuteType(function *object.JBuiltInFunction, args []object.JValue) (object.JValue, error) {
	arg := args[0]

	return object.NewJString(object.GetJValueType(arg)), nil
}

func ExecuteIsNumber(function *object.JBuiltInFunction, args []object.JValue) (object.JValue, error) {
	if _, ok := args[0].(*object.JNumber); !ok {
		return FALSE, nil
	}

	return TRUE, nil
}

func ExecuteIsString(function *object.JBuiltInFunction, args []object.JValue) (object.JValue, error) {
	if _, ok := args[0].(*object.JString); !ok {
		return FALSE, nil
	}

	return TRUE, nil
}

func ExecuteIsList(function *object.JBuiltInFunction, args []object.JValue) (object.JValue, error) {
	if _, ok := args[0].(*object.JList); !ok {
		return FALSE, nil
	}

	return TRUE, nil
}

func ExecuteIsFunction(function *object.JBuiltInFunction, args []object.JValue) (object.JValue, error) {
	if _, ok := args[0].(*object.JFunction); !ok {
		return FALSE, nil
	}

	return TRUE, nil
}

func ExecutePrint(function *object.JBuiltInFunction, args []object.JValue) (object.JValue, error) {
	fmt.Print(args[0])

	return nil, nil
}

func ExecutePrintln(function *object.JBuiltInFunction, args []object.JValue) (object.JValue, error) {
	fmt.Println(args[0])

	return nil, nil
}

func ExecuteInput(function *object.JBuiltInFunction, args []object.JValue) (object.JValue, error) {
	textBytes, _, _ := bufio.NewReader(os.Stdin).ReadLine()

	return object.NewJString(string(textBytes)), nil
}

func ExecuteInputNumber(function *object.JBuiltInFunction, args []object.JValue) (object.JValue, error) {
	textBytes, _, _ := bufio.NewReader(os.Stdin).ReadLine()
	text := string(textBytes)

	var (
		number interface{}
		err    error
	)

	number, err = strconv.Atoi(text)
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

func ExecuteRun(function *object.JBuiltInFunction, args []object.JValue) (object.JValue, error) {
	filename, ok := args[0].GetValue().(string)
	if !ok {
		return nil, errors.Wrap(&common.JRunTimeError{
			JError: &common.JError{
				StartPos: function.StartPos,
				EndPos:   function.EndPos,
			},
			Context: function.GetContext(),
			Details: "First arguments must be a string",
		}, "failed to call run")
	}

	scriptFile, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(&common.JRunTimeError{
			JError: &common.JError{
				StartPos: function.StartPos,
				EndPos:   function.EndPos,
			},
			Context: function.GetContext(),
			Details: "Failed to load script " + filename + ", error: " + err.Error(),
		}, "failed to call run")
	}

	bytes, err := io.ReadAll(scriptFile)
	if err != nil {
		return nil, errors.Wrap(&common.JRunTimeError{
			JError: &common.JError{
				StartPos: function.StartPos,
				EndPos:   function.EndPos,
			},
			Context: function.GetContext(),
			Details: "Failed to load script " + filename + ", error: " + err.Error(),
		}, "failed to call run")
	}

	if _, err = Run(filename, string(bytes)); err != nil {
		return nil, errors.Wrap(&common.JRunTimeError{
			JError: &common.JError{
				StartPos: function.StartPos,
				EndPos:   function.EndPos,
			},
			Context: function.GetContext(),
			Details: "Failed to finish executing script " + filename + "\n" + err.Error(),
		}, "failed to call run")
	}

	return nil, nil
}

func ExecuteRunShell(function *object.JBuiltInFunction, args []object.JValue) (object.JValue, error) {
	shellStr := args[0].GetValue().(string)
	cmd := exec.Command("/bin/sh", "-c", shellStr)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	if err := cmd.Run(); err != nil || errBuf.Len() > 0 {
		return nil, errors.Wrap(&common.JRunTimeError{
			JError: &common.JError{
				StartPos: function.StartPos,
				EndPos:   function.EndPos,
			},
			Context: function.GetContext(),
		}, errBuf.String())
	}

	return object.NewJString(outBuf.String()), nil
}
