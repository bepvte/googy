//TRUE MEANS BLOCKED
package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"
	"upper.io/db.v3"
)

const (
	channel = iota + 1
	user
	role
	guild
)

type permrule struct { //where is where the perm takes place, what is on what command it takes place
	Priority int    `db:"priority"`
	State    bool   `db:"state"`
	Where    string `db:"where"`
	What     string `db:"what"`
	Type     int    `db:"type"`
	Guild    string `db:"guild"`
}

var channelregex = regexp.MustCompile("<#(\\d+)>")
var roleregex = regexp.MustCompile("<@&(\\d+)>")

func permCheck(s *discordgo.Session, m *discordgo.MessageCreate, what string) bool {
	c, _ := s.State.Channel(m.ChannelID)
	member, _ := s.State.Member(c.GuildID, m.Author.ID)

	roleconds := make([]db.Compound, 0)
	for _, trole := range member.Roles {
		roleconds = append(roleconds, db.Cond{
			"where": trole,
			"what":  what,
			"type":  role,
		})
	}
	roleconds = append(roleconds,
		db.Cond{
			"where": m.ChannelID,
			"what":  what,
			"type":  channel,
		},
		db.Cond{
			"where": m.Author.ID,
			"what":  what,
			"type":  user,
		},
		db.Cond{
			"where": c.GuildID,
			"what":  what,
			"type":  guild,
		})
	conds := db.Or(roleconds...)

	res := database.Collection("perms").Find(db.And(db.Cond{"guild": c.GuildID}, conds))

	res.OrderBy("-priority")

	var state bool
	var rule permrule

	for res.Next(&rule) {
		state = rule.State
	}

	return state
}

func permWrap(s *discordgo.Session, m *discordgo.MessageCreate, what string, callback func(*discordgo.Session, *discordgo.MessageCreate)) {
	if !permCheck(s, m, what) {
		callback(s, m)
	}
}

const permusage = "Couldnt understand that\n`&add command [enabled|disabled] [priority] [channel-mention|role-mention|user-mention|leave blank for serverwide]`"

//												  1	         2                3    		4
func permAdd(s *discordgo.Session, m *discordgo.MessageCreate) {
	userperm, err := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if !(userperm&discordgo.PermissionManageRoles != 0 || m.Author.ID == "147077474222604288") {
		return
	}
	if err != nil {
		panic(err)
	}
	sv := csv.NewReader(strings.NewReader(m.Content))
	sv.Comma = ' '
	parsed, err := sv.Read()
	if err != nil || len(parsed) < 3 {
		s.ChannelMessageSend(m.ChannelID, permusage)
		return
	}
	var state bool
	{
		switch parsed[2] {
		case "enabled":
			state = false
		case "disabled":
			state = true
		default:
			s.ChannelMessageSend(m.ChannelID, permusage)
			return
		}
	}

	priority, err := strconv.Atoi(parsed[3])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, permusage)
		return
	}
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Println("[PERMADD] wack time finding channel: ", err)
		s.ChannelMessageSend(m.ChannelID, "Something went wrong getting the servers id, try again?")
		return
	}
	perm := permrule{
		Priority: priority,
		State:    state,
		What:     parsed[1],
		Guild:    c.GuildID,
	}
	channelMatches := channelregex.FindStringSubmatch(m.Content)
	roleMatches := roleregex.FindStringSubmatch(m.Content)
	switch {
	case len(channelMatches) >= 2:
		perm.Where = channelMatches[1]
		perm.Type = channel
	case len(m.Mentions) >= 1:
		perm.Where = m.Mentions[0].ID
		perm.Type = user
	case len(roleMatches) >= 2:
		perm.Where = roleMatches[1]
		perm.Type = role
	default:
		perm.Where = c.GuildID
		perm.Type = guild
	}
	if _, err := database.Collection("perms").Insert(perm); err != nil {
		log.Println("[PERMADD] ", err)
		s.ChannelMessageSend(m.ChannelID, "Couldnt access database, try again later")
	}
	s.ChannelMessageSend(m.ChannelID, "Done")
}

func permList(s *discordgo.Session, m *discordgo.MessageCreate) {
	userperm, _ := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if !(userperm&discordgo.PermissionManageRoles != 0 || m.Author.ID == "147077474222604288") {
		return
	}
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Println("[PERMLIST] ", err)
		return
	}
	res := database.Collection("perms").Find("guild", c.GuildID)

	var buf bytes.Buffer
	writ := tabwriter.NewWriter(&buf, 15, 0, 2, ' ', tabwriter.AlignRight)
	fmt.Fprintln(writ, "```csv\nTarget\tCommand\tState\tPriority\tType\t")

	var rule permrule
	for res.Next(&rule) {
		var stype string
		switch rule.Type {
		case channel:
			stype = "channel"
		case guild:
			stype = "guild"
		case role:
			stype = "role"
		case user:
			stype = "user"
		}
		var sstate string
		if rule.State {
			sstate = "disabled"
		} else {
			sstate = "enabled"
		}
		fmt.Fprintf(writ, "%v\t%v\t%v\t%v\t%v\t\n", rule.Where, rule.What, sstate, rule.Priority, stype)
	}
	fmt.Fprintln(writ, "\n```")
	writ.Flush()
	str := string(buf.Bytes())
	for _, x := range split(str, 2000) {
		s.ChannelMessageSend(m.ChannelID, x)
	}
}

func permDel(s *discordgo.Session, m *discordgo.MessageCreate) {
	userperm, _ := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if !(userperm&discordgo.PermissionManageRoles != 0 || m.Author.ID == "147077474222604288") {
		return
	}
	t := strings.Split(m.Content, " ")
	if len(t) < 2 {
		s.ChannelMessageSend(m.ChannelID, "&del <command>")
		return
	}
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Println("[PERMDEL] ", err)
		return
	}
	res, err := database.DeleteFrom("perms").Where("guild = ? AND what = ?", c.GuildID, strings.Join(t[1:], " ")).Exec()
	if err != nil {
		log.Println("[PERMDEL] ", err)
		s.ChannelMessageSend(m.ChannelID, "Database wouldnt accept that...")
		return
	}
	count, err := res.RowsAffected()
	if err != nil {
		log.Println("[PERMDEL] ", err)
		s.ChannelMessageSend(m.ChannelID, "We did it, but the database wouldnt tell us how many were deleted.")
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Deleted %v permrules", count))
}

func split(buf string, lim int) []string {
	var chunk string
	chunks := make([]string, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}
