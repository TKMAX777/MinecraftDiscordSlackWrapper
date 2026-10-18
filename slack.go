package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/TKMAX777/MinecraftDiscordSlackWrapper/mcheads"
	"github.com/TKMAX777/MinecraftDiscordSlackWrapper/minecraft"
	"github.com/TKMAX777/MinecraftDiscordSlackWrapper/slack_webhook"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	scm "github.com/slack-go/slack/socketmode"
)

// SlackHandler handles Slack conversations
type SlackHandler struct {
	api *slack.Client
	scm *scm.Client

	sendChannel          chan CommandContent
	lastMessageTS        string
	rosterHistoryChecked bool

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

			if !(s.settings.SendOption.All || s.settings.SendOption.JoinLeft) {
				return nil
			}

			onlineUserNum++
			if s.settings.AddOnlineNumber {
				switch onlineUserNum {
				case 0, 1:
					sMessage.Text = fmt.Sprintf(":%s: `%s joined the game`\nActive: %d player", s.settings.Reaction.Join, message.User, onlineUserNum)
				default:
					sMessage.Text = fmt.Sprintf(":%s: `%s joined the game`\nActive: %d players", s.settings.Reaction.Join, message.User, onlineUserNum)
				}
			} else {
				sMessage.Text = fmt.Sprintf(":%s: `%s joined the game`", s.settings.Reaction.Join, message.User)
			}
		case minecraft.MessageTypeLeft:
			if s.settings.SendJoinStateMessage {
				err := s.sendUserState(message.Type)
				if err != nil {
					return err
				}
			}

			if !(s.settings.SendOption.All || s.settings.SendOption.JoinLeft) {
				return nil
			}

			onlineUserNum--
			if s.settings.AddOnlineNumber {
				switch onlineUserNum {
				case 0, 1:
					sMessage.Text = fmt.Sprintf(":%s: `%s left the game`\nActive: %d player", s.settings.Reaction.Left, message.User, onlineUserNum)
				default:
					sMessage.Text = fmt.Sprintf(":%s: `%s left the game`\nActive: %d players", s.settings.Reaction.Left, message.User, onlineUserNum)
				}
			} else {
				sMessage.Text = fmt.Sprintf(":%s: `%s left the game`", s.settings.Reaction.Left, message.User)
			}
		case minecraft.MessageTypeThreadINFO:
			if !(s.settings.SendOption.All || s.settings.SendOption.ThreadINFO) {
				return nil
			}

			sMessage.Text = message.Message
		case minecraft.MessageTypeVillagerDeath:
			if !(s.settings.SendOption.All || s.settings.SendOption.VillagerDeath) {
				return nil
			}

			content := message.Content.(minecraft.MessageContentVillagerDeath)
			sMessage.Text = fmt.Sprintf(
				":%s: Villager: %s (ID: %d) died\n%s",
				s.settings.Reaction.VillagerDeath,
				content.Job,
				content.ID,
				content.DiedMessage,
			)
		case minecraft.MessageTypeDeath:
			if !(s.settings.SendOption.All || s.settings.SendOption.Death) {
				return nil
			}

			sMessage.Text = fmt.Sprintf(":%s: %s", s.settings.Reaction.Death, message.Message)
		case minecraft.MessageTypeReachedTheAdvancement:
			if !(s.settings.SendOption.All || s.settings.SendOption.ReachedTheAdvancement) {
				return nil
			}

			sMessage.Text = fmt.Sprintf(":%s: %s", s.settings.Reaction.Advancement, message.Message)
		case minecraft.MessageTypeMessage:
			if !(s.settings.SendOption.All || s.settings.SendOption.Message) {
				return nil
			}
			sMessage.Username = message.User
			if avatarURI, err := mcheads.GetCurrentAvaterURI(message.User); err == nil {
				sMessage.IconURL = avatarURI
			} else {
				log.Printf("Slack chat avatar for %s: %v", message.User, err)
			}
			sMessage.Text = message.Message
		case minecraft.MessageTypeServermessage:
			if !(s.settings.SendOption.All || s.settings.SendOption.Message) {
				return nil
			}

			sMessage.Text = message.Message
		case minecraft.MessageTypeDifficultySet:
			if !(s.settings.SendOption.All || s.settings.SendOption.DifficultySet) {
				return nil
			}

			sMessage.Text = fmt.Sprintf(":%s: %s", s.settings.Reaction.DifficultySet, message.Message)
		case minecraft.MessageTypeOther:
			if !s.settings.SendOption.All {
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

	userSettings, err := GetUsersSettings()
	if err != nil {
		fmt.Fprintf(os.Stderr, errors.Wrap(err, "GetUsersSettings").Error())
		return
	}

	user := userSettings.GetUser(ev.User, ServiceTypeSlack)

	var name = slackUser.Profile.DisplayName
	if name == "" {
		name = slackUser.RealName
	}

	user.Name = name

	for _, text := range strings.Split(ev.Text, "\n") {
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
		switch s.serverType {
		case "vanilla", "":
			command.Command = "/" + msg[0]
		}

		if len(msg) < 2 {
			return
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
	return s.RefreshUserState()
}

// RefreshUserState updates the existing Slack roster without reposting it.
func (s *SlackHandler) RefreshUserState() error {
	if !s.settings.SendJoinStateMessage {
		return nil
	}
	state := s.joinState.Snapshot()
	marker := fmt.Sprintf("MinecraftuserStateMessage,%s", s.webhook.Identity.UserID)
	if !s.rosterHistoryChecked {
		s.rosterHistoryChecked = true
		if s.lastMessageTS == "" {
			messages, err := s.webhook.GetMessages(s.settings.ChannelID, "", 100)
			if err == nil {
				for _, msg := range messages {
					if msg.Text == marker {
						s.lastMessageTS = msg.TS
						break
					}
				}
			}
		}
	}
	if len(state) == 0 {
		if s.lastMessageTS != "" {
			_, err := s.webhook.Remove(s.settings.ChannelID, s.lastMessageTS)
			var apiErr *slack_webhook.APIError
			if err != nil && !(errors.As(err, &apiErr) && apiErr.Code == "message_not_found") {
				return err
			}
			s.lastMessageTS = ""
		}
		return nil
	}
	avatarURIs := make(map[string]string, len(state))
	for name := range state {
		if uri, err := mcheads.GetCurrentAvaterURI(name); err == nil {
			avatarURIs[name] = uri
		} else {
			log.Printf("Slack roster avatar for %s: %v", name, err)
		}
	}
	message := slack_webhook.Message{
		AsUser: false, Channel: s.settings.ChannelID, Username: s.settings.UserName,
		IconURL: s.settings.AvaterURI, Blocks: buildSlackUserStateBlocks(state, time.Now(), avatarURIs), Text: marker,
		TS: s.lastMessageTS,
	}
	var ts string
	var err error
	if message.TS == "" {
		ts, err = s.webhook.Send(message)
	} else {
		ts, err = s.webhook.Update(message)
		var apiErr *slack_webhook.APIError
		if errors.As(err, &apiErr) && apiErr.Code == "message_not_found" {
			s.lastMessageTS = ""
			message.TS = ""
			ts, err = s.webhook.Send(message)
		}
	}
	if err != nil {
		return err
	}
	if ts == "" {
		return errors.New("RefreshSlackUserState: missing message timestamp")
	}
	s.lastMessageTS = ts
	return nil
}

func buildSlackUserStateBlocks(state map[string]UserState, now time.Time, avatarURIs map[string]string) []slack_webhook.BlockBase {
	header := slack_webhook.SectionBlock()
	header.Text = slack_webhook.MrkdwnElement(fmt.Sprintf("*Active Users*  ·  %d online", len(state)))
	blocks := []slack_webhook.BlockBase{header}
	names := make([]string, 0, len(state))
	for name := range state {
		names = append(names, name)
	}
	sort.Strings(names)
	elements := []slack_webhook.BlockElement{}
	for _, name := range names {
		if uri := avatarURIs[name]; uri != "" {
			elements = append(elements, slack_webhook.ImageElement(uri, name))
		}
		elements = append(elements,
			slack_webhook.BlockElement{Type: "plain_text", Text: fmt.Sprintf("%s  ·  ⏱ %s", name, formatLoginDuration(state[name].JoinedAt, now))},
		)
		if len(elements) >= 8 {
			blocks = append(blocks, slack_webhook.ContextBlock(elements...))
			elements = nil
		}
	}
	if len(elements) > 0 {
		blocks = append(blocks, slack_webhook.ContextBlock(elements...))
	}
	blocks = append(blocks, slack_webhook.DividerBlock())
	return blocks
}
