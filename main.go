package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/mikkelstb/feedfetcher/config"
	"github.com/mikkelstb/feedfetcher/feed"
)

var config_file string
var sq_config config.RepositoryConfig
var db *SQLite

func init() {

	// To do: Make config relative

	var err error
	cfg, err := config.Read("/Users/mikkel/feedfetcher/config.json")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("error: config file not read, aborting")
		os.Exit(1)
	}

	for r := range cfg.Repositories {
		if cfg.Repositories[r].Type == "sqlite3" {
			sq_config = cfg.Repositories[r]
		}
	}

	db, err = NewSQLite(sq_config.Address)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	fmt.Println("Hello World")
	//getLatestArticles(10)

	resources_fileserver := http.FileServer(http.Dir("./resources"))

	http.HandleFunc("/", listArticles)
	http.HandleFunc("/feed/", listArticlesJson)
	http.HandleFunc("/test/", listArticlesTest)
	http.HandleFunc("/resources/", http.StripPrefix("/resources", resources_fileserver).ServeHTTP)
	http.ListenAndServe(":3001", nil)

}

func listArticlesTest(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(404)
	}

	languages, ok := r.Form["lang"]

	if !ok {
		languages = make([]string, 0)
	}

	var p Page
	p.setChecked(languages)

	p.Title = "HeadlineTracker"
	p.Today = time.Now().Format(time.ANSIC)
	p.Articles = getLatestArticles(20, languages)

	t, _ := template.ParseFiles(
		"./html_templates/list_test.html",
		"./html_templates/article.html",
		"./html_templates/header.html",
		"./html_templates/dropdown.html",
	)

	t.Execute(w, p)
}

func listArticlesJson(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(404)
	}

	languages, ok := r.Form["lang[]"]
	if !ok {
		languages = make([]string, 0)
	}

	fmt.Println(r.Form)
	fmt.Println(languages)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Encoding", "utf-8")
	articles := getLatestArticles(20, languages)
	jsondata, err := json.Marshal(articles)
	if err != nil {
		w.WriteHeader(404)
	}
	w.Write(jsondata)
}

func listArticles(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(404)
	}

	languages, ok := r.Form["lang"]

	if !ok {
		languages = make([]string, 0)
	}

	var p Page
	p.setChecked(languages)

	p.Title = "HeadlineTracker"
	p.Today = time.Now().Format(time.ANSIC)
	p.Articles = getLatestArticles(20, languages)

	t, _ := template.ParseFiles(
		"./html_templates/list.html",
		"./html_templates/article.html",
		"./html_templates/header.html",
		"./html_templates/dropdown.html",
	)

	t.Execute(w, p)
}

// TODO: Skrive et skript som leser de 10 siste artikler fra databasefil sorteret på docdate

func getLatestArticles(nb_articles int, languages []string) []feed.NewsItem {
	var articles []feed.NewsItem

	articles, err := db.GetArticles(nb_articles, languages)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return articles
}
