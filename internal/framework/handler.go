package framework

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	Prefix = "!"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Validate message is not itself
	if m.Author.Bot {
		return
	}

	// Ignore messages without a prefix
	if !strings.HasPrefix(m.Content, Prefix) {
		return
	}

	// Trim command
	content := strings.TrimPrefix(m.Content, Prefix)
	parts := strings.Fields(content)

	if len(parts) == 0 {
		return
	}

	cmdName := parts[0]
	args := parts[1:]

	// Validate if command exists
	commands := GetCommands()
	cmd, exists := commands[cmdName]
	if !exists {
		s.ChannelMessageSend(m.ChannelID, "Unknown command")
		return
	}

	cmd.Execute(h, s, m, args)
}
