package main

import (
	"fmt"

	"github.com/mjbozo/mdx"
)

func main() {
	config := &mdx.GeneratorConfig{InputFilename: "sample.mdx", OutputFilename: "sample.html"}
	n, err := mdx.Generate(config)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d bytes written to sample.html\n", n)
}
