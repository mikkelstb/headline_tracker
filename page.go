package main

import "github.com/mikkelstb/feedfetcher/feed"

type Page struct {
	Today      string
	Title      string
	SourceName string
	Articles   []feed.NewsItem
}
