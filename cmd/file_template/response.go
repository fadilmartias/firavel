package file_template

var ResponseTemplate = `package responses

type {{ .Name }} struct {
	ID        *int    ` + "`json:\"id\"`" + `
	Name      *string ` + "`json:\"name\"`" + `
	CreatedAt *string ` + "`json:\"created_at\"`" + `
	UpdatedAt *string ` + "`json:\"updated_at\"`" + `
}
`
