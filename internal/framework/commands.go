package framework

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// CommandHandler is implemented by anything that handles a slash command.
type CommandHandler interface {
	// Definition returns the ApplicationCommand definition sent to Discord.
	Definition() *discordgo.ApplicationCommand
	// Handle is called when the command is invoked.
	Handle(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// CommandRegistry maps command names to their handlers.
type CommandRegistry struct {
	handlers map[string]CommandHandler
	// Preserve insertion order for syncing.
	order []string
}

func newCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		handlers: make(map[string]CommandHandler),
	}
}

// Add registers a CommandHandler. Panics if a handler with the same name is
// registered twice — this is intentional so misconfiguration is caught at
// startup, not silently at runtime.
func (r *CommandRegistry) Add(h CommandHandler) {
	name := h.Definition().Name
	if _, exists := r.handlers[name]; exists {
		panic(fmt.Sprintf("framework: duplicate command registered: %q", name))
	}
	r.handlers[name] = h
	r.order = append(r.order, name)
}

func (r *CommandRegistry) dispatch(s *discordgo.Session, i *discordgo.InteractionCreate) {
	name := i.ApplicationCommandData().Name
	h, ok := r.handlers[name]
	if !ok {
		log.Printf("framework: no handler for command %q", name)
		return
	}
	h.Handle(s, i)
}

// sync uploads all registered command definitions to Discord.
// App ID is derived from the session. Guild ID is empty for global commands.
// To scope to a guild during development, replace "" with your guild ID.
func (r *CommandRegistry) sync(s *discordgo.Session) error {
	appID := s.State.User.ID
	guildID := "" // replace with a guild ID to test commands instantly in a single server

	defs := make([]*discordgo.ApplicationCommand, 0, len(r.order))
	for _, name := range r.order {
		defs = append(defs, r.handlers[name].Definition())
	}

	for _, def := range defs {
		if _, err := s.ApplicationCommandCreate(appID, guildID, def); err != nil {
			return fmt.Errorf("framework: failed to register command %q: %w", def.Name, err)
		}
		log.Printf("framework: registered command /%s", def.Name)
	}
	return nil
}
