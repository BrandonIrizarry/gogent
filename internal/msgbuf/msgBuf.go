package msgbuf

import "google.golang.org/genai"

type MsgBuf struct {
	Messages []*genai.Content
}

func NewMsgBuf() MsgBuf {
	buf := make([]*genai.Content, 0)

	return MsgBuf{buf}
}

func (msgBuf *MsgBuf) AddMessage(message *genai.Content) {
	msgBuf.Messages = append(msgBuf.Messages, message)
}

func (msgBuf *MsgBuf) AddText(text string) {
	content := genai.NewContentFromText(text, "user")

	msgBuf.Messages = append(msgBuf.Messages, content)
}

func (msgBuf *MsgBuf) AddToolPart(part *genai.Part) {
	content := genai.NewContentFromParts([]*genai.Part{part}, "tool")

	msgBuf.Messages = append(msgBuf.Messages, content)
}
