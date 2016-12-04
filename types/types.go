package types

import (
	"errors"
	"strings"
)

type Kind int

const (
	Null        Kind = iota
	Scalar      Kind = iota
	Enum        Kind = iota
	Obj         Kind = iota
	Interface   Kind = iota
	List        Kind = iota
	NonNullable Kind = iota
	Union       Kind = iota
	Temp        Kind = iota
)

// these are the canonical names for the fundamental scalar types
const (
	T_String  = "String"
	T_Int     = "Int"
	T_Float   = "Float"
	T_Boolean = "Boolean"
	T_ID      = "ID"
	T_Null    = "*NULL*"
)

// Type is the structure that contains information about a GraphQL Type.
// The Kind of type determines how it stores information. An attempt
// was made to use Go's type system to do this, by defining an interface
// that Type would implement, but it started getting pretty big, and the
// code to use it was getting cumbersome. Instead, we're going to
// just switch on the Kind value instead and we can reuse a lot of things
// between the different Kinds.
// Here's the behavior of different types:
//
// Kind         Subtypes/Fields
// -----------  --------------------------------------------------------
// Null         nil
// Scalar       nil
// Enum         nil, with Values of this type stored in the value registry
// Obj          nil, Fields contains an ordered list of fields
// Interface    just like obj
// List         Subtypes[0] contains the type of the list; an empty list may have no subtype.
// NonNullable  Subtypes[0] contains the nonnullable type
// Union        Subtypes contains a list of the types in the union
// Temp         nil
type Type struct {
	Name     string
	Kind     Kind
	Subtypes []Type
	Fields   []Field
}

// A field is a Name and a Type, used in Obj (and Interface) types,
// as well as in the Arg list for a field.
type Field struct {
	N    string
	Args []Field
	T    Type
}

// Key returns the unique name of this type
func (t Type) Key() string {
	switch t.Kind {
	case Union:
		stnames := make([]string, 0)
		for _, st := range t.Subtypes {
			stnames = append(stnames, st.Key())
		}
		return TypeNameFor(t.Kind, stnames...)
	case List, NonNullable:
		return TypeNameFor(t.Kind, t.Subtypes[0].Key())
	default:
		return TypeNameFor(t.Kind, t.Name)
	}
}

func (t Type) String() string {
	return t.Key()
}

func (t Type) Is(other Type) bool {
	if t.Kind != other.Kind || t.Name != other.Name {
		return false
	}
	if len(t.Subtypes) != len(other.Subtypes) {
		return false
	}
	switch t.Kind {
	case List, Union, NonNullable:
		for i := range t.Subtypes {
			if !t.Subtypes[i].Is(other.Subtypes[i]) {
				return false
			}
		}
	}
	return false
}

func (t Type) HasField(name string) bool {
	if t.Kind != Obj {
		return false
	}
	for _, f := range t.Fields {
		if f.N == name {
			return true
		}
	}
	return false
}

func (t Type) GetField(name string) (Field, error) {
	if t.Kind != Obj {
		return Field{}, errors.New("GetField on non-object.")
	}
	for _, f := range t.Fields {
		if f.N == name {
			return f, nil
		}
	}
	return Field{}, errors.New("Field not found: " + name)
}

// TypeNameFor is a free function that takes a kind and a list of names and
// returns a unique name given that combination
func TypeNameFor(k Kind, names ...string) string {
	switch k {
	case List:
		return "[" + names[0] + "]"
	case NonNullable:
		return names[0] + "!"
	case Union:
		return strings.Join(names, "|")
	default:
		return names[0]
	}
}
