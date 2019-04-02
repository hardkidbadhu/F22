package route

import (
	"net/http"

	"F22/handlers"

	"github.com/gorilla/mux"
)

func NewRouter(p *handlers.Provider) *mux.Router {

	router := mux.NewRouter()

	//Anonymous function for wrapping middleware
	fn := func(handler http.HandlerFunc) http.Handler {
		return Adapt(http.HandlerFunc(handler), recoverHandler(), loggingHandler(p.Logger()), sessionAuthenticator(p.DB(), p.Config()))
	}

	ws := func(handler http.HandlerFunc) http.Handler {
		return Adapt(http.HandlerFunc(handler), recoverHandler(), loggingHandler(p.Logger()))
	}

	router.Handle("/", ws(p.RedirectToHome)).Methods("GET")
	router.Handle("/login", ws(p.Login)).Methods("GET")
	router.Handle("/signUp", ws(p.SignUp)).Methods("GET")
	router.Handle("/authenticate", ws(p.Authenticate)).Methods("POST")

	router.Handle("/newArticle", fn(p.NewArticle)).Methods("GET")
	router.Handle("/createNewArticle", fn(p.SaveArticle)).Methods("POST")

	router.Handle("/article/{id}", fn(p.OpenArticle)).Methods("GET")

	router.Handle("/article/{id}/postComment", fn(p.PostComment)).Methods("POST")

	router.Handle("/comment/{id}/reply", fn(p.Reply)).Methods("GET")
	router.Handle("/article/{id}/vote", fn(p.VoteArticle)).Methods("GET")

	return router
}
