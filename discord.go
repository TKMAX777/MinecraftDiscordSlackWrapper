package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/TKMAX777/MinecraftDiscordSlackWrapper/discord_webhook"
	"github.com/TKMAX777/MinecraftDiscordSlackWrapper/mcheads"
	"github.com/TKMAX777/MinecraftDiscordSlackWrapper/minecraft"
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
	joinState  *JoinState

	lastMessageID string
	stateMu       sync.Mutex
	closed        bool
	avatars       map[string]cachedStateAvatar

	settings DiscordSetting
}

type cachedStateAvatar struct {
	joinedAt time.Time
	head     image.Image
}

func NewDiscordHandler(settings DiscordSetting, joinState *JoinState) *DiscordHandler {
	var handler = &DiscordHandler{
		webhook:   discord_webhook.New(settings.Token),
		settings:  settings,
		joinState: joinState,
		idRegExp:  regexp.MustCompile(`<@!(\d+)>`),
	}

	var err = handler.webhook.SetHookURI(settings.HookURI)
	if err != nil {
		log.Println("Failed to set custom webhook uri", err.Error())
	}

	return handler
}

func (d *DiscordHandler) Close() error {
	d.stateMu.Lock()
	defer d.stateMu.Unlock()
	d.closed = true
	if d.session != nil {
		d.session.Close()
	}
	if d.settings.SendJoinStateMessage && d.lastMessageID != "" {
		return d.webhook.Delete(d.settings.ChannelID, d.lastMessageID)
	}
	return nil
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
	session, err := discordgo.New("Bot " + d.settings.Token)
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
		var dMessage = discord_webhook.Message{
			UserName:  d.settings.UserName,
			AvaterURL: d.settings.AvaterURI,
		}
		switch message.Type {
		case minecraft.MessageTypeJoin:
			if d.settings.SendJoinStateMessage {
				err := d.sendUserState(message.Type)
				if err != nil {
					return err
				}
			}

			if !d.settings.SendOption.All && !d.settings.SendOption.JoinLeft {
				return nil
			}

			onlineUserNum++
			if d.settings.AddOnlineNumber {
				switch onlineUserNum {
				case 0, 1:
					dMessage.Content = fmt.Sprintf(":%s: `%s joined the game`\nActive: %d player", d.settings.Reaction.Join, message.User, onlineUserNum)
				default:
					dMessage.Content = fmt.Sprintf(":%s: `%s joined the game`\nActive: %d players", d.settings.Reaction.Join, message.User, onlineUserNum)
				}
			} else {
				dMessage.Content = fmt.Sprintf(":%s: `%s joined the game`", d.settings.Reaction.Join, message.User)
			}
		case minecraft.MessageTypeLeft:
			if d.settings.SendJoinStateMessage {
				err := d.sendUserState(message.Type)
				if err != nil {
					return err
				}
			}

			if !(d.settings.SendOption.All || d.settings.SendOption.JoinLeft) {
				return nil
			}

			onlineUserNum--
			if d.settings.AddOnlineNumber {
				switch onlineUserNum {
				case 0, 1:
					dMessage.Content = fmt.Sprintf(":%s: `%s left the game`\nActive: %d player", d.settings.Reaction.Left, message.User, onlineUserNum)
				default:
					dMessage.Content = fmt.Sprintf(":%s: `%s left the game`\nActive: %d players", d.settings.Reaction.Left, message.User, onlineUserNum)
				}
			} else {
				dMessage.Content = fmt.Sprintf(":%s: `%s left the game`", d.settings.Reaction.Left, message.User)
			}
		case minecraft.MessageTypeThreadINFO:
			if !(d.settings.SendOption.All || d.settings.SendOption.ThreadINFO) {
				return nil
			}

			dMessage.Content = message.Message
		case minecraft.MessageTypeVillagerDeath:
			if !(d.settings.SendOption.All || d.settings.SendOption.VillagerDeath) {
				return nil
			}

			content := message.Content.(minecraft.MessageContentVillagerDeath)
			dMessage.Content = fmt.Sprintf(
				":%s: Villager: %s (ID: %d) died\n%s",
				d.settings.Reaction.VillagerDeath,
				content.Job,
				content.ID,
				content.DiedMessage,
			)
		case minecraft.MessageTypeDeath:
			if !(d.settings.SendOption.All || d.settings.SendOption.Death) {
				return nil
			}

			dMessage.Content = fmt.Sprintf(":%s: %s", d.settings.Reaction.Death, message.Message)
		case minecraft.MessageTypeReachedTheAdvancement:
			if !(d.settings.SendOption.All || d.settings.SendOption.ReachedTheAdvancement) {
				return nil
			}

			dMessage.Content = fmt.Sprintf(":%s: %s", d.settings.Reaction.Advancement, message.Message)
		case minecraft.MessageTypeMessage:
			if !(d.settings.SendOption.All || d.settings.SendOption.Message) {
				return nil
			}
			dMessage.UserName = message.User
			if avatarURI, err := mcheads.GetCurrentAvaterURI(message.User); err == nil {
				dMessage.AvaterURL = avatarURI
			} else {
				log.Printf("Discord chat avatar for %s: %v", message.User, err)
			}
			dMessage.Content = message.Message
		case minecraft.MessageTypeServermessage:
			if !(d.settings.SendOption.All || d.settings.SendOption.Message) {
				return nil
			}

			dMessage.Content = message.Message
		case minecraft.MessageTypeOther:
			if !d.settings.SendOption.All {
				return nil
			}

			dMessage.Content = message.Message
		}

		_, err := d.webhook.Send(d.settings.ChannelID, dMessage, false, nil)

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

	userSettings, err := GetUsersSettings()
	if err != nil {
		fmt.Fprintf(os.Stderr, errors.Wrap(err, "GetUsersSettings").Error())
		return
	}

	user := userSettings.GetUser(m.Author.ID, ServiceTypeDiscord)

	user.Name = m.Member.Nick
	if user.Name == "" {
		user.Name = m.Author.Username
	}

	for _, text := range strings.Split(m.Message.Content, "\n") {
		var command CommandContent
		var msg = strings.Split(strings.TrimSpace(text), " ")

		var permissions = userSettings.GetPermissions(user.Groups)

		hasPrefixSay, err := permissions.Verify(strings.TrimSpace(text))
		if err != nil {
			fmt.Fprintf(os.Stderr, errors.Wrap(err, "Verify").Error())
			return
		}

		// check the message has prefix: "say"
		// if there is the prefix, not escape "@" to "at_"
		if !hasPrefixSay {
			ok, err := permissions.Verify("say")
			if err != nil {
				fmt.Fprintf(os.Stderr, errors.Wrap(err, "Verify").Error())
				return
			}
			if !ok {
				// nothing to do
				return
			}

			// send as message
			msg = append([]string{"say"}, msg...)
		}

		command.Command = msg[0]

		// vanilla needs slash to execute commands
		switch d.serverType {
		case "vanilla", "":
			command.Command = "/" + msg[0]
		}

		if len(msg) < 2 {
			return
		}

		switch command.Command {
		case "/msg", "/say", "msg", "say":
			msg[1] = fmt.Sprintf("[%s]%s", user.Name, msg[1])
			command.Options = strings.Join(msg[1:], " ")

			for _, id := range d.idRegExp.FindAllStringSubmatch(command.Options, -1) {
				// replace Slack User ids
				if len(id) < 2 {
					continue
				}

				mem, err := s.GuildMember(m.GuildID, id[1])
				if err != nil {
					continue
				}

				var idName = mem.Nick
				if idName == "" {
					idName = mem.User.Username
				}
				command.Options = strings.Join(strings.Split(command.Options, "!"+id[1]), idName)
			}

			if !hasPrefixSay {
				// escape "@" (target selector)
				command.Options = strings.ReplaceAll(command.Options, "@", "at_")
			}
		default:
			command.Options = strings.Join(msg[1:], " ")
		}

		d.sendChannel <- command
	}
}

func (d *DiscordHandler) sendUserState(event minecraft.MessageType) error {
	return d.updateUserState(event == minecraft.MessageTypeJoin || event == minecraft.MessageTypeLeft)
}

// RefreshUserState edits the existing roster so minute updates do not post new messages.
func (d *DiscordHandler) RefreshUserState() error {
	return d.updateUserState(false)
}

func (d *DiscordHandler) updateUserState(repost bool) error {
	d.stateMu.Lock()
	defer d.stateMu.Unlock()
	if d.closed || !d.settings.SendJoinStateMessage {
		return nil
	}
	const LogonPngName = "Logon Users.png"

	var dFiles []discord_webhook.File

	switch len(d.joinState.Snapshot()) {
	case 0:
		if d.lastMessageID != "" {
			if err := d.webhook.Delete(d.settings.ChannelID, d.lastMessageID); err != nil {
				return err
			}
		}
		d.lastMessageID = ""
		d.avatars = nil
	default:
		r, err := d.makeUserStateImage()
		if err != nil {
			return errors.Wrap(err, "MakeReactionImage")
		}
		dFiles = []discord_webhook.File{
			{
				FileName:    LogonPngName,
				Reader:      r,
				ContentType: "image/png",
			},
		}

		var message = &discord_webhook.Message{
			UserName:  d.settings.UserName,
			AvaterURL: d.settings.AvaterURI,
		}
		previousID := d.lastMessageID
		if d.lastMessageID == "" || repost {
			message, err = d.webhook.Send(d.settings.ChannelID, *message, true, dFiles)
		} else {
			message, err = d.webhook.Edit(d.settings.ChannelID, d.lastMessageID, *message, dFiles)
			var apiErr *discord_webhook.HTTPError
			if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound {
				// The roster may have been manually deleted; recreate it once.
				d.lastMessageID = ""
				return d.refreshDeletedUserState()
			}
		}
		if err != nil {
			return errors.Wrap(err, "Send")
		}

		if message.ID == "" {
			return errors.New("SendUserState: missing message ID")
		}
		d.lastMessageID = message.ID
		if repost && previousID != "" && previousID != message.ID {
			return d.webhook.Delete(d.settings.ChannelID, previousID)
		}
	}

	return nil
}

// stateMu is already held by RefreshUserState.
func (d *DiscordHandler) refreshDeletedUserState() error {
	r, err := d.makeUserStateImage()
	if err != nil {
		return err
	}
	message, err := d.webhook.Send(d.settings.ChannelID, discord_webhook.Message{
		UserName: d.settings.UserName, AvaterURL: d.settings.AvaterURI,
	}, true, []discord_webhook.File{{FileName: "Logon Users.png", Reader: r, ContentType: "image/png"}})
	if err != nil {
		return err
	}
	if message.ID == "" {
		return errors.New("SendUserState: missing message ID")
	}
	d.lastMessageID = message.ID
	return nil
}

func (d *DiscordHandler) makeUserStateImage() (io.Reader, error) {
	state := d.joinState.Snapshot()
	names := make([]string, 0, len(state))
	for name := range state {
		names = append(names, name)
	}
	if d.avatars == nil {
		d.avatars = make(map[string]cachedStateAvatar)
	}
	for name := range d.avatars {
		if _, online := state[name]; !online {
			delete(d.avatars, name)
		}
	}
	sort.Strings(names)
	players := make([]userStatePlayer, 0, len(names))
	for _, name := range names {
		player := userStatePlayer{Name: name, JoinedAt: state[name].JoinedAt}
		if cached, ok := d.avatars[name]; ok && cached.joinedAt.Equal(player.JoinedAt) {
			player.Head = cached.head
			players = append(players, player)
			continue
		}
		r, err := mcheads.GetAvater(name)
		if err == nil {
			player.Head, _, err = image.Decode(r)
			r.Close()
		}
		if err != nil {
			log.Printf("makeUserStateImage: %s: %s", name, err)
		} else {
			d.avatars[name] = cachedStateAvatar{joinedAt: player.JoinedAt, head: player.Head}
		}
		players = append(players, player)
	}
	canvas, err := renderUserStateImage(players, time.Now())
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	if err := png.Encode(&out, canvas); err != nil {
		return nil, err
	}
	return &out, nil
}
