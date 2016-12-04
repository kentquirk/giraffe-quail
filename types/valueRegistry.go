package types

import "errors"

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
