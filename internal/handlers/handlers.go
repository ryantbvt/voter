package handlers

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// PingButtonHandler handles the "ping:clicked" button.
type PingButtonHandler struct{}

func (h *PingButtonHandler) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "You clicked the button! 🎉",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

// FeedbackModalHandler handles the "modal:feedback" modal submission.
type FeedbackModalHandler struct{}

func (h *FeedbackModalHandler) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()

	var text string
	for _, row := range data.Components {
		row, ok := row.(*discordgo.ActionsRow)
		if !ok {
			continue
		}
		for _, comp := range row.Components {
			input, ok := comp.(*discordgo.TextInput)
			if !ok {
				continue
			}
			if input.CustomID == "feedback_text" {
				text = input.Value
			}
		}
	}

	// In a real app you'd persist this somewhere.
	log.Printf("feedback from %s: %s", i.Member.User.Username, text)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Thanks for your feedback! We received: *%s*", text),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
