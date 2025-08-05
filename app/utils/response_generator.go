package utils

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strings"
	"text/template"
)

func GenerateResponseStructToFile(filename, structName, outputPath string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.AllErrors)
	if err != nil {
		return err
	}

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec := spec.(*ast.TypeSpec)
			if typeSpec.Name.Name != structName {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			fields := []Field{}

			for _, field := range structType.Fields.List {
				if len(field.Names) == 0 {
					continue
				}

				name := field.Names[0].Name
				typeStr := resolveType(field.Type)
				jsonTag := fmtJsonTag(name)

				if !strings.HasPrefix(typeStr, "*") {
					typeStr = "*" + typeStr
				}
				fields = append(fields, Field{
					Name:    name,
					Type:    typeStr,
					JSONTag: jsonTag,
				})
			}

			return writeToFile(outputPath, structName+"Response", fields)
		}
	}

	return nil
}

type Field struct {
	Name    string
	Type    string
	JSONTag string
}

func resolveType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return resolveType(t.X)
	case *ast.SelectorExpr:
		return resolveType(t.X) + "." + t.Sel.Name
	default:
		return "any"
	}
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func fmtJsonTag(name string) string {
	if name == "ID" {
		return "id"
	}
	snake := matchFirstCap.ReplaceAllString(name, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func writeToFile(path, structName string, fields []Field) error {
	tmpl := `package responses

import "time"

type {{.StructName}} struct {
{{- range .Fields }}
	{{ .Name }} {{ .Type }} ` + "`json:\"{{ .JSONTag }},omitempty\"`" + `
{{- end }}
}
`

	data := struct {
		StructName string
		Fields     []Field
	}{
		StructName: structName,
		Fields:     fields,
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	t := template.Must(template.New("response").Parse(tmpl))
	return t.Execute(f, data)
}
