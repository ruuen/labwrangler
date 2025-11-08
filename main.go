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
	// Config settings; priority: cmdline->env vars->config file
	// TODO: add a wrapper struct for app's config settings and get this out of main
	// var guildId string
	// var appId string
	// var appToken string
	// flag.StringVar(&guildId, "guildId", "", "Guild ID to monitor")
	// flag.Parse()
	//
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Unable to load environment variables: %v", err)
	// }
	// appId = os.Getenv("DISCORD_APP_ID")
	// appToken = os.Getenv("DISCORD_APP_TOKEN")
	// if guildId == "" {
	// 	guildId = os.Getenv("DISCORD_GUILD_ID")
	// }

	config := config.New()
	err := config.Init()
	if err != nil {
		log.Fatalln(err)
	}

	// TODO: this needs later improvement with the config wrapper to dynamically check and return missing settings. im the only dunce who would forget right now so it's fine :)
	// switch {
	// case appId == "":
	// case appToken == "":
	// case guildId == "":
	// 	log.Fatalf("You haven't provided a configuration setting. I'll leave it up to you to guess which one. Good luck!")
	// }

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
