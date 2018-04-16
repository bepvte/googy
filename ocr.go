//was forced  to last second switc hfrom gosseract to tesseract.v1

package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/GeertJohan/go.leptonica.v1"
	"gopkg.in/GeertJohan/go.tesseract.v1"
)

const ocrTimeout = time.Second * 10

func ocrInit() *tesseract.Tess {
	tessdata_prefix := os.Getenv("TESSDATA_PREFIX")
	if tessdata_prefix == "" {
		tessdata_prefix = "/usr/share/tesseract-ocr/tessdata"
	}
	ocrcl, err := tesseract.NewTess(tessdata_prefix, "eng")
	if err != nil {
		panic("Error while initializing Tess: " + err.Error())
	}
	// setup a whitelist of all basic ascii
	err = ocrcl.SetVariable("tessedit_char_whitelist", ` !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_abcdefghijklmnopqrstuvwxyz{|}~`+"`")
	if err != nil {
		panic("Failed to SetVariable: " + err.Error())
	}

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
	wd, _ := os.Getwd()

	tmpfile, err := ioutil.TempFile(wd, "googy-")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpfile.Name())
	io.Copy(tmpfile, resp)

	pix, err := leptonica.NewPixFromFile(filepath.Join(wd, tmpfile.Name()))
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "that image couldnt process....")
		return
	}
	defer pix.Close() // remember to cleanup

	ocrcl := ocrInit()
	defer ocrcl.Close()

	ocrcl.SetImagePix(pix)
	if checkTimeout(c, m, s) {
		return
	}

	if err != nil {
		log.Println("[OCR] error: ", err)
		s.ChannelMessageSend(m.ChannelID, "OCR failed with error\n```"+err.Error()+"\n```")
		return
	}
	t := ocrcl.Text()
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
