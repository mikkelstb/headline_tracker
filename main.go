package main

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mikkelstb/headline_tracker/config"
	"github.com/mikkelstb/headline_tracker/database"
	"github.com/mikkelstb/headline_tracker/obj"
	"github.com/mikkelstb/ir_models/boolean"
	"github.com/mikkelstb/ir_models/ipop"

	"github.com/mikkelstb/simplelog"
)




var sourcepage struct{
	Today string
	Title string
	SourceName string
	Articles []obj.NewsItem
	Category [][]string
}

var db *database.Database
var cfg *config.General
var termdictionary *boolean.TermDictionary
var term_dic_database *ipop.Database

func main () {
	
	simplelog.InitLoggers("./ht.log")

	cfg = config.Read("./headline_tracker_config.json")

	db = new(database.Database)
	db.Init(cfg.Username, cfg.Password, cfg.Db_name)

	//termdictionary.Init()

	http.HandleFunc("/", mainHandlerFunc)
	http.HandleFunc("/search/", searchHandlerFunc)
	http.HandleFunc("/source/", listhandlerFunc)
	http.HandleFunc("/resources/", getResources)
	http.ListenAndServe(":3000", nil)
}

func listhandlerFunc(w http.ResponseWriter, r *http.Request) {

	simplelog.Info.Println("Listhandlerfunc: " + r.URL.Path)
	options := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	p := sourcepage
	p.Title = "Headline tracker. Latest news"
	p.Today = time.Now().Local().Format("02/01 2006")

	switch len(options) {
	case 1:
		p.Articles = db.GetArticles(time.Now())
		break
	case 2:
		source_id, err := strconv.Atoi(options[1])
		if err != nil {
			simplelog.Error.Println( options[1] + " is Not a valid source_id")
		}
		p.Articles = db.GetArticlesFromSource(time.Now(), source_id)
	}
	
	t, _ := template.ParseFiles("./html_templates/list.html", "./html_templates/header.html", "./html_templates/article.html")
	t.Execute(w, p)
}


func mainHandlerFunc(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./html_templates/index.html", "./html_templates/header.html")
	p := sourcepage
	t.Execute(w, p)
}

func searchHandlerFunc(w http.ResponseWriter, r *http.Request) {

}

func getResources(w http.ResponseWriter, r *http.Request) {
	simplelog.Info.Println("calling getResources " + r.RequestURI)
	http.ServeFile(w, r, "." + r.RequestURI)
}

