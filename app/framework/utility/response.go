package utility

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	ResponseWriter *http.ResponseWriter
}

type CustomResponse struct {
	Status  int8   `json:"status"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (r Response) Send(payload any, status_code int) {
	w := *r.ResponseWriter
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status_code)

	res, err := json.Marshal(&payload)

	if err != nil {
		log.Fatal("Decode error: ", err)
	}

	w.Write(res)
}
