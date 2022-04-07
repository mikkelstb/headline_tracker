package main

import (
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mikkelstb/feedfetcher/feed"
)

type SQLite struct {
	db *sql.DB
}

func NewSQLite(filename string) (*SQLite, error) {
	sq := new(SQLite)

	db_config, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}

	sq.db = db_config
	return sq, nil
}

func (r *SQLite) String() string {
	return "sqlite db"
}

func (r *SQLite) WriteSingle(a feed.NewsItem) (string, error) {
	_, err := r.db.Exec("INSERT INTO newsitem (docdate, id, source, headline, story, url) values(?,?,?,?,?,?)", a.Docdate, a.GetId(), a.FeedId, a.Headline, a.Story, a.Url)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (r *SQLite) GetArticles(nb_articles int, languages []string) ([]feed.NewsItem, error) {
	var articles []feed.NewsItem = make([]feed.NewsItem, 0, nb_articles)

	base_query := `select newsitem.docdate, newsitem.id, newsitem.headline, newsitem.story, newsitem.url, source.screen_name from newsitem left join source on newsitem.source=source.id`
	lang_query := `where source.language = '` + strings.Join(languages, `' or source.language = '`) + `'`
	sort_query := `order by docdate desc limit ?`

	var query string

	if len(languages) > 0 {
		query = strings.Join([]string{base_query, lang_query, sort_query}, " ")
	} else {
		query = strings.Join([]string{base_query, sort_query}, " ")
	}

	//fmt.Println(query)

	//rows, err := r.db.Query("select newsitem.docdate, newsitem.id, newsitem.headline, newsitem.story, newsitem.url, source.screen_name from newsitem left join source on newsitem.source=source.id order by docdate desc limit ?", nb_articles)
	rows, err := r.db.Query(query, nb_articles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, docdate, headline, story, url, source_name sql.NullString
		err := rows.Scan(&docdate, &id, &headline, &story, &url, &source_name)
		if err != nil {
			return nil, err
		}
		articles = append(articles, feed.NewsItem{Id: id.String, Docdate: docdate.String, Source: source_name.String, Headline: headline.String, Story: story.String, Url: url.String})
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return articles, nil
}

func (r *SQLite) Close() error {
	return r.db.Close()
}
