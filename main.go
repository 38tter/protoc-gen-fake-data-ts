package main

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ProtoMessage struct {
	messages []protoreflect.FullName
}

func (protoMessage *ProtoMessage) Generate(plugin *protogen.Plugin) error {
	protoFiles := plugin.Files
	for _, file := range protoFiles {
		if len(file.Services) == 0 {
			continue
		}
	}

	return nil
}

func main() {
	var g = ProtoMessage{}
	protogen.Options{
		ParamFunc: nil,
	}.Run(g.Generate)
}
