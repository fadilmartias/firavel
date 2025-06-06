package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

// Grup perintah make
var makeCmd = &cobra.Command{
	Use:   "make",
	Short: "Commands to create new files",
}

var makeControllerCmd = &cobra.Command{
	Use:   "controller [name]",
	Short: "Create a new controller file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if !strings.HasSuffix(name, "Controller") {
			name += "Controller"
		}
		data := struct{ Name string }{Name: name}
		createFileFromTemplate("app/http/controllers", name+".go", controllerTemplate, data)
	},
}

func init() {
	makeCmd.AddCommand(makeControllerCmd)
}

// Helper untuk membuat file dari template
func createFileFromTemplate(dir, filename, tmplContent string, data interface{}) {
	path := filepath.Join(dir, filename)
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("Error: File %s already exists.\n", path)
		return
	}

	// Buat direktori jika belum ada
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	tmpl, err := template.New("file").Parse(tmplContent)
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		return
	}

	err = tmpl.Execute(file, data)
	if err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		return
	}

	fmt.Printf("File created successfully: %s\n", path)
}

// Isi template untuk controller
const controllerTemplate = `package controllers

import "github.com/gofiber/fiber/v2"

type {{.Name}} struct {
	BaseController
}

func New{{.Name}}() *{{.Name}} {
	return &{{.Name}}{}
}

// Index example method
func (ctrl *{{.Name}}) Index(c *fiber.Ctx) error {
	return ctrl.SuccessResponse(c, SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "Hello from {{.Name}}",
		Data:    nil,
	})
}
`
