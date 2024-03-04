package migration

import "coralscale/app/schema"

func AutoMigrate() {
	schema.UserSchema()
	schema.ProfileSchema()
}