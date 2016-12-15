package main

import (
	"fmt"
	"os"

	"github.com/kentquirk/giraffe-quail/gql"
)

func main() {
	gq, err := gql.FromFile("../../parser/tests/starwars.schema")
	if err != nil {
		fmt.Println("Error parsing schema." + err.Error())
		os.Exit(1)
	}
	query, err := gq.Types.Get("Query")
	for _, f := range query.Fields {
		fmt.Println("f = ", f.N)
	}
	fmt.Println("s = ", gq.Scope)

	qs := `query HeroNameQuery {
              hero {
                name
              }
            }
         `
	ops, err := gq.ParseString("querytest", qs)
	if err != nil {
		fmt.Println("Error parsing query." + err.Error())
		os.Exit(1)
	}
	gq.Register("hero", HeroHandler{"Luke"})
	for _, op := range ops {
		if err := gq.DoOp(op); err != nil {
			fmt.Println(err)
		}
	}
}
