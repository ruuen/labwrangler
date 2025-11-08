package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/ruuen/labwrangler/commands"
	"github.com/ruuen/labwrangler/config"
)

func main() {
	config := config.New()
	err := config.Init()
	if err != nil {
		log.Fatalln(err)
	}

	s, err := discordgo.New("Bot " + config.Settings["AppToken"])
	if err != nil {
		log.Fatalf("Unable to open Discord session: %v", err)
	}

	// Command registration
	targetCommands := []*discordgo.ApplicationCommand{
		&commands.ChannelClearCommand,
	}

	createdCommands, err := s.ApplicationCommandBulkOverwrite(config.Settings["AppId"], config.Settings["GuildId"], targetCommands)
	if err != nil {
		log.Fatalf("Failed to create command: %v", err)
	}

	for _, cmd := range createdCommands {
		log.Printf("Created command: %v", cmd.Name)
	}

	// Handler registration
	s.AddHandler(commands.ChannelClearHandler)

	err = s.Open()
	if err != nil {
		log.Fatalf("Failed to open websocket connection with Discord: %v", err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sig

	s.Close()
}
