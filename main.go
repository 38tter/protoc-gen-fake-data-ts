package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const SpaceCharacterNum = 4

type ProtoMessage struct {
	messages []protoreflect.FullName
}

func (protoMessage *ProtoMessage) Generate(plugin *protogen.Plugin) error {
	protoFiles := plugin.Files
	for _, file := range protoFiles {
		fileName := strings.Replace(file.Desc.Path(), ".proto", "", 1)

		generatedFileName := fileName + "_fake.ts"
		generatedFilePath := protogen.GoImportPath(file.Desc.Path())

		t := plugin.NewGeneratedFile(generatedFileName, generatedFilePath)

		var code []string

		for _, message := range file.Messages {
			code = append(code, protoMessage.GenerateFakeDataClass(message)...)
		}
		t.P(strings.Join(code[:], "\n"))
	}

	return nil
}

func (protoMessage *ProtoMessage) GenerateFakeDataClass(message *protogen.Message) []string {
	var code = []string{
		fmt.Sprintf("export const %s = {", strcase.ToLowerCamel(string(message.Desc.Name()))),
	}
	for _, field := range message.Fields {
		code = append(
			code,
			fmt.Sprintf("%s%s", strings.Repeat(" ", SpaceCharacterNum), strcase.ToLowerCamel(field.GoName)),
		)
		_ = protoMessage.GenerateStructForFaker(message)
	}
	code = append(code, "}")

	return code
}

func (protoMessage *ProtoMessage) GenerateStructForFaker(message *protogen.Message) reflect.Value {
	var fields []reflect.StructField

	for _, field := range message.Fields {
		fields = append(fields, reflect.StructField{
			Name: field.GoName,
			//Type: reflect.TypeOf(field.Desc.Kind().String()),
			Type: reflect.TypeOf(""),
			Tag:  "",
		})
	}
	structDef := reflect.New(reflect.StructOf(fields))

	return structDef.Elem()
}

func main() {
	var g = ProtoMessage{}
	protogen.Options{
		// TODO: Add some reguralization
		ParamFunc: nil,
	}.Run(g.Generate)
}
