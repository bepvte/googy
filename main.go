package main

import (
	"database/sql"
	"github.com/bwmarrin/discordgo"
	"github.com/davecgh/go-spew/spew"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
)

var s *discordgo.Session

var db *sql.DB

var banned = map[string]bool{}

func main() {
	//token, err := ioutil.ReadFile("token")
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatalln("MAKE A `token` FILE")
	}
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS store (myid text PRIMARY KEY)"); err != nil {
		panic(err)
	}
	{
		rows, err := db.Query("SELECT myid FROM store")
		if err != nil {
			panic(err)
		}
		for rows.Next() {
			var x string
			if err := rows.Scan(&x); err != nil {
				panic(err)
			}
			banned[x] = true
		}
	}
	//bannedFile, err = os.OpenFile("banned", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	//if err != nil {
	//	panic(err)
	//}
	//scanner := bufio.NewScanner(bannedFile)
	//for scanner.Scan() {
	//	banned[scanner.Text()] = true
	//}
	s, err = discordgo.New("Bot " + string(token))
	if err != nil {
		panic(err)
	}
	if err := s.Open(); err != nil {
		panic(err)
	}
	s.AddHandler(messageCreate)

	ocrInit()

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
		var trimmed string
		for x := range prefixes {
			if strings.HasPrefix(strings.ToLower(m.Content), prefixes[x]) {
				trimmed = strings.TrimSpace(strings.TrimPrefix(m.Content[len(prefixes[x]):], ",")) // trimmed it wowow
				log.Println(trimmed)
				defer func() {
					e := recover()
					if e != nil {
						log.Println("Panic caught: ", spew.Sprint(e))
						s.ChannelMessageSend(m.ChannelID, "Error uhh ahhh ahh uuhhhh\n```\n"+spew.Sprint(e)+"\n``` ahh uhhh ahh ahh")
					}
				}()
				result, err := google(trimmed)
				if err != nil {
					panic(err)
				} else {
					msg := result[0]
					var resultSanitized []string
					for _, x := range result[1:] {
						resultSanitized = append(resultSanitized, "<"+x.url+">")
					}
					_, err := s.ChannelMessageSend(m.ChannelID, msg.url+" - ```"+msg.desc+"```"+"\n**See also:**\n"+strings.Join(resultSanitized, "\n"))
					if err != nil {
						log.Println(err)
					}
					break
				}
			}
		}
	case strings.ToLower(m.Content) == "$pacman":
		s.ChannelMessageSend(m.ChannelID, "<:pacman:324163173596790786>")
	case strings.HasPrefix(strings.ToLower(m.Content), "$botban"):
		permissions, err := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
		if err != nil {
			log.Println(err)
			return
		}
		if (permissions&discordgo.PermissionBanMembers) > 0 || m.Author.ID == "147077474222604288" {
			if len(m.Mentions) != 1 {
				s.ChannelMessageSend(m.ChannelID, "$botban <usermention>")
				return
			}
			if _, err := db.Exec("INSERT INTO store VALUES ($1)", m.Mentions[0].ID); err != nil {
				s.ChannelMessageSend(m.ChannelID, "failed to ban user. maybe they are already banned?")
				return
			}
			banned[m.Mentions[0].ID] = true
		}
	case strings.HasPrefix(strings.ToLower(m.Content), "$ocr"):
		ocr(s, m)
	}

}
