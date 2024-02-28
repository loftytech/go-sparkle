package utility

type Request struct {
	method string
	path   string
	params map[string]string
}

func HandleSetRequest(method string, path string, params map[string]string) Request {
	return Request{
		method: method,
		params: params,
		path:   path,
	}
}
