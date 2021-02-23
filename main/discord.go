package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Message put message data
type Message struct {
	Content string
	Name    string
}

// Say handles say function
type Say struct {
	sendChannel chan Message
}

// New make send message from text and discord ID
func (m *Message) New(text, discordID string) error {
	var err error

	users, err := ReadNameDict()
	if err != nil {
		return err
	}

	user, ok := users.findNameFromDiscordID(discordID)
	if !ok {
		return fmt.Errorf("Undefined user")
	}

	m.Content = text
	m.Name = user

	return nil
}

// Handler handle say commands sent to discord
func (say Say) Handler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.GuildID != Settings.Discord.GuildID {
		return
	}

	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	if m.Message.Content == "" {
		return
	}

	var text string
	{
		var msg = strings.Split(
			strings.TrimSpace(m.Message.Content), "say ",
		)

		if len(msg) < 2 || msg[0] != "" {
			return
		}
		text = msg[1]
	}

	var message Message

	err := message.New(text, m.Author.ID)
	if err != nil {
		return
	}
	fmt.Printf("%#v\n", m.Message.Content)

	say.sendChannel <- message

	return
}
