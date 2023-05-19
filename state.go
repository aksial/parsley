package parsley

type StateFunc func(*Lexer) StateFunc
