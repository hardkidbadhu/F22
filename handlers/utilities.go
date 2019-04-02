package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"F22/config"
)


func renderTemplate(w http.ResponseWriter, cfg *config.Config, fileName string, arg map[string]interface{}) {
	htmlStr := LoadTemplate(cfg, fileName, arg)
	fmt.Fprint(w, htmlStr)
}

func LoadTemplate(cfg *config.Config, templatePath string, args map[string]interface{}) (templateString string) {

	_, templateName := filepath.Split(templatePath)

	temp, err := template.New(templateName).Delims(cfg.DelimsL, cfg.DelimsR).ParseFiles(templatePath)
	if err != nil {
		log.Println("Error in parse files:", err.Error())
		return
	}

	b := bytes.Buffer{}

	err = temp.Delims(cfg.DelimsL, cfg.DelimsR).Execute(&b, args)
	if err != nil {
		log.Println("Error in execute template:", err.Error())
		return
	}

	templateString = b.String()
	return
}


func renderJson(w http.ResponseWriter, status int, res interface{}) {
	resByte, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(resByte)
}