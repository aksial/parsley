package parsley

var EOF Item = Item{"eof", nil}

type Item struct {
	Type  string
	Value []byte
}
