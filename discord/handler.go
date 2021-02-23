package discord

import "fmt"

// Handler is a discord webhook handler
type Handler struct {
	hookURI      string
	errorHookURI string
	avaterURI    string
	userName     string
	errMessage   errMessage
}

type errMessage struct {
	avaterURI string
	userName  string
}

// NewHandler create new discord message handler
func NewHandler(hookURI string) *Handler {
	return &Handler{hookURI: hookURI}
}

// SetErrorHookURI put error hook URI
func (h *Handler) SetErrorHookURI(hookURI string) *Handler {
	h.errorHookURI = hookURI
	return h
}

// SetProfile set defalt message profile
func (h *Handler) SetProfile(avaterURI, userName string) *Handler {
	h.avaterURI = avaterURI
	h.userName = userName

	return h
}

// SetErrorProfile set defalt error message profile
func (h *Handler) SetErrorProfile(avaterURI, userName string) *Handler {
	h.errMessage = errMessage{avaterURI, userName}
	return h
}

// SendMessage send message to discord by using default settings
func (h *Handler) SendMessage(text string) error {
	var m Message = Message{
		Text:   text,
		Name:   h.userName,
		Avatar: h.avaterURI,
	}
	return m.Send(h.hookURI)
}

// SendError send error to discord
func (h *Handler) SendError(e error) error {
	if e == nil {
		return fmt.Errorf("NilErrorError: Error is nil pointer")
	}

	var m Message = Message{
		Text:   fmt.Sprintf("以下の問題が発生しました。\n%s", e.Error()),
		Name:   h.errMessage.userName,
		Avatar: h.errMessage.avaterURI,
	}

	if m.Name == "" {
		m.Name = h.userName
	}

	if m.Avatar == "" {
		m.Avatar = h.avaterURI
	}

	if h.errorHookURI == "" {
		return m.Send(h.hookURI)
	}

	return m.Send(h.errorHookURI)
}
