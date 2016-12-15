package typeschema

import (
	"fmt"
	"io"

	"github.com/kentquirk/giraffe-quail/types"
)

var TR *types.TypeRegistry
var GlobalScope *types.Scope

// LoadString parses a schema string and turns it into type and value registries.
// It resets the TR and GlobalScope global pointers to new values and returns them.
// if error is non-nil, it's an errList from the parser.
func LoadString(s string) (*types.TypeRegistry, *types.Scope, error) {
	TR = types.NewTypeRegistry()
	GlobalScope = types.NewScope()

	_, err := Parse("string", []byte(s))
	return TR, GlobalScope, err
}

// LoadReader parses a schema string and turns it into type and value registries.
// It resets the TR and GlobalScope global pointers to new values and returns them.
// if error is non-nil, it's an errList from the parser.
// The name parameter is used in error messages
func LoadReader(name string, r io.Reader) (*types.TypeRegistry, *types.Scope, error) {
	TR = types.NewTypeRegistry()
	GlobalScope = types.NewScope()

	_, err := ParseReader(name, r)
	return TR, GlobalScope, err
}

// LoadFile parses a schema from a file and turns it into type and value registries.
// It resets the TR and GlobalScope global pointers to new values and returns them.
// Could return a normal error if file couldn't be read, or an errList from the parser
// if the file was good but had errors.
func LoadFile(filename string) (*types.TypeRegistry, *types.Scope, error) {
	TR = types.NewTypeRegistry()
	GlobalScope = types.NewScope()

	_, err := ParseFile(filename)
	return TR, GlobalScope, err
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
