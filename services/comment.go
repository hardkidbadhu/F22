package services

import (
	"F22/models"

	"github.com/globalsign/mgo/bson"
)

type Comment interface {
	Find(bson.ObjectId) (*models.Comment, error)
	Save(*models.Comment) error
	List(bson.M) ([]models.Comment, error)
}
