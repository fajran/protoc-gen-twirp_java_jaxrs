package main

import (
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic("error reading input")
	}

	req := &plugin.CodeGeneratorRequest{}
	err = proto.Unmarshal(data, req)
	if err != nil {
		panic("error parsing input proto")
	}

	g := newGenerator(req)

	err = g.Generate()
	if err != nil {
		panic("error generating output")
	}

	data, err = proto.Marshal(g.Response)
	if err != nil {
		panic("error marshaling output proto")
	}

	_, err = os.Stdout.Write(data)
	if err != nil {
		panic("error writing output")
	}
}
