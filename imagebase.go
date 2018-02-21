package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func getImage(m *discordgo.MessageCreate, s *discordgo.Session, modname, cmdname string) io.ReadCloser {
		defer func() {
			e := recover()
			if e != nil {
				log.Println(fmt.Sprintf("[%v] ", modname), e)
			}
		}()
		var earliestUrl string
		if len(m.Attachments) != 0 {
			earliestUrl = getUrl(m.Message)
		} else if len(strings.Split(m.Content, " ")) == 2 {
			earliestUrl = strings.Split(m.Content, " ")[1]
			if _, err := url.ParseRequestURI(earliestUrl); err != nil {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("The url you put next to '$%v' is invalid.", cmdname))
				return nil
			}
		} else {
			msgs, err := s.ChannelMessages(m.ChannelID, 100, m.ID, "", "")
			if err != nil {
				panic(err)
			}
			for _, x := range msgs {
				u := getUrl(x)
				if u != "" {
					earliestUrl = u
					log.Println(u)
					break
				}
			}
		}
		if earliestUrl == "" {
			s.ChannelMessageSend(m.ChannelID, "Your command didnt include a url or attachment!")
			return nil
		}

		//alright time to go


		resp, err := http.Get(earliestUrl)
		if err != nil {
			log.Println(fmt.Sprintf("[%v] Couldnt get: %v", modname, earliestUrl))
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