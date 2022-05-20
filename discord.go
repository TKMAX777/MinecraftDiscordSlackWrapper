package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math"
	"regexp"
	"strings"

	"golang.org/x/image/draw"

	"github.com/TKMAX777/MinecraftDiscordSlackWrapper/discord_webhook"
	"github.com/TKMAX777/MinecraftDiscordSlackWrapper/mcheads"
	"github.com/TKMAX777/MinecraftDiscordSlackWrapper/minecraft"
	"github.com/bwmarrin/discordgo"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"

	"github.com/pkg/errors"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/math/fixed"
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

	settings DiscordSetting
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

			if d.settings.SendOption&(SendSettingJoinLeft|SendSettingAll) == 0 {
				return nil
			}

			onlineUserNum++
			if d.settings.AddOnlineNumber {
				switch onlineUserNum {
				case 0, 1:
					dMessage.Content = fmt.Sprintf("%s `%s joined the game`\nActive: %d player", d.settings.Reaction.Join, message.User, onlineUserNum)
				default:
					dMessage.Content = fmt.Sprintf("%s `%s joined the game`\nActive: %d players", d.settings.Reaction.Join, message.User, onlineUserNum)
				}
			} else {
				dMessage.Content = fmt.Sprintf("%s `%s joined the game`", d.settings.Reaction.Join, message.User)
			}
		case minecraft.MessageTypeLeft:
			if d.settings.SendJoinStateMessage {
				err := d.sendUserState(message.Type)
				if err != nil {
					return err
				}
			}

			if d.settings.SendOption&(SendSettingJoinLeft|SendSettingAll) == 0 {
				return nil
			}

			onlineUserNum--
			if d.settings.AddOnlineNumber {
				switch onlineUserNum {
				case 0, 1:
					dMessage.Content = fmt.Sprintf("%s `%s left the game`\nActive: %d player", d.settings.Reaction.Left, message.User, onlineUserNum)
				default:
					dMessage.Content = fmt.Sprintf("%s `%s left the game`\nActive: %d players", d.settings.Reaction.Left, message.User, onlineUserNum)
				}
			} else {
				dMessage.Content = fmt.Sprintf("%s `%s left the game`", d.settings.Reaction.Left, message.User)
			}
		case minecraft.MessageTypeThreadINFO:
			if d.settings.SendOption&(SendSettingThreadINFO|SendSettingAll) == 0 {
				return nil
			}

			dMessage.Content = message.Message
		case minecraft.MessageTypeDeath:
			if d.settings.SendOption&(SendSettingDeath|SendSettingAll) == 0 {
				return nil
			}

			dMessage.Content = fmt.Sprintf("%s %s", d.settings.Reaction.Death, message.Message)
		case minecraft.MessageTypeReachedTheAdvancement:
			if d.settings.SendOption&(SendSettingReachedTheAdvancement|SendSettingAll) == 0 {
				return nil
			}

			dMessage.Content = fmt.Sprintf("%s %s", d.settings.Reaction.Advancement, message.Message)
		case minecraft.MessageTypeMessage:
			if d.settings.SendOption&(SendSettingMessage|SendSettingAll) == 0 {
				return nil
			}
			dMessage.UserName = message.User
			dMessage.AvaterURL = mcheads.GetAvaterURI(message.User)
			dMessage.Content = message.Message
		case minecraft.MessageTypeServermessage:
			if d.settings.SendOption&(SendSettingMessage|SendSettingAll) == 0 {
				return nil
			}

			dMessage.Content = message.Message
		case minecraft.MessageTypeOther:
			if d.settings.SendOption&SendSettingAll == 0 {
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

	userDict, _ := ReadNameDict()
	if userDict == nil {
		userDict = Users{}
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
			user.Name = m.Author.Username
		}
	}

	for _, text := range strings.Split(m.Message.Content, "\n") {
		var command CommandContent
		var msg = strings.Split(text, " ")

		var permissions = GetPermissions(user.PermissionCode)

		// check the message has prefix: "say"
		// if there is the prefix, not escape "@" to "at_"
		var hasPrefixSay = true

		command.Command, ok = permissions[msg[0]]
		if !ok {
			_, ok = permissions["say"]
			if !ok || !user.SendAllMessages {
				return
			}
			msg = append([]string{"say"}, msg...)
			command.Command = "/say"
			hasPrefixSay = false
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

			if !hasPrefixSay {
				// escape "@" (target selector)
				command.Options = strings.ReplaceAll(command.Options, "@", "at_")
			}
		case "/difficulty", "difficulty":
			switch msg[1] {
			case "p", "peaceful":
				if d.settings.Difficulty&DifficultyPeaceful == 0 {
					return
				}
			case "e", "easy":
				if d.settings.Difficulty&DifficultyEasy == 0 {
					return
				}
			case "n", "normal":
				if d.settings.Difficulty&DifficultyNormal == 0 {
					return
				}
			case "h", "hard":
				if d.settings.Difficulty&DifficultyHard == 0 {
					return
				}
			}

			command.Options = msg[1]
		default:
			command.Options = strings.Join(msg[1:], " ")
		}

		d.sendChannel <- command
	}
}

func (d *DiscordHandler) sendUserState(event minecraft.MessageType) error {
	const LogonPngName = "Logon Users.png"

	var dFiles []discord_webhook.File

	if d.lastMessageID != "" {
		var err = d.webhook.Delete(d.settings.ChannelID, d.lastMessageID)
		if err != nil {
			log.Printf("sendUserState: Delete: %s\n", err.Error())
		}
	}

	switch len(d.joinState.State) {
	case 0:
		d.lastMessageID = ""
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
		message, err = d.webhook.Send(d.settings.ChannelID, *message, true, dFiles)
		if err != nil {
			return errors.Wrap(err, "Send")
		}

		d.lastMessageID = message.ID
	}

	return nil
}

func (d *DiscordHandler) makeUserStateImage() (r io.Reader, err error) {
	const UserNameSize = 30
	const imageHeight = 40
	const imageWidth = 800

	const imageMarginSide = 5
	const imageMarginLine = 5
	const userMargin = 5

	const HeadSize = imageHeight - imageMarginLine

	ft, err := truetype.Parse(gomono.TTF)
	if err != nil {
		return nil, errors.Wrap(err, "FontParseError")
	}

	var laneFrames = make([]*image.RGBA, 0)
	var frame = image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))

	var X = imageMarginSide
	for username := range d.joinState.State {
		r, err := mcheads.GetAvater(username)
		if err != nil {
			log.Println("makeUserStateImage: GetAvater: ", err.Error())
			continue
		}

		headImage, _, err := image.Decode(r)
		if err != nil {
			log.Println("makeUserStateImage: Decode: ", err.Error())
			continue
		}

		headImage = d.resize(headImage, HeadSize)

		var errConunt int
	again:
		var dr = &font.Drawer{
			Dst: frame,
			Src: image.Black,
			Face: truetype.NewFace(
				ft,
				&truetype.Options{
					Size: UserNameSize,
				},
			),
			Dot: fixed.Point26_6{},
		}

		// Confirm that the width occupied by the user does not exceed the image size.
		if imageWidth < X+headImage.Bounds().Dx()+dr.MeasureString(username).Ceil()+imageMarginSide-userMargin {
			laneFrames = append(laneFrames, frame)
			frame = image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))

			if errConunt > 1 {
				log.Println("INFO: makeUserStateImage: UserName too long:", username)
				continue
			}

			errConunt++
			X = imageMarginSide

			goto again
		}

		var imgPoint = image.Point{X, imageMarginLine}
		draw.Copy(frame, imgPoint, headImage, headImage.Bounds(), draw.Over, nil)

		dr.Dot.X = fixed.I(headImage.Bounds().Dx() + userMargin + X)
		dr.Dot.Y = fixed.I(imageHeight)

		dr.DrawString(username)

		X += X + headImage.Bounds().Dx() + dr.MeasureString(username).Ceil() + userMargin
	}

	laneFrames = append(laneFrames, frame)
	frame = image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight*(len(laneFrames)+1)+imageMarginLine))
	frame = d.fillFrame(frame, color.White)

	var dr = &font.Drawer{
		Dst: frame,
		Src: image.Black,
		Face: truetype.NewFace(
			ft,
			&truetype.Options{
				Size: UserNameSize,
			},
		),
		Dot: fixed.Point26_6{},
	}

	dr.Dot.X = fixed.I(imageMarginSide)
	dr.Dot.Y = fixed.I(imageHeight)

	dr.DrawString("<Logon Users>")

	for i := range laneFrames {
		var imgPoint = image.Point{0, imageHeight * (i + 1)}
		draw.Copy(frame, imgPoint, laneFrames[i], laneFrames[i].Bounds(), draw.Over, nil)
	}

	var encodePNG = new(bytes.Buffer)
	err = png.Encode(encodePNG, frame)
	if err != nil {
		return nil, errors.Wrap(err, "Encode")
	}

	return encodePNG, nil
}

func (d DiscordHandler) fillFrame(frame *image.RGBA, c color.Color) *image.RGBA {
	var rect = frame.Rect
	var newFrame = &image.RGBA{}

	*newFrame = *frame

	for h := rect.Min.Y; h < rect.Max.Y; h++ {
		for v := rect.Min.X; v < rect.Max.X; v++ {
			newFrame.Set(v, h, c)
		}
	}

	return newFrame
}

func (d DiscordHandler) resize(srcImage image.Image, MaxSize int) *image.RGBA {
	// resize png, jpg
	var width, height = float64(srcImage.Bounds().Size().X), float64(srcImage.Bounds().Size().Y)

	var ratio float64
	if width > height {
		ratio = float64(MaxSize) / width
	} else {
		ratio = float64(MaxSize) / height
	}
	srcImage = resize.Resize(
		uint(math.Floor(width*ratio)),
		uint(math.Floor(height*ratio)),
		srcImage, resize.Lanczos3,
	)
	resizedImage := image.NewRGBA(srcImage.Bounds())
	draw.FloydSteinberg.Draw(resizedImage, srcImage.Bounds(), srcImage, image.Point{})

	return resizedImage
}
