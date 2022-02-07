package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/TKMAX777/MinecraftDiscordWrapper/discord_webhook"
	"github.com/TKMAX777/MinecraftDiscordWrapper/minecraft"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

// DiscordHandler handles Discord conversations
type DiscordHandler struct {
	session *discordgo.Session

	sendChannel chan CommandContent
	idRegExp    *regexp.Regexp

	webhook *discord_webhook.Handler

	serverType string

	settings DiscordSetting
}

func NewDiscordHandler(settings DiscordSetting) *DiscordHandler {
	var handler = &DiscordHandler{
		webhook:  discord_webhook.New(settings.Token),
		settings: settings,
		idRegExp: regexp.MustCompile(`<@!(\d+)>`),
	}

	var err = handler.webhook.SetHookURI(settings.HookURI)
	if err != nil {
		log.Println("Failed to set custom webhook uri", err.Error())
	}

	return handler
}

func (d *DiscordHandler) SetCommandInput(stdin chan CommandContent) *DiscordHandler {
	d.sendChannel = stdin
	return d
}

func (d *DiscordHandler) SetServerType(serverType string) *DiscordHandler {
	d.serverType = serverType
	return d
}

func (d *DiscordHandler) StartSession() error {
	session, err := discordgo.New(d.settings.Token)
	if err != nil {
		return errors.Wrap(err, "StartSession")
	}
	d.session = session

	d.session.AddHandler(d.getMessage)
	return errors.Wrap(d.session.Open(), "OpeningDiscordSession")
}

func (d *DiscordHandler) SendMessageFunction() MessageSender {
	var onlineUserNum int

	return MessageSender(func(message minecraft.Message) error {
		var content string
		switch message.Type {
		case minecraft.MessageTypeJoin:
			if d.settings.SendOption&(SendSettingJoinLeft|SendSettingAll) == 0 {
				return nil
			}

			onlineUserNum++
			if d.settings.AddOnlineNumber {
				switch onlineUserNum {
				case 0, 1:
					content = fmt.Sprintf("%s `%s joined the game`\nActive: %d player", d.settings.Reaction.Join, message.User, onlineUserNum)
				default:
					content = fmt.Sprintf("%s `%s joined the game`\nActive: %d players", d.settings.Reaction.Join, message.User, onlineUserNum)
				}
			} else {
				content = fmt.Sprintf("%s `%s joined the game`", d.settings.Reaction.Join, message.User)
			}
		case minecraft.MessageTypeLeft:
			if d.settings.SendOption&(SendSettingJoinLeft|SendSettingAll) == 0 {
				return nil
			}

			onlineUserNum--
			if d.settings.AddOnlineNumber {
				switch onlineUserNum {
				case 0, 1:
					content = fmt.Sprintf("%s `%s left the game`\nActive: %d player", d.settings.Reaction.Left, message.User, onlineUserNum)
				default:
					content = fmt.Sprintf("%s `%s left the game`\nActive: %d players", d.settings.Reaction.Left, message.User, onlineUserNum)
				}
			} else {
				content = fmt.Sprintf("%s `%s left the game`", d.settings.Reaction.Left, message.User)
			}
		case minecraft.MessageTypeThreadINFO:
			if d.settings.SendOption&(SendSettingThreadINFO|SendSettingAll) == 0 {
				return nil
			}

			content = message.Message
		case minecraft.MessageTypeDeath:
			if d.settings.SendOption&(SendSettingDeath|SendSettingAll) == 0 {
				return nil
			}

			content = fmt.Sprintf("%s %s", d.settings.Reaction.Death, message.Message)
		case minecraft.MessageTypeReachedTheAdvancement:
			if d.settings.SendOption&(SendSettingReachedTheAdvancement|SendSettingAll) == 0 {
				return nil
			}

			content = fmt.Sprintf("%s %s", d.settings.Reaction.Advancement, message.Message)
		case minecraft.MessageTypeMessage:
			if d.settings.SendOption&(SendSettingMessage|SendSettingAll) == 0 {
				return nil
			}

			content = message.Message
		case minecraft.MessageTypeOther:
			if d.settings.SendOption&SendSettingAll == 0 {
				return nil
			}

			content = message.Message
		}

		_, err := d.webhook.Send(d.settings.ChannelID, discord_webhook.Message{
			UserName:  d.settings.UserName,
			AvaterURL: d.settings.AvaterURI,
			Content:   content,
		}, false, nil)

		return err
	})
}

// Handler handle say commands sent to discord
func (d *DiscordHandler) getMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.GuildID != d.settings.GuildID || m.ChannelID != d.settings.ChannelID {
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

	var user User
	user, ok := userDict.findUserFromDiscordID(m.Author.ID)
	if !ok {
		if d.settings.Permissions != 0 {
			user = User{
				PermissionCode:  d.settings.Permissions,
				SendAllMessages: d.settings.SendAllMessages,
			}
		} else {
			user, ok = userDict.findUserFromDiscordID("Default")
			if !ok {
				return
			}
		}

		user.Name = m.Member.Nick
		if user.Name == "" {
			user.Name = m.Member.User.Username
		}
	}

	for _, text := range strings.Split(m.Message.Content, "\n") {
		var command CommandContent
		var msg = strings.Split(text, " ")

		var permissions = GetPermissions(user.PermissionCode)
		command.Command, ok = permissions[msg[0]]
		if !ok {
			_, ok = permissions["say"]
			if !ok || !user.SendAllMessages {
				return
			}
			msg = append([]string{"say"}, msg...)
			command.Command = "/say"
		}

		if d.serverType == "paper" {
			command.Command = strings.TrimPrefix(command.Command, "/")
		}

		if len(msg) < 2 {
			return
		}

		if msg[1] == ";" {
			msg = msg[:1]
		}

		switch command.Command {
		case "/msg", "/say", "msg", "say":
			msg[1] = fmt.Sprintf("[%s]%s", user.Name, msg[1])
			command.Options = strings.Join(msg[1:], " ")

			for _, match := range d.idRegExp.FindAllStringSubmatch(command.Options, -1) {
				u, ok := userDict.findUserFromDiscordID(match[1])
				if !ok {
					continue
				}
				command.Options = strings.ReplaceAll(command.Options, "!"+u.DiscordID, u.Name)
			}

		default:
			command.Options = strings.Join(msg[1:], " ")
		}

		fmt.Printf("[Discord]%v\n", text)

		d.sendChannel <- command
	}
}
