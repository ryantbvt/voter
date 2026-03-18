package framework

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// ModalHandler handles a modal submit interaction.
//
// Like buttons, modals are matched by custom_id — exact or prefix.
type ModalHandler interface {
	Handle(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// ModalHandlerFunc is a convenience adapter so plain functions satisfy ModalHandler.
type ModalHandlerFunc func(s *discordgo.Session, i *discordgo.InteractionCreate)

func (f ModalHandlerFunc) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	f(s, i)
}

// ModalRegistry maps modal custom IDs (exact or prefix) to handlers.
type ModalRegistry struct {
	exact  map[string]ModalHandler
	prefix []modalPrefixEntry
}

type modalPrefixEntry struct {
	prefix  string
	handler ModalHandler
}

func newModalRegistry() *ModalRegistry {
	return &ModalRegistry{
		exact: make(map[string]ModalHandler),
	}
}

// Add registers a handler for an exact modal custom_id.
func (r *ModalRegistry) Add(customID string, h ModalHandler) {
	if _, exists := r.exact[customID]; exists {
		panic(fmt.Sprintf("framework: duplicate modal ID registered: %q", customID))
	}
	r.exact[customID] = h
}

// AddFunc is a convenience wrapper around Add for plain functions.
func (r *ModalRegistry) AddFunc(customID string, h ModalHandlerFunc) {
	r.Add(customID, h)
}

// AddPrefix registers a handler that matches any custom_id starting with prefix.
func (r *ModalRegistry) AddPrefix(prefix string, h ModalHandler) {
	r.prefix = append(r.prefix, modalPrefixEntry{prefix: prefix, handler: h})
}

func (r *ModalRegistry) dispatch(s *discordgo.Session, i *discordgo.InteractionCreate) {
	id := i.ModalSubmitData().CustomID

	if h, ok := r.exact[id]; ok {
		h.Handle(s, i)
		return
	}

	for _, entry := range r.prefix {
		if len(id) >= len(entry.prefix) && id[:len(entry.prefix)] == entry.prefix {
			entry.handler.Handle(s, i)
			return
		}
	}

	log.Printf("framework: no handler for modal custom_id %q", id)
}
