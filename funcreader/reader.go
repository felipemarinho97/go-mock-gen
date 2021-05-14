package funcreader

import (
	"go/doc"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

// FuncPathAndName Get the name and path of a func
func FuncPathAndName(f reflect.Value) string {
	return runtime.FuncForPC(f.Pointer()).Name()
}

// FuncName Get the name of a func (with package path)
func FuncName(f reflect.Value) string {
	splitFuncName := strings.Split(FuncPathAndName(f), ".")
	return splitFuncName[len(splitFuncName)-1]
}

// FuncDescription Get description of a func
func FuncDescription(f reflect.Value) string {
	fileName, _ := runtime.FuncForPC(f.Pointer()).FileLine(0)
	funcName := FuncName(f)
	fset := token.NewFileSet()

	dir := filepath.Dir(fileName)

	// Parse src
	parsedAst, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	importPath, _ := filepath.Abs("/")

	for _, pkg := range parsedAst {

		myDoc := doc.New(pkg, importPath, doc.AllDecls)
		for _, theFunc := range myDoc.Funcs {
			if theFunc.Name == funcName {
				return theFunc.Doc
			}
		}
		for _, t := range myDoc.Types {
			for _, theFunc := range t.Methods {
				if theFunc.Name == funcName {
					return theFunc.Doc
				}
			}
		}
	}
	return ""
}
