package handlers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"F22/err"
)

func parseJson(params io.ReadCloser, data interface{}) (bool, error) {
	b, _ := ioutil.ReadAll(params)
	error := json.Unmarshal(b, data)
	if error == nil {
		return true, nil
	}
	return false, err.UIError{
		error,
		"Error in parsing JSON!.",
		http.StatusBadRequest,
	}
}


func renderJson(w http.ResponseWriter, status int, res interface{}) {
	resByte, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(resByte)
}