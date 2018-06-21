//was forced  to last second switc hfrom gosseract to tesseract.v1

package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/otiai10/gosseract"
)

const ocrTimeout = time.Second * 10

func ocrInit() *gosseract.Client {
	ocrcl := gosseract.NewClient()

	tessdata_prefix := os.Getenv("TESSDATA_PREFIX")
	if tessdata_prefix == "" {
		p := "/usr/share/tessdata/" //this sucks fuck u golang
		ocrcl.TessdataPrefix = &p
	}

	// setup a whitelist of all basic ascii
	ocrcl.SetWhitelist(` !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_abcdefghijklmnopqrstuvwxyz{|}~` + "`")
	return ocrcl
}

func ocr(s *discordgo.Session, m *discordgo.MessageCreate) {
	c := time.NewTimer(ocrTimeout)
	resp := getImage(m, s, "OCR", "ocr")
	if resp == nil {
		return
	}
	defer resp.Close()
	if checkTimeout(c, m, s) {
		return
	}

	var buf bytes.Buffer

	io.Copy(&buf, resp)

	ocrcl := ocrInit()
	defer ocrcl.Close()

	ocrcl.SetImageFromBytes(buf.Bytes())
	if checkTimeout(c, m, s) {
		return
	}

	t, err := ocrcl.Text()
	if err != nil {
		log.Println("[OCR] error: ", err)
		s.ChannelMessageSend(m.ChannelID, "OCR failed with error\n```"+err.Error()+"\n```")
		return
	}
	if t == "" {
		s.ChannelMessageSend(m.ChannelID, "nothing found")

	} else if len(t) >= 400 {
		s.ChannelMessageSend(m.ChannelID, t[:400]+"\nThe rest of the result was too long to display")
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
