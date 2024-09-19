package response

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

func RespondWithJSON(rw http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling json: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(code)
	rw.Write(data)
}

func RespondWithError(rw http.ResponseWriter, msg string, code int) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	log.Println(msg)
	RespondWithJSON(rw, code, errorResponse{
		Error: msg,
	})
}

func RespondWithHTML(rw http.ResponseWriter, html string, code int) {
	data := []byte(html)
	rw.Header().Set("Content-Type", "text/html")
	rw.WriteHeader(code)
	rw.Write(data)
}

func RespondWithHTMLError(rw http.ResponseWriter, error string, code int) {
	errMsg := "<p>" + error + "</p>"
	log.Println(errMsg)
	RespondWithHTML(rw, errMsg, code)
}

func RespondWithTemplate(rw http.ResponseWriter, templateName string, templateData interface{}) {
	tmpl := template.Must(template.ParseGlob(filepath.Join("views", "templates", templateName)))
	err := tmpl.Execute(rw, templateData)
	if err != nil {
		RespondWithHTMLError(rw, err.Error(), http.StatusInternalServerError)
	}
}

func RespondWithPartialTemplate(rw http.ResponseWriter, partialFolderName, templateName string, templateData interface{}, funcMap template.FuncMap) {
	tmpl := template.New(templateName).Funcs(funcMap)
	tmpl = template.Must(tmpl.ParseGlob(filepath.Join("views", "templates", partialFolderName, templateName)))
	err := tmpl.Execute(rw, templateData)
	if err != nil {
		RespondWithHTMLError(rw, err.Error(), http.StatusInternalServerError)
	}
}

