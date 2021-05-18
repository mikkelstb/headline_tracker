package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mikkelstb/headline_tracker/obj"
	"github.com/mikkelstb/simplelog"

	_ "github.com/go-sql-driver/mysql"
)



type Database struct {
	login_credentials string
	initialized bool
}


type source struct {
	name string
	id int
}


func (db *Database) Init(username string, password string, db_name string) {
	db.login_credentials = username + ":" + password + "@/" + db_name
	db.initialized = true
	simplelog.InitLoggers("./ht.log")
}



func (db *Database) GetArticlesFromSource(date time.Time, source_id int) []obj.NewsItem {

	source_name := db.GetSourceName(source_id)
	db_connection, err := sql.Open("mysql", db.login_credentials)

	if err != nil {
		simplelog.Error.Println(err.Error())
	}
	defer db_connection.Close()
	db_connection.SetConnMaxLifetime(time.Second * 5)
	db_connection.SetMaxOpenConns(10)
	db_connection.SetConnMaxIdleTime(10)

	const MySQLDateType string = "2006-01-02 15:04:05"

	query := `select id, headline, story, pubdate, create_time, url from article where pubdate like? and source=? order by pubdate DESC limit 10`
	datestring := date.Format("2006-01-02") + "%"
	rows, err := db_connection.Query(query, datestring, source_id)
	if err != nil {
		simplelog.Error.Println(err.Error())
	}
	defer rows.Close()

	var articles = []obj.NewsItem{}

	for rows.Next() {
		var article = obj.NewsItem{}
		var pubdate string
		var createtime string
		if err := rows.Scan(&article.LocalID, &article.Headline, &article.Story, &pubdate, &createtime, &article.Url); err != nil {
			simplelog.Error.Println(err.Error())
		}
		article.Pubdate, err = time.Parse(MySQLDateType, pubdate)
		article.CollectTime, err = time.Parse(MySQLDateType, createtime)
		article.Source = source_name
		articles = append(articles, article)
	}
	return articles
}




func (db *Database) GetArticles(date time.Time) []obj.NewsItem {

	db_connection, err := sql.Open("mysql", db.login_credentials)

	if err != nil {
		simplelog.Error.Println(err.Error())
	}
	defer db_connection.Close()
	db_connection.SetConnMaxLifetime(time.Second * 5)
	db_connection.SetMaxOpenConns(10)
	db_connection.SetConnMaxIdleTime(10)

	const MySQLDateType string = "2006-01-02 15:04:05"

	query := `select id, headline, story, pubdate, create_time, url, name from articlelist where pubdate like? order by pubdate DESC limit 10`
	datestring := date.Format("2006-01-02") + "%"
	rows, err := db_connection.Query(query, datestring)
	if err != nil {
		simplelog.Error.Println(err.Error())
	}
	defer rows.Close()

	var articles = []obj.NewsItem{}

	for rows.Next() {
		var article = obj.NewsItem{}
		var pubdate string
		var createtime string
		if err := rows.Scan(&article.ID, &article.Headline, &article.Story, &pubdate, &createtime, &article.Url, &article.Source); err != nil {
			simplelog.Error.Println(err.Error())
		}
		article.Pubdate, err = time.Parse(MySQLDateType, pubdate)
		article.CollectTime, err = time.Parse(MySQLDateType, createtime)
		articles = append(articles, article)
	}
	return articles
}

func (db *Database) GetSourceName(id int) string {

	db_connection, err := sql.Open("mysql", db.login_credentials)

	if err != nil {
		simplelog.Error.Println(err.Error())
	}
	defer db_connection.Close()
	db_connection.SetConnMaxLifetime(time.Second * 5)
	db_connection.SetMaxOpenConns(10)
	db_connection.SetConnMaxIdleTime(10)

	const query = `select name from source where id=?`
	var name string

	row := db_connection.QueryRow(query, id)
	err = row.Scan(&name)
	if err != nil {
		log.Println("No name found")
	}
	return name
}

func (db *Database) GetArticlesFromIDs(docIDs []int) []obj.NewsItem {

	db_connection, err := sql.Open("mysql", db.login_credentials)
	if err != nil {
		simplelog.Error.Println(err.Error())
	}
	defer db_connection.Close()
	db_connection.SetConnMaxLifetime(time.Second * 5)
	db_connection.SetMaxOpenConns(10)
	db_connection.SetConnMaxIdleTime(10)

	const MySQLDateType string = "2006-01-02 15:04:05"

	//Convert []int -> new string{2, 5, 6, 7}
	ids_string := strings.Trim(strings.Join(strings.Split(fmt.Sprint(docIDs), " "), ", "), "[]")
	
	const query = `id, headline, story, pubdate, create_time, url from article where id in (?)`
	rows, err := db_connection.Query(query, ids_string)

	if err != nil {
		simplelog.Error.Println(err.Error())
	}
	defer rows.Close()

	var articles = []obj.NewsItem{}

	for rows.Next() {
		var article = obj.NewsItem{}
		var pubdate string
		var createtime string
		if err := rows.Scan(&article.ID, &article.Headline, &article.Story, &pubdate, &createtime, &article.Url, &article.Source); err != nil {
			simplelog.Error.Println(err.Error())
		}
		article.Pubdate, err = time.Parse(MySQLDateType, pubdate)
		article.CollectTime, err = time.Parse(MySQLDateType, createtime)
		articles = append(articles, article)
	}
	return articles
}
