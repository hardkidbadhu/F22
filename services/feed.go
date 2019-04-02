package services

import (
	"F22/models"

	"github.com/globalsign/mgo/bson"
)

type Article interface {
	Find(bson.ObjectId) (*models.Article, error)
	List(bson.M) ([]models.Article, error)
	Save(models.Article) error
	Update(bson.M, bson.M) error
}
