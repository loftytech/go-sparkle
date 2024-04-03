package utility

import (
	"net/http"
)

type Request struct {
	Method         string
	Body          []byte
	Path           string
	Params         map[string]string
	ResponseWriter *http.ResponseWriter
}

func HandleSetRequest(method string, path string, params map[string]string, w *http.ResponseWriter, body []byte) Request {
	request := Request{
		Method:         method,
		Params:         params,
		Path:           path,
		ResponseWriter: w,
		Body: body,
	}

	return request
}
