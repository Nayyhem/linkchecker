package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"LinkChecker/checker"
)

var tmpl = template.Must(template.ParseFiles("templates/index.html"))

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		link := r.FormValue("url")
		if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") {
			tmpl.Execute(w, map[string]interface{}{
				"Error": "URL invalide - doit commencer par http:// ou https://",
			})
			return
		}

		deadLinks := checker.CheckLinkPage(link)

		tmpl.Execute(w, map[string]interface{}{
			"Analyzed":  true,
			"DeadLinks": deadLinks,
		})
	} else {
		tmpl.Execute(w, nil)
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Serveur lancé sur http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
