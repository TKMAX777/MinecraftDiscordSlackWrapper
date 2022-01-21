package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

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

	sendChannel chan CommandContent

	regExp struct {
		UserID  *regexp.Regexp
		Channel *regexp.Regexp
		URI     *regexp.Regexp
	}

	webhook *slack_webhook.Handler

	messageUnescaper *strings.Replacer

	settings SlackSetting
}

func NewSlackHandler(settings SlackSetting) *SlackHandler {
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

	return &handler
}

func (s *SlackHandler) SetCommandInput(stdin chan CommandContent) *SlackHandler {
	s.sendChannel = stdin
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
		var content string
		switch message.Type {
		case minecraft.MessageTypeJoin:
			if s.settings.SendOption&(SendSettingJoinLeft|SendSettingAll) == 0 {
				return nil
			}

			onlineUserNum++
			if s.settings.AddOnlineNumber {
				switch onlineUserNum {
				case 0, 1:
					content = fmt.Sprintf("%s `%s joined the game`\nActive: %d player", s.settings.Reaction.Join, message.User, onlineUserNum)
				default:
					content = fmt.Sprintf("%s `%s joined the game`\nActive: %d players", s.settings.Reaction.Join, message.User, onlineUserNum)
				}
			} else {
				content = fmt.Sprintf("%s `%s joined the game`", s.settings.Reaction.Join, message.User)
			}
		case minecraft.MessageTypeLeft:
			if s.settings.SendOption&(SendSettingJoinLeft|SendSettingAll) == 0 {
				return nil
			}

			onlineUserNum--
			if s.settings.AddOnlineNumber {
				switch onlineUserNum {
				case 0, 1:
					content = fmt.Sprintf("%s `%s left the game`\nActive: %d player", s.settings.Reaction.Left, message.User, onlineUserNum)
				default:
					content = fmt.Sprintf("%s `%s left the game`\nActive: %d players", s.settings.Reaction.Left, message.User, onlineUserNum)
				}
			} else {
				content = fmt.Sprintf("%s `%s left the game`", s.settings.Reaction.Left, message.User)
			}
		case minecraft.MessageTypeThreadINFO:
			if s.settings.SendOption&(SendSettingThreadINFO|SendSettingAll) == 0 {
				return nil
			}

			content = message.Message
		case minecraft.MessageTypeOther:
			if s.settings.SendOption&SendSettingAll == 0 {
				return nil
			}

			content = message.Message
		}

		_, err := s.webhook.Send(slack_webhook.Message{
			AsUser:   false,
			Channel:  s.settings.ChannelID,
			Username: s.settings.UserName,
			IconURL:  s.settings.AvaterURI,
			Text:     content,
		})

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
		command.Command, ok = permissions[msg[0]]

		if !ok {
			_, ok = permissions["say"]
			if !ok || !s.settings.SendAllMessages {
				return
			}
			msg = append([]string{"say"}, msg...)
			command.Command = "/say"
		}

		if len(msg) < 2 {
			return
		}

		if msg[1] == ";" {
			msg = msg[:1]
		}

		switch command.Command {
		case "/msg", "/say":
			msg[1] = fmt.Sprintf("[%s]%s", user.Name, msg[1])
			command.Options = strings.Join(msg[1:], " ")

			text, err = s.escapeMessage(text)
			if err != nil {
				continue
			}

		default:
			command.Options = strings.Join(msg[1:], " ")
		}

		fmt.Printf("[Slack]%v\n", text)

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
		content = strings.Join(strings.Split(content, id[0]), "`@"+repl+"`")
	}

	for _, ch := range s.regExp.Channel.FindAllStringSubmatch(content, -1) {
		if len(ch) < 3 {
			continue
		}

		content = strings.Join(strings.Split(content,
			fmt.Sprintf("<#%s|%s>", ch[1], ch[2])),
			fmt.Sprintf("`#%s`(URI: <%sarchives/%s>)", ch[2], s.webhook.Identity.WorkspaceURI, ch[1]),
		)
	}

	for _, uri := range s.regExp.URI.FindAllStringSubmatch(content, -1) {
		if len(uri) < 3 {
			continue
		}

		if uri[1] == uri[2] {
			content = strings.Join(strings.Split(content,
				fmt.Sprintf("<%s|%s>", uri[1], uri[2])),
				fmt.Sprintf("<%s>)", uri[1]),
			)
		}
		content = strings.Join(strings.Split(content,
			fmt.Sprintf("<%s|%s>", uri[1], uri[2])),
			fmt.Sprintf("%s(URI: <%s>)", uri[2], uri[1]),
		)
	}

	return s.messageUnescaper.Replace(content), nil
}