package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// Value contains actual values manipulated by GraphQL.
type Value struct {
	N string      // the name
	T Type        // the type of the data
	V interface{} // the actual data being represented by the value
}

func (v Value) String() string {
	return fmt.Sprintf("%s:%s %v", v.T.String(), v.N, v.V)
}

// Variables are a variant of value where the V parameter is always another Value object
type Variable Value

func (v Value) AsInt() (int, error) {
	switch {
	case typereg.IsInt(v.T):
		return v.V.(int), nil
	case typereg.IsFloat(v.T):
		return int(v.V.(float64)), nil
	case typereg.IsBool(v.T):
		if v.V.(bool) {
			return -1, nil
		} else {
			return 0, nil
		}
	case typereg.IsStr(v.T):
		i, err := strconv.ParseInt(v.V.(string), 10, 32)
		return int(i), err
	default:
		return 0, errors.New("Unable to convert to Int")
	}
}

func (v Value) AsFloat() (float64, error) {
	switch {
	case typereg.IsFloat(v.T):
		return v.V.(float64), nil
	case typereg.IsInt(v.T):
		return float64(v.V.(int64)), nil
	default:
		return 0, errors.New("Unable to convert to Float")
	}
}

func (v Value) AsBool() (bool, error) {
	switch {
	case typereg.IsBool(v.T):
		return v.V.(bool), nil
	case typereg.IsFloat(v.T):
		return v.V.(float64) != 0, nil
	case typereg.IsInt(v.T):
		return float64(v.V.(int64)) != 0, nil
	case typereg.IsStr(v.T):
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
	switch {
	case typereg.IsBool(v.T):
		if v.V.(bool) {
			return "true", nil
		} else {
			return "false", nil
		}
	case typereg.IsFloat(v.T):
		return fmt.Sprintf("%f", v.V), nil
	case typereg.IsInt(v.T):
		return fmt.Sprintf("%d", v.V), nil
	case typereg.IsStr(v.T):
		return v.V.(string), nil
	default:
		return "", errors.New("Unable to convert to String")
	}
}

func (v Value) AsList() ([]Value, error) {
	switch {
	case typereg.IsList(v.T):
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
	if len(item) == 0 {
		return
	}
	if !typereg.IsList(v.T) {
		panic("Attempt to append to a non-list!")
	}
	// see if we need to coerce our value type into whatever the list contains
	list := v.V.([]Value)
	for _, i := range item {
		if v.T.Is(i.T) {
		}
	}
	v.V = append(list, item...)
}

// An ObjectField is a Value where the N field is the name
// of the field, and the V field is a nested Value object
func (v Value) SetField(val Value) {
	if !typereg.IsObj(v.T) {
		panic("Attempt to SetField on a non-object!")
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
	if !typereg.IsObj(v.T) {
		panic("Attempt to Get from a non-object!")
	}
	items := v.V.([]Value)
	for i := range items {
		if items[i].N == key {
			return items[i], true
		}
	}
	return typereg.MakeNull(), false
}
