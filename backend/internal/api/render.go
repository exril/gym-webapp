package api

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

var tmplEmbed embed.FS

func renderFiles(tmpl string, w http.ResponseWriter, d interface{}) {
	t, err := template.ParseFS(tmplEmbed, fmt.Sprintf("tmpl/%s.html", tmpl))
	if err != nil {
		log.Fatal(err)
	}

	if err := t.Execute(w, d); err != nil {
		log.Fatal(err)
	}
}
