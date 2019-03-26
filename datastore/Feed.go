package datastore

import (
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
func NewArticle(db *mgo.Session, cfg *config.Config) services.Article {
	return &articleStore{
		session:    db,
		collection: db.DB(cfg.DatabaseName).C("article"),
		conf:       cfg,
	}
}

func (as *articleStore) Find (userId bson.ObjectId) (*models.Article, error) {
	article := &models.Article{}
	if err := as.collection.FindId(userId).One(&article); err != nil {
		log.Printf("Error - Datastore - FindByID - %s", err.Error())
		return nil, err
	}

	return article, nil
}

func (as *articleStore) List(id bson.ObjectId) ([]*models.Article, error) {

	articles := []*models.Article{}
	if err := as.collection.Find(bson.M{"author": id}).All(&articles); err != nil {
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