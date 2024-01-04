package main

import (
	"strings"

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

		fileName := strings.Replace(file.Desc.Path(), ".proto", "", 1)

		generatedFileName := fileName + "_fake.ts"
		generatedFilePath := protogen.GoImportPath(file.Desc.Path())

		t := plugin.NewGeneratedFile(generatedFileName, generatedFilePath)

		t.P("hoge\nfuga")
	}

	return nil
}

func main() {
	var g = ProtoMessage{}
	protogen.Options{
		ParamFunc: nil,
	}.Run(g.Generate)
}
