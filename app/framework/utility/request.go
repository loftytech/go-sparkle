package utility

import (
	"net/http"
)

type Request struct {
	Method         string
	Path           string
	Params         map[string]string
	ResponseWriter *http.ResponseWriter
}

func HandleSetRequest(method string, path string, params map[string]string, w *http.ResponseWriter) Request {
	request := Request{
		Method:         method,
		Params:         params,
		Path:           path,
		ResponseWriter: w,
	}

	return request
}
