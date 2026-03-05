package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/ryantbvt/voter/internal/framework"
)

func main() {
	// load configs
	conf := framework.LoadEnv()

	// initialize the discord bot
	discordServer, err := discordgo.New("Bot " + conf.DiscToken)
	if err != nil {
		log.Fatal("Error starting Discord Bot")
	}

	// add handlers
	handler := framework.NewHandler()
	discordServer.AddHandler(handler.MessageHandler)

	// open discord bot
	if err := discordServer.Open(); err != nil {
		log.Fatal("Error opening Discord connection", err)
	}

	defer discordServer.Close()

	log.Println("Discord bot is now running")

	// Wait for ctrl + C or killsig
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	log.Println("Shutting down bot")
}
