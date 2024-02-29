package main

import (
	// core "coralscale/app/framework/core"
	"coralscale/app/framework/db"
	"coralscale/app/framework/migration"
	"coralscale/routes"
)

func main() {
	routes.Register()

	db.Init()
	migration.AutoMigrate()
	// core.InitWebService(core.Handler)
}
