package types

import (
	"errors"
	"strconv"
)

type TypeRegistry struct {
	Types map[string]Type
}

func NewTypeRegistry() *TypeRegistry {
	tr := &TypeRegistry{Types: make(map[string]Type)}
	// these are the fundamental types we understand at start
	tr.Register(T_String, Scalar)
	tr.Register(T_Int, Scalar)
	tr.Register(T_Float, Scalar)
	tr.Register(T_Boolean, Scalar)
	tr.Register(T_ID, Scalar)
	tr.Register(T_Null, Null)
	return tr
}

// Register adds a type by name to the type registry. It is an error if the type
// already exists.
func (tr *TypeRegistry) Register(name string, k Kind, subtypes ...Type) (Type, error) {
	return tr.RegisterWithFields(name, k, nil, subtypes...)
}

// Register adds a type by name to the type registry. It is an error if the type
// already exists.
func (tr *TypeRegistry) RegisterWithFields(name string, k Kind, fields []Field, subtypes ...Type) (Type, error) {
	newtype := Type{Name: name, Kind: k, Subtypes: subtypes, Fields: fields}
	if t, found := tr.Types[name]; found {
		if t.Kind != Temp {
			return t, errors.New("Type " + name + " already defined.")
		}
	}
	tr.Types[name] = newtype
	return newtype, nil
}

// Update modifies a type by looking it up by name in the type registry and then
// replacing it.
// It is an error if the type does not exist.
func (tr *TypeRegistry) Update(u Type) (Type, error) {
	name := u.Key()
	if t, found := tr.Types[name]; !found {
		return t, errors.New("Type " + name + " not defined but attempted update.")
	}
	tr.Types[name] = u
	return u, nil
}

// Get retrieves a type from the registry by name. It returns an error
// if the type does not exist
func (tr *TypeRegistry) Get(name string) (Type, error) {
	t, ok := tr.Types[name]
	if ok {
		return t, nil
	}
	return t, errors.New("Type " + name + " was not found in the type registry.")
}

// MustGet retrieves a type from the registry by name. It panics if the type
// does not exist. Mainly intended for internal use.
func (tr *TypeRegistry) MustGet(name string) Type {
	if t, ok := tr.Types[name]; ok {
		return t
	}
	panic("Type " + name + " was not found in the type registry.")
}

// MaybeGet retrieves a type from the registry by name. If the type does
// not exist, it creates a placeholder type that can be updated later.
func (tr *TypeRegistry) MaybeGet(name string) (Type, error) {
	t, ok := tr.Types[name]
	if !ok {
		return tr.Register(name, Temp)
	}
	return t, nil
}

// Int is a convenience method for retrieving the standard Int type
func (tr *TypeRegistry) Int() Type {
	return tr.MustGet(T_Int)
}

// Float is a convenience method for retrieving the standard Float type
func (tr *TypeRegistry) Float() Type {
	return tr.MustGet(T_Float)
}

// Str is a convenience method for retrieving the standard String type
func (tr *TypeRegistry) Str() Type {
	return tr.MustGet(T_String)
}

// Bool is a convenience method for retrieving the standard Bool type
func (tr *TypeRegistry) Bool() Type {
	return tr.MustGet(T_Boolean)
}

// Null is a convenience method for retrieving the standard Null type
func (tr *TypeRegistry) Null() Type {
	return tr.MustGet(T_Null)
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

// MakeEmptyList constructs a Value that has a type that is a List
// where the type's Subtypes list is empty.
// It also contains an empty array of Value objects.
func (tr *TypeRegistry) MakeEmptyList() Value {
	return Value{T: Type{Kind: List, Subtypes: []Type{}}, V: make([]Value, 0)}
}

// MakeListOf constructs a Value that has a type that is a List
// and contains an empty array of Value objects.
func (tr *TypeRegistry) MakeListOf(listType Type) Value {
	return Value{T: Type{Kind: List, Subtypes: []Type{listType}}, V: make([]Value, 0)}
}

func (tr *TypeRegistry) MakeNamelessObj() Value {
	// This makes a Value for a nameless, empty ObjType
	return Value{T: Type{Kind: Obj, Fields: make([]Field, 0)}, V: make(map[string]Value)}
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
