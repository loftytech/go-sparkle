package routes

import (
	router "coralscale/app/framework/core"
	"coralscale/app/framework/utility"
	"fmt"
)

func Register() {

	router.Post("/v1/login", func(res *utility.Response, req *utility.Request) {

		fmt.Println("accessed path: " + req.Path)

		res.Send(utility.CustomResponse{
			Status:  1,
			Message: "Login attempt initiated",
			Data: string(req.Body),
		}, 200)

	})

	router.Get("/v1/profile/:username/fetch", func(res *utility.Response, req *utility.Request) {
		fmt.Println("accessed path: " + req.Path)
		res.Send(utility.CustomResponse{
			Status:  1,
			Message: "Profile fetched successfully",
			Data:    req.Params["username"],
		}, 200)
	})

}
