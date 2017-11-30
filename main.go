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
	"os"
	"bufio"
)

var s *discordgo.Session

var prefixes = []string{
	"ok google",
	"okay google",
	"hey google",
	"$google", //order prevents double send
	"$g",
}

type result struct {
	url, desc string
}

var banned = map[string]bool{}
var bannedFile *os.File

func main() {
	token, err := ioutil.ReadFile("token")
	if err != nil || string(token) == "" {
		log.Fatalln("YOU FUCKED IT\nMAKE A `token` FILE")
	}
	bannedFile, err = os.OpenFile("banned", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(bannedFile)
	for scanner.Scan() {
		banned[scanner.Text()] = true
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
	if m.Author.ID == s.State.User.ID || m.Author.Bot || banned[m.Author.ID] {
		return
	}
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
				break
			}
		}
	}
	if strings.ToLower(m.Content) == "$pacman" {
		s.ChannelMessageSend(m.ChannelID, "<:pacman:324163173596790786>")
		return
	}
	if strings.ToLower(m.Content) == "$botban" {
		log.Println(m.Author.ID)
		permissions, err := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
		if err != nil {
			log.Println(err)
			return
		}
		if (permissions & discordgo.PermissionBanMembers) > 0 || m.Author.ID == "147077474222604288"{
			log.Println("hi")
			if len(m.Mentions) != 1 {
				s.ChannelMessageSend(m.ChannelID, "$botban <usermention>")
				return
			}
			if _, err := bannedFile.WriteString(m.Mentions[0].ID + "\n"); err != nil {
				s.ChannelMessageSend(m.ChannelID, "failed to write file. contact bot author.")
				return
			}
			banned[m.Mentions[0].ID] = true
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
