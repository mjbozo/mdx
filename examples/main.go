package main

import (
	"fmt"

	"github.com/mjbozo/mdx"
)

func main() {
	config := &mdx.GeneratorConfig{
		Title:          "MDX Sample",
		InputFilename:  "sample.mdx",
		OutputFilename: "sample.html",
		Links: []map[string]string{{
			"rel":  "stylesheet",
			"href": "sample.css",
		}},
	}

	n, err := mdx.Generate(config)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d bytes written to sample.html\n", n)
}
