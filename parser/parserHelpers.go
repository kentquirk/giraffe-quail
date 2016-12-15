package parser

import (
	"errors"
	"fmt"
)

func (c *current) NewError(s string, code string) error {
	return errors.New(fmt.Sprintf("Parser error %s (%s) @ %d:%d",
		s, code, c.pos.line, c.pos.col))
}

func (c *current) WrapError(e error, code string) error {
	return errors.New(fmt.Sprintf("Parser error %s (%s) @ %d:%d",
		e.Error(), code, c.pos.line, c.pos.col))
}

func DumpErrors(err error) {
	list := err.(errList)
	for _, err := range list {
		pe := err.(*parserError)
		fmt.Printf("%+v\n", pe)
	}
}

// Str tries to interpret an interface{} as a string;
// if it succeeds, returns the string, if it fails, returns ""
func Str(s interface{}) string {
	if st, ok := s.(string); ok {
		return st
	}
	return ""
}
