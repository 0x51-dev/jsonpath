package ir

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/0x51-dev/upeg/parser"
	"strings"
)

// Boolean = true / false
// true    = %x74.72.75.65
// false   = %x66.61.6c.73.65
type Boolean bool

func (s Boolean) String() string {
	if s {
		return "true"
	}
	return "false"
}

func (s Boolean) Value(_ any) (any, error) {
	return s, nil
}

func (s Boolean) argument() {}

func (s Boolean) comparable() {}

func (s Boolean) literal() {}

type Literal interface {
	// literal = number / string-literal / true / false / null
	literal()

	// Comparable = literal / ...
	Comparable

	// FunctionArgument = literal / ...
	FunctionArgument
}

// ParseLiteral parses a literal node, which is the result of parsing grammar.Literal.
func ParseLiteral(n *parser.Node) (Literal, error) {
	name := "Literal"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	if len(n.Children()) != 1 {
		return nil, NewInvalidNodeStructureError(name, n)
	}

	switch n := n.Children()[0]; n.Name {
	case "Number":
		var number string
		for _, n := range n.Children() {
			number += n.Value()
		}
		lit := Number(number)
		return &lit, nil
	case "StringLiteral":
		return ParseStringLiteral(n)
	case "True":
		b := Boolean(true)
		return &b, nil
	case "False":
		b := Boolean(false)
		return &b, nil
	case "Null":
		return new(Null), nil
	default:
		return nil, NewInvalidNodeStructureError(name, n)
	}
}

// Null = %x6e.75.6c.6c
type Null struct{}

func (s Null) String() string {
	return "null"
}

func (s Null) Value(_ any) (any, error) {
	return nil, nil
}

func (s Null) argument() {}

func (s Null) comparable() {}

func (s Null) literal() {}

// Number = (int / "-0") [ frac ] [ exp ] ; decimal number
// int    = "0" / (["-"] DIGIT1 *DIGIT)      ; - optional
// DIGIT  = %x30-39              ; 0-9
// DIGIT1 = %x31-39                    ; 1-9 non-zero digit
// frac   = "." 1*DIGIT                  ; decimal fraction
// exp    = "e" [ "-" / "+" ] 1*DIGIT    ; decimal exponent
type Number json.Number

func (s Number) String() string {
	return string(s)
}

func (s Number) Value(_ any) (any, error) {
	jn := json.Number(s)
	if strings.Contains(string(jn), ".") || strings.Contains(string(jn), "e") {
		return jn.Float64()
	}
	return jn.Int64()
}

func (s Number) argument() {}

func (s Number) comparable() {}

func (s Number) literal() {}

// String         = %x22 *double-quoted %x22 / %x27 *single-quoted %x27       ; 'string'
// double-quoted  = unescaped / %x27 / ESC %x22 / ESC escapable
// single-quoted  = unescaped / %x22 / ESC %x27 / ESC escapable
// ESC            = %x5C
// unescaped      = %x20-21 / %x23-26 / %x28-5B / %x5D-D7FF / %xE000-10FFFF
// escapable      = %x62 / %x66 / %x6E / %x72 / %x74 / "/" / "\" / (%x75 hexchar)
// hexchar        = non-surrogate / (high-surrogate "\" %x75 low-surrogate)
// non-surrogate  = ((DIGIT / "A"/"B"/"C" / "E"/"F") 3HEXDIG) / ("D" %x30-37 2HEXDIG )
// high-surrogate = "D" ("8"/"9"/"A"/"B") 2HEXDIG
// low-surrogate  = "D" ("C"/"D"/"E"/"F") 2HEXDIG
// HEXDIG         = DIGIT / "A" / "B" / "C" / "D" / "E" / "F"
type String string

func ParseStringLiteral(n *parser.Node) (*String, error) {
	name := "StringLiteral"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}

	var str string
	for _, n := range n.Children() {
		switch n.Name {
		case "UnescapedDq", "UnescapedSq":
			str += n.Value()
		case "EscapableDq", "EscapableSq":
			if v := n.Value(); v != "" {
				switch v {
				case "\"", "'", "/", "\\":
					str += v
				case "b":
					str += "\b"
				case "f":
					str += "\f"
				case "n":
					str += "\n"
				case "r":
					str += "\r"
				case "t":
					str += "\t"
				default:
					return nil, NewInvalidNodeStructureError(name, n)
				}
			} else {
				if len(n.Children()) != 1 {
					return nil, NewInvalidNodeStructureError(name, n)
				}
				n := n.Children()[0]
				if n.Name != "Hexchar" {
					return nil, NewInvalidNodeStructureError(name, n)
				}
				raw, err := hex.DecodeString(n.Value())
				if err != nil {
					return nil, NewInvalidNodeStructureError(name, n)
				}
				str += string(raw)
			}
		default:
			return nil, NewInvalidNodeStructureError(name, n)
		}
	}
	s := String(str)
	return &s, nil
}

func (s String) String() string {
	return fmt.Sprintf("'%s'", string(s))
}

func (s String) Value(_ any) (any, error) {
	return string(s), nil
}

func (s String) argument() {}

func (s String) comparable() {}

func (s String) literal() {}
