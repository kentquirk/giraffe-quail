package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kentquirk/giraffe-quail/gql"
)

func main() {
	gq, err := gql.FromFile("../../parser/tests/status.schema")
	if err != nil {
		fmt.Println("Error parsing schema." + err.Error())
		os.Exit(1)
	}
	query, err := gq.Types.Get("Query")
	for _, f := range query.Fields {
		fmt.Println("f = ", f.N)
	}
	fmt.Println("s = ", gq.Scope)

	qs := `{
              status {
                Name
                revision
              }
            }
         `
	ops, err := gq.ParseString("querytest", qs)
	if err != nil {
		fmt.Println("Error parsing query." + err.Error())
		os.Exit(1)
	}
	h := NewStatusHandler()
	gq.Register("status", h)
	err = gq.DoOps(ops)
	if err != nil {
		fmt.Println("Error performing operations." + err.Error())
		os.Exit(1)
	}

	data, err := json.MarshalIndent(h.Data, "", "  ")
	// buf := &bytes.Buffer{}
	// json.NewEncoder(buf).Encode(h.Data)
	fmt.Println(string(data))
}
