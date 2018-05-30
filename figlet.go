package main

import (
	"github.com/bwmarrin/discordgo"
	"strings"
	"os/exec"
	"bytes"
	"log"
)

func figlet(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(strings.Replace(m.Content, " ", "", -1)) > 10 || len(m.Content) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Usage: "+prefix+"figlet [less then 10 letters]")
		return
	}
	c := exec.Command("figlet", m.Content)
	var out bytes.Buffer
	c.Stdout = &out
	if err := c.Run(); err != nil {
		log.Printf("[FIGLET] Error: %v\n", err)
		return
	}
	s.ChannelMessageSend(m.ChannelID, "```\n"+strings.Replace(out.String(), "`", "\\`", -1)+"\n```")
}