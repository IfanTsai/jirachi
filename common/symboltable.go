package common

import "github.com/IfanTsai/jirachi/pkg/safemap"

type JSymbolTable struct {
	Symbols *safemap.SafeMap[any]
	Parent  *JSymbolTable
}

func NewJSymbolTable(parent *JSymbolTable) *JSymbolTable {
	return &JSymbolTable{
		Symbols: safemap.NewSafeMap[any](),
		Parent:  parent,
	}
}

func (st *JSymbolTable) Get(name any) any {
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

func (st *JSymbolTable) Set(name, value any) *JSymbolTable {
	st.Symbols.Set(name, value)

	return st
}

func (st *JSymbolTable) Remove(name any) *JSymbolTable {
	st.Symbols.Del(name)

	return st
}
