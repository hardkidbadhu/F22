package datastore

import (
	"F22/config"
	"F22/models"
	"F22/services"
	"log"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

var _ services.User = &userStore{}

type userStore struct {
	session    *mgo.Session
	collection *mgo.Collection
	conf       *config.Config
}

func NewUser(db *mgo.Session, cfg *config.Config) services.User {
	return &userStore{
		session:    db,
		collection: db.DB(cfg.DatabaseName).C("user"),
		conf:       cfg,
	}
}

func (us *userStore) FindByID(userId bson.ObjectId) (*models.User, error) {
	userIns := &models.User{}
	if err := us.collection.FindId(userId).One(&userIns); err != nil {
		log.Printf("Error - Datastore - FindByID - %s", err.Error())
		return nil, err
	}

	return userIns, nil
}

func (us *userStore) FindByCredentials(un, pw string) (*models.User, error) {
	userIns := &models.User{}

	if err := us.collection.Find(bson.M{"username": bson.RegEx{`^` + un + `$`, `i`}, "password": pw}).One(&userIns); err != nil {
		log.Printf("Error - Datastore - FindByID - %s", err.Error())
		return nil, err
	}

	return userIns, nil
}
