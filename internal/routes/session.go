package routes

import (
	"log"
	"net/http"
	"strings"
	"time"
)

func SetSession(rw http.ResponseWriter, req *http.Request) {
	symbol := strings.TrimPrefix(req.URL.Path, "/products/")
	log.Println("in set session", symbol)
	time.Sleep(5 * time.Second)
}
