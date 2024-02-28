package core

import (
	"log"
	"net/http"
)

func InitWebService(handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
