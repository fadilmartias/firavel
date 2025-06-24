package commands

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/fadilmartias/firavel/app/utils"
	"github.com/spf13/cobra"
)

var generateAll bool

var makeResponseCmd = &cobra.Command{
	Use:   "make:response [ModelName]",
	Short: "Generate response struct from a model or all models using --all",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if generateAll {
			files, err := os.ReadDir("app/models")
			if err != nil {
				fmt.Println("❌ Failed to read model directory:", err)
				return
			}

			for _, file := range files {
				if strings.HasSuffix(file.Name(), ".go") {
					modelFile := filepath.Join("app/models", file.Name())
					modelName := getStructNameFromFile(modelFile)
					if modelName == "" {
						continue
					}
					outDir := "app/responses"
					os.MkdirAll(outDir, 0755)
					outPath := filepath.Join(outDir, utils.SnakeCase(modelName)+"_response.go")
					err := utils.GenerateResponseStructToFile(modelFile, modelName, outPath)
					if err != nil {
						fmt.Printf("❌ Failed to generate %s: %v\n", modelName, err)
					} else {
						fmt.Printf("✅ Generated: %s\n", outPath)
					}
				}
			}
			return
		}

		if len(args) == 0 {
			fmt.Println("❌ Model name required if --all not set")
			return
		}

		modelName := args[0]
		modelFile := findModelFile(modelName)
		if modelFile == "" {
			fmt.Printf("❌ Model file for %s not found.\n", modelName)
			return
		}

		outDir := "app/responses"
		os.MkdirAll(outDir, 0755)
		outPath := filepath.Join(outDir, utils.SnakeCase(modelName)+"_response.go")

		err := utils.GenerateResponseStructToFile(modelFile, modelName, outPath)
		if err != nil {
			fmt.Printf("❌ Failed to generate response struct: %v\n", err)
			return
		}
		fmt.Printf("✅ Generated: %s\n", outPath)
	},
}

func findModelFile(modelName string) string {
	files, err := os.ReadDir("app/models")
	if err != nil {
		fmt.Println("Failed to read model directory:", err)
		return ""
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".go") &&
			strings.Contains(strings.ToLower(file.Name()), strings.ToLower(modelName)) {
			return filepath.Join("app/models", file.Name())
		}
	}
	return ""
}

func getStructNameFromFile(filename string) string {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.AllErrors)
	if err != nil {
		return ""
	}
	for _, decl := range node.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if _, ok := typeSpec.Type.(*ast.StructType); ok {
						return typeSpec.Name.Name
					}
				}
			}
		}
	}
	return ""
}

func init() {
	makeResponseCmd.Flags().BoolVar(&generateAll, "all", false, "Generate response structs for all models")
}
