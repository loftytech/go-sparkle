package routes

import (
	router "coralscale/app/framework/core"
	"coralscale/app/framework/utility"
	"fmt"
	"net/http"
)

func Register() {

	router.Post("/v1/login", func(w *http.ResponseWriter, request *utility.Request) {

		fmt.Fprintf(*w, "Current route login")

	})

	router.Get("/v1/profile/:username/fetch", func(w *http.ResponseWriter, req *utility.Request) {
		request := &req

		fmt.Fprintf(*w, "Current route profile, %s", request)

	})

}
