package lexer

type JPosition struct {
	Index    int
	Col      int
	Ln       int64
	Filename string
}

func NewJPosition(index, col int, ln int64, filename string) *JPosition {
	return &JPosition{
		Index:    index,
		Col:      col,
		Ln:       ln,
		Filename: filename,
	}
}

func (p *JPosition) Advance(text []byte) {
	p.Index++
	p.Col++

	if text[p.Index] == '\n' {
		p.Ln++
		p.Col = 0
	}
}
