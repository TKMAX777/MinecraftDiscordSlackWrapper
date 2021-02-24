package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// CommandContent put minecraft command content
type CommandContent struct {
	Command string
	Options string
}

// MinecraftCommand handles minecraft function
type MinecraftCommand struct {
	sendChannel chan CommandContent
}

// Handler handle say commands sent to discord
func (c MinecraftCommand) Handler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.GuildID != Settings.Discord.GuildID {
		return
	}

	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	if m.Message.Content == "" {
		return
	}

	userDict, err := ReadNameDict()
	if err != nil {
		return
	}

	user, ok := userDict.findUserFromDiscordID(m.Author.ID)
	if !ok {
		user, ok = userDict.findUserFromDiscordID("Default")
		if !ok {
			return
		}
	}

	for _, text := range strings.Split(m.Message.Content, "\n") {
		var command CommandContent
		var msg = strings.Split(text, " ")

		if len(msg) < 2 {
			return
		}

		var permissions = GetPermissions(user.PermissionCode)
		command.Command, ok = permissions[msg[0]]
		if !ok {
			return
		}

		switch command.Command {
		case "/msg", "/say":
			msg[1] = fmt.Sprintf("[%s]%s", user.Name, msg[1])
		}

		if msg[1] == ";" {
			msg = msg[:1]
		}

		command.Options = strings.Join(msg[1:], " ")

		fmt.Printf("[Discord]%v\n", text)

		c.sendChannel <- command
	}

	return
}
