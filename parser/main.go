package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
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
		for _, q := range queries {
			fmt.Println("parsing ", a)
			fmt.Printf("%s\n", q)
			_, err := Parse(a, q)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
}
