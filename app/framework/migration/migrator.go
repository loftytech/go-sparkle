package migration

import (
	"coralscale/app/framework/schema"
	"coralscale/app/models"
)

func AutoMigrate() {
	schema.CreateModel(models.User{})
	schema.CreateModel(models.Profile{})
}
