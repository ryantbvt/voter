package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"github.com/ryantbvt/voter/internal/commands"
	"github.com/ryantbvt/voter/internal/framework"
	"github.com/ryantbvt/voter/internal/handlers"
)

func main() {
	cfg := framework.LoadEnv()

	s, err := discordgo.New("Bot " + cfg.DiscToken)
	if err != nil {
		log.Fatalf("error creating session: %v", err)
	}

	// ── Build the router ─────────────────────────────────────────────────────

	router := framework.NewRouter()

	// Slash commands
	router.Commands().Add(&commands.PingCommand{})
	router.Commands().Add(&commands.FeedbackCommand{})

	// Buttons — exact ID
	router.Buttons().Add("ping:clicked", &handlers.PingButtonHandler{})

	// Modals — prefix match so "modal:feedback" and "modal:feedback:123" both hit this handler
	router.Modals().AddPrefix("modal:feedback", &handlers.FeedbackModalHandler{})

	// ── Wire up and open ──────────────────────────────────────────────────────

	s.AddHandler(router.Handler())
	s.Identify.Intents = discordgo.IntentsGuilds

	if err := s.Open(); err != nil {
		log.Fatalf("error opening connection: %v", err)
	}
	defer s.Close()

	// Sync commands — app ID is derived from the session automatically.
	// To scope to a single guild during dev, edit guildID in framework/commands.go.
	if err := router.Sync(s); err != nil {
		log.Fatalf("error syncing commands: %v", err)
	}

	log.Println("Bot is running. Press CTRL+C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}
