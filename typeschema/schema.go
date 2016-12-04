package typeschema

import (
	"fmt"
	"io/ioutil"

	"github.com/kentquirk/giraffe-quail/types"
)

var TR *types.TypeRegistry
var VR *types.ValueRegistry

// ParseString parses a schema string and turns it into type and value registries.
// It resets the TR and VR global pointers to new values and returns them.
// if error is non-nil, it's an errList from the parser.
func LoadSchemaFromString(s string) (*types.TypeRegistry, *types.ValueRegistry, error) {
	TR = types.NewTypeRegistry()
	VR = types.NewValueRegistry()

	_, err := Parse("string", []byte(s))
	return TR, VR, err
}

// ParseFile parses a schema from a file and turns it into type and value registries.
// It resets the TR and VR global pointers to new values and returns them.
// Could return a normal error if file couldn't be read, or an errList from the parser
// if the file was good but had errors.
func LoadSchemaFromFile(filename string) (*types.TypeRegistry, *types.ValueRegistry, error) {
	TR = types.NewTypeRegistry()
	VR = types.NewValueRegistry()

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}
	_, err = Parse(filename, b)
	return TR, VR, err
}

func DumpErrors(err error) {
	switch e := err.(type) {
	case errList:
		for _, err := range e {
			pe := err.(*parserError)
			fmt.Printf("%+v\n", pe)
		}
	default:
		fmt.Println(e)
	}
}
