package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/TKMAX777/MinecraftDiscordWrapper/mcheads"
	"github.com/TKMAX777/MinecraftDiscordWrapper/minecraft"
	"github.com/TKMAX777/MinecraftDiscordWrapper/slack_webhook"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	scm "github.com/slack-go/slack/socketmode"
)

// SlackHandler handles Slack conversations
type SlackHandler struct {
	api *slack.Client
	scm *scm.Client

	sendChannel   chan CommandContent
	lastMessageTS string

	regExp struct {
		UserID  *regexp.Regexp
		Channel *regexp.Regexp
		URI     *regexp.Regexp
	}

	serverType string
	joinState  *JoinState
	webhook    *slack_webhook.Handler

	messageUnescaper *strings.Replacer

	settings SlackSetting
}

func NewSlackHandler(settings SlackSetting, joinState *JoinState) *SlackHandler {
	var handler SlackHandler

	handler.api = slack.New(
		settings.Token,
		slack.OptionAppLevelToken(settings.EventToken),
	)
	slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags))

	handler.regExp.UserID = regexp.MustCompile(`<@(\S+)>`)
	handler.regExp.Channel = regexp.MustCompile(`<#(\S+)\|(\S+)>`)
	handler.regExp.URI = regexp.MustCompile(`<(https??://\S+)\|(\S+)>`)

	handler.settings = settings

	handler.messageUnescaper = strings.NewReplacer(
		"&amp;", "&",
		"&lt;", "<",
		"&gt;", ">",
	)

	handler.webhook = slack_webhook.New(settings.Token)

	handler.joinState = joinState

	return &handler
}

func (s *SlackHandler) SetCommandInput(stdin chan CommandContent) *SlackHandler {
	s.sendChannel = stdin
	return s
}

func (s *SlackHandler) SetServerType(serverType string) *SlackHandler {
	s.serverType = serverType
	return s
}

func (s *SlackHandler) StartSession() error {
	s.scm = scm.New(s.api)
	go func() {
		var err = s.scm.Run()
		if err != nil {
			log.Println(errors.Wrap(err, "OpeningSlackSession"))
		}
	}()

	go func() {
		for ev := range s.scm.Events {
			switch ev.Type {
			case scm.EventTypeConnected:
				fmt.Printf("Start websocket connection with Slack\n")
			case scm.EventTypeEventsAPI:
				s.scm.Ack(*ev.Request)

				evp, ok := ev.Data.(slackevents.EventsAPIEvent)
				if !ok {
					continue
				}
				switch evp.Type {
				case slackevents.CallbackEvent:
					switch evi := evp.InnerEvent.Data.(type) {
					case *slackevents.AppMentionEvent:
					case *slackevents.MessageEvent:
						s.getMessage(evi)
					}
				}
			}
		}
	}()
	return nil
}

func (s *SlackHandler) SendMessageFunction() MessageSender {
	var onlineUserNum int

	return MessageSender(func(message minecraft.Message) error {
		var sMessage = slack_webhook.Message{
			AsUser:   false,
			Channel:  s.settings.ChannelID,
			Username: s.settings.UserName,
			IconURL:  s.settings.AvaterURI,
		}
		switch message.Type {
		case minecraft.MessageTypeJoin:
			if s.settings.SendJoinStateMessage {
				err := s.sendUserState(message.Type)
				if err != nil {
					return err
				}
			}

			if s.settings.SendOption&(SendSettingJoinLeft|SendSettingAll) == 0 {
				return nil
			}

			onlineUserNum++
			if s.settings.AddOnlineNumber {
				switch onlineUserNum {
				case 0, 1:
					sMessage.Text = fmt.Sprintf("%s `%s joined the game`\nActive: %d player", s.settings.Reaction.Join, message.User, onlineUserNum)
				default:
					sMessage.Text = fmt.Sprintf("%s `%s joined the game`\nActive: %d players", s.settings.Reaction.Join, message.User, onlineUserNum)
				}
			} else {
				sMessage.Text = fmt.Sprintf("%s `%s joined the game`", s.settings.Reaction.Join, message.User)
			}
		case minecraft.MessageTypeLeft:
			if s.settings.SendJoinStateMessage {
				err := s.sendUserState(message.Type)
				if err != nil {
					return err
				}
			}

			if s.settings.SendOption&(SendSettingJoinLeft|SendSettingAll) == 0 {
				return nil
			}

			onlineUserNum--
			if s.settings.AddOnlineNumber {
				switch onlineUserNum {
				case 0, 1:
					sMessage.Text = fmt.Sprintf("%s `%s left the game`\nActive: %d player", s.settings.Reaction.Left, message.User, onlineUserNum)
				default:
					sMessage.Text = fmt.Sprintf("%s `%s left the game`\nActive: %d players", s.settings.Reaction.Left, message.User, onlineUserNum)
				}
			} else {
				sMessage.Text = fmt.Sprintf("%s `%s left the game`", s.settings.Reaction.Left, message.User)
			}
		case minecraft.MessageTypeThreadINFO:
			if s.settings.SendOption&(SendSettingThreadINFO|SendSettingAll) == 0 {
				return nil
			}

			sMessage.Text = message.Message
		case minecraft.MessageTypeDeath:
			if s.settings.SendOption&(SendSettingDeath|SendSettingAll) == 0 {
				return nil
			}

			sMessage.Text = fmt.Sprintf("%s %s", s.settings.Reaction.Death, message.Message)
		case minecraft.MessageTypeReachedTheAdvancement:
			if s.settings.SendOption&(SendSettingReachedTheAdvancement|SendSettingAll) == 0 {
				return nil
			}

			sMessage.Text = fmt.Sprintf("%s %s", s.settings.Reaction.Advancement, message.Message)
		case minecraft.MessageTypeMessage:
			if s.settings.SendOption&(SendSettingMessage|SendSettingAll) == 0 {
				return nil
			}
			sMessage.Username = message.User
			sMessage.IconURL = mcheads.GetAvaterURI(message.User)
			sMessage.Text = message.Message
		case minecraft.MessageTypeServermessage:
			if s.settings.SendOption&(SendSettingMessage|SendSettingAll) == 0 {
				return nil
			}

			sMessage.Text = message.Message
		case minecraft.MessageTypeDifficultySet:
			if s.settings.SendOption&(SendSettingDifficultySet|SendSettingAll) == 0 {
				return nil
			}

			sMessage.Text = fmt.Sprintf("%s %s", s.settings.Reaction.DifficultySet, message.Message)
		case minecraft.MessageTypeOther:
			if s.settings.SendOption&SendSettingAll == 0 {
				return nil
			}

			sMessage.Text = message.Message
		}

		_, err := s.webhook.Send(sMessage)

		return err
	})
}

// getMessage handles messages posted to Slack
func (s *SlackHandler) getMessage(ev *slackevents.MessageEvent) {
	if ev.Channel != s.settings.ChannelID {
		return
	}

	if ev.Text == "" {
		return
	}

	slackUser, err := s.api.GetUserInfo(ev.User)
	if err != nil {
		return
	}

	var name = slackUser.Profile.DisplayName
	if name == "" {
		name = slackUser.RealName
	}

	var user User
	user.Name = name

	for _, text := range strings.Split(ev.Text, "\n") {
		var command CommandContent
		var msg = strings.Split(text, " ")

		var permissions = GetPermissions(s.settings.Permissions)

		var ok bool

		// check the message has prefix: "say"
		// if there is the prefix, not escape "@" to "at_"
		var hasPrefixSay = true

		command.Command, ok = permissions[msg[0]]
		if !ok {
			_, ok = permissions["say"]
			if !ok || !s.settings.SendAllMessages {
				return
			}
			msg = append([]string{"say"}, msg...)
			command.Command = "/say"
			hasPrefixSay = false
		}

		// if server uses paperMC, commands do not contain "/""
		if s.serverType == "paper" {
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
			command.Options = strings.Join(msg[1:], " ")
			command.Options, err = s.escapeMessage(command.Options)
			if err != nil {
				continue
			}

			if !hasPrefixSay {
				// escape "@" (target selector)
				command.Options = strings.ReplaceAll(command.Options, "@", "at_")
			}

			command.Options = fmt.Sprintf("[%s]%s", user.Name, command.Options)
		case "/difficulty", "difficulty":
			switch msg[1] {
			case "p", "peaceful":
				if s.settings.Difficulty&DifficultyPeaceful == 0 {
					return
				}
			case "e", "easy":
				if s.settings.Difficulty&DifficultyEasy == 0 {
					return
				}
			case "n", "normal":
				if s.settings.Difficulty&DifficultyNormal == 0 {
					return
				}
			case "h", "hard":
				if s.settings.Difficulty&DifficultyHard == 0 {
					return
				}
			}

			command.Options = msg[1]
		default:
			command.Options = strings.Join(msg[1:], " ")
		}

		s.sendChannel <- command
	}
}

func (s *SlackHandler) escapeMessage(content string) (output string, err error) {
	for _, id := range s.regExp.UserID.FindAllStringSubmatch(content, -1) {
		if len(id) < 2 {
			continue
		}

		u, err := s.api.GetUserInfo(id[1])
		if err != nil {
			return "", err
		}

		var repl = u.Profile.DisplayName
		if repl == "" {
			repl = u.RealName
		}
		content = strings.Join(strings.Split(content, id[0]), "@"+repl+"")
	}

	for _, ch := range s.regExp.Channel.FindAllStringSubmatch(content, -1) {
		if len(ch) < 3 {
			continue
		}

		content = strings.Join(strings.Split(content,
			fmt.Sprintf("<#%s|%s>", ch[1], ch[2])),
			fmt.Sprintf("#%s", ch[2]),
		)
	}

	for _, uri := range s.regExp.URI.FindAllStringSubmatch(content, -1) {
		if len(uri) < 3 {
			continue
		}

		if uri[1] == uri[2] {
			content = strings.Join(strings.Split(content,
				fmt.Sprintf("<%s|%s>", uri[1], uri[2])),
				uri[1],
			)
		}
		content = strings.Join(strings.Split(content,
			fmt.Sprintf("<%s|%s>", uri[1], uri[2])),
			fmt.Sprintf("%s(URI: <%s>)", uri[2], uri[1]),
		)
	}

	return s.messageUnescaper.Replace(content), nil
}

func (s *SlackHandler) sendUserState(event minecraft.MessageType) error {
	var VoiceStateMessageText = fmt.Sprintf("MinecraftuserStateMessage,%s", s.webhook.Identity.UserID)
	var message = slack_webhook.Message{
		AsUser:   false,
		Channel:  s.settings.ChannelID,
		Username: s.settings.UserName,
		IconURL:  s.settings.AvaterURI,
		Blocks:   s.buildUserStateBlock(),
		Text:     VoiceStateMessageText,
	}

	var ts = s.lastMessageTS
	switch {
	case len(s.joinState.State) < 1:
		// there are no player
		if ts == "" {
			// if not found a last message, find from message history
			messages, err := s.webhook.GetMessages(s.settings.ChannelID, "", 100)
			if err == nil {
				for _, msg := range messages {
					if msg.Text == VoiceStateMessageText {
						// *repost user messages contains DummyURIs
						ts = msg.TS
						break
					}
				}
			}

			if ts == "" {
				return nil
			}
		}
		s.lastMessageTS = ""
		s.webhook.Remove(message.Channel, ts)
	default:
		// there are some players
		var ts = s.lastMessageTS
		if ts == "" {
			// if not found a last message, find from message history
			messages, err := s.webhook.GetMessages(s.settings.ChannelID, "", 100)
			if err == nil {
				for _, msg := range messages {
					if msg.Text == VoiceStateMessageText {
						ts = msg.TS
						break
					}
				}
			}
		}
		if ts != "" {
			s.webhook.Remove(message.Channel, ts)
		}

		message.TS = ts
		ts, err := s.webhook.Send(message)
		if err != nil {
			return err
		}

		s.lastMessageTS = ts
	}

	return nil
}

func (s *SlackHandler) buildUserStateBlock() []slack_webhook.BlockBase {
	var blocks = []slack_webhook.BlockBase{}

	var channelText = "Logon Users"
	var channelNameElement = slack_webhook.MrkdwnElement(channelText)

	blocks = append(
		blocks,
		slack_webhook.ContextBlock(channelNameElement),
	)

	var userCount int
	var elements = []slack_webhook.BlockElement{}

	for username := range s.joinState.State {
		var userImage = mcheads.GetAvaterURI(username)
		var imageElm = slack_webhook.ImageElement(userImage, username)
		var userElm = slack_webhook.MrkdwnElement(username)

		elements = append(elements, imageElm, userElm)

		userCount++
		if userCount%4 == 0 {
			var block = slack_webhook.ContextBlock(elements...)
			blocks = append(blocks, block)

			elements = []slack_webhook.BlockElement{}
		}
	}

	if userCount%4 > 0 {
		var block = slack_webhook.ContextBlock(elements...)
		blocks = append(blocks, block)
	}

	blocks = append(blocks, slack_webhook.DividerBlock())

	return blocks
}
