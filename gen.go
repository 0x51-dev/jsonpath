package jsonpath

//go:generate go run github.com/0x51-dev/upeg/cmd/abnf --in=internal/grammar/grammar.abnf --out=internal/grammar/grammar.go --ignore=name-first,name-char,DIGIT,DIGIT1,ALPHA,B,S,LCALPHA,function-name-first,function-name-char,unescaped,logical-or-expr,filter-query,singular-query --package=grammar
