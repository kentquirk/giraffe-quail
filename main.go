package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kentquirk/giraffe-quail/parser"
)

func main() {
	// limit := flag.Int("limit", 0, "The max number of items to post.")
	// importid := flag.String("importid", "", "Set an import ID string to be added to the tags of each passage imported in this session.")
	flag.Parse()

	for _, a := range flag.Args() {
		b, err := ioutil.ReadFile(a)
		if err != nil {
			fmt.Println("couldn't read ", a)
			continue
		}

		queries := bytes.Split(b, []byte("\n\n"))
		hadErrors := 0
		for _, q := range queries {
			fmt.Println("parsing ", a)
			fmt.Printf("%s\n", q)
			_, err := parser.Parse(a, q)
			if err != nil {
				parser.DumpErrors(err)
				hadErrors = 1
			}
		}
		os.Exit(hadErrors)
	}
}
