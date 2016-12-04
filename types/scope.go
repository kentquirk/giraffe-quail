package types

import "errors"

type Scope struct {
	Values map[string]Value
}

func NewScope() *Scope {
	sc := &Scope{Values: make(map[string]Value)}
	return sc
}

// Create sets up a value in the scope. It is an error if the name
// already exists.
func (sc *Scope) Create(name string, v Value) (Value, error) {
	if _, found := sc.Values[name]; found {
		return v, errors.New("Value " + name + " already defined; cannot be overridden.")
	}
	sc.Values[name] = v
	return v, nil
}

// Get retrieves a value from the registry by name. It returns an error
// if the value does not exist
func (sc *Scope) Get(name string) (Value, error) {
	t, ok := sc.Values[name]
	if ok {
		return t, nil
	}
	return t, errors.New("Value " + name + " was not found in the scope.")
}

// MustGet retrieves a value from the registry by name. It panics if the value
// does not exist. Mainly intended for internal use.
func (sc *Scope) MustGet(name string) Value {
	if t, ok := sc.Values[name]; ok {
		return t
	}
	panic("Value " + name + " was not found in the scope.")
}
