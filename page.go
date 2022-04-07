package main

import "github.com/mikkelstb/feedfetcher/feed"

type Page struct {
	Today     string
	Title     string
	Articles  []feed.NewsItem
	Languages map[string]Language
}

type Language struct {
	Iso     string
	Name    string
	Checked bool
}

func (p *Page) setChecked(checked_langs []string) {
	p.Languages = make(map[string]Language)
	p.Languages["dan"] = Language{Iso: "dan", Name: "Danish", Checked: false}
	p.Languages["kor"] = Language{Iso: "kor", Name: "Korean", Checked: false}
	p.Languages["eng"] = Language{Iso: "eng", Name: "English", Checked: false}
	p.Languages["nob"] = Language{Iso: "nob", Name: "Norwegian", Checked: false}

	for _, lang := range checked_langs {
		if lan, ok := p.Languages[lang]; ok {
			lan.Checked = true
			p.Languages[lang] = lan
		}
	}
}
