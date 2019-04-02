package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"F22/datastore"
	"F22/internal"
	"F22/models"

	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func (p *Provider) NewArticle(rw http.ResponseWriter, r *http.Request) {

	args := make(map[string]interface{})

	token := r.URL.Query().Get("token")

	args["url"] = fmt.Sprintf("%s/createNewArticle?token=%s", p.cfg.Endpoint, r.URL.Query().Get("token"))

	if val := context.Get(r, "userId"); val == nil {
		log.Println("Error - NewArticle - Invalid user identifier!.")
		renderTemplate(rw, p.cfg, "views/login.html", args)
		return
	}

	args["login"] = fmt.Sprintf("%s/login=%s", p.cfg.Endpoint, token)
	args["newArticle"] = fmt.Sprintf("%s/newArticle?token=%s", p.cfg.Endpoint, token)
	args["home"] = fmt.Sprintf("%s?token=%s", p.cfg.Endpoint, token)

	renderTemplate(rw, p.cfg, "views/newArticle.html", args)

}

func (p *Provider) SaveArticle(rw http.ResponseWriter, r *http.Request) {

	article := struct {
		Title string
		Url   string
		Text  string
	}{
		r.FormValue("title"),
		r.FormValue("url"),
		r.FormValue("message"),
	}

	arg := make(map[string]interface{})

	var (
		userId_s   string
		ok         bool
		articleIns models.Article
		error      error
	)

	if userId_s, ok = context.Get(r, "userId").(string); !ok && bson.IsObjectIdHex(userId_s) {
		log.Print("Error - Invalid user id supplied!.")
		arg["hasError"] = true
		arg["message"] = "Invalid user name or password!."
		renderTemplate(rw, p.cfg, "views/login.html", arg)
		return
	}

	articleIns.Id = bson.NewObjectId()
	articleIns.Title = article.Title
	articleIns.Description = article.Text
	articleIns.Url = article.Url
	articleIns.Author = bson.ObjectIdHex(userId_s)
	articleIns.CreatedDate = time.Now().UTC()

	token := r.URL.Query().Get("token")
	arg["url"] = fmt.Sprintf("%s/createNewArticle?token=%s", p.cfg.Endpoint, token)
	if error = datastore.NewArticle(p.db, p.cfg).Save(articleIns); error != nil {
		log.Print("Error - Invalid user id supplied!.")
		arg["message"] = "Something went wrong please try after sometime!"
		renderTemplate(rw, p.cfg, "views/newArticle.html", arg)
		return
	}

	arg["login"] = fmt.Sprintf("%s/login?token=%s", p.cfg.Endpoint, token)
	arg["newArticle"] = fmt.Sprintf("%s/newArticle=%s", p.cfg.Endpoint, token)
	arg["home"] = fmt.Sprintf("%s?token=%s", p.cfg.Endpoint, token)

	articles, _ := datastore.NewArticle(p.db, p.cfg).List(bson.M{})
	arg["articles"] = internal.ArticleListJSON(articles, p.db, p.cfg, token)

	renderTemplate(rw, p.cfg, "views/home.html", arg)
}

func (p *Provider) OpenArticle(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	articleId_s := vars["id"]

	token := r.URL.Query().Get("token")

	if ok := bson.IsObjectIdHex(articleId_s); !ok {
		http.RedirectHandler(fmt.Sprintf("%s?token=%s", p.cfg.Endpoint, token), http.StatusSeeOther)
		return
	}

	//fetch the article from the database

	articleIns, error := datastore.NewArticle(p.db, p.cfg).Find(bson.ObjectIdHex(articleId_s))
	if error != nil {
		p.log.Printf("Error - OpenArticle - %s", error.Error())
		http.RedirectHandler(fmt.Sprintf("%s?token=%s", p.cfg.Endpoint, token), http.StatusSeeOther)
		return
	}

	if len(articleIns.Url) > 3 {
		http.Redirect(rw, r, articleIns.Url, http.StatusMovedPermanently)
		return
	}

	args := make(map[string]interface{})

	args["login"] = fmt.Sprintf("%s/login?token=%s", p.cfg.Endpoint, token)
	args["newArticle"] = fmt.Sprintf("%s/newArticle?token=%s", p.cfg.Endpoint, token)
	args["home"] = fmt.Sprintf("%s?token=%s", p.cfg.Endpoint, token)

	if name, ok := context.Get(r, "sessionUser").(string); ok {
		args["name"] = name
	}

	articleJson := internal.GetArticleJson(articleIns, p.db, p.cfg, token)

	args["article"] = *articleJson
	log.Print(args["article"])
	args["url"] = fmt.Sprintf("%s/article/%s/postComment?token=%s", p.cfg.Endpoint, articleId_s, token)

	renderTemplate(rw, p.cfg, "views/addComment.html", args)

}

func (p *Provider) PostComment(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	articleId_s := vars["id"]

	token := r.URL.Query().Get("token")
	commentId := r.URL.Query().Get("commentId")

	if !bson.IsObjectIdHex(articleId_s) {
		http.RedirectHandler(fmt.Sprintf("%s?token=%s", p.cfg.Endpoint, token), http.StatusTemporaryRedirect)
		return
	}

	//fetch the article from the database

	var (
		error      error
		articleIns *models.Article
		commentIns *models.Comment
		userIns    *models.User
		args       map[string]interface{}
	)

	if articleIns, error = datastore.NewArticle(p.db, p.cfg).Find(bson.ObjectIdHex(articleId_s)); error != nil {
		p.log.Printf("Error - OpenArticle - %s", error.Error())
		http.RedirectHandler(fmt.Sprintf("%s?token=%s", p.cfg.Endpoint, token), http.StatusTemporaryRedirect)
		return
	}

	//Comment service instance
	cmtSrv := datastore.NewComment(p.db, p.cfg)

	userId, ok := context.Get(r, "userId").(string)
	if !ok {
		p.log.Print("Error - Handler - Interface conversion...")
		p.RedirectToHome(rw, r)
		return
	}

	parentId := bson.NewObjectId()
	firstComment := true

	if bson.IsObjectIdHex(commentId) {
		if commentIns, _ = cmtSrv.Find(bson.ObjectIdHex(commentId)); commentIns != nil {
			parentId = commentIns.Id
			firstComment = false
		}
	}

	newComment := models.Comment{
		parentId,
		parentId,
		articleIns.Id,
		bson.ObjectIdHex(userId),
		firstComment,
		"",
		r.FormValue("message"),
		0,
		0,
		time.Now().UTC(),
		"",
		"",
	}

	args = make(map[string]interface{})

	args["login"] = fmt.Sprintf("%s/login=%s", p.cfg.Endpoint, token)
	args["newArticle"] = fmt.Sprintf("%s/newArticle?token=%s", p.cfg.Endpoint, token)
	args["home"] = fmt.Sprintf("%s?token=%s", p.cfg.Endpoint, token)

	if error = cmtSrv.Save(&newComment); error != nil {
		p.log.Printf("Error - handler - PostComment - %s", error.Error())

		if userIns, error = datastore.NewUser(p.db, p.cfg).FindByToken(token); userIns != nil && error == nil {
			args["name"] = userIns.Name
		}

		articles, _ := datastore.NewArticle(p.db, p.cfg).List(bson.M{})
		args["articles"] = internal.ArticleListJSON(articles, p.db, p.cfg, token)
		renderTemplate(rw, p.cfg, "views/home.html", args)
		return
	}

	articleJson := internal.GetArticleJson(articleIns, p.db, p.cfg, token)

	args["article"] = *articleJson
	args["url"] = fmt.Sprintf("%s/article/%s/postComment?token=%s", p.cfg.Endpoint, articleId_s, token)

	renderTemplate(rw, p.cfg, "views/addComment.html", args)
}

func (p *Provider) Reply(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	commentId := vars["id"]

	token := r.URL.Query().Get("token")

	if !bson.IsObjectIdHex(commentId) {
		p.log.Printf("Invalid objectId - %s", commentId)
		http.RedirectHandler(fmt.Sprintf("%s?token=%s", p.cfg.Endpoint, token), http.StatusTemporaryRedirect)
		return
	}

	cmtIns, error := datastore.NewComment(p.db, p.cfg).Find(bson.ObjectIdHex(commentId))
	if error != nil {
		p.log.Printf("Error - Handler - Reply - %s", error.Error())
		http.RedirectHandler(fmt.Sprintf("%s?token=%s", p.cfg.Endpoint, token), http.StatusTemporaryRedirect)
		return
	}

	args := make(map[string]interface{})
	args["comment"] = *cmtIns

	args["url"] = fmt.Sprintf("%s/article/%s/postComment?token=%s&commentId=%s", p.cfg.Endpoint, cmtIns.ArticleId.Hex(), token, commentId)

	renderTemplate(rw, p.cfg, "views/reply.html", args)
}

func (p *Provider) VoteArticle(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	articleId := vars["id"]

	token := r.URL.Query().Get("token")
	vote := r.URL.Query().Get("vote")

	var (
		articles   []models.Article
		articleIns *models.Article
		error      error
		args       map[string]interface{}
	)

	args = make(map[string]interface{})

	articleSrv := datastore.NewArticle(p.db, p.cfg)

	if articles, error = datastore.NewArticle(p.db, p.cfg).List(bson.M{}); error != nil {
		log.Printf("Error - Handler - %s", error.Error())
		renderTemplate(rw, p.cfg, "views/home.html", args)
		return
	}
	args["login"] = fmt.Sprintf("%s/login=%s", p.cfg.Endpoint, token)
	args["newArticle"] = fmt.Sprintf("%s/newArticle?token=%s", p.cfg.Endpoint, token)
	args["home"] = fmt.Sprintf("%s?token=%s", p.cfg.Endpoint, token)
	if name, ok := context.Get(r, "sessionUser").(string); ok {
		args["name"] = name
	}

	args["articles"] = internal.ArticleListJSON(articles, p.db, p.cfg, token)
	if articleIns, error = articleSrv.Find(bson.ObjectIdHex(articleId)); error != nil {
		p.log.Printf("Error - Handler - %s", error.Error())
		renderTemplate(rw, p.cfg, "views/home.html", args)
		return
	}

	switch vote {
	case "like":
		if error = articleSrv.Update(bson.M{"_id" : articleIns.Id}, bson.M{"$inc" : bson.M{"likes": 1}}); error != nil {
			p.log.Printf("Error - Handler - %s", error.Error())
		}
		renderTemplate(rw, p.cfg, "views/home.html", args)
		return
	case "dlike":
		if error = articleSrv.Update(bson.M{"_id" : articleIns.Id}, bson.M{"$inc" : bson.M{"dLikes": 1}}); error != nil {
			p.log.Printf("Error - Handler - %s", error.Error())
		}
		renderTemplate(rw, p.cfg, "views/home.html", args)
		return
	}

	renderTemplate(rw, p.cfg, "views/home.html", args)
}
