package datastore

import (
	"F22/db"
	"log"

	"F22/config"
	"F22/models"
	"F22/services"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

//Throws error in compile time if articleStore is not a type of services.Article
//Lets compiler checks whether the interface implementation is done properly
var _ services.Article = &articleStore{}

type articleStore struct {
	session    *mgo.Session
	collection *mgo.Collection
	conf       *config.Config
}

//Interface implementation
//Can return articleStore as type services.Article coz it implements all the service interface methods
func NewArticle(session *mgo.Session, cfg *config.Config) services.Article {
	return &articleStore{
		session:    session,
		collection: session.DB(cfg.DatabaseName).C(db.Article),
		conf:       cfg,
	}
}

func (as *articleStore) Find(userId bson.ObjectId) (*models.Article, error) {
	article := &models.Article{}
	if err := as.collection.FindId(userId).One(&article); err != nil {
		log.Printf("Error - Datastore - FindByID - %s", err.Error())
		return nil, err
	}

	return article, nil
}

func (as *articleStore) List(m bson.M) ([]models.Article, error) {

	articles := []models.Article{}
	if err := as.collection.Find(m).All(&articles); err != nil {
		log.Printf("Error - Datastore - FindByID - %s", err.Error())
		return nil, err
	}

	return articles, nil

}

func (as *articleStore) Save(article models.Article) error {

	err := as.collection.Insert(&article)
	if err != nil {
		log.Printf("ERROR: Save(%s) - %q\n", article.Id, err)
		retry := as.conf.RetryDBInsert
		for retry > 0 {
			log.Printf("RETRY: Save(%s) - %q\n", article.Id, retry)
			if mgo.IsDup(err) {
				article.Id = bson.NewObjectId()
				if err = as.collection.Insert(article); err == nil {
					break
				}
			}

			if err.Error() == "EOF" {
				as.session.Refresh()
			}

			if err = as.collection.Insert(article); err == nil {
				break
			}
			retry--
		}

		if err != nil {
			log.Printf("ERROR: Save(%s) - %q\n", article.Id, err)
			return err
		}
	}

	return nil

}

func (as *articleStore) Update(selector, updater bson.M) error {

	err := as.collection.Update(selector, updater)
	if err != nil {
		log.Printf("ERROR: Update(%+v, %+v) - %q\n", selector, updater, err)
		retry := as.conf.RetryDBInsert

		for retry > 0 {
			log.Printf("ERROR: Update(%+v, %+v) - %q\n", selector, updater, err)
			if mgo.IsDup(err) {
				break
			}
			if err.Error() == "EOF" {
				as.session.Refresh()
			}
			if err = as.collection.Update(updater, selector); err == nil {
				break
			}
			retry--
		}
		if err != nil {
			log.Printf("ERROR: Update(%+v, %+v) - %q\n", selector, updater, err)
			return err
		}
	}

	return nil

}
