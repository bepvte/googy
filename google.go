package main

import (
	"errors"
	"github.com/anaskhan96/soup"
	"golang.org/x/net/html"
	"log"
	"net/url"
	"strings"
)

var prefixes = []string{
	"ok google",
	"okay google",
	"hey google",
	"&google",
	"&g",
	"ok googy",
	"okay googy",
	"hey googy",
}

type result struct {
	url, desc string
}

func google(s string) ([]result, error) {
	var resp string
	var err error
	defer func() {
		e := recover()
		if e != nil {
			//log.Println("contents: "+ resp)
			panic(e)
		}
	}()
	if s == "panictest" {
		panic(errors.New("fof"))
	}
	resp, err = soup.Get("https://www.google.com/search?q=" + url.QueryEscape(s))
	if err != nil {
		return []result{}, errors.New("failed to reach google")
	}
	var results = []result{}
	root := soup.HTMLParse(resp)
	for _, x := range root.FindAll("div", "class", "g") {
		if len(results) > 3 {
			break
		}
		//var buf bytes.Buffer
		//html.Render(&buf, x.Pointer)
		//log.Println(buf.String())
		if x.Attrs()["class"] != "g" {
			continue
		}
		linkMom := x.Find("h3", "class", "r")
		if linkMom.Error != nil {
			log.Println(linkMom.Error)
			continue
		}
		linkTarget := linkMom.Find("a")
		descMom := x.Find("span", "class", "st")
		if descMom.Error != nil {
			log.Println(descMom.Error)
			continue
		}
		descTarget := descMom.Pointer
		var f func(*html.Node)
		descList := []string{}
		f = func(n *html.Node) {
			if n != nil && n.Type == html.TextNode {
				if n.Data != "\n" {
					descList = append(descList, n.Data)
				}
			}
			if n != nil {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					f(c)
				}
			}
		}
		f(descTarget)

		if strings.HasPrefix(linkTarget.Attrs()["href"], "/url?") {
			toTrim, _ := url.Parse("https://www.google.com" + linkTarget.Attrs()["href"])
			results = append(results, result{toTrim.Query().Get("q"), strings.Join(descList, " ")})
		}
	}
	return results, nil
}
