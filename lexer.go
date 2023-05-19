package parsley

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

type Lexer struct {
	state StateFunc
	input []byte
	items chan Item
	start int
	pos   int
	width int
}

func Lex(s StateFunc, input []byte) *Lexer {
	l := new(Lexer)
	l.state = s
	l.input = input
	l.items = make(chan Item, 2)
	return l
}

func (l *Lexer) Emit(itemType string) {
	l.items <- Item{itemType, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *Lexer) AcceptFunc(f func(rune) bool) (bool, error) {
	r, err := l.Next()
	if f(r) {
		return true, err
	}
	l.Backup()
	return false, nil
}

func (l *Lexer) Accept(valid string) (bool, error) {
	r, err := l.Next()
	if strings.ContainsRune(valid, r) {
		return true, err
	}
	l.Backup()
	return false, err
}

func (l *Lexer) AcceptRunFunc(f func(rune) bool) (n int, err error) {
	for err == nil {
		r, newErr := l.Next()
		if f(r) {
			err = newErr
			n++
		} else {
			break
		}
	}
	l.Backup()
	return
}

func (l *Lexer) AcceptRun(valid string) (n int, err error) {
	for err == nil {
		r, newErr := l.Next()
		if strings.ContainsRune(valid, r) {
			err = newErr
			n++
		} else {
			break
		}
	}
	l.Backup()
	return
}

func (l *Lexer) Next() (r rune, err error) {
	if l.pos >= len(l.input)-1 {
		err = io.EOF
	}
	r, l.width = utf8.DecodeRune(l.input[l.pos:])
	l.pos += l.width
	return
}

func (l *Lexer) NextItem() (item Item) {
	for l.state != nil {
		select {
		case item = <-l.items:
			return
		default:
			l.state = l.state(l)
		}
	}
	return EOF
}

func (l *Lexer) Run() []byte {
	return l.input[l.start:l.pos]
}

func (l *Lexer) Ignore() {
	l.start = l.pos
}

func (l *Lexer) Backup() {
	l.pos -= l.width
}

func (l *Lexer) Peek() (r rune, err error) {
	r, err = l.Next()
	l.Backup()
	return
}

func (l *Lexer) EOF() bool {
	return l.pos >= len(l.input)-1
}

func (l *Lexer) Errorf(next StateFunc, format string, args ...any) StateFunc {
	l.items <- Item{
		"error",
		[]byte(fmt.Sprintf(format, args...)),
	}
	return next
}

func (l *Lexer) Fatalf(format string, args ...any) StateFunc {
	return l.Errorf(nil, format, args...)
}
