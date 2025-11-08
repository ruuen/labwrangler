package config

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Settings map[string]string
}

type setting struct {
	name string
	cmdKey string
	cmdDescription string
	envKey string
	configKey string
	required bool
	value string
}

func New() (s *AppConfig) {
	return &AppConfig{
		Settings: make(map[string]string),
	}
}

func (c *AppConfig) Init() error {
	settings := []*setting{
		{
			name: "AppId",
			envKey: "DISCORD_APP_ID",
			configKey: "appId",
			required: true,
		},
		{
			name: "AppToken",
			envKey: "DISCORD_APP_TOKEN",
			configKey: "appToken",
			required: true,
		},
		{
			name: "GuildId",
			cmdKey: "guildId",
			cmdDescription: "Application Token of Discord bot",
			envKey: "DISCORD_GUILD_ID",
			configKey: "guildId",
			required: true,
		},
	}

	// cmdline params
	for _, v := range settings {
		if v.cmdKey != "" {
			flag.StringVar(&v.value, v.cmdKey, "", v.cmdDescription)
		}
	}
	flag.Parse()

	// env vars
	godotenv.Load()
	for _, v := range settings {
		if v.value == "" && v.envKey != "" {
			env := os.Getenv(v.envKey)
			v.value = env
		}
	}

	// config file (TBC)

	// app defaults (TBC)

	for _, v := range settings {
		c.Settings[v.name] = v.value
	}

	return nil
}
