package interpreter_test

import (
	"testing"

	"github.com/IfanTsai/jirachi/interpreter"

	"github.com/IfanTsai/jirachi/common"
	"github.com/pkg/errors"

	"github.com/IfanTsai/jirachi/interpreter/object"
	"github.com/stretchr/testify/require"
)

func TestExecuteLen(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		args        []object.JValue
		checkResult func(t *testing.T, resValue object.JValue, err error)
	}{
		{
			name: "len string",
			args: []object.JValue{object.NewJString("hello world")},
			checkResult: func(t *testing.T, resValue object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, object.NewJNumber(nil), resValue)
				require.IsType(t, 0, resValue.GetValue())
				require.Equal(t, 11, resValue.GetValue())
			},
		},
		{
			name: "len list",
			args: []object.JValue{object.NewJList([]object.JValue{
				object.NewJString("hello"),
				object.NewJString("world"),
				object.NewJNumber(1),
				object.NewJNumber(1.2),
				object.NewJNumber(object.NewJList([]object.JValue{})),
			})},
			checkResult: func(t *testing.T, resValue object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, object.NewJNumber(nil), resValue)
				require.IsType(t, 0, resValue.GetValue())
				require.Equal(t, 5, resValue.GetValue())
			},
		},
		{
			name: "wrong variable type",
			args: []object.JValue{object.NewJNumber(0)},
			checkResult: func(t *testing.T, resValue object.JValue, err error) {
				t.Helper()
				require.Error(t, err)
				require.IsType(t, &common.JRunTimeError{}, errors.Cause(err))
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			resValue, err := interpreter.ExecuteLen(interpreter.Len, testCase.args)
			testCase.checkResult(t, resValue, err)
		})
	}
}

func TestExecuteType(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		args        []object.JValue
		checkResult func(t *testing.T, resValue object.JValue, err error)
	}{
		{
			name: "type number",
			args: []object.JValue{object.NewJNumber(0)},
			checkResult: func(t *testing.T, resValue object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, object.NewJString(""), resValue)
				require.Equal(t, object.Number, resValue.String())
			},
		},
		{
			name: "type string",
			args: []object.JValue{object.NewJString("hello")},
			checkResult: func(t *testing.T, resValue object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, object.NewJString(""), resValue)
				require.Equal(t, object.String, resValue.String())
			},
		},
		{
			name: "type list",
			args: []object.JValue{object.NewJList([]object.JValue{})},
			checkResult: func(t *testing.T, resValue object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, object.NewJString(""), resValue)
				require.Equal(t, object.List, resValue.String())
			},
		},
		{
			name: "type function",
			args: []object.JValue{object.NewJFunction("test_func", []string{}, nil)},
			checkResult: func(t *testing.T, resValue object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, object.NewJString(""), resValue)
				require.Equal(t, object.Function, resValue.String())
			},
		},
		{
			name: "type built-in function",
			args: []object.JValue{interpreter.Type},
			checkResult: func(t *testing.T, resValue object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, object.NewJString(""), resValue)
				require.Equal(t, object.BuiltInFunction, resValue.String())
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			resValue, err := interpreter.ExecuteType(interpreter.Type, testCase.args)
			testCase.checkResult(t, resValue, err)
		})
	}
}
