package utility

type Request struct {
	Method string
	Path   string
	Params map[string]string
}


func HandleSetRequest(method string, path string, params map[string]string) Request {
	request := Request{
		Method: method,
		Params: params,
		Path:   path,
	}


	return request
}
