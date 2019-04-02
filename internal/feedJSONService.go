package internal

import (
	"F22/config"
	"F22/datastore"
	"F22/models"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"log"
)

func GetArticleJson(article *models.Article, dbSession *mgo.Session, cfg *config.Config, token string) *models.Article {

	cmtSrv := datastore.NewComment(dbSession, cfg)

	cmtMap := make(map[models.Comment][]models.Comment)

	usrSrv := datastore.NewUser(dbSession, cfg)
	userMap := make(map[bson.ObjectId]string)
	users, _ := usrSrv.FindAll(bson.M{})
	for u := range users {
		userMap[users[u].Id] = users[u].Name
	}

	if comments, error := cmtSrv.List(bson.M{"articleId": article.Id, "firstComment": true}); error == nil && len(comments) > 0 {

		//Group comments in the way the parent to child comments i.e main comments and reply to that
		for c := range comments {

			comments[c].ReplyLink = fmt.Sprintf("%s/comment/%s/reply?token=%s", cfg.Endpoint, comments[c].Id.Hex(), token)
			if authName, ok := userMap[comments[c].AuthorId]; ok {
				log.Print(authName)
				comments[c].AuthorName = authName
			}

			if replies, error := cmtSrv.List(bson.M{"articleId": article.Id, "ParentId": comments[c].Id, "_id": bson.M{"$ne": comments[c].Id}}); error == nil {

				for r := range replies {

					if authName, ok := userMap[replies[r].AuthorId]; ok {
						replies[r].AuthorName = authName
					}
				}

				cmtMap[comments[c]] = replies
				continue
			}

			cmtMap[comments[c]] = []models.Comment{}
		}

	}

	article.Comments = cmtMap

	log.Printf("Article %+v", article)
	return article
}

func ArticleListJSON(articles []models.Article, dbSession *mgo.Session, cfg *config.Config, token string) []models.Article {

	usrSrv := datastore.NewUser(dbSession, cfg)

	userMap := make(map[bson.ObjectId]string)
	users, _ := usrSrv.FindAll(bson.M{})
	for u := range users {
		userMap[users[u].Id] = users[u].Name
	}

	for a := range articles {
		if userIns, error := usrSrv.FindByID(articles[a].Author); error == nil {
			articles[a].AuthorName = userIns.Name
		}

		articles[a].UpVoteUrl = fmt.Sprintf("%s/article/%s/vote?token=%s&vote=like", cfg.Endpoint, articles[a].Id.Hex(), token)
		articles[a].DownVoteUrl = fmt.Sprintf("%s/article/%s/vote?token=%s&vote=dlike", cfg.Endpoint, articles[a].Id.Hex(), token)
		articles[a].RedirectUrl = fmt.Sprintf("%s/article/%s?token=%s", cfg.Endpoint, articles[a].Id.Hex(), token)
		articles[a].CreatedDateStr = articles[a].CreatedDate.Format("01-Jan-2006")

		log.Print(articles[a].UpVoteUrl, articles[a].DownVoteUrl)
	}

	return articles
}
