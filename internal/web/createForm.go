package web

import (
	"fmt"
	"net/http"
	"text/template"
)

func createForm(rw http.ResponseWriter, _ *http.Request) {
	fmt.Println("create form hit")
	tmpl := template.Must(template.ParseFiles("views/createForm.html"))
	err := tmpl.Execute(rw, nil)
	if err != nil {
		http.Error(rw, "error rendering template", http.StatusInternalServerError)
		return
	}
}
