package framework

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Name        string
	Description string
	Execute     func(h *Handler, s *discordgo.Session, m *discordgo.MessageCreate, args []string)
}

func GetCommands() map[string]Command {
	return map[string]Command{
		"ping": {
			Name:        "ping",
			Description: "Pong",
			Execute:     pingPong,
		},
	}
}

func pingPong(h *Handler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	log.Println("ping/pong request received")
	s.ChannelMessageSend(m.ChannelID, "Pong")
}
