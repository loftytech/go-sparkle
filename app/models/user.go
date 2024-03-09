package models

type User struct {
	Id         int `orm:"primary"`
	Username   string
	Password   string
	First_name string
	Last_name  string
	Email      string
}
