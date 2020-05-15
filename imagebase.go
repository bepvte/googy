package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func getImage(m *discordgo.MessageCreate, s *discordgo.Session, modname, cmdname string) io.ReadCloser {
	defer func() {
		e := recover()
		if e != nil {
			log.Println(fmt.Sprintf("[%v] ", modname), e)
		}
	}()
	var earliestURL string
	if len(m.Attachments) != 0 {
		earliestURL = getURL(m.Message)
	} else if len(strings.Split(m.Content, " ")) == 2 {
		earliestURL = strings.Split(m.Content, " ")[1]
		if _, err := url.ParseRequestURI(earliestURL); err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("The url you put next to '$%v' is invalid.", cmdname))
			return nil
		}
	} else {
		msgs, err := s.ChannelMessages(m.ChannelID, 100, m.ID, "", "")
		if err != nil {
			panic(err)
		}
		for _, x := range msgs {
			u := getURL(x)
			if u != "" {
				earliestURL = u
				break
			}
		}
	}
	if earliestURL == "" {
		s.ChannelMessageSend(m.ChannelID, "Your command didnt include a url or attachment!")
		return nil
	}

	//alright time to go

	resp, err := http.Get(earliestURL)
	if err != nil {
		log.Println(fmt.Sprintf("[%v] Couldnt get: %v", modname, earliestURL))
		s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
		return nil
	}
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "image/") {
		log.Println(resp.Request.URL)
		s.ChannelMessageSend(m.ChannelID, "The link wont work...")
		return nil
	}
	return resp.Body
}

func getURL(x *discordgo.Message) string {
	if len(x.Attachments) != 0 {
		return x.Attachments[0].URL
	}
	if x.Content != "" {
		if _, err := url.ParseRequestURI(x.Content); err == nil {
			if h, err := http.Head(x.Content); err == nil && strings.HasPrefix(h.Header.Get("Content-Type"), "image/") {
				return x.Content
			}
		}
	}
	return ""
}

func reverse(lst []*discordgo.Message) chan *discordgo.Message {
	ret := make(chan *discordgo.Message)
	go func() {
		for i := range lst {
			ret <- lst[len(lst)-1-i]
		}
		close(ret)
	}()
	return ret
}
