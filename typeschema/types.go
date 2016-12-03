package typeschema

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
// List         Subtypes[0] contains the type of the list
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

// TypeNameFor takes a kind and a list of names and returns a unique
// name given that combination
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