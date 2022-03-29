package main

import (
	"fmt"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/mikkelstb/feedfetcher/config"
	"github.com/mikkelstb/feedfetcher/feed"
)

var sq_config config.RepositoryConfig
var db *SQLite

func init() {

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
	getLatestArticles(10)

	resources_fileserver := http.FileServer(http.Dir("./resources"))

	http.HandleFunc("/", listArticles)
	http.HandleFunc("/resources/", http.StripPrefix("/resources", resources_fileserver).ServeHTTP)
	http.ListenAndServe(":3001", nil)

}

func listArticles(w http.ResponseWriter, r *http.Request) {

	var p Page
	p.Title = "HeadlineTracker"
	p.Today = time.Now().Format(time.ANSIC)
	p.Articles = getLatestArticles(10)
	p.SourceName = "test"

	t, _ := template.ParseFiles(
		"./html_templates/list.html",
		"./html_templates/article.html",
		"./html_templates/header.html",
	)
	t.Execute(w, p)
}

// TODO: Skrive et skript som leser de 10 siste artikler fra databasefil sorteret på docdate

func getLatestArticles(nb_articles int) []feed.NewsItem {
	var articles []feed.NewsItem

	articles, err := db.GetArticles(nb_articles)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return articles
}
