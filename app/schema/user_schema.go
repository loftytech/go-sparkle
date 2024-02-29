package schema

import (
	 "coralscale/app/framework/schema"
)

func UserSchema () {
	userSchema := schema.Schema{
		TableName: "users",
	}
	userSchema.Id()
	userSchema.String("username", 64)
	userSchema.String("password", 64)
	userSchema.String("first_name", 64)
	userSchema.String("last_name", 64)
	userSchema.String("email", 64)
	userSchema.Create()
}