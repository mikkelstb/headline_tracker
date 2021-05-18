package obj

import "time"

//NewsItem is a general structure for all articles in the system
type NewsItem struct {
	ID int
	LocalID int
	Source string
	Headline string
	Intro string
	Story string
	Bylines []Byline
	Pubdate time.Time
	CollectTime time.Time
	Categories []Category
	Url string
}

type Byline struct {
	Name string
	Title string
	Email string
}

type Category struct {
	Name string
	id int
}


func (ni *NewsItem) LocalPubdate() string {
	return ni.Pubdate.Format("2001 02-03 15:04:05")
}
