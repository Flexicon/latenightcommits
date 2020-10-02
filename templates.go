package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

// Template implements the echo.Renderer interface
type Template struct {
	templates *template.Template
}

// NewTemplateRenderer constructor
func NewTemplateRenderer() echo.Renderer {
	return &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
}

// Render the given template with data
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
