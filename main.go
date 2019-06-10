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

const prefix = "$"

var prefixes = []string{
	"ok google",
	"okay google",
	"hey google",
	prefix + "google",
	prefix + "g",
	"ok googy",
	"okay googy",
	"hey googy",
}

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
	switch {
	default:
		// var trimmed string
		for x := range prefixes {
			if strings.HasPrefix(strings.ToLower(m.Content), prefixes[x]+" ") {
				s.ChannelMessageSend(m.ChannelID, "use scruntddg you freak!!!!!!!! good bye")

				// 		trimmed = strings.TrimSpace(strings.TrimPrefix(m.Content[len(prefixes[x]+" "):], ",")) // trimmed it wowow
				// 		log.Println(trimmed)
				// 		defer func() {
				// 			e := recover()
				// 			if e != nil {
				// 				log.Println("Panic caught: ", spew.Sprint(e))
				// 				s.ChannelMessageSend(m.ChannelID, "Error uhh ahhh ahh uuhhhh\n```\n"+spew.Sprint(e)+"\n``` ahh uhhh ahh ahh")
				// 			}
				// 		}()
				// 		result, err := google(trimmed)
				// 		if err != nil {
				// 			panic(err)
				// 		} else {
				// 			msg := result[0]
				// 			var resultSanitized []string
				// 			for _, x := range result[1:] {
				// 				resultSanitized = append(resultSanitized, "<"+x.Link+">")
				// 			}
				// 			var desc string
				// 			if msg.Description != "" {
				// 				desc = " ```" + msg.Description + "```"
				// 			}
				// 			_, err := s.ChannelMessageSend(m.ChannelID, msg.Link+desc+"\n**See also:**\n"+strings.Join(resultSanitized, "\n"))
				// 			if err != nil {
				// 				log.Println(err)
				// 			}
				// }
				break

			}
		}
	case strings.ToLower(m.Content) == prefix+"pacman":
		s.ChannelMessageSend(m.ChannelID, "<:pacman:324163173596790786>")
	case isCommand(m.Content, "botban"):
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
	case isCommand(m.Content, "ocr"):
		ocr(s, m)
	case isCommand(m.Content, "help"):
		s.ChannelMessageSend(m.ChannelID, "yerm")
	case isCommand(m.Content, "say"):
		if m.Author.ID == os.Getenv("OWNER") {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
			s.ChannelMessageSend(m.ChannelID, strings.TrimPrefix(m.Content, prefix+"say"))
		}
		// case isCommand(m.Content, "add"):
		// 	permAdd(s, m)
		// case isCommand(m.Content, "perms"):
		// 	permList(s, m)
		// case isCommand(m.Content, "del"):
		// 	permDel(s, m)
		//case strings.HasPrefix(strings.ToLower(m.Content), prefix+"magick"):
		//	permWrap(s, m, "magick", magick)
		//case strings.HasPrefix(strings.ToLower(m.Content), prefix+"squish"):
		//	permWrap(s, m, "magick", squish)
		//case strings.HasPrefix(strings.ToLower(m.Content), prefix+"squosh"):
		//	permWrap(s, m, "magick", squosh)
	case isCommand(m.Content, "knuckles"):
		s.ChannelMessageSend(m.ChannelID, "CHUCKLES")
	case isCommand(m.Content, "figlet"):
		figlet(s, m)
	case isCommand(m.Content, "nick"):
		c, err := s.Channel(m.ChannelID)
		if err != nil {
			return
		}
		if err := s.GuildMemberNickname(c.GuildID, "@me", strings.TrimPrefix(m.Content, prefix+"nick ")); err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
			log.Println("[NICK] Error:" + err.Error())
		}
	case isCommand(m.Content, "tickle"):
		s.ChannelMessageSend(m.ChannelID, "HEHEHEHEHEHEHE!!!")
	}
}
func isCommand(test, command string) bool {
	return strings.HasPrefix(strings.ToLower(test), prefix+command)
}
