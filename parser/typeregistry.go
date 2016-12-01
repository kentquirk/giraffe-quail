package parser

import (
	"errors"
	"strconv"
	"strings"
)

// A Type description contains the information about the type of a data element.
// Type descriptions consist of some series (possibly 0) wrapper types,
// with a BaseType at the bottom of their Type description.

// Type interface supports:
// * A query to request the name of the BaseType.
// * Asking if the Type supports null values.
// * Getting a string representation of the type.
type Type interface {
	BaseType() string
	Nullable() bool
	String() string
	Is(t Type) bool
}

// TypeRegistry manages a collection of types by name. There are some predefined
// type names: String, Int, Float, Bool, and ID (ID is just a string that
// is not necessarily human-readable) are all names accessible to users.
// The name of the null type is "*NULL*", which cannot be created by
// a user because * is not a valid character in the name set.
type TypeRegistry map[string]Type

func CreateTypeRegistry() TypeRegistry {
	tr := TypeRegistry{}
	// these are the fundamental types we understand at start
	tr.Register("String", BaseType{"String"})
	tr.Register("Int", BaseType{"Int"})
	tr.Register("Float", BaseType{"Float"})
	tr.Register("Bool", BaseType{"Bool"})
	tr.Register("ID", BaseType{"ID"})
	tr.Register("*NULL*", BaseType{"*NULL*"})
	return tr
}

// Register adds a type by name to the type registry. It is an error if the type
// already exists.
func (tr TypeRegistry) Register(name string, t Type) error {
	if _, found := tr[name]; found {
		return errors.New("Type " + name + " already defined.")
	}
	tr[name] = t
	return nil
}

// Get retrieves a type from the registry by name. It returns an error
// if the type does not exist
func (tr TypeRegistry) Get(name string) (Type, error) {
	t, ok := tr[name]
	if ok {
		return t, nil
	}
	return t, errors.New("Type " + name + " was not found in the type registry.")
}

// MustGet retrieves a type from the registry by name. It panics if the type
// does not exist. Mainly intended for internal use.
func (tr TypeRegistry) MustGet(name string) Type {
	if t, ok := tr[name]; ok {
		return t
	}
	panic("Type " + name + " was not found in the type registry.")
}

// Int is a convenience method for retrieving the standard Int type
func (tr TypeRegistry) Int() Type {
	return tr.MustGet("Int")
}

// Float is a convenience method for retrieving the standard Float type
func (tr TypeRegistry) Float() Type {
	return tr.MustGet("Float")
}

// Str is a convenience method for retrieving the standard String type
func (tr TypeRegistry) Str() Type {
	return tr.MustGet("String")
}

// Bool is a convenience method for retrieving the standard Bool type
func (tr TypeRegistry) Bool() Type {
	return tr.MustGet("Bool")
}

// Null is a convenience method for retrieving the standard Null type
func (tr TypeRegistry) Null() Type {
	return tr.MustGet("*NULL*")
}

// ParseInt constructs an Int Value object from a string formatted as an integer.
func (tr TypeRegistry) ParseInt(s string) (Value, error) {
	v, err := strconv.ParseInt(s, 10, 32)
	val := tr.MakeInt(int(v))
	return val, err
}

// ParseFloat constructs a Float Value object from a string formatted as a float.
func (tr TypeRegistry) ParseFloat(s string) (Value, error) {
	v, err := strconv.ParseFloat(s, 64)
	val := tr.MakeFloat(v)
	return val, err
}

func (tr TypeRegistry) MakeInt(v int) Value {
	return Value{T: tr.Int(), V: v, N: ""}
}

func (tr TypeRegistry) MakeFloat(v float64) Value {
	return Value{T: tr.Float(), V: v, N: ""}
}

func (tr TypeRegistry) MakeStr(v string) Value {
	return Value{T: tr.Str(), V: v, N: ""}
}

func (tr TypeRegistry) MakeBool(v bool) Value {
	return Value{T: tr.Bool(), V: v, N: ""}
}

func (tr TypeRegistry) MakeNull() Value {
	return Value{T: tr.Null(), V: nil, N: ""}
}

func (tr TypeRegistry) MakeList() Value {
	// This makes a Value for a ListType that contains only Null
	// Its type can be set explicitly
	return Value{T: ListType{Contains: tr.Null()}, V: nil, N: ""}
}

func (tr TypeRegistry) MakeListOf(t Type) Value {
	// This makes a Value for a ListType that contains a specific type
	return Value{T: ListType{Contains: t}, V: nil, N: ""}
}

func (tr TypeRegistry) IsList(t Type) bool {
	if !t.Nullable() {
		t = t.(NonNullType).Contains
	}
	_, ok := t.(ListType)
	return ok
}

func (tr TypeRegistry) IsObj(t Type) bool {
	if !t.Nullable() {
		t = t.(NonNullType).Contains
	}
	_, ok := t.(ObjType)
	return ok
}

func (tr TypeRegistry) IsBaseType(t Type, typename string) bool {
	if !t.Nullable() {
		t = t.(NonNullType).Contains
	}
	bt, ok := t.(BaseType)
	if !ok {
		return false
	}
	return bt.Name == typename
}

func (tr TypeRegistry) IsInt(t Type) bool {
	return tr.IsBaseType(t, "Int")
}

func (tr TypeRegistry) IsFloat(t Type) bool {
	return tr.IsBaseType(t, "Float")
}

func (tr TypeRegistry) IsBool(t Type) bool {
	return tr.IsBaseType(t, "Bool")
}

func (tr TypeRegistry) IsNull(t Type) bool {
	return tr.IsBaseType(t, "*NULL*")
}

func (tr TypeRegistry) IsStr(t Type) bool {
	return tr.IsBaseType(t, "String")
}

func (tr TypeRegistry) IsID(t Type) bool {
	return tr.IsBaseType(t, "ID")
}

func (tr TypeRegistry) MakeNamelessObj() Value {
	// This makes a Value for a nameless, empty ObjType
	return Value{T: ObjType{}, V: nil, N: ""}
}

// searches the type registry for enums and then searches all enums for
// the appropriate value. It may be better to do this a) with an interface
// and b) with a more efficient search method
func (tr TypeRegistry) FindEnum(v string) (Value, error) {
	for _, t := range tr {
		if typ, ok := t.(EnumType); ok {
			for _, val := range typ.Values {
				if v == val {
					return Value{T: t, V: v, N: ""}, nil
				}
			}
		}
	}
	return Value{}, errors.New("'" + v + "' was not found in any enum.")
}

// BaseType is the fundamental unit of type.
// However, most cases won't store a BaseType, but an instance of the Type interface.
type BaseType struct {
	Name string
}

func (t BaseType) BaseType() string {
	return t.Name
}

func (t BaseType) Nullable() bool {
	return true
}

func (t BaseType) String() string {
	return t.BaseType()
}

func (t BaseType) Is(o Type) bool {
	return t.String() == o.String()
}

// ListType is a wrapper type to contain a list of other types.
type ListType struct {
	Contains Type
}

func (t ListType) BaseType() string {
	return t.Contains.BaseType()
}

// A ListType is nullable by default but could be wrapped in a NonNullType.
// It can also contain NonNullTypes.
func (t ListType) Nullable() bool {
	return true
}

func (t ListType) String() string {
	return "[" + t.BaseType() + "]"
}

func (t ListType) Is(o Type) bool {
	return t.String() == o.String()
}

// NonNullType is a wrapper that indicates that a type cannot be null.
type NonNullType struct {
	Contains Type
}

func (t NonNullType) BaseType() string {
	return t.Contains.BaseType()
}

func (t NonNullType) Nullable() bool {
	return false
}

func (t NonNullType) String() string {
	return t.BaseType() + "!"
}

func (t NonNullType) Is(o Type) bool {
	return t.String() == o.String()
}

// ObjType is a wrapper that contains a set of object fields, each of which
// is also a type. The fields are kept in an array (rather than a map) so
// that field order can be maintained.
type ObjType struct {
	Name     string
	Contains []Type
}

func (t ObjType) BaseType() string {
	return t.Name
}

func (t ObjType) Nullable() bool {
	return true
}

func (t ObjType) String() string {
	return "obj " + t.Name
}

func (t ObjType) Is(o Type) bool {
	return t.String() == o.String()
}

func (t ObjType) Add(element Type) {
	t.Contains = append(t.Contains, element)
}

// EnumType is a wrapper that contains a set of enum values, each of which
// is represented by its string name. The enum values are kept in an array.
type EnumType struct {
	Name   string
	Values []string
}

func (t EnumType) BaseType() string {
	return t.Name
}

func (t EnumType) Nullable() bool {
	return true
}

func (t EnumType) String() string {
	return "enum " + t.Name + "{" + strings.Join(t.Values, ", ") + "}"
}

func (t EnumType) Is(o Type) bool {
	return t.String() == o.String()
}

func (t EnumType) Add(element string) {
	t.Values = append(t.Values, element)
}
