package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type User struct {
	Id          bson.ObjectId `json:"id" bson:"_id"`
	Name        string        `json:"name" bson:"name"`
	UserName    string        `json:"username" bson:"username"`
	Password    string        `json:"password" bson:"password"`
	CreatedDate time.Time     `json:"createdDate" bson:"createdDate"`
}

//Stringer implementation returns only the name when the structure is printed
func (u *User) String() string {
	return u.Name
}