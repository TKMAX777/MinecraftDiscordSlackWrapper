package minecraft

import (
	"bufio"
	"io"
	"log"
	"regexp"
	"strings"

	"github.com/TKMAX777/MinecraftDiscordWrapper/process"
	"github.com/pkg/errors"
)

// Handler handles minecraft server process
type Handler struct {
	settings Setting
	process  process.Process

	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
}

type Setting struct {
	ThreadInfoRegExp string
	JAVA             string
	CustomBinaryPath string
	Options          []string
}

type MessageType int

const (
	MessageTypeJoin MessageType = iota
	MessageTypeLeft
	MessageTypeThreadINFO
	MessageTypeOther
)

type Message struct {
	Type    MessageType
	Message string

	// when Message type is Join / Left, User will have Join or Left user name
	User string
}

// NewHandler makes new minecraft handler
func NewHandler(settings Setting) *Handler {
	return &Handler{settings: settings, process: process.Process{}}
}

// Start starts minectaft server
func (m *Handler) Start() (chan Message, error) {
	var opt = []string{}
	if m.settings.CustomBinaryPath != "" {
		err := m.process.Start(m.settings.CustomBinaryPath, opt...)
		if err != nil {
			return nil, errors.Wrap(err, "StartMinecraft")
		}
	} else {
		opt = append(m.settings.Options, "-jar", "server.jar", "nogui")
		if m.settings.JAVA == "" {
			m.settings.JAVA = "java"
		}

		var err = m.process.Start(m.settings.JAVA, opt...)
		if err != nil {
			return nil, errors.Wrap(err, "StartMinecraft")
		}
	}

	m.stdin, m.stdout, m.stderr, _ = m.process.Pipes()

	var cMessage = m.sendMessages()
	return cMessage, nil
}

func (m *Handler) sendMessages() chan Message {
	var joinedOrLeftRegExp = regexp.MustCompile(`\]: (\S+) (joined|left) (the game)$`)

	var infoTextRegExp *regexp.Regexp
	switch m.settings.ThreadInfoRegExp {
	case "", "default":
		infoTextRegExp = regexp.MustCompile(`\[\d{2}:\d{2}:\d{2}\] \[Server thread/INFO\]: (.+)`)
	case "paper":
		infoTextRegExp = regexp.MustCompile(`\[\d{2}:\d{2}:\d{2} INFO\]: (.+)`)
	default:
		infoTextRegExp = regexp.MustCompile(m.settings.ThreadInfoRegExp)
	}

	var rdr = bufio.NewReaderSize(m.stdout, bufio.MaxScanTokenSize)

	var cMessage = make(chan Message)

	go func() {
		defer m.stdout.Close()

		for {
			text, err := rdr.ReadString('\n')
			if err == io.EOF {
				return
			} else if err != nil {
				if err.Error() == "read |0: file already closed" {
					return
				}
				log.Printf("ErrorAtReadOutput: %s", err.Error())
			}

			var message Message
			text = strings.TrimSpace(text)

			switch {
			case joinedOrLeftRegExp.Match([]byte(text)):
				switch joinedOrLeftRegExp.FindStringSubmatch(text)[2] {
				case "joined":
					message.Type = MessageTypeJoin
				case "left":
					message.Type = MessageTypeLeft
				}

				message.User = joinedOrLeftRegExp.FindStringSubmatch(text)[1]

				text = joinedOrLeftRegExp.ReplaceAllString(text, "$1 $2 $3")

			case infoTextRegExp.Match([]byte(text)):
				message.Type = MessageTypeThreadINFO
				text = infoTextRegExp.ReplaceAllString(text, `$1`)
			default:
				message.Type = MessageTypeOther
			}
			message.Message = text
			cMessage <- message
		}
	}()

	return cMessage
}

// Interrupt send keyboard interrupt to the minecraft server
func (m *Handler) Interrupt() {
	m.process.Interrupt()
}

// Kill send kill the minecraft server immediately
func (m *Handler) Kill() {
	m.process.Kill()
}

// Pipes returns minecrafts process stdin / stderr pipes
//    these pipes are automatically closed when the process killed by Handler.Stop()
func (m *Handler) Pipes() (stdin io.WriteCloser, stderr io.ReadCloser) {
	return m.stdin, m.stderr
}
