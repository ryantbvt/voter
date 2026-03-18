package commands

import (
	"github.com/bwmarrin/discordgo"
)

// PingCommand responds to /ping with a simple pong message.
// It also sends back a button so you can see button wiring in action.
type PingCommand struct{}

func (c *PingCommand) Definition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Replies with Pong!",
	}
}

func (c *PingCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong! 🏓",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Click me",
							Style:    discordgo.PrimaryButton,
							CustomID: "ping:clicked",
						},
					},
				},
			},
		},
	})
}
