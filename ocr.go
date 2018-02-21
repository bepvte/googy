//was forced  to last second switc hfrom gosseract to tesseract.v1

package main

import (
	"github.com/bwmarrin/discordgo"
	"gopkg.in/GeertJohan/go.leptonica.v1"
	"gopkg.in/GeertJohan/go.tesseract.v1"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var ocrcl *tesseract.Tess

func ocrInit() {
	tessdata_prefix := os.Getenv("TESSDATA_PREFIX")
	if tessdata_prefix == "" {
		tessdata_prefix = "/usr/share/tesseract-ocr/tessdata"
	}
	var err error
	ocrcl, err = tesseract.NewTess(tessdata_prefix, "eng")
	if err != nil {
		log.Fatalf("Error while initializing Tess: %s\n", err)
	}
	// setup a whitelist of all basic ascii
	err = ocrcl.SetVariable("tessedit_char_whitelist", ` !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_abcdefghijklmnopqrstuvwxyz{|}~`+"`")
	if err != nil {
		log.Fatalf("Failed to SetVariable: %s\n", err)
	}

}

func ocr(s *discordgo.Session, m *discordgo.MessageCreate) {
	defer func() {
		e := recover()
		if e != nil {
			log.Println("[OCR] ", e)
		}
	}()
	var earliestUrl string
	if len(m.Attachments) != 0 {
		earliestUrl = getUrl(m.Message)
	} else if len(strings.Split(m.Content, " ")) == 2 {
		earliestUrl = strings.Split(m.Content, " ")[1]
		if _, err := url.ParseRequestURI(earliestUrl); err != nil {
			//this block gets copy pasted like twenty times bc im lazy and bad
			if _, err := s.ChannelMessageSend(m.ChannelID, "The url you put next to '$ocr' is invalid."); err != nil {
				panic(err)
			}
			return
		}
		log.Println(earliestUrl)
	} else {
		msgs, err := s.ChannelMessages(m.ChannelID, 100, m.ID, "", "")
		if err != nil {
			panic(err)
		}
		for _, x := range msgs {
			u := getUrl(x)
			if u != "" {
				earliestUrl = u
				break
			}
		}
	}
	if earliestUrl == "" {
		if _, err := s.ChannelMessageSend(m.ChannelID, "I cant see any images posted as direct links or attachments in the last 30 messages, and your command didnt include a url or attachment!"); err != nil {
			panic(err)
		}
		return
	}

	//alright time to go
	wd, _ := os.Getwd()
	tmpfile, err := ioutil.TempFile(wd, "googy-")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpfile.Name())
	resp, err := http.Get(earliestUrl)
	if err != nil {
		log.Println("[OCR] Couldnt get: ", earliestUrl)
		if err := s.MessageReactionAdd(m.ChannelID, m.ID, "❌"); err != nil {
			panic(err)
		}
		return
	}
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "image/") {
		log.Println(resp.Request.URL)
		if _, err := s.ChannelMessageSend(m.ChannelID, "The link wont work..."); err != nil {
			panic(err)
		}
		return
	}
	io.Copy(tmpfile, resp.Body)
	resp.Body.Close()

	pix, err := leptonica.NewPixFromFile(filepath.Join(wd, tmpfile.Name()))
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "that image couldnt process....")
	}
	defer pix.Close() // remember to cleanup

	ocrcl.SetImagePix(pix)

	if err != nil {
		log.Println("[OCR] error with ocr: ", err)
		if _, err := s.ChannelMessageSend(m.ChannelID, "OCR failed with error\n```"+err.Error()+"\n```"); err != nil {
			panic(err)
		}
	}
		s.ChannelMessageSend(m.ChannelID, ocrcl.Text())
}

func getUrl(x *discordgo.Message) string {
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
		for i, _ := range lst {
			ret <- lst[len(lst)-1-i]
		}
		close(ret)
	}()
	return ret
}