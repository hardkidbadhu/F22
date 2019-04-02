package handlers

import (
	"fmt"
	"github.com/pborman/uuid"
	"log"
	"net/http"
	"time"

	"F22/datastore"
	"F22/internal"
	"F22/models"

	"github.com/globalsign/mgo/bson"
)

func (p *Provider) Login(rw http.ResponseWriter, r *http.Request) {

	args := make(map[string]interface{})

	token := r.URL.Query().Get("token")

	userIns, err := datastore.NewUser(p.db, p.cfg).FindByToken(token)
	if err != nil {
		args["authenticate"] = fmt.Sprintf("%s/authenticate", p.cfg.Endpoint)
		args["signUp"] = fmt.Sprintf("%s/signUp", p.cfg.Endpoint)
		args["name"] = `Guest`
		renderTemplate(rw, p.cfg, "views/login.html", args)
		return
	}

	args["login"] = fmt.Sprintf("%s/login?token=%s", p.cfg.Endpoint, token)
	args["newArticle"] = fmt.Sprintf("%s/newArticle?token=%s", p.cfg.Endpoint, token)
	args["home"] = fmt.Sprintf("%s?token=%s", p.cfg.Endpoint, token)
	args["name"] = userIns.Name

	renderTemplate(rw, p.cfg, "views/home.html", args)

}

func (p *Provider) SignUp(rw http.ResponseWriter, r *http.Request) {

	args := make(map[string]interface{})

	args["authenticate"] = fmt.Sprintf("%s/authenticate", p.cfg.Endpoint)
	renderTemplate(rw, p.cfg, "views/signUp.html", args)
}

//Authenticate validates the user auth
func (p *Provider) Authenticate(rw http.ResponseWriter, r *http.Request) {

	action := r.FormValue("methodAction")

	log.Print(action)

	arg := make(map[string]interface{})

	switch action {

	case "login":

		credentials := struct {
			username string `json:"username"`
			Password string `json:"password"`
		}{
			r.FormValue("un"),
			r.FormValue("pw"),
		}

		//Parse the credentials from request and convert to the credentials type
		var (
			error   error
			userIns *models.User
		)

		if userIns, error = datastore.NewUser(p.db, p.cfg).FindByCredentials(credentials.username, credentials.Password); error != nil && userIns == nil {
			p.log.Printf("Error - Handler - Authenticate - %s", error.Error())
			arg["hasError"] = true
			arg["message"] = "Invalid user name or password!."
			renderTemplate(rw, p.cfg, "views/login.html", arg)
			return
		}

		log.Println("Tok ", userIns.AccessToken)

		arg["newArticle"] = fmt.Sprintf("%s/newArticle?token=%s", p.cfg.Endpoint, userIns.AccessToken)
		arg["login"] = fmt.Sprintf("%s/login", p.cfg.Endpoint)
		articles, _ := datastore.NewArticle(p.db, p.cfg).List(bson.M{})
		arg["articles"] = internal.ArticleListJSON(articles, p.db, p.cfg, userIns.AccessToken)
		arg["name"] = userIns.Name

		renderTemplate(rw, p.cfg, "views/home.html", arg)
		return

	case "signUp":

		userDetails := struct {
			UserName          string
			Password          string
			confirmedPassword string
			Name              string
		}{
			r.FormValue("un"),
			r.FormValue("pw"),
			r.FormValue("cpw"),
			r.FormValue("name"),
		}

		uuid := uuid.NewRandom()
		loginToken := uuid.String()

		log.Println("Tok ", loginToken)

		u := models.User{
			bson.NewObjectId(),
			userDetails.Name,
			userDetails.UserName,
			userDetails.Password,
			loginToken,
			time.Now().UTC(),
		}

		arg := make(map[string]interface{})

		if err := datastore.NewUser(p.db, p.cfg).SaveUser(u); err != nil {
			log.Printf("Error - SignUp - %s", err.Error())
			arg["hasError"] = true
			arg["message"] = "Something went wrong please try after sometimes!."
			renderTemplate(rw, p.cfg, "views/login.html", arg)
			return
		}

		arg["login"] = fmt.Sprintf("%s/login?token=%s", p.cfg.Endpoint, loginToken)
		arg["newArticle"] = fmt.Sprintf("%s/newArticle?token=%s", p.cfg.Endpoint, loginToken)
		arg["home"] = fmt.Sprintf("%s", p.cfg.Endpoint)
		arg["name"] = u.Name

		articles, _ := datastore.NewArticle(p.db, p.cfg).List(bson.M{})
		arg["articles"] = internal.ArticleListJSON(articles, p.db, p.cfg, loginToken)
		renderTemplate(rw, p.cfg, "views/home.html", arg)
		return

	default:
		arg["hasError"] = true
		arg["message"] = "Something went wrong please try after sometimes!."
		renderTemplate(rw, p.cfg, "views/login.html", arg)
		return
	}

}

func (p *Provider) RedirectToHome(rw http.ResponseWriter, r *http.Request) {

	args := make(map[string]interface{})

	token := r.URL.Query().Get("token")

	args["login"] = fmt.Sprintf("%s/login", p.cfg.Endpoint)
	args["newArticle"] = fmt.Sprintf("%s/newArticle", p.cfg.Endpoint)
	args["home"] = fmt.Sprintf("%s", p.cfg.Endpoint)
	args["name"] = "Guest user"

	var (
		userIns  *models.User
		error    error
		articles []models.Article
	)

	if len(token) > 3 {

		if userIns, error = datastore.NewUser(p.db, p.cfg).FindByToken(token); userIns != nil && error == nil {
			args["name"] = userIns.Name
		}

		args["login"] = fmt.Sprintf("%s/login=%s", p.cfg.Endpoint, token)
		args["newArticle"] = fmt.Sprintf("%s/newArticle?token=%s", p.cfg.Endpoint, token)
		args["home"] = fmt.Sprintf("%s?token=%s", p.cfg.Endpoint, token)
	}

	//Home page contains list of articles
	//Query to list all the articles

	if articles, error = datastore.NewArticle(p.db, p.cfg).List(bson.M{}); error != nil {
		log.Printf("Error - Handler - %s", error.Error())
		renderTemplate(rw, p.cfg, "views/home.html", args)
		return
	}

	args["articles"] = internal.ArticleListJSON(articles, p.db, p.cfg, token)

	log.Print(args["articles"])
	renderTemplate(rw, p.cfg, "views/home.html", args)

}

