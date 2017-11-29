package main

import (
	"errors"
	"golang.org/x/net/html"

	"github.com/anaskhan96/soup"
	"github.com/bwmarrin/discordgo"

	"io/ioutil"
	"log"
	"net/url"
	"strings"
)

var s *discordgo.Session

var prefixes = []string{
	"ok google",
	"okay google",
	"hey google",
	"!google",
	"!g",
}

type result struct {
	url, desc string
}

func main() {
	token, err := ioutil.ReadFile("token")
	if err != nil || string(token) == "" {
		log.Fatalln("YOU FUCKED IT\nMAKE A `token` FILE")
	}
	s, err = discordgo.New("Bot " + string(token))
	if err != nil {
		panic("aww shit")
	}
	if err := s.Open(); err != nil {
		panic(err)
	}
	s.AddHandler(messageCreate)

	c := make(chan interface{})
	<-c
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var trimmed string
	for x := range prefixes {
		if strings.HasPrefix(strings.ToLower(m.Content), prefixes[x]) {
			trimmed = strings.TrimSpace(strings.TrimPrefix(m.Content[len(prefixes[x]):], ",")) // trimmed it wowow
			log.Println(trimmed)
			result, err := google(trimmed)
			if err != nil {
				log.Println(err)
			} else {
				msg := result[0]
				var resultSanitized []string
				for _, x := range result[1:] {
					resultSanitized = append(resultSanitized, "<"+x.url+">")
				}
				_, err := s.ChannelMessageSend(m.ChannelID, msg.url+" - ```"+msg.desc+"```"+"\n**See also:**\n"+strings.Join(resultSanitized, "\n"))
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func google(s string) ([]result, error) {
	resp, err := soup.Get("https://www.google.com/search?q=" + url.QueryEscape(s))
	if err != nil {
		return []result{}, errors.New("failed to reach google")
	}
	var results = []result{}
	root := soup.HTMLParse(resp)
	for _, x := range root.FindAll("div", "class", "g") {
		if len(results) > 3 {
			break
		}
		linkTarget := x.Find("h3", "class", "r").Find("a")
		descTarget := x.Find("span", "class", "st").Pointer
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
