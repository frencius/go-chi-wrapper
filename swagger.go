package server

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"strings"

	"github.com/swaggo/swag"
)

// Path is a default swagger doc path
var Path = "docs/swagger.yaml"

// Swagger is a swagger instance
type Swagger struct{}

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
}

func init() {
	swag.Register(swag.Name, &Swagger{})
}

// ReadDoc got from `swag init` implementation with adjustment in reading swagger json or yaml
func (s *Swagger) ReadDoc() string {
	if doc, err := ioutil.ReadFile(Path); err != nil {
		return err.Error()
	} else {
		doc := string(doc)
		sInfo := SwaggerInfo
		sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

		t, err := template.New("swagger_info").Funcs(template.FuncMap{
			"marshal": func(v interface{}) string {
				a, _ := json.Marshal(v)
				return string(a)
			},
		}).Parse(doc)
		if err != nil {
			return doc
		}

		var tpl bytes.Buffer
		if err := t.Execute(&tpl, sInfo); err != nil {
			return doc
		}

		return tpl.String()
	}
}
