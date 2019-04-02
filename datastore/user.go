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

func (us *userStore) FindByToken(token string) (*models.User, error) {
	userIns := &models.User{}
	if err := us.collection.Find(bson.M{"accessToken": token}).One(&userIns); err != nil {
		log.Printf("Error - Datastore - FindByID - %s", err.Error())
		return nil, err
	}

	return userIns, nil
}


func (as *userStore) SaveUser(user models.User) error {

	err := as.collection.Insert(&user)
	if err != nil {
		log.Printf("ERROR: Save(%s) - %q\n", user.Id, err)
		retry := as.conf.RetryDBInsert
		for retry > 0 {
			log.Printf("RETRY: Save(%s) - %q\n", user.Id, retry)
			if mgo.IsDup(err) {
				user.Id = bson.NewObjectId()
				if err = as.collection.Insert(user); err == nil {
					break
				}
			}

			if err.Error() == "EOF" {
				as.session.Refresh()
			}

			if err = as.collection.Insert(user); err == nil {
				break
			}
			retry--
		}

		if err != nil {
			log.Printf("ERROR: Save(%s) - %q\n", user.Id, err)
			return err
		}
	}

	return nil

}

func (us *userStore) FindAll(m bson.M) ([]models.User, error) {
	userIns := []models.User{}
	if err := us.collection.Find(m).All(&userIns); err != nil {
		log.Printf("Error - Datastore - FindByID - %s", err.Error())
		return nil, err
	}

	return userIns, nil
}