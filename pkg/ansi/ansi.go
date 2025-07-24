package ansi

import (
	"bytes"
	"fmt"
	"text/template"
)

type Data struct {
	Fields []Field
}

type Field struct {
	Name         string
	DisplayValue string
	Color        string
}

func OutputKeyVal(data Data) string {
	// This uses the ANSI tags specified in tview: https://github.com/rivo/tview/wiki/ANSI
	tmplStr := `{{range .Fields}}{{if eq .Color "red"}}{{printf "[red]%s:[-] " .Name}}{{printf "[red]%s[-]" .DisplayValue}}{{else if eq .Color "green"}}{{printf "[green]%s:[-] " .Name}}{{printf "[green]%s[-]" .DisplayValue}}{{else}}{{printf "%s:%s" .Name .DisplayValue}}{{end}}
{{end}}`

	tmpl, err := template.New("output").Parse(tmplStr)
	if err != nil {
		panic(err)
	}

	var output bytes.Buffer
	if err := tmpl.Execute(&output, data); err != nil {
		panic(err)
	}

	return output.String()
}

func NewField(name string, value string, color string) Field {
	return Field{
		Name:         name,
		DisplayValue: value,
		Color:        color,
	}
}

func CreateField(name string, value any) Field {
	var displayValue string
	var color string

	switch v := value.(type) {
	case string:
		if v == "" {
			displayValue = "-"
			color = "red"
		} else {
			displayValue = v
			color = "green"
		}
	case bool:
		if v {
			displayValue = "true"
			color = "green"
		} else {
			displayValue = "false"
			color = "red"
		}
	default:
		displayValue = fmt.Sprint(v)
		color = "default"
	}

	return NewField(name, displayValue, color)
}

func SimpleString(value any) string {
	field := CreateField("", value)
	return fmt.Sprintf("[%s]%s[-]", field.Color, field.DisplayValue)
}
