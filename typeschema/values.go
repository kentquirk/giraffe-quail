package typeschema

import (
	"errors"
	"strconv"
)

type Value struct {
	T Type
	V interface{}
}

type ValueRegistry struct {
	Values map[string]Value
}

func NewValueRegistry() *ValueRegistry {
	vr := &ValueRegistry{Values: make(map[string]Value)}
	return vr
}

// Register sets up a value in the value registry. It is an error if the name
// already exists.
func (vr *ValueRegistry) Register(name string, v Value) (Value, error) {
	if _, found := vr.Values[name]; found {
		return v, errors.New("Value " + name + " already defined; cannot be overridden.")
	}
	vr.Values[name] = v
	return v, nil
}

func (tr *TypeRegistry) MakeInt(v int) Value {
	return Value{T: tr.Int(), V: v}
}

func (tr *TypeRegistry) MakeFloat(v float64) Value {
	return Value{T: tr.Float(), V: v}
}

func (tr *TypeRegistry) MakeStr(v string) Value {
	return Value{T: tr.Str(), V: v}
}

func (tr *TypeRegistry) MakeBool(v bool) Value {
	return Value{T: tr.Bool(), V: v}
}

func (tr *TypeRegistry) MakeNull() Value {
	return Value{T: tr.Null(), V: nil}
}

// ParseInt constructs an Int Value object from a string formatted as an integer.
func (tr *TypeRegistry) ParseInt(s string) (Value, error) {
	v, err := strconv.ParseInt(s, 10, 32)
	val := tr.MakeInt(int(v))
	return val, err
}

// ParseFloat constructs a Float Value object from a string formatted as a float.
func (tr *TypeRegistry) ParseFloat(s string) (Value, error) {
	v, err := strconv.ParseFloat(s, 64)
	val := tr.MakeFloat(v)
	return val, err
}
