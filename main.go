package main

import (
	core "coralscale/app/framework/core"
	"coralscale/routes"
)

func main() {
	routes.Register()
	core.InitWebService(core.Handler)
}
