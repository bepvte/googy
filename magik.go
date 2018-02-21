// +build ignore

package main

import (
	"bytes"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/gographics/imagick.v2/imagick"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func magickInit() {
	imagick.Initialize()
}

func magick(m *discordgo.MessageCreate, s *discordgo.Session) {
	defer func() {
		e := recover()
		if e != nil {
			log.Println("[MAGICK] ", e)
		}
	}()
	var earliestUrl string
	if len(m.Attachments) != 0 {
		earliestUrl = getUrl(m.Message)
	} else if len(strings.Split(m.Content, " ")) == 2 {
		earliestUrl = strings.Split(m.Content, " ")[1]
		if _, err := url.ParseRequestURI(earliestUrl); err != nil {
			s.ChannelMessageSend(m.ChannelID, "The url you put next to '$magick' is invalid.")
			return
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
		return
	}

	var buf bytes.Buffer
	//alright time to go
	resp, err := http.Get(earliestUrl)
	if err != nil {
		log.Println("[MAGICK] Couldnt get: ", earliestUrl)
		s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
		return
	}
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "image/") {
		log.Println(resp.Request.URL)
		s.ChannelMessageSend(m.ChannelID, "The link wont work...")
		return
	}
	io.Copy(&buf, resp.Body)
	resp.Body.Close()

	mw := imagick.NewMagickWand()

	if err := mw.ReadImageBlob(buf.Bytes()); err != nil {
		s.ChannelMessageSend(m.ChannelID, "couldnt parse image")
		return
	}

	// Get original logo size
	//width := mw.GetImageWidth()
	//height := mw.GetImageHeight()

	// Calculate half the size
	//hWidth := uint(width / 2)
	//hHeight := uint(height / 2)

}
