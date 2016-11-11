package parser

import (
	"errors"
	"fmt"
)

func (c *current) NewError(s string, code string) error {
	return errors.New(fmt.Sprintf("Parser error %s (%s) @ %d:%d",
		s, code, c.pos.line, c.pos.col))
}

func DumpErrors(err error) {
	list := err.(errList)
	for _, err := range list {
		pe := err.(*parserError)
		fmt.Printf("%+v\n", pe)
	}
}
