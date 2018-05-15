// +build ignore

//ITS BAD
package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/fogleman/gg"
)

func meme(s *discordgo.Session, m *discordgo.MessageCreate) {
	sv := csv.NewReader(strings.NewReader(m.Content))
	sv.Comma = ' '
	parsed, err := sv.Read()
	if err != nil || len(parsed) < 3 {
		s.ChannelMessageSend(m.ChannelID, permusage)
		return
	}

	defer func() {
		e := recover()
		if e != nil {
			log.Println("[MEME] ", e)
		}
	}()
	var earliestUrl string
	if len(m.Attachments) != 0 {
		earliestUrl = getUrl(m.Message)
	} else if len(parsed) == 4 {
		earliestUrl = parsed[3]
		if _, err := url.ParseRequestURI(earliestUrl); err != nil {
			s.ChannelMessageSend(m.ChannelID, "The url you put next to '&meme' is invalid.")
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

	//alright time to go

	resp, err := http.Get(earliestUrl)
	if err != nil {
		log.Println(fmt.Sprintf("[MEME] Couldnt get: %v", earliestUrl))
		s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
		return
	}
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "image/") {
		log.Println(resp.Request.URL)
		s.ChannelMessageSend(m.ChannelID, "The link wont work...")
		return
	}

	i, _, err := image.Decode(resp.Body)
	if err != nil {
		log.Println(resp.Request.URL)
		s.ChannelMessageSend(m.ChannelID, "The link wont work...")
		return
	}
	resp.Body.Close()

	dc := gg.NewContextForImage(i)
	ppep := i.Bounds().Size()
	if err := dc.LoadFontFace("font/impact.ttf", 96); err != nil {
		panic(err)
	}
	dc.SetRGB(0, 0, 0)
	n := 6 // "stroke" size
	for dy := -n; dy <= n; dy++ {
		for dx := -n; dx <= n; dx++ {
			if dx*dx+dy*dy >= n*n {
				// give it rounded corners
				continue
			}
			x := float64(ppep.X)/2 + float64(dx)
			y := float64(ppep.Y)/2 + float64(dy)
			dc.DrawStringAnchored(parsed[1], x, y, 0.5, 0.5)
		}
	}
	dc.SetRGB(1, 1, 1)
	dc.DrawStringAnchored(parsed[1], float64(ppep.X)/2, float64(ppep.Y)/2, 0.5, 0.5)
	var buf bytes.Buffer
	dc.EncodePNG(&buf)
	s.ChannelFileSend(m.ChannelID, "meme.png", buf.Bytes())
}
