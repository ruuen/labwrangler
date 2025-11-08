package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/ruuen/labwrangler/commands"
)

func main() {
	// TODO: load from env files, using flags for now
	appId := flag.String("appId", "", "Application ID for Discord")
	appToken := flag.String("appToken", "", "Application token for Discord")
	// this is fine though but also give option to load from env/config file in some priority
	guildId := flag.String("guildId", "", "Guild ID to monitor")
	flag.Parse()

	s, err := discordgo.New("Bot " + *appToken)
	if err != nil {
		log.Fatalf("Unable to open Discord session: %v", err)
	}

	targetCommands := []*discordgo.ApplicationCommand{
		&commands.ChannelClearCommand,
	}

	createdCommands, err := s.ApplicationCommandBulkOverwrite(*appId, *guildId, targetCommands)
	if err != nil {
		log.Fatalf("Failed to create command: %v", err)
	}

	for _, cmd := range createdCommands {
		log.Printf("Created command: %v", cmd.Name)
	}

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
