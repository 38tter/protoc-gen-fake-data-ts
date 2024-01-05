package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-faker/faker/v4"
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

	str := protoMessage.GenerateStructForFaker(message).Interface()
	err := faker.FakeData(&str)

	if err != nil {
		fmt.Println(err)
	}

	strValue := reflect.ValueOf(str)
	strType := reflect.TypeOf(str)

	for i := 0; i < strValue.NumField(); i++ {
		field := strType.Field(i)
		fieldValue := strValue.Field(i)

		// TODO: Fix condition
		if !fieldValue.CanInt() && !fieldValue.CanFloat() && !fieldValue.CanUint() {
			v := fmt.Sprintf("\"%s\"", fieldValue.String())
			code = append(code, fmt.Sprintf("%s%s: %v,", strings.Repeat(" ", SpaceCharacterNum), field.Name, v))
		} else {
			code = append(code, fmt.Sprintf("%s%s: %v,", strings.Repeat(" ", SpaceCharacterNum), field.Name, fieldValue))
		}
	}

	code = append(code, "}")

	return code
}

func (protoMessage *ProtoMessage) GenerateStructForFaker(message *protogen.Message) reflect.Value {
	var fields []reflect.StructField

	for _, field := range message.Fields {
		str, err := mapProtoKindToGoTypes(field)
		if err != nil {
			// TODO: error handling
			continue
		}
		fields = append(fields, str)
	}
	structDef := reflect.New(reflect.StructOf(fields))

	return structDef.Elem()
}

func mapProtoKindToGoTypes(field *protogen.Field) (reflect.StructField, error) {
	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		return reflect.StructField{
			Name: field.GoName,
			Type: reflect.TypeOf(true),
			Tag:  "",
		}, nil
	case protoreflect.Int64Kind, protoreflect.Int32Kind:
		return reflect.StructField{
			Name: field.GoName,
			Type: reflect.TypeOf(0),
			Tag:  "",
		}, nil
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return reflect.StructField{
			Name: field.GoName,
			Type: reflect.TypeOf(0.0),
			Tag:  "",
		}, nil
	case protoreflect.StringKind:
		return reflect.StructField{
			Name: field.GoName,
			Type: reflect.TypeOf(""),
			Tag:  "",
		}, nil
	// TODO: Add more kinds
	default:
		return reflect.StructField{}, errors.New("invalid protoreflect type deteceted")
	}
}

func main() {
	var g = ProtoMessage{}
	protogen.Options{
		// TODO: Add some reguralization
		ParamFunc: nil,
	}.Run(g.Generate)
}
