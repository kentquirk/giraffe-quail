package parser

import (
	"io/ioutil"

	"github.com/kentquirk/giraffe-quail/types"
)

var TR *types.TypeRegistry
var GlobalScope *types.Scope

// ParseString parses a query string.
// if error is non-nil, it's an errList from the parser.
func LoadQueryFromString(s string, scope *types.Scope) error {
	_, err := Parse("string", []byte(s))
	return err
}

// ParseFile parses a query from a file.
// Could return a normal error if file couldn't be read, or an errList from the parser
// if the file was good but had errors.
func LoadQueryFromFile(filename string) error {

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	_, err = Parse(filename, b)
	return err
}
