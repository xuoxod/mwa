package models

import (
	"github.com/xuoxod/mwa/internal/forms"
)

// Holds data sent from handler to template
type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
	Flash           string
	Warning         string
	Error           string
	Form            *forms.Form
	IsAuthenticated int
	IsAdmin         int
}
