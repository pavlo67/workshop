package scriptor

import (
	"regexp"
)

var reInfix *regexp.Regexp

func init() {
	infixesStr := "^["
	for sign := range Infixes {
		if sign == "\\" || sign == "-" {
			infixesStr += "\\"
		}
		infixesStr += sign
	}
	infixesStr += "]"

	reInfix = regexp.MustCompile(infixesStr)
}

type Func func(a, b interface{}) interface{}

type Infix struct {
	Priority  int
	Signatura [3]Type
	Func
}

var Infixes = map[string][]Infix{
	"+": {
		{0, [3]Type{TypeInt, TypeInt, TypeInt}, AddInt},
		{0, [3]Type{TypeFloat, TypeFloat, TypeFloat}, AddFloat},
	},
	"-": {
		{0, [3]Type{TypeInt, TypeInt, TypeInt}, SubInt},
		{0, [3]Type{TypeFloat, TypeFloat, TypeFloat}, SubFloat},
	},
	"*": {
		{10, [3]Type{TypeInt, TypeInt, TypeInt}, MultInt},
		{10, [3]Type{TypeFloat, TypeFloat, TypeFloat}, MultFloat},
	},
	"/": {
		{10, [3]Type{TypeInt, TypeInt, TypeInt}, DivInt},
		{10, [3]Type{TypeFloat, TypeFloat, TypeFloat}, DivFloat},
	},
}

// + --------------------------------------------------------------------

var _ Func = AddInt

func AddInt(a, b interface{}) interface{} {
	aInt, _ := a.(int64)
	bInt, _ := b.(int64)
	return aInt + bInt
}

var _ Func = AddFloat

func AddFloat(a, b interface{}) interface{} {
	aFloat, _ := a.(float64)
	bFloat, _ := b.(float64)
	return aFloat + bFloat
}

// + --------------------------------------------------------------------

var _ Func = SubInt

func SubInt(a, b interface{}) interface{} {
	aInt, _ := a.(int64)
	bInt, _ := b.(int64)
	return aInt - bInt
}

var _ Func = SubFloat

func SubFloat(a, b interface{}) interface{} {
	aFloat, _ := a.(float64)
	bFloat, _ := b.(float64)
	return aFloat - bFloat
}

// * --------------------------------------------------------------------

var _ Func = MultInt

func MultInt(a, b interface{}) interface{} {
	aInt, _ := a.(int64)
	bInt, _ := b.(int64)
	return aInt * bInt
}

var _ Func = MultFloat

func MultFloat(a, b interface{}) interface{} {
	aFloat, _ := a.(float64)
	bFloat, _ := b.(float64)
	return aFloat * bFloat
}

// * --------------------------------------------------------------------

var _ Func = DivInt

func DivInt(a, b interface{}) interface{} {
	aInt, _ := a.(int64)
	bInt, _ := b.(int64)
	return aInt / bInt
}

var _ Func = DivFloat

func DivFloat(a, b interface{}) interface{} {
	aFloat, _ := a.(float64)
	bFloat, _ := b.(float64)
	return aFloat / bFloat
}
