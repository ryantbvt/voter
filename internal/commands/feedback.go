package commands

import (
	"github.com/bwmarrin/discordgo"
)

// FeedbackCommand opens a modal when invoked.
// The modal submit is handled separately in handlers/feedback_modal.go.
type FeedbackCommand struct{}

func (c *FeedbackCommand) Definition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "feedback",
		Description: "Submit feedback via a form",
	}
}

func (c *FeedbackCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "modal:feedback",
			Title:    "Submit Feedback",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "feedback_text",
							Label:       "Your feedback",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "Tell us what you think...",
							Required:    true,
							MaxLength:   1000,
						},
					},
				},
			},
		},
	})
}
