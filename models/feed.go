package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type Article struct {
	Id             bson.ObjectId         `json:"id" bson:"_id"`
	Title          string                `json:"title" bson:"title"`
	Description    string                `json:"description" bson:"description"`
	Url            string                `json:"url" bson:"url"`
	Author         bson.ObjectId         `json:"author" bson:"author"`
	AuthorName     string                `json:"authorName" bson:"-"`
	Comments       map[Comment][]Comment `json:"comments" json:"comments"`
	HasReplies     bool                  `json:"hasComments" bson:"hasComments"`
	CreatedDate    time.Time             `json:"-" bson:"createdDate"`
	CreatedDateStr string                `json:"createdDate" bson:"-"`
	Likes          int                   `json:"likes" bson:"likes"`
	DLikes         int                   `json:"dLikes" bson:"dLikes"`
	RedirectUrl    string                `json:"redirectUrl"`
	UpVoteUrl      string                `json:"upVoteUrl"`
	DownVoteUrl    string                `json:"downVoteUrl"`
}

type Comment struct {
	Id             bson.ObjectId `json:"id" bson:"_id"`
	ParentId       bson.ObjectId `json:"ParentId" bson:"ParentId"`
	ArticleId      bson.ObjectId `json:"articleId" bson:"articleId"`
	AuthorId       bson.ObjectId `json:"authorId" bson:"authorId"`
	FirstComment   bool          `json:"-" bson:"firstComment"` //maintaining this boolean to identify this is parent since $where comparison doesn't work on objectId's
	AuthorName     string        `json:"authorName" bson:"-"`
	CommentString  string        `json:"commentStr" bson:"commentStr"`
	Likes          int           `json:"likes" bson:"likes"`
	DLikes         int           `json:"dLikes" bson:"dLikes"`
	CreatedDate    time.Time     `json:"-" bson:"createdDate"`
	CreatedDateStr string        `json:"createdDate" bson:"-"`
	ReplyLink      string        `json:"reply_link"`
}
