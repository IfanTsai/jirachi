package common

type JContext struct {
	Name           string
	Parent         *JContext
	ParentEntryPos *JPosition
	SymbolTable    *JSymbolTable
}

func NewJContext(name string, symbolTable *JSymbolTable, parent *JContext, parentEntryPos *JPosition) *JContext {
	return &JContext{
		Name:           name,
		SymbolTable:    symbolTable,
		Parent:         parent,
		ParentEntryPos: parentEntryPos,
	}
}
