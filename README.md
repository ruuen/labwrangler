# LabWrangler

WIP Discord bot using the Gateway API which helps manage my homelab Discord server and learn Golang.

## Planned Features

- Bulk delete messages in a channel (as Discord doesn't provide bulk delete in-app)
    - Filtering to control deletion of messages with certain reactions
- Monitor Diun image update notification channel and provide buttons to trigger container updates in different ways
    - Option to create PRs in Github that bump versions for later manual/automated action
    - Option to trigger a workflow for runners to update containers
- Manage hosted game server actions (ie. start/shutdown/restarting game server processes, performing and applying backups, game-specific commands)
    - Vanilla / Modded Minecraft
    - Space Engineers
- Restrict commands and actions by roles

## Hosting

LabWrangler is built using Discord's Gateway API - no inbound port forwarding is required.

The bot opens a HTTPS WebSocket connection to Discord for sending & receiving events.

A container image will be provided from this repo soon.

## Development Environment

- Go 1.25.4+ installed
- Discord application and bot user created using [Discord's developer site.](https://discord.com/developers/applications)

1. Clone the repo

```bash
git clone git@github.com:ruuen/labwrangler.git
```

2. Set following environment variables by either creating a `.env` file at project root, or defining them in your shell:

| Variable | Description |
| -------------- | --------------- |
| DISCORD_APP_ID | Application ID of your bot from Discord's developer site. |
| DISCORD_APP_TOKEN | Application token of your bot from Discord's developer site. |
| DISCORD_GUILD_ID | Guild ID of your Discord server. |

3. Run the app

```bash
go run main.go
```

## License

This project is licensed under the terms of the MIT license.
