package types

import (
	"errors"
	"strconv"
)

type Constructor func() Value
type TypeRegistry struct {
	Types        map[string]Type
	Constructors map[string]Constructor
}

func NewTypeRegistry() *TypeRegistry {
	tr := &TypeRegistry{
		Types:        make(map[string]Type),
		Constructors: make(map[string]Constructor),
	}
	// these are the fundamental types we understand at start
	tr.Register(T_String, Scalar)
	tr.RegisterConstructor(T_String, func() Value { return tr.MakeStr("") })
	tr.Register(T_Int, Scalar)
	tr.RegisterConstructor(T_Int, func() Value { return tr.MakeInt(0) })
	tr.Register(T_Float, Scalar)
	tr.RegisterConstructor(T_Float, func() Value { return tr.MakeFloat(0) })
	tr.Register(T_Boolean, Scalar)
	tr.RegisterConstructor(T_Boolean, func() Value { return tr.MakeBool(false) })
	tr.Register(T_ID, Scalar)
	tr.RegisterConstructor(T_ID, func() Value { return tr.MakeID("") })
	tr.Register(T_Null, Null)
	tr.RegisterConstructor(T_Null, func() Value { return tr.MakeNull() })
	return tr
}

// Register adds a type by name to the type registry. It is an error if the type
// already exists.
func (tr *TypeRegistry) Register(name string, k Kind, subtypes ...Type) (Type, error) {
	return tr.RegisterWithFields(name, k, nil, subtypes...)
}

// RegisterConstructor adds a type constructor to the type registry. The type must exist
// in the registry.
func (tr *TypeRegistry) RegisterConstructor(name string, f Constructor) error {
	if _, found := tr.Types[name]; !found {
		return errors.New("Type " + name + " not defined when adding Constructor.")
	}
	tr.Constructors[name] = f
	return nil
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

// ID is a convenience method for retrieving the standard ID type
func (tr *TypeRegistry) ID() Type {
	return tr.MustGet(T_ID)
}

// Bool is a convenience method for retrieving the standard Bool type
func (tr *TypeRegistry) Bool() Type {
	return tr.MustGet(T_Boolean)
}

// Null is a convenience method for retrieving the standard Null type
func (tr *TypeRegistry) Null() Type {
	return tr.MustGet(T_Null)
}

func (tr *TypeRegistry) ListType(t Type) Type {
	return Type{Kind: List, Subtypes: []Type{t}}
}

func (tr *TypeRegistry) NonNullableType(t Type) Type {
	return Type{Kind: List, Subtypes: []Type{t}}
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

func (tr *TypeRegistry) MakeID(v string) Value {
	return Value{T: tr.ID(), V: v}
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

func (tr *TypeRegistry) MakeObj(t Type) Value {
	// This makes a Value for a specific ObjType
	return Value{T: t, V: make(map[string]Value)}
}

func (tr *TypeRegistry) MakeValueOf(t Type) (Value, error) {
	if constr, ok := tr.Constructors[t.Key()]; ok {
		return constr(), nil
	}
	switch t.Kind {
	case Null, Scalar:
		return tr.MakeNull(), errors.New("Constructor for Scalar type " + t.Name + " is missing.")
	case Obj, Interface, Union:
		return tr.MakeObj(t), nil
	case List:
		return tr.MakeListOf(t), nil
	default:
		return tr.MakeNull(), errors.New("Don't know how to make variable of type " + t.Name)
	}
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
