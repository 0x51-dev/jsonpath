package ir

import (
	"fmt"
	"github.com/0x51-dev/upeg/parser"
	"strconv"
	"strings"
)

func parseInt(n *parser.Node) (int, error) {
	if n.Name != "Int" {
		return 0, NewInvalidNodeStructureError("Int", n)
	}
	return strconv.Atoi(n.Value())
}

type AbsSingularQuery struct {
	Segments []SingularQuerySegment
}

func ParseAbsSingularQuery(n *parser.Node) (*AbsSingularQuery, error) {
	name := "AbsSingularQuery"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	cs := n.Children()
	if len(cs) != 2 {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	if cs[0].Name != "RootIdentifier" {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	segments, err := ParseSingularQuerySegments(cs[1])
	if err != nil {
		return nil, err
	}
	return &AbsSingularQuery{
		Segments: segments,
	}, nil
}

func (s AbsSingularQuery) String() string {
	str := "$"
	for _, s := range s.Segments {
		str += s.String()
	}
	return str
}

func (s AbsSingularQuery) Value(ref any) (any, error) {
	current := ref
	for _, segment := range s.Segments {
		c, err := segment.Value(current)
		if err != nil {
			return nil, err
		}
		current = c
	}
	return current, nil
}

func (s AbsSingularQuery) comparable() {}

type BasicExpr interface {
	fmt.Stringer

	basicExpr()
}

func ParseBasicExpr(n *parser.Node) (BasicExpr, error) {
	name := "BasicExpr"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	switch n := n.Children()[0]; n.Name {
	case "ParenExpr":
		return ParseParenExpr(n)
	case "ComparisonExpr":
		return ParseComparisonExpr(n)
	case "TestExpr":
		return ParseTestExpr(n)
	default:
		return nil, NewInvalidNodeStructureError(name, n)
	}
}

type BracketedSelection struct {
	Selectors []Selector
}

func ParseBracketedSelection(n *parser.Node) (*BracketedSelection, error) {
	name := "BracketedSelection"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	var selectors []Selector
	for _, n := range n.Children() {
		selector, err := ParseSelector(n)
		if err != nil {
			return nil, err
		}
		selectors = append(selectors, selector)
	}
	return &BracketedSelection{
		Selectors: selectors,
	}, nil
}

func (s BracketedSelection) String() string {
	var str []string
	for _, s := range s.Selectors {
		str = append(str, s.String())
	}
	return "[" + strings.Join(str, ", ") + "]"
}

func (s BracketedSelection) childSegment() {}

func (s BracketedSelection) segment() {}

type ChildSegment interface {
	Segment

	childSegment()
}

func ParseChildSegment(n *parser.Node) (ChildSegment, error) {
	name := "ChildSegment"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	switch n := n.Children()[0]; n.Name {
	case "BracketedSelection":
		s, err := ParseBracketedSelection(n)
		if err != nil {
			return nil, err
		}
		return s, nil
	case "WildcardSelector":
		return new(WildcardSelector), nil
	case "MemberNameShorthand":
		return &MemberNameShorthand{
			Name: n.Value(),
		}, nil
	default:
		return nil, NewInvalidNodeStructureError(name, n)
	}
}

type Comparable interface {
	fmt.Stringer

	Value(ref any) (any, error)
	comparable()
}

func ParseComparable(n *parser.Node) (Comparable, error) {
	name := "Comparable"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	switch n := n.Children()[0]; n.Name {
	case "Literal":
		return ParseLiteral(n)
	case "RelSingularQuery":
		return ParseRelSingularQuery(n)
	case "AbsSingularQuery":
		return ParseAbsSingularQuery(n)
	case "FunctionExpr":
		return ParseFunctionExpr(n)
	default:
		return nil, NewInvalidNodeStructureError(name, n)
	}
}

type ComparisonExpr struct {
	Left, Right Comparable
	Op          string
}

func ParseComparisonExpr(n *parser.Node) (*ComparisonExpr, error) {
	name := "ComparisonExpr"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	cs := n.Children()
	if len(cs) != 3 {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	left, err := ParseComparable(cs[0])
	if err != nil {
		return nil, err
	}
	if typ := typeOfArgument(left); typ != valueType && typ != unknownType {
		return nil, fmt.Errorf("invalid arg type %s for %s", left, name)
	}
	op := cs[1].Value()
	right, err := ParseComparable(cs[2])
	if err != nil {
		return nil, err
	}
	if typ := typeOfArgument(left); typ != valueType && typ != unknownType {
		return nil, fmt.Errorf("invalid arg type %s for %s", left, name)
	}
	return &ComparisonExpr{
		Left:  left,
		Op:    op,
		Right: right,
	}, nil
}

func (s ComparisonExpr) String() string {
	return fmt.Sprintf("%s %s %s", s.Left, s.Op, s.Right)
}

func (s ComparisonExpr) basicExpr() {}

type DescendantSegment struct {
	Segment Segment
}

func ParseDescendantSegment(n *parser.Node) (*DescendantSegment, error) {
	name := "DescendantSegment"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	switch n := n.Children()[0]; n.Name {
	case "BracketedSelection":
		segment, err := ParseBracketedSelection(n)
		if err != nil {
			return nil, err
		}
		return &DescendantSegment{
			Segment: segment,
		}, nil
	case "WildcardSelector":
		return &DescendantSegment{
			Segment: new(WildcardSelector),
		}, nil
	case "MemberNameShorthand":
		return &DescendantSegment{
			Segment: &MemberNameShorthand{
				Name: n.Value(),
			},
		}, nil
	default:
		return nil, NewInvalidNodeStructureError(name, n)
	}
}

func (s DescendantSegment) String() string {
	switch s := s.Segment.(type) {
	case *MemberNameShorthand:
		return fmt.Sprintf("..%s", s.Name)
	}
	return fmt.Sprintf("..%s", s.Segment)
}

func (s DescendantSegment) segment() {}

type FilterSelector struct {
	LogicalExpr *LogicalExpr
}

func ParseFilterSelector(n *parser.Node) (*FilterSelector, error) {
	name := "FilterSelector"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	logicalExpr, err := ParseLogicalExpr(n.Children()[0])
	if err != nil {
		return nil, err
	}
	return &FilterSelector{
		LogicalExpr: logicalExpr,
	}, nil
}

func (s FilterSelector) String() string {
	return fmt.Sprintf("?%s", s.LogicalExpr.String())
}

func (s FilterSelector) selector() {}

type FunctionArgument interface {
	fmt.Stringer

	argument()
}

func ParseFunctionArgument(n *parser.Node) (FunctionArgument, error) {
	name := "FunctionArgument"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	switch n := n.Children()[0]; n.Name {
	case "Literal":
		return ParseLiteral(n)
	case "RelQuery":
		return ParseRelQuery(n)
	case "JsonpathQuery":
		return ParseJSONPathQuery(n)
	case "LogicalExpr":
		return ParseLogicalExpr(n)
	case "FunctionExpr":
		return ParseFunctionExpr(n)
	default:
		return nil, NewInvalidNodeStructureError(name, n)
	}
}

type FunctionExpr struct {
	Name      string
	Arguments []FunctionArgument
}

func ParseFunctionExpr(n *parser.Node) (*FunctionExpr, error) {
	name := "FunctionExpr"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	cs := n.Children()
	functionName := cs[0].Value()
	var args []FunctionArgument
	for _, n := range cs[1:] {
		if n.Name != "FunctionArgument" {
			return nil, NewInvalidNodeStructureError(name, n)
		}
		arg, err := ParseFunctionArgument(n)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	switch functionName {
	case "length":
		if l := len(args); l != 1 {
			return nil, fmt.Errorf("invalid number of arguments (%d) for %s", l, functionName)
		}
		if arg := args[0]; typeOfArgument(arg) != valueType {
			return nil, fmt.Errorf("invalid arg type %s for %s", arg, functionName)
		}
	case "count":
		if l := len(args); l != 1 {
			return nil, fmt.Errorf("invalid number of arguments (%d) for %s", l, functionName)
		}
		if arg := args[0]; typeOfArgument(arg) != nodesType {
			return nil, fmt.Errorf("invalid arg type %s for %s", arg, functionName)
		}
	case "value":
		if l := len(args); l != 1 {
			return nil, fmt.Errorf("invalid number of arguments (%d) for %s", l, functionName)
		}
		if typ := typeOfArgument(args[0]); typ != nodesType && typ != valueType {
			return nil, fmt.Errorf("invalid arg type %s for %s", args[0], functionName)
		}
	case "match", "search":
		if l := len(args); l != 2 {
			return nil, fmt.Errorf("invalid number of arguments (%d) for %s", l, functionName)
		}
		if arg := args[0]; typeOfArgument(arg) != valueType {
			return nil, fmt.Errorf("invalid arg type %s for %s", arg, functionName)
		}
		if arg := args[1]; typeOfArgument(arg) != valueType {
			return nil, fmt.Errorf("invalid arg type %s for %s", arg, functionName)
		}
	}

	return &FunctionExpr{
		Name:      functionName,
		Arguments: args,
	}, nil
}

func (s FunctionExpr) String() string {
	var str []string
	for _, a := range s.Arguments {
		str = append(str, a.String())
	}
	return fmt.Sprintf("%s(%s)", s.Name, strings.Join(str, ", "))
}

func (s FunctionExpr) Value(ref any) (any, error) {
	panic("should be implemented by the caller of this function")
}

func (s FunctionExpr) argument() {}

func (s FunctionExpr) comparable() {}

func (s FunctionExpr) testExpr() {}

type IndexSegment struct {
	Selector *IndexSelector
}

func ParseIndexSegment(n *parser.Node) (*IndexSegment, error) {
	name := "IndexSegment"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	selector, err := ParseIndexSelector(n.Children()[0])
	if err != nil {
		return nil, err
	}
	return &IndexSegment{
		Selector: selector,
	}, nil
}

func (s IndexSegment) String() string {
	return fmt.Sprintf("[%s]", s.Selector.String())
}

func (s IndexSegment) Value(ref any) (any, error) {
	switch ref := ref.(type) {
	case []any:
		if len(ref) <= s.Selector.Index {
			return nil, nil
		}
		return ref[s.Selector.Index], nil
	default:
		return nil, fmt.Errorf("unsupported ref type: %T", ref)
	}
}

func (s IndexSegment) singularQuerySegment() {}

type IndexSelector struct {
	Index int
}

func ParseIndexSelector(n *parser.Node) (*IndexSelector, error) {
	name := "IndexSelector"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	idx, err := parseInt(n.Children()[0])
	if err != nil {
		return nil, err
	}
	return &IndexSelector{Index: idx}, nil
}

func (s IndexSelector) String() string {
	return fmt.Sprintf("%d", s.Index)
}

func (s IndexSelector) selector() {}

type JSONPathQuery struct {
	Segments []Segment
}

func ParseJSONPathQuery(n *parser.Node) (*JSONPathQuery, error) {
	name := "JsonpathQuery"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	cs := n.Children()
	if len(cs) != 2 {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	if cs[0].Name != "RootIdentifier" {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	segments, err := ParseSegments(cs[1])
	if err != nil {
		return nil, err
	}
	return &JSONPathQuery{
		Segments: segments,
	}, nil
}

func (q JSONPathQuery) String() string {
	str := "$"
	for _, s := range q.Segments {
		str += s.String()
	}
	return str
}

func (q JSONPathQuery) argument() {}

func (q JSONPathQuery) testExpr() {}

type LogicalAndExpr struct {
	Expressions []BasicExpr
}

func ParseLogicalAndExpr(n *parser.Node) (*LogicalAndExpr, error) {
	name := "LogicalAndExpr"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	var expressions []BasicExpr
	for _, n := range n.Children() {
		expr, err := ParseBasicExpr(n)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expr)
	}
	return &LogicalAndExpr{
		Expressions: expressions,
	}, nil
}

func (s LogicalAndExpr) String() string {
	var str []string
	for _, e := range s.Expressions {
		str = append(str, e.String())
	}
	return strings.Join(str, "&&")
}

type LogicalExpr struct {
	Expressions []*LogicalAndExpr
}

func ParseLogicalExpr(n *parser.Node) (*LogicalExpr, error) {
	name := "LogicalExpr"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	var expressions []*LogicalAndExpr
	for _, n := range n.Children() {
		expr, err := ParseLogicalAndExpr(n)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expr)
	}
	return &LogicalExpr{
		Expressions: expressions,
	}, nil
}

func (s LogicalExpr) String() string {
	var str []string
	for _, e := range s.Expressions {
		str = append(str, e.String())
	}
	return strings.Join(str, " || ")
}

func (s LogicalExpr) argument() {}

type MemberNameShorthand struct {
	Name string
}

func (s MemberNameShorthand) String() string {
	return fmt.Sprintf("['%s']", s.Name)
}

func (s MemberNameShorthand) childSegment() {}

func (s MemberNameShorthand) segment() {}

type NameSegment struct {
	Name string
}

func ParseNameSegment(n *parser.Node) (*NameSegment, error) {
	name := "NameSegment"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	switch n := n.Children()[0]; n.Name {
	case "NameSelector":
		selector, err := ParseNameSelector(n)
		if err != nil {
			return nil, err
		}
		return &NameSegment{
			Name: selector.Name,
		}, nil
	case "MemberNameShorthand":
		return &NameSegment{
			Name: n.Value(),
		}, nil
	default:
		return nil, NewInvalidNodeStructureError(name, n)
	}
}

func (s NameSegment) String() string {
	return fmt.Sprintf("[%s]", s.Name)
}

func (s NameSegment) Value(ref any) (any, error) {
	switch ref := ref.(type) {
	case map[string]any:
		return ref[s.Name], nil
	default:
		return nil, nil
	}
}

func (s NameSegment) singularQuerySegment() {}

type NameSelector struct {
	Name string
}

func ParseNameSelector(n *parser.Node) (*NameSelector, error) {
	name := "NameSelector"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	str, err := ParseStringLiteral(n.Children()[0])
	if err != nil {
		return nil, err
	}
	return &NameSelector{
		Name: string(*str),
	}, nil
}

func (s NameSelector) String() string {
	name := strings.Replace(s.Name, "'", "\\'", -1)
	return fmt.Sprintf("'%s'", name)
}

func (s NameSelector) selector() {}

type ParenExpr struct {
	Negation    bool
	LogicalExpr *LogicalExpr
}

func ParseParenExpr(n *parser.Node) (*ParenExpr, error) {
	name := "ParenExpr"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	var negation bool
	for _, n := range n.Children() {
		switch n.Name {
		case "LogicalNotOp":
			negation = true
		case "LogicalExpr":
			logicalExpr, err := ParseLogicalExpr(n)
			if err != nil {
				return nil, err
			}
			return &ParenExpr{
				Negation:    negation,
				LogicalExpr: logicalExpr,
			}, nil
		default:
			return nil, NewInvalidNodeStructureError(name, n)
		}
	}
	return nil, NewInvalidNodeStructureError(name, n)
}

func (s ParenExpr) String() string {
	if s.Negation {
		return fmt.Sprintf("!(%s)", s.LogicalExpr.String())
	}
	return fmt.Sprintf("(%s)", s.LogicalExpr.String())
}

func (s ParenExpr) basicExpr() {}

type RelQuery struct {
	Segments []Segment
}

func ParseRelQuery(n *parser.Node) (*RelQuery, error) {
	name := "RelQuery"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	cs := n.Children()
	if len(cs) != 2 {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	if cs[0].Name != "CurrentNodeIdentifier" {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	segments, err := ParseSegments(cs[1])
	if err != nil {
		return nil, err
	}
	return &RelQuery{
		Segments: segments,
	}, nil
}

func (q RelQuery) String() string {
	str := "@"
	for _, s := range q.Segments {
		str += s.String()
	}
	return str
}

func (q RelQuery) argument() {}

func (q RelQuery) testExpr() {}

type RelSingularQuery struct {
	Segments []SingularQuerySegment
}

func ParseRelSingularQuery(n *parser.Node) (*RelSingularQuery, error) {
	name := "RelSingularQuery"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	cs := n.Children()
	if len(cs) != 2 {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	if cs[0].Name != "CurrentNodeIdentifier" {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	segments, err := ParseSingularQuerySegments(cs[1])
	if err != nil {
		return nil, err
	}
	return &RelSingularQuery{
		Segments: segments,
	}, nil
}

func (s RelSingularQuery) String() string {
	str := "@"
	for _, s := range s.Segments {
		str += s.String()
	}
	return str
}

func (s RelSingularQuery) Value(ref any) (any, error) {
	current := ref
	for _, segment := range s.Segments {
		c, err := segment.Value(current)
		if err != nil {
			return nil, err
		}
		current = c
	}
	return current, nil
}

func (s RelSingularQuery) comparable() {}

type Segment interface {
	fmt.Stringer

	segment()
}

func ParseSegment(n *parser.Node) (Segment, error) {
	name := "Segment"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	switch n := n.Children()[0]; n.Name {
	case "ChildSegment":
		return ParseChildSegment(n)
	case "DescendantSegment":
		return ParseDescendantSegment(n)
	default:
		return nil, NewInvalidNodeStructureError(name, n)
	}
}

func ParseSegments(n *parser.Node) ([]Segment, error) {
	name := "Segments"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	var segments []Segment
	for _, n := range n.Children() {
		segment, err := ParseSegment(n)
		if err != nil {
			return nil, err
		}
		segments = append(segments, segment)
	}
	return segments, nil
}

type Selector interface {
	fmt.Stringer

	selector()
}

func ParseSelector(n *parser.Node) (Selector, error) {
	name := "Selector"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	switch n := n.Children()[0]; n.Name {
	case "NameSelector":
		return ParseNameSelector(n)
	case "WildcardSelector":
		return new(WildcardSelector), nil
	case "IndexSelector":
		return ParseIndexSelector(n)
	case "SliceSelector":
		return ParseSliceSelector(n)
	case "FilterSelector":
		return ParseFilterSelector(n)
	default:
		return nil, NewInvalidNodeStructureError(name, n)
	}
}

type SingularQuerySegment interface {
	fmt.Stringer

	Value(ref any) (any, error)
	singularQuerySegment()
}

func ParseSingularQuerySegments(n *parser.Node) ([]SingularQuerySegment, error) {
	name := "SingularQuerySegments"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	var segments []SingularQuerySegment
	for _, n := range n.Children() {
		switch n.Name {
		case "NameSegment":
			segment, err := ParseNameSegment(n)
			if err != nil {
				return nil, err
			}
			segments = append(segments, segment)
		case "IndexSegment":
			segment, err := ParseIndexSegment(n)
			if err != nil {
				return nil, err
			}
			segments = append(segments, segment)
		default:
			return nil, NewInvalidNodeStructureError(name, n)
		}
	}
	return segments, nil
}

type SliceSelector struct {
	Start, End, Step int
}

func ParseSliceSelector(n *parser.Node) (*SliceSelector, error) {
	name := "SliceSelector"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	start, end, step := 0, -1, 1
	for _, n := range n.Children() {
		switch n.Name {
		case "Start":
			idx, err := parseInt(n.Children()[0])
			if err != nil {
				return nil, err
			}
			start = idx
		case "End":
			idx, err := parseInt(n.Children()[0])
			if err != nil {
				return nil, err
			}
			end = idx
		case "Step":
			idx, err := parseInt(n.Children()[0])
			if err != nil {
				return nil, err
			}
			step = idx
		default:
			return nil, NewInvalidNodeStructureError(name, n)
		}
	}
	return &SliceSelector{
		Start: start,
		End:   end,
		Step:  step,
	}, nil
}

func (s SliceSelector) String() string {
	var str string
	if 0 < s.Start {
		str += fmt.Sprintf("%d", s.Start)
	}
	str += ":"
	if 0 < s.End {
		str += fmt.Sprintf("%d", s.End)
	}
	if 0 != s.Step && 1 != s.Step {
		str += fmt.Sprintf(":%d", s.Step)
	}
	return str
}

func (s SliceSelector) selector() {}

type TestExpr struct {
	Negation bool
	TestExpr TestExpression
}

func ParseTestExpr(n *parser.Node) (*TestExpr, error) {
	name := "TestExpr"
	if n.Name != name {
		return nil, NewInvalidNodeStructureError(name, n)
	}
	var negation bool
	for _, n := range n.Children() {
		switch n.Name {
		case "LogicalNotOp":
			negation = true
		case "FunctionExpr":
			q, err := ParseFunctionExpr(n)
			if err != nil {
				return nil, err
			}
			if typ := typeOfArgument(q); typ != logicalType {
				return nil, fmt.Errorf("invalid arg type %s for %s", q, name)
			}
			return &TestExpr{
				Negation: negation,
				TestExpr: q,
			}, nil
		case "RelQuery":
			q, err := ParseRelQuery(n)
			if err != nil {
				return nil, err
			}
			return &TestExpr{
				Negation: negation,
				TestExpr: q,
			}, nil
		case "JsonpathQuery":
			q, err := ParseJSONPathQuery(n)
			if err != nil {
				return nil, err
			}
			return &TestExpr{
				Negation: negation,
				TestExpr: q,
			}, nil
		default:
			return nil, NewInvalidNodeStructureError(name, n)
		}

	}
	return nil, NewInvalidNodeStructureError(name, n)
}

func (s TestExpr) String() string {
	if s.Negation {
		return fmt.Sprintf("!%s", s.TestExpr)
	}
	return s.TestExpr.String()
}

func (s TestExpr) basicExpr() {}

type TestExpression interface {
	fmt.Stringer

	testExpr()
}

type WildcardSelector struct{}

func (s WildcardSelector) String() string {
	return "*"
}

func (s WildcardSelector) childSegment() {}

func (s WildcardSelector) segment() {}

func (s WildcardSelector) selector() {}

type argumentType int

const (
	unknownType argumentType = iota
	valueType
	logicalType
	nodesType
)

func typeOfArgument(arg any) argumentType {
	switch arg := arg.(type) {
	case *JSONPathQuery:
		for _, segment := range arg.Segments {
			if typ := typeOfArgument(segment); typ != valueType {
				return typ
			}
		}
		return valueType
	case *FunctionExpr:
		switch arg.Name {
		case "length", "count", "value":
			return valueType
		case "match", "search":
			return logicalType
		default:
			return unknownType
		}
	case *WildcardSelector:
		return nodesType
	case *DescendantSegment:
		return typeOfArgument(arg.Segment)
	case *RelSingularQuery, *AbsSingularQuery, *MemberNameShorthand:
		return valueType
	case *RelQuery:
		for _, segment := range arg.Segments {
			if typ := typeOfArgument(segment); typ != valueType {
				return typ
			}
		}
		return valueType
	case *Boolean, *Null, *Number, *String:
		return valueType
	default:
		panic(fmt.Sprintf("unsupported arg type %T", arg))
	}
}
