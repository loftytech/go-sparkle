package schema

import (
	 "coralscale/app/framework/schema"
)

func ProfileSchema () {
	userSchema := schema.Schema{
		TableName: "profile",
	}
	userSchema.Integer("profile_id").AutoIncrement().Primary()
	userSchema.Double("balance")
	userSchema.String("work", 64)
	userSchema.String("state", 64)
	userSchema.String("city", 64)
	userSchema.String("country", 64)
	userSchema.Text("about")
	userSchema.Create()
}