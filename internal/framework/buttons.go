package framework

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// ButtonHandler handles a button component interaction.
//
// Button custom IDs can be either exact strings or prefixes. For example,
// a paginator might use IDs like "page:next:42" and "page:prev:42" — you'd
// register a single handler under the prefix "page:" using AddPrefix.
type ButtonHandler interface {
	Handle(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// ButtonHandlerFunc is a convenience adapter so plain functions satisfy ButtonHandler.
type ButtonHandlerFunc func(s *discordgo.Session, i *discordgo.InteractionCreate)

func (f ButtonHandlerFunc) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	f(s, i)
}

// ButtonRegistry maps button custom IDs (exact or prefix) to handlers.
type ButtonRegistry struct {
	exact  map[string]ButtonHandler
	prefix []prefixEntry
}

type prefixEntry struct {
	prefix  string
	handler ButtonHandler
}

func newButtonRegistry() *ButtonRegistry {
	return &ButtonRegistry{
		exact: make(map[string]ButtonHandler),
	}
}

// Add registers a handler for an exact custom_id.
func (r *ButtonRegistry) Add(customID string, h ButtonHandler) {
	if _, exists := r.exact[customID]; exists {
		panic(fmt.Sprintf("framework: duplicate button ID registered: %q", customID))
	}
	r.exact[customID] = h
}

// AddFunc is a convenience wrapper around Add for plain functions.
func (r *ButtonRegistry) AddFunc(customID string, h ButtonHandlerFunc) {
	r.Add(customID, h)
}

// AddPrefix registers a handler that matches any custom_id starting with prefix.
// Useful for dynamic IDs that carry embedded state (e.g. "confirm:delete:123").
func (r *ButtonRegistry) AddPrefix(prefix string, h ButtonHandler) {
	r.prefix = append(r.prefix, prefixEntry{prefix: prefix, handler: h})
}

func (r *ButtonRegistry) dispatch(s *discordgo.Session, i *discordgo.InteractionCreate) {
	id := i.MessageComponentData().CustomID

	// Exact match first.
	if h, ok := r.exact[id]; ok {
		h.Handle(s, i)
		return
	}

	// Prefix match (first registered wins).
	for _, entry := range r.prefix {
		if len(id) >= len(entry.prefix) && id[:len(entry.prefix)] == entry.prefix {
			entry.handler.Handle(s, i)
			return
		}
	}

	log.Printf("framework: no handler for button custom_id %q", id)
}
