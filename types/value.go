package types

import (
	"errors"
	"fmt"
)

// Value is a container for a value in the system
// It always has a type, and the V member will vary based on the type of the
// value being manipulated.
// T        V
// List     []<listType>
// Obj      map[string]Value, where the string is the field name
//          Objs are always rendered in order of the field names found in the Type
type Value struct {
	T Type
	V interface{}
}

type NamedValue struct {
	N string
	V Value
}

func (v *Value) String() string {
	return fmt.Sprintf("%v <%s>", v.V, v.T.String())
}

func (v *Value) Append(items ...Value) error {
	if len(items) == 0 {
		return nil
	}
	if v.T.Kind != List {
		return errors.New("Attempt to append to a non-list!")
	}
	// see if we need to convert the type of our list
	if len(v.T.Subtypes) == 0 {
		v.T.Subtypes = append(v.T.Subtypes, items[0].T)
	}

	// see if we need to coerce our value type into whatever the list contains
	list := v.V.([]Value)
	for _, i := range items {
		if v.T.Subtypes[0].Is(i.T) {
			list = append(list, i)
		} else {
			return errors.New(
				fmt.Sprintf("List requires '%s' but '%s' is not the proper type.",
					v.T.Subtypes[0].String(),
					v.String()))
		}
	}
	v.V = list
	return nil
}

func (v *Value) SetField(name string, val Value) error {
	if v.T.Kind != Obj {
		return errors.New("Attempt to SetField on a non-object!")
	}
	m := v.V.(map[string]Value)
	// if it's an anonymous object, just set the field
	if v.T.Name == "" {
		m[name] = val
		return nil
	}
	// but if we have an actual type, validate the type
	f, err := v.T.GetField(name)
	if err != nil {
		return err
	}
	if !f.T.Is(val.T) {
		return errors.New(fmt.Sprintf("Incompatible types: <%s> and <%s>",
			f.T.String(), val.T.String()))
	}
	m[name] = val
	return nil
}

// Set sets a value from another value, if it's compatible
func (v *Value) Set(other Value) error {
	if v.T.Key() == other.T.Key() {
		v.V = other.V
		return nil
	}

	// needs a lot more here

	return errors.New(fmt.Sprintf("Incompatible types: <%s> and <%s>",
		v.T.String(), other.T.String()))
}
