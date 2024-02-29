package core

import (
	"coralscale/app/framework/utility"
	"fmt"
	"net/http"
	"strings"
)

type RequestRoute struct {
	method             string
	path               string
	filteredParamsPath string
	params             []string
	paramLength        int
}

var requestData = RequestRoute{}

func (r *RequestRoute) resolveRequestRoute(w *http.ResponseWriter) {

	for _, v := range route_list {
		fmt.Println(v)

		splited_path := strings.Split(r.path[1:], "/")
		req_path_len := len(splited_path)

		fmt.Println("r.method: ", r.method)
		fmt.Println("v.method: ", v.method)

		if req_path_len != v.pathArrLen || r.method != v.method {
			continue
		}

		param_map := map[string]string{}
		filteredParamsPath := ""
		// params := []string{}
		for idx, value := range splited_path {
			param, ok := v.params[idx]
			if ok {
				param_map[param] = value
			} else {
				filteredParamsPath += "/" + value
			}
		}

		fmt.Println("req filteredParamsPath: ", filteredParamsPath)
		fmt.Println("route filteredParamsPath: ", v.filteredParamsPath)

		if filteredParamsPath != v.filteredParamsPath {
			continue
		}

		// for idx, value := range v.params {
		// 	param_map[value] = splited_path[idx]
		// }

		request := utility.HandleSetRequest(v.method, v.path, param_map, w)
		response := utility.Response{ResponseWriter: w}

		v.handler(&response, &request)

		// fmt.Println("request match: ", request)
	}
}

type Route struct {
	method             string
	path               string
	pathArrLen         int
	filteredParamsPath string
	params             map[int]string
	handler            func(*utility.Response, *utility.Request)
}

func (r *Route) resolveRoute() {
	path := r.path
	fmt.Println("resolveRoute path: ", path)
	splited_path := strings.Split(path[1:], "/")
	filteredParamsPath := ""
	params := map[int]string{}
	for idx, v := range splited_path {
		if len(v) > 0 && v[0:1] == ":" {
			if len(v) > 0 {
				params[idx] = v[1:]
			} else {
				params[idx] = ""
			}
		} else {
			filteredParamsPath += "/" + v
		}
	}

	r.pathArrLen = len(splited_path)
	r.params = params
	// r.paramLength = len(params)
	r.filteredParamsPath = filteredParamsPath
}

var route_list = []Route{}

func Get(path string, handler func(response *utility.Response, request *utility.Request)) {
	validateRoute("GET", path, handler)
}

func Post(path string, handler func(response *utility.Response, request *utility.Request)) {
	validateRoute("POST", path, handler)
}

func Patch(path string, handler func(response *utility.Response, request *utility.Request)) {
	validateRoute("PATCH", path, handler)
}

func Put(path string, handler func(response *utility.Response, request *utility.Request)) {
	validateRoute("PUT", path, handler)
}

func Delete(path string, handler func(response *utility.Response, request *utility.Request)) {
	validateRoute("DRLETE", path, handler)
}

func validateRoute(method string, path string, handler func(response *utility.Response, request *utility.Request)) {
	route := Route{method: method, path: path, handler: handler}
	route.resolveRoute()
	route_list = append(route_list, route)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	requestData = RequestRoute{
		method:             r.Method,
		path:               r.URL.Path,
		filteredParamsPath: "",
		params:             []string{},
		paramLength:        0,
	}
	requestData.resolveRequestRoute(&w)

	// fmt.Println("requestData: ", requestData)

	// fmt.Fprintf(w, "Current route %s!", r.URL.Path[1:])

}
