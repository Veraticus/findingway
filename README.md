# murult

Bot for Materia Ultimate Raid (MUR) discord server.

## Usage

Every configuration is provided via the environment variable.
This is not the most secure way to store secrets but leaking a discord token isn't that destructive I think?
Anyways, you can run it by:

```shell
export DISCORD_TOKEN=<discord bot token>
SLEEP=<how many minutes to sleep between job>
go run cmd/main.go
```
