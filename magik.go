package main

import (
	"bytes"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/gographics/imagick.v2/imagick"
	"io"
	"log"
)

func magickInit() {
	imagick.Initialize()
}

func magick(s *discordgo.Session, m *discordgo.MessageCreate) {
	r := getImage(m, s, "[MAGICK]", "$magick")
	if r == nil {
		return
	}
	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()

	mw := imagick.NewMagickWand()

	if err := mw.ReadImageBlob(buf.Bytes()); err != nil {
		s.ChannelMessageSend(m.ChannelID, "couldnt parse image")
		return
	}
	width := mw.GetImageWidth()
	height := mw.GetImageHeight()

	hWidth := uint(width / 2)
	hHeight := uint(height / 2)
	msg, _ := s.ChannelMessageSend(m.ChannelID, "Working...")
	if err := mw.LiquidRescaleImage(hHeight, hWidth, 0, 0); err != nil {
		log.Println("[MAGICK] ", err)
		return
	}

	res := mw.GetImageBlob()
	s.ChannelMessageDelete(m.ChannelID, msg.ID)
	if _, err := s.ChannelFileSend(m.ChannelID, "magick."+mw.GetImageFormat(), bytes.NewBuffer(res)); err != nil {
		log.Println("[MAGICK] ", err)
		return
	}
}
func squish(s *discordgo.Session, m *discordgo.MessageCreate) {
	r := getImage(m, s, "[MAGICK]", "$magick")
	if r == nil {
		return
	}
	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()

	mw := imagick.NewMagickWand()

	if err := mw.ReadImageBlob(buf.Bytes()); err != nil {
		s.ChannelMessageSend(m.ChannelID, "couldnt parse image")
		return
	}
	width := mw.GetImageWidth()
	height := mw.GetImageHeight()

	hWidth := uint(width / 2)
	msg, _ := s.ChannelMessageSend(m.ChannelID, "Working...")
	if err := mw.LiquidRescaleImage(height, hWidth, 0, 0); err != nil {
		log.Println("[MAGICK] ", err)
		return
	}

	res := mw.GetImageBlob()
	s.ChannelMessageDelete(m.ChannelID, msg.ID)
	if _, err := s.ChannelFileSend(m.ChannelID, "magick."+mw.GetImageFormat(), bytes.NewBuffer(res)); err != nil {
		log.Println("[MAGICK] ", err)
		return
	}
}
func squosh(s *discordgo.Session, m *discordgo.MessageCreate) {
	r := getImage(m, s, "[MAGICK]", "$magick")
	if r == nil {
		return
	}
	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()

	mw := imagick.NewMagickWand()

	if err := mw.ReadImageBlob(buf.Bytes()); err != nil {
		s.ChannelMessageSend(m.ChannelID, "couldnt parse image")
		return
	}
	width := mw.GetImageWidth()
	height := mw.GetImageHeight()

	hHeight := uint(height / 2)
	msg, _ := s.ChannelMessageSend(m.ChannelID, "Working...")
	if err := mw.LiquidRescaleImage(hHeight, width, 0, 0); err != nil {
		log.Println("[MAGICK] ", err)
		return
	}

	res := mw.GetImageBlob()
	s.ChannelMessageDelete(m.ChannelID, msg.ID)
	if _, err := s.ChannelFileSend(m.ChannelID, "magick."+mw.GetImageFormat(), bytes.NewBuffer(res)); err != nil {
		log.Println("[MAGICK] ", err)
		return
	}
}