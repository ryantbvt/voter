package framework

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// Router holds all registered handlers and dispatches incoming interactions.
type Router struct {
	commands *CommandRegistry
	buttons  *ButtonRegistry
	modals   *ModalRegistry
}

// NewRouter creates a new Router with empty registries.
func NewRouter() *Router {
	return &Router{
		commands: newCommandRegistry(),
		buttons:  newButtonRegistry(),
		modals:   newModalRegistry(),
	}
}

// Commands returns the slash command registry (used to register handlers and sync with Discord).
func (r *Router) Commands() *CommandRegistry {
	return r.commands
}

// Buttons returns the button interaction registry.
func (r *Router) Buttons() *ButtonRegistry {
	return r.buttons
}

// Modals returns the modal submit registry.
func (r *Router) Modals() *ModalRegistry {
	return r.modals
}

// Handler returns a discordgo event handler function to pass to session.AddHandler.
func (r *Router) Handler() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			r.commands.dispatch(s, i)
		case discordgo.InteractionMessageComponent:
			data := i.MessageComponentData()
			switch data.ComponentType {
			case discordgo.ButtonComponent:
				r.buttons.dispatch(s, i)
			default:
				log.Printf("unhandled component type: %v", data.ComponentType)
			}
		case discordgo.InteractionModalSubmit:
			r.modals.dispatch(s, i)
		default:
			log.Printf("unhandled interaction type: %v", i.Type)
		}
	}
}

// Sync registers all slash commands with Discord globally.
// Call this after Open() so that s.State.User.ID is available.
// To scope commands to a specific guild during development, see the commented
// guildID parameter in commands.sync.
func (r *Router) Sync(s *discordgo.Session) error {
	return r.commands.sync(s)
}
