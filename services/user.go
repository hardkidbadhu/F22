package services

import (
	"F22/models"

	"github.com/globalsign/mgo/bson"
)

type User interface {
	FindByID(bson.ObjectId) (*models.User, error)
	FindByToken(string) (*models.User, error)
	FindByCredentials(string, string) (*models.User, error)
	SaveUser(models.User) error
	FindAll(bson.M) ([]models.User, error)
}

