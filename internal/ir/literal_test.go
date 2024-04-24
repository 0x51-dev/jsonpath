package ir

import (
	"fmt"
	"github.com/0x51-dev/jsonpath/internal/grammar"
	"github.com/0x51-dev/upeg/parser"
	"testing"
)

func TestParseLiteral(t *testing.T) {
	t.Run("Number", func(t *testing.T) {
		for _, test := range []struct {
			input          string
			representation string
		}{
			{"-0", "0"},
			{"0", "0"}, // int
			{"-1", "-1"},
			{"-10", "-10"},
			{"99", "99"},
			{"0.0", "0"}, // frac
			{"-1.01", "-1.01"},
			{"1.33", "1.33"},
			{"-0.1", "-0.1"},
			{"0.0e1", "0"}, // exp
			{"0e0", "0"},
			{"0e-1", "0"},
			{"0e+1", "0"},
			{"1e1", "10"},
			{"1e-1", "0.1"},
			{"1e+1", "10"},
			{"1.1e1", "11"},
		} {
			p, err := parser.New([]rune(test.input))
			if err != nil {
				t.Fatal(err)
			}
			n, err := p.ParseEOF(grammar.Literal)
			if err != nil {
				t.Fatal(err)
			}
			lit, err := ParseLiteral(n)
			if err != nil {
				t.Fatal(err)
			}
			v, err := lit.Value(nil)
			if err != nil {
				t.Fatal(err)
			}
			if fmt.Sprintf("%v", v) != test.representation {
				t.Fatalf("expected %q, got %q", test.representation, v)
			}
		}
	})
	t.Run("StringLiteral", func(t *testing.T) {
		for _, test := range []struct {
			input          string
			representation string
		}{
			{`"hello"`, "hello"},
			{`'hello'`, "hello"},
			{`'he"llo'`, `he"llo`},
			{`"he'llo'"`, `he'llo'`},
			{`"\/\\\b\f\n\r\t\uF09F\u8DBA"`, "/\\\b\f\n\r\t🍺"},
			{`"he\"llo"`, `he"llo`},
			{`'he\'llo'`, `he'llo`},
		} {
			p, err := parser.New([]rune(test.input))
			if err != nil {
				t.Fatal(err)
			}
			n, err := p.ParseEOF(grammar.Literal)
			if err != nil {
				t.Fatal(err)
			}
			lit, err := ParseLiteral(n)
			if err != nil {
				t.Fatal(err)
			}
			v, err := lit.Value(nil)
			if err != nil {
				t.Fatal(err)
			}
			if v != test.representation {
				t.Fatalf("expected %q, got %q", test.representation, v)
			}
		}
	})
	t.Run("Boolean", func(t *testing.T) {
		for _, boolean := range []string{"true", "false"} {
			p, err := parser.New([]rune(boolean))
			if err != nil {
				t.Fatal(err)
			}
			n, err := p.ParseEOF(grammar.Literal)
			if err != nil {
				t.Fatal(err)
			}
			lit, err := ParseLiteral(n)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := lit.Value(nil); err != nil {
				t.Fatal(err)
			}
		}
	})
	t.Run("Null", func(t *testing.T) {
		p, err := parser.New([]rune("null"))
		if err != nil {
			t.Fatal(err)
		}
		n, err := p.ParseEOF(grammar.Literal)
		if err != nil {
			t.Fatal(err)
		}
		lit, err := ParseLiteral(n)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := lit.Value(nil); err != nil {
			t.Fatal(err)
		}
	})
}
