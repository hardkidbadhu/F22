package datastore

import (
	"F22/config"
	"F22/db"
	"F22/models"
	"F22/services"
	"github.com/globalsign/mgo/bson"
	"log"

	"github.com/globalsign/mgo"
)

type commentStore struct {
	session    *mgo.Session
	collection *mgo.Collection
	conf       *config.Config
}

func NewComment(session *mgo.Session, cfg *config.Config) services.Comment {
	return &commentStore{
		session:    session,
		collection: session.DB(cfg.DatabaseName).C(db.Comments),
		conf:       cfg,
	}
}

func (c *commentStore) Save(comment *models.Comment) error {

	err := c.collection.Insert(&comment)
	if err != nil {
		log.Printf("ERROR: Save(%s) - %q\n", comment.Id, err)
		retry := c.conf.RetryDBInsert
		for retry > 0 {
			log.Printf("RETRY: Save(%s) - %q\n", comment.Id, retry)
			if mgo.IsDup(err) {
				comment.Id = bson.NewObjectId()
				if err = c.collection.Insert(comment); err == nil {
					break
				}
			}

			if err.Error() == "EOF" {
				c.session.Refresh()
			}

			if err = c.collection.Insert(comment); err == nil {
				break
			}
			retry--
		}

		if err != nil {
			log.Printf("ERROR: Save(%s) - %q\n", comment.Id, err)
			return err
		}
	}

	return nil


}

func (c * commentStore) List(m bson.M) ([]models.Comment, error) {

	comments:= []models.Comment{}

	if err := c.collection.Find(m).All(&comments); err != nil {
		log.Printf("Error - Datastore - FindByID - %s", err.Error())
		return nil, err
	}

	return comments, nil

}

func (c *commentStore) Find(cmtId bson.ObjectId) (*models.Comment, error) {

	commentIns := &models.Comment{}

	if err := c.collection.FindId(cmtId).One(&commentIns); err != nil {
		log.Printf("Error - Datastore - FindByID - %s", err.Error())
		return nil, err
	}

	return commentIns, nil
}