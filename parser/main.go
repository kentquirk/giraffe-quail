package main

import (
	"flag"
	"fmt"
	"io/ioutil"
)

func main() {
	// limit := flag.Int("limit", 0, "The max number of items to post.")
	// importid := flag.String("importid", "", "Set an import ID string to be added to the tags of each passage imported in this session.")
	flag.Parse()

	for _, a := range flag.Args() {
		bytes, err := ioutil.ReadFile(a)
		if err != nil {
			fmt.Println("couldn't read ", a)
			continue
		}
		fmt.Println("parsing ", a)
		r, err := Parse(a, bytes)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%#v\n", r)
	}
}
