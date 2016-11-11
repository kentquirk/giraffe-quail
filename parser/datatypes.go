package parser

import (
	"fmt"
	"strconv"
)

type ValueType int

const (
	VT_Var ValueType = iota
	VT_Int
	VT_Float
	VT_Str
	VT_Bool
	VT_List
	VT_Obj
	VT_Enum
	VT_Null
	VT_Field
)

var VT_Names = []string{
	"VT_Var",
	"VT_Int",
	"VT_Float",
	"VT_Str",
	"VT_Bool",
	"VT_List",
	"VT_Obj",
	"VT_Enum",
	"VT_Null",
	"VT_Field",
}

type Value struct {
	T ValueType
	V interface{}
	N string
}

func (v Value) String() string {
	return fmt.Sprintf("%s:%s %v", VT_Names[v.T], v.N, v.V)
}

func MakeVar(name string) Value {
	return Value{T: VT_Var, V: nil, N: name}
}

func MakeInt(v int64) Value {
	return Value{T: VT_Int, V: v, N: ""}
}

func ParseInt(s string) (Value, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	val := Value{T: VT_Int, V: v, N: ""}
	return val, err
}

func MakeFloat(v float64) Value {
	return Value{T: VT_Float, V: v, N: ""}
}

func ParseFloat(s string) (Value, error) {
	v, err := strconv.ParseFloat(s, 64)
	val := Value{T: VT_Float, V: v, N: ""}
	return val, err
}

func MakeStr(v string) Value {
	return Value{T: VT_Str, V: v, N: ""}
}

func MakeBool(v bool) Value {
	return Value{T: VT_Bool, V: v, N: ""}
}

func MakeNull() Value {
	return Value{T: VT_Null, V: nil, N: ""}
}

func MakeEnum(name string) Value {
	return Value{T: VT_Enum, V: nil, N: name}
}

func MakeList() Value {
	return Value{T: VT_List, V: make([]Value, 0), N: ""}
}

func MakeObj() Value {
	return Value{T: VT_List, V: make(map[string]Value), N: ""}
}

func MakeField(name string, v Value) Value {
	return Value{T: VT_Field, V: v, N: name}
}

func (v Value) AsInt64() int64 {
	switch v.T {
	case VT_Int:
		return v.V.(int64)
	case VT_Float:
		return int64(v.V.(float64))
	default:
		return 0
	}
}

func (v Value) AsInt() int {
	return int(v.AsInt64())
}

func (v Value) AsFloat64() float64 {
	switch v.T {
	case VT_Float:
		return v.V.(float64)
	case VT_Int:
		return float64(v.V.(int64))
	default:
		return 0
	}
}

func (v Value) AsBool() bool {
	switch v.T {
	case VT_Bool:
		return v.V.(bool)
	default:
		return false
	}
}

func (v Value) AsStr() string {
	switch v.T {
	case VT_Str:
		return v.V.(string)
	default:
		return ""
	}
}

func (v Value) AsList() []Value {
	switch v.T {
	case VT_List:
		return v.V.([]Value)
	default:
		return []Value{}
	}
}

func (v Value) Append(item ...Value) {
	if v.T != VT_List {
		panic("Attempt to append to a non-list!")
	}
	list := v.V.([]Value)
	v.V = append(list, item...)
}

func (v Value) Set(key string, val Value) {
	if v.T != VT_Obj {
		panic("Attempt to Set to a non-object!")
	}
	m := v.V.(map[string]Value)
	m[key] = val
}

func (v Value) SetField(f Value) {
	if v.T != VT_Obj {
		panic("Attempt to SetField to a non-object!")
	}
	if f.T != VT_Field {
		panic("Attempt to SetField from a non-field!")
	}
	m := v.V.(map[string]Value)
	m[f.N] = f.V.(Value)
}

func (v Value) Get(key string) (Value, bool) {
	if v.T != VT_Obj {
		panic("Attempt to Get from a non-object!")
	}
	m := v.V.(map[string]Value)
	val, ok := m[key]
	return val, ok
}
