package main

import (
	"database/sql"

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

func (r *SQLite) GetArticles(nb_articles int) ([]feed.NewsItem, error) {
	var articles []feed.NewsItem = make([]feed.NewsItem, 0, nb_articles)

	rows, err := r.db.Query("select * from newsitem order by docdate desc limit ?", nb_articles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, source, docdate, headline, story, url string
		err := rows.Scan(&docdate, &id, &source, &headline, &story, &url)
		if err != nil {
			return nil, err
		}
		articles = append(articles, feed.NewsItem{Id: id, Docdate: docdate, Source: source, Headline: headline, Story: story, Url: url})
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return articles, nil
}

func (r *SQLite) Close() error {
	return r.db.Close()
}
