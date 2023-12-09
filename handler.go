package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	urlverifier "github.com/davidmytton/url-verifier"
	"github.com/go-chi/chi/v5"
)

type UrlData struct {
	Url                   string
	Valid                 bool
	AmountOfUrlsGenerated int
	InternalError         bool
	UrlEmpty              bool
}

func (a *App) shortenHandler(w http.ResponseWriter, r *http.Request) {
	var hasInternalError bool = false
	var hasValidUrl bool = true
	var IsUrlEmpty bool = false
	var short string
	var AmountOfUrlsGenerated int

	formUrl := r.PostFormValue("inputUrl")

	if len(formUrl) == 0 {
		log.Println("url is empty")
		IsUrlEmpty = true
	}

	verifier := urlverifier.NewVerifier()
	ret, err := verifier.Verify(formUrl)

	if err != nil || !ret.IsURL {
		fmt.Println("url isnt valid")
		hasValidUrl = false
	}
	if hasValidUrl {

		short, err = a.shortenUrl(formUrl)
		if err != nil {
			log.Println("error generating a short code")
			hasInternalError = true
		}

		AmountOfUrlsGenerated, err = a.getAmountOfLinks()
		if err != nil {
			log.Println("error getting the amount of urls generated")
			// we dont crash here, just leave it at 0
		}
	}

	u := UrlData{Url: strings.Join([]string{a.RootUrl, "r/", short}, ""), Valid: hasValidUrl, AmountOfUrlsGenerated: AmountOfUrlsGenerated, InternalError: hasInternalError, UrlEmpty: IsUrlEmpty}

	tmpl := template.Must(template.ParseFS(html, "pages/shortened.html"))
	tmpl.Execute(w, u)
}

func (a *App) redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "id")
	url, err := a.getFromDatabase(shortCode)
	if err != nil {
		log.Printf("couldn't find %s in the database\n", shortCode)
		tmpl := template.Must(template.ParseFS(html, "base.html", "pages/404.html"))
		tmpl.Execute(w, nil)
		return
	}
	if !strings.HasPrefix(url, "http://") || !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}
	http.Redirect(w, r, url, http.StatusPermanentRedirect)
}
