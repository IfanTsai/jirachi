package common

type JPosition struct {
	Index    int
	Col      int
	Ln       int64
	Filename string
	Text     string
}

func NewJPosition(index, col int, ln int64, filename, text string) *JPosition {
	return &JPosition{
		Index:    index,
		Col:      col,
		Ln:       ln,
		Filename: filename,
		Text:     text,
	}
}

func (p *JPosition) Advance(text []byte) *JPosition {
	p.Col++

	if text != nil && p.Index >= 0 && p.Index < len(text) && text[p.Index] == '\n' {
		p.Ln++
		p.Col = 0
	}

	p.Index++

	return p
}

func (p *JPosition) Back(text []byte) *JPosition {
	p.Index--
	p.Col--

	if text != nil && p.Index >= 0 && text[p.Index] == '\n' {
		p.Ln--
		p.Col = len(text) - 1
	}

	return p
}

func (p *JPosition) Copy() *JPosition {
	return &JPosition{
		Index:    p.Index,
		Col:      p.Col,
		Ln:       p.Ln,
		Filename: p.Filename,
		Text:     p.Text,
	}
}
