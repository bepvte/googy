//was forced  to last second switc hfrom gosseract to tesseract.v1

package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/otiai10/gosseract"
)

const ocrTimeout = time.Second * 10

var ocrRegex = regexp.MustCompile(`\$ocr(\w{3})`)

var ocrLangs = []string{
	"eng",
	"ara",
	"jpn",
	"rus",
	"kor",
	"ell",
	"fra",
}

func ocrInit(lang string) *gosseract.Client {
	ocrcl := gosseract.NewClient()

	tessdataPrefix := os.Getenv("TESSDATA_PREFIX")
	if tessdataPrefix == "" {
		p := "/usr/share/tessdata/" //this sucks fuck u golang
		ocrcl.TessdataPrefix = &p
	}

	// setup a whitelist of all basic ascii
	if lang == "eng" {
		ocrcl.SetWhitelist(` !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_abcdefghijklmnopqrstuvwxyz{|}~` + "`")
	}
	ocrcl.SetLanguage(lang)
	return ocrcl
}

func ocr(s *discordgo.Session, m *discordgo.MessageCreate) {
	c := time.NewTimer(ocrTimeout)
	var lang string
	{
		if len(m.Content) >= 7 {
			langmatch := ocrRegex.FindStringSubmatch(m.Content[0:7])
			if len(langmatch) >= 2 {
				for x := range ocrLangs {
					if langmatch[1] == ocrLangs[x] {
						lang = ocrLangs[x]
					}
				}
			}
		}
	}
	if lang == "" {
		lang = "eng"
		m.Content = strings.Replace(m.Content, "ocr", "ocreng", 1)
	}
	resp := getImage(m, s, "OCR", "ocr"+lang)
	if resp == nil {
		return
	}
	defer resp.Close()
	if checkTimeout(c, m, s) {
		return
	}

	var buf bytes.Buffer

	io.Copy(&buf, resp)

	ocrcl := ocrInit(lang)
	defer ocrcl.Close()

	ocrcl.SetImageFromBytes(buf.Bytes())
	if checkTimeout(c, m, s) {
		return
	}

	t, err := ocrcl.Text()
  t = strings.Replace(t, "@", "@​", -1)
	if err != nil {
		log.Println("[OCR] error: ", err)
		s.ChannelMessageSend(m.ChannelID, "OCR failed with error\n```"+err.Error()+"\n```")
		return
	}
	if t == "" {
		s.ChannelMessageSend(m.ChannelID, "nothing found")

	} else if len(t) >= 900 {
		s.ChannelMessageSend(m.ChannelID, t[:900]+"\nThe rest of the result was too long to display")
	} else {
		s.ChannelMessageSend(m.ChannelID, t)
	}
}

func checkTimeout(time *time.Timer, m *discordgo.MessageCreate, s *discordgo.Session) bool {
	select {
	case <-time.C:
		s.ChannelMessageSend(m.ChannelID, "timed out")
		return true
	default:
		return false
	}
}
