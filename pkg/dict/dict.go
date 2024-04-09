package dict

type Dict struct {
	Words map[string]string
}

func NewDict() *Dict {
	return &Dict{
		Words: make(map[string]string),
	}
}

func (d *Dict) Define(word, definition string) bool {
	d.Words[word] = definition
	return true
}
