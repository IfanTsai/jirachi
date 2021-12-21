package common

import "github.com/IfanTsai/jirachi/pkg/safemap"

type JSymbolTable struct {
	Symbols *safemap.SafeMap
	Parent  *JSymbolTable
}

func NewJSymbolTable(parent *JSymbolTable) *JSymbolTable {
	return &JSymbolTable{
		Symbols: safemap.NewSafeMap(),
		Parent:  parent,
	}
}

func (st *JSymbolTable) Get(name interface{}) interface{} {
	symbolTable := st

	for symbolTable != nil {
		if value, ok := symbolTable.Symbols.Get(name); ok {
			return value
		} else {
			symbolTable = symbolTable.Parent
		}
	}

	return nil
}

func (st *JSymbolTable) Set(name, value interface{}) *JSymbolTable {
	st.Symbols.Set(name, value)

	return st
}

func (st *JSymbolTable) Remove(name interface{}) *JSymbolTable {
	st.Symbols.Del(name)

	return st
}
