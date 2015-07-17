package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

var tmplt = template.Must(template.ParseFiles(fmt.Sprintf("%s/static/test.html", os.Getenv("ROOT_PATH"))))

type Page struct {
	Title string
	Body  []byte
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func a() { fmt.Println("a") }
