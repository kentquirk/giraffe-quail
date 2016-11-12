package parser

import (
	"errors"
	"fmt"
	"regexp"
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
	VT_NonNull
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
	"VT_NonNull",
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

// A variable is a ValueType that has a name in N but a Value type for its V parameter.
// The V param has the information for the type of the Variable
func MakeVar(name string) Value {
	return Value{T: VT_Var, V: nil, N: name}
}

func MakeInt(v int) Value {
	return Value{T: VT_Int, V: v, N: ""}
}

func MakeFloat(v float64) Value {
	return Value{T: VT_Float, V: v, N: ""}
}

func ParseInt(s string) (Value, error) {
	v, err := strconv.ParseInt(s, 10, 32)
	val := Value{T: VT_Int, V: v, N: ""}
	return val, err
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

func MakeNonNull(inner Value) Value {
	return Value{T: VT_NonNull, V: inner, N: ""}
}

func MakeEnum(name string) Value {
	return Value{T: VT_Enum, V: nil, N: name}
}

func MakeList() Value {
	return Value{T: VT_List, V: make([]Value, 0), N: ""}
}

// an object variable uses the name field to contain the name of
// the object itself, and the V field is an array of other values
// in the order desired. The N fields of the array items indicate
// their individual names
func MakeObj() Value {
	return Value{T: VT_List, V: make([]Value, 0), N: ""}
}

func MakeField(name string, v Value) Value {
	return Value{T: VT_Field, V: v, N: name}
}

func (v Value) AsInt() (int, error) {
	switch v.T {
	case VT_Int:
		return v.V.(int), nil
	case VT_Float:
		return int(v.V.(float64)), nil
	case VT_Bool:
		if v.V.(bool) {
			return -1, nil
		} else {
			return 0, nil
		}
	case VT_Str:
		i, err := strconv.ParseInt(v.V.(string), 10, 32)
		return int(i), err
	default:
		return 0, errors.New("Unable to convert to Int")
	}
}

func (v Value) AsFloat() (float64, error) {
	switch v.T {
	case VT_Float:
		return v.V.(float64), nil
	case VT_Int:
		return float64(v.V.(int64)), nil
	default:
		return 0, errors.New("Unable to convert to Float")
	}
}

func (v Value) AsBool() (bool, error) {
	switch v.T {
	case VT_Bool:
		return v.V.(bool), nil
	case VT_Float:
		return v.V.(float64) != 0, nil
	case VT_Int:
		return float64(v.V.(int64)) != 0, nil
	case VT_Str:
		s := v.V.(string)
		m, _ := regexp.Match("(?i:^(yes|true|1|y|t)$)", []byte(s))
		if m {
			return true, nil
		}
		m, _ = regexp.Match("(?i:^(no|false|0|n|f)$)", []byte(s))
		if m {
			return false, nil
		}
		fallthrough
	default:
		return false, errors.New("Unable to convert to Bool")
	}
}

func (v Value) AsStr() (string, error) {
	switch v.T {
	case VT_Bool:
		if v.V.(bool) {
			return "true", nil
		} else {
			return "false", nil
		}
	case VT_Float:
		return fmt.Sprintf("%f", v.V), nil
	case VT_Int:
		return fmt.Sprintf("%d", v.V), nil
	case VT_Str:
		return v.V.(string), nil
	default:
		return "", errors.New("Unable to convert to String")
	}
}

func (v Value) AsList() ([]Value, error) {
	switch v.T {
	case VT_List:
		return v.V.([]Value), nil
	default:
		return []Value{}, errors.New("Unable to convert to List")
	}
}

// func (v Value) Assign(x Value) error {
//     if x == nil && v.NonNull {
//         return errors.New(fmt.Sprintf("Attempt to assign null to a non-null value (%s)", v.N))
//     }
// }

func (v Value) Append(item ...Value) {
	if v.T != VT_List {
		panic("Attempt to append to a non-list!")
	}
	list := v.V.([]Value)
	v.V = append(list, item...)
}

func (v Value) Set(val Value) {
	if v.T != VT_Obj {
		panic("Attempt to Set to a non-object!")
	}
	items := v.V.([]Value)
	for i := range items {
		if items[i].N == val.N {
			items[i].V = val.V
			return
		}
	}
	v.V = append(items, val)
}

func (v Value) Get(key string) (Value, bool) {
	if v.T != VT_Obj {
		panic("Attempt to Get from a non-object!")
	}
	items := v.V.([]Value)
	for i := range items {
		if items[i].N == key {
			return items[i], true
		}
	}
	return MakeNull(), false
}
