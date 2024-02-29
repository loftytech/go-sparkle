package routes

import (
	router "coralscale/app/framework/core"
	"coralscale/app/framework/utility"
)

func Register() {

	router.Post("/v1/login", func(res *utility.Response, req *utility.Request) {

		res.Send(utility.CustomResponse{
			Status:  1,
			Message: "Login attempt initiated",
		}, 200)

	})

	router.Get("/v1/profile/:username/fetch", func(res *utility.Response, req *utility.Request) {
		res.Send(utility.CustomResponse{
			Status:  1,
			Message: "Profile fetched successfully",
			Data:    req.Params["username"],
		}, 200)
	})

}
