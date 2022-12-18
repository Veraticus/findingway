# murult

Bot for Materia Ultimate Raid (MUR) discord server.

## Usage

Every configuration is provided via the environment variable.
This is not the most secure way to store secrets but leaking a discord token isn't that destructive I think?
Anyways, you can run it by:

```shell
export DISCORD_TOKEN=<discord bot token>
export GUILD_ID=<guild ID to get emojis from>
export CHANNEL_ID=<channel ID to write to>
export WORLD=<world to look for>
SLEEP=<how many minutes to sleep between job>
go run cmd/main.go
```

It also accepts 1 command line argument,

1. `--once`, which causes the server to run the loop once and exits.
   Mainly useful for debugging purposes or setting this up as a job from a CRON.
