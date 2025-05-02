package templates

import (
	"embed"
	"html/template"
)

//go:embed "*.gohtml"
var templateFS embed.FS

type Templates interface {
	GetEventListTemplate() (*template.Template, error)
	GetEventDetailsTemplate() (*template.Template, error)
	GetCreateEventTemplate() (*template.Template, error)
}

type templates struct {
}

func New() Templates {
	return &templates{}
}

func (t *templates) GetEventListTemplate() (*template.Template, error) {
	return template.ParseFS(templateFS, "events-list.gohtml")
}

func (t *templates) GetEventDetailsTemplate() (*template.Template, error) {
	return template.ParseFS(templateFS, "event-details.gohtml")
}

func (t *templates) GetCreateEventTemplate() (*template.Template, error) {
	return template.ParseFS(templateFS, "create-event.gohtml")
}
