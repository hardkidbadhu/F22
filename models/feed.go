package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type Article struct {
	Id          bson.ObjectId `json:"id" bson:"_id"`
	Title       string        `json:"title" bson:"title"`
	Description string        `json:"description" bson:"description"`
	Url         string        `json:"url" bson:"url"`
	Author      bson.ObjectId `json:"author" bson:"author"`
	comments    Comment       `json:"comments" json:"comments"`
	CreatedDate time.Time     `json:"createdDate" bson:"createdDate"`
	Likes       int           `json:"likes" bson:"likes"`
	DLikes      int           `json:"dLikes" bson:"dLikes"`
}

type Comment struct {
	Id            bson.ObjectId `json:"id" bson:"_id"`
	ParentId      bson.ObjectId `json:"ParentId" bson:"ParentId"`
	CommentString string        `json:"commentStr" bson:"commentStr"`
	Likes         int           `json:"likes" bson:"likes"`
	DLikes        int           `json:"dLikes" bson:"dLikes"`
	CreatedDate   time.Time     `json:"createdDate" bson:"createdDate"`
}
