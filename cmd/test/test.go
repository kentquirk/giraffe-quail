package main

import (
	"fmt"
	"os"

	"github.com/kentquirk/giraffe-quail/parser"
	"github.com/kentquirk/giraffe-quail/types"
	"github.com/kentquirk/giraffe-quail/typeschema"
)

func main() {
	// parser.TR = types.NewTypeRegistry()
	// parser.GlobalScope = types.NewScope()
	var err error

	tr, gs, err := typeschema.LoadSchemaFromFile("../../parser/tests/starwars.schema")
	if err != nil {
		fmt.Println("Error parsing schema." + err.Error())
		os.Exit(1)
	}
	query, err := tr.Get("Query")
	for _, f := range query.Fields {
		fmt.Println("f = ", f.N)
	}
	fmt.Println("gs = ", gs)

	qs := `query HeroNameQuery {
              hero {
                name
              }
            }
         `
	qi, err := parser.Parse("querytest", []byte(qs))
	if err != nil {
		fmt.Println("Error parsing query." + err.Error())
		os.Exit(1)
	}
	qia := qi.([]interface{})
	for _, q := range qia {
		ss := q.(types.SelectionSet)
		fmt.Println("ss = ", ss.Fields[0].Name)
	}
}
