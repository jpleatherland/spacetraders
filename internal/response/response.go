package response

import (
	"encoding/json"
	"log"
	"net/http"
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
