package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/davecgh/go-spew/spew"
)

var s *discordgo.Session

var bannedFile *os.File
var banned = map[string]bool{}

func main() {
	//token, err := ioutil.ReadFile("token")
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatalln("set env var TOKEN")
	}
	var err error
	bannedFile, err = os.OpenFile("banned", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(bannedFile)
	for scanner.Scan() {
		banned[scanner.Text()] = true
	}
	s, err = discordgo.New("Bot " + string(token))
	if err != nil {
		panic(err)
	}
	if err := s.Open(); err != nil {
		panic(err)
	}
	s.AddHandler(messageCreate)
	s.UpdateStatus(0, "with god.")

	log.Println("We goin")
	c := make(chan interface{})
	<-c
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Author.Bot || banned[m.Author.ID] {
		return
	}
	var prefix = "$"
	if m.GuildID == "449701194881826819" {
		prefix = "."
	}
	switch {
	case strings.ToLower(m.Content) == prefix+"pacman":
		s.ChannelMessageSend(m.ChannelID, "<:pacman:324163173596790786>")
	case isCommand(m.Content, "botban", prefix):
		permissions, err := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
		if err != nil {
			log.Println(err)
			return
		}
		if (permissions&discordgo.PermissionBanMembers) > 0 || m.Author.ID == "147077474222604288" {
			if len(m.Mentions) != 1 {
				s.ChannelMessageSend(m.ChannelID, prefix+"botban <usermention>")
				return
			}
			if _, err := io.WriteString(bannedFile, m.Mentions[0].ID+"\n"); err != nil {
				s.ChannelMessageSend(m.ChannelID, "Error uhh ahhh ahh uuhhhh\n```\n"+spew.Sprint(err)+"\n``` ahh uhhh ahh ahh")
			}
			banned[m.Mentions[0].ID] = true
		}
	case isCommand(m.Content, "ocr", prefix):
		ocr(s, m)
	case isCommand(m.Content, "help", prefix):
		s.ChannelMessageSend(m.ChannelID, "yerm")
	case isCommand(m.Content, "say", prefix):
		if m.Author.ID == os.Getenv("OWNER") {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
			s.ChannelMessageSend(m.ChannelID, strings.TrimPrefix(m.Content, prefix+"say"))
		}
		//case strings.HasPrefix(strings.ToLower(m.Content), prefix+"magick"):
		//	permWrap(s, m, "magick", magick)
		//case strings.HasPrefix(strings.ToLower(m.Content), prefix+"squish"):
		//	permWrap(s, m, "magick", squish)
		//case strings.HasPrefix(strings.ToLower(m.Content), prefix+"squosh"):
		//	permWrap(s, m, "magick", squosh)
	case isCommand(m.Content, "knuckles", prefix):
		s.ChannelMessageSend(m.ChannelID, "CHUCKLES")
	case isCommand(m.Content, "figlet", prefix):
		figlet(s, m)
	case isCommand(m.Content, "nick", prefix):
		c, err := s.Channel(m.ChannelID)
		if err != nil {
			return
		}
		if err := s.GuildMemberNickname(c.GuildID, "@me", strings.TrimPrefix(m.Content, prefix+"nick ")); err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
			log.Println("[NICK] Error:" + err.Error())
		}
	case isCommand(m.Content, "tickle", prefix):
		s.ChannelMessageSend(m.ChannelID, "HEHEHEHEHEHEHE!!!")
	}
}
func isCommand(test, command, prefix string) bool {
	return strings.HasPrefix(strings.ToLower(test), prefix+command)
}
