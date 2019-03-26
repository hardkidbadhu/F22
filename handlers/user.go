package handlers

import "net/http"

//Authenticate validates the user auth
func (p *Provider) Authenticate (rw http.ResponseWriter, r *http.Request) {

	credentials := struct {
		username string `json:"username"`
		Password string  `json:"password"`
	}{}

	//Parse the credentials from request and convert to the credentials type
	if ok, err := parseJson(r.Body, &credentials); !ok && err != nil {
		renderJson(rw, http.StatusBadRequest, err)
		return
	}



}
