package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fadilmartias/firavel/app/utils"
	"github.com/fadilmartias/firavel/cmd/file_template"
	"github.com/jinzhu/inflection"
	"github.com/spf13/cobra"
)

var makeFileCmd = &cobra.Command{
	Use:   "make:file",
	Short: "Generate files (controller, model, seeder, etc)",
	Run: func(cmd *cobra.Command, args []string) {
		var choices = []string{
			"Controller",
			"Model",
			"Migration",
			"Seeder",
			"Factory",
			"Response",
		}
		selected := []string{}
		prompt := &survey.MultiSelect{
			Message: "Select file types to generate:",
			Options: choices,
		}
		survey.AskOne(prompt, &selected)

		var name string
		survey.AskOne(&survey.Input{
			Message: "Enter base name (without suffix):",
		}, &name)

		lowercaseName := strings.ToLower(name)
		for _, item := range selected {
			switch item {
			case "Controller":
				createFileFromTemplate("app/http/controllers/v1", utils.SnakeCase(lowercaseName)+"_controller.go", file_template.ControllerTemplate, map[string]string{
					"Name":       utils.StudlyCase(name) + "Controller",
					"LowerName":  lowercaseName,
					"SnakeName":  utils.SnakeCase(lowercaseName),
					"PluralName": inflection.Plural(utils.SnakeCase(lowercaseName)),
				})
			case "Model":
				createFileFromTemplate("app/models", utils.SnakeCase(lowercaseName)+"_model.go", file_template.ModelTemplate, map[string]string{
					"Name":       utils.StudlyCase(name),
					"LowerName":  lowercaseName,
					"SnakeName":  utils.SnakeCase(lowercaseName),
					"PluralName": inflection.Plural(utils.SnakeCase(lowercaseName)),
				})
			case "Seeder":
				createFileFromTemplate("database/seeders", utils.SnakeCase(lowercaseName)+"_seeder.go", file_template.SeederTemplate, map[string]string{
					"Name":       utils.StudlyCase(name),
					"LowerName":  lowercaseName,
					"SnakeName":  utils.SnakeCase(lowercaseName),
					"PluralName": inflection.Plural(utils.SnakeCase(lowercaseName)),
				})

			case "Migration":
				migrationName := fmt.Sprintf("create_%s_table", inflection.Plural(utils.SnakeCase(lowercaseName)))
				timestampStr := utils.Timestamp()
				filename := fmt.Sprintf("%s_%s.go", timestampStr, migrationName)
				createFileFromTemplate("database/migrations", filename, file_template.MigrationTemplate, map[string]string{
					"MigrationName": migrationName,
					"Timestamp":     timestampStr,
					"Name":          utils.StudlyCase(name),
					"LowerName":     lowercaseName,
					"SnakeName":     utils.SnakeCase(lowercaseName),
					"PluralName":    inflection.Plural(utils.SnakeCase(lowercaseName)),
				})
			case "Factory":
				createFileFromTemplate("database/factories", utils.SnakeCase(lowercaseName)+"_factory.go", file_template.FactoryTemplate, map[string]string{
					"Name":       utils.StudlyCase(name),
					"LowerName":  lowercaseName,
					"SnakeName":  utils.SnakeCase(lowercaseName),
					"PluralName": inflection.Plural(utils.SnakeCase(lowercaseName)),
				})
			case "Response":
				createFileFromTemplate("app/responses", utils.SnakeCase(lowercaseName)+"_response.go", file_template.ResponseTemplate, map[string]string{
					"Name":       utils.StudlyCase(name) + "Response",
					"LowerName":  lowercaseName,
					"SnakeName":  utils.SnakeCase(lowercaseName),
					"PluralName": inflection.Plural(utils.SnakeCase(lowercaseName)),
				})
			}
		}
	},
}

func createFileFromTemplate(dir, filename, tmplContent string, data any) {
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
