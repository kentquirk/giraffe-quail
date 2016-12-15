package main

import (
	"fmt"

	gqtypes "github.com/kentquirk/giraffe-quail/types"
)

type HeroHandler struct {
	Name string
}

func (h HeroHandler) Operate(global, local *gqtypes.Scope) error {
	fmt.Println("Hero: " + h.Name)
	return nil
}
