package main

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/felipemarinho97/go-mock-gen/funcreader"
)

// Interface ...
type Interface struct {
	Name string
}

func main() {
	code := GenerateMockCode()

	fmt.Println(code)

}

func GenerateMockCode() string {
	typeToMock := s3.Client{}
	replacer := Interface{
		Name: "S3Client",
	}

	t := reflect.TypeOf(&typeToMock)
	var builder strings.Builder
	var b bytes.Buffer

	in := CreateMockInterface(t)
	struc := CreateMockStruct(t)
	fns := CreateMockedFunctions(t)

	fmt.Fprintf(&builder, "package clients\n\n")
	fmt.Fprintf(&builder, "/*\n* CODE GENERATED AUTOMATICALLY WITH github.com/felipemarinho97/go-mock-gen\n*/\n\n")
	fmt.Fprintf(&builder, in)
	fmt.Fprintf(&builder, struc)
	fmt.Fprintf(&builder, fns)

	template.Must(template.New("template").
		Parse(builder.String())).
		Execute(&b, replacer)

	return b.String()
}

// CreateMockInterface ...
func CreateMockInterface(t reflect.Type) string {
	var b strings.Builder
	b.WriteString("// I{{.Name}} generic client\ntype I{{.Name}} interface {\n")

	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		doc := funcreader.FuncDescription(m.Func)
		if doc != "" {
			re := regexp.MustCompile(`\n`)
			s := re.ReplaceAllString(doc, "\n\t// ")
			b.WriteString("\t// " + s + "\n")
		}
		fmt.Fprintf(&b, "\t%v%v %v\n", m.Name, GetInputSignature(m), GetOutputSignature(m))
	}

	b.WriteString("}\n")

	return b.String()
}

// CreateMockStruct ...
func CreateMockStruct(t reflect.Type) string {
	var b strings.Builder
	b.WriteString("// {{.Name}}Mock generic client mock\ntype {{.Name}}Mock struct {\n")

	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		fmt.Fprintf(&b, "\t%vMock\tfunc%v %v\n", m.Name, GetInputSignature(m), GetOutputSignature(m))
	}

	b.WriteString("}\n")

	return b.String()
}

// CreateMockedFunctions ...
func CreateMockedFunctions(t reflect.Type) string {
	var b strings.Builder

	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		b.WriteString(CreateMockedFunction(m) + "\n")
	}

	return b.String()
}

// CreateMockedFunction ...
func CreateMockedFunction(m reflect.Method) string {
	var b strings.Builder
	fmt.Fprintf(&b, "// %v mocked function\nfunc (m {{.Name}}Mock) %v%v%v {\n",
		m.Name, m.Name, GetInputSignature(m), GetOutputSignature(m))

	fmt.Fprintf(&b, "\treturn m.%vMock%v\n}\n", m.Name, GetInputParams(m))

	return b.String()
}

//GetInputSignature returns the input signature from a given method
func GetInputSignature(m reflect.Method) string {
	var b strings.Builder
	b.WriteString("(")

	for i := 1; i < m.Type.NumIn(); i++ {
		t := m.Type.In(i)

		if m.Type.IsVariadic() && m.Type.NumIn()-1 == i {
			fmt.Fprintf(&b, "arg%v ...%v", i, t.Elem())
		} else {
			fmt.Fprintf(&b, "arg%v %v", i, t.String())
		}

		if m.Type.NumIn()-1 != i {
			fmt.Fprintf(&b, ", ")
		}

	}
	b.WriteString(")")

	return b.String()
}

//GetOutputSignature returns the output signature from a given method
func GetOutputSignature(m reflect.Method) string {
	var b strings.Builder
	b.WriteString("(")

	for i := 0; i < m.Type.NumOut(); i++ {
		t := m.Type.Out(i)

		fmt.Fprintf(&b, t.String())

		if m.Type.NumOut()-1 != i {
			fmt.Fprintf(&b, ", ")
		}

	}
	b.WriteString(")")

	return b.String()
}

//GetInputParams returns the input signature from a given method
func GetInputParams(m reflect.Method) string {
	var b strings.Builder
	b.WriteString("(")

	for i := 1; i < m.Type.NumIn(); i++ {

		if m.Type.IsVariadic() && m.Type.NumIn()-1 == i {
			fmt.Fprintf(&b, "arg%v...", i)
		} else {
			fmt.Fprintf(&b, "arg%v,", i)
		}

	}
	b.WriteString(")")

	return b.String()
}
