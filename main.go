package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading input")
		os.Exit(1)
	}

	req := &plugin.CodeGeneratorRequest{}
	err = proto.Unmarshal(data, req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing input proto")
		os.Exit(1)
	}

	g := newGenerator(req)

	err = g.Generate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error generating output: %s\n", err)
		os.Exit(1)
	}

	data, err = proto.Marshal(g.Response)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error marshaling output proto")
		os.Exit(1)
	}

	_, err = os.Stdout.Write(data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error writing output")
		os.Exit(1)
	}
}
