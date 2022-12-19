package main

import (
	"log"
	"os"
	"strconv"
	"time"

	murult "github.com/yuyuriyuri/trappingway-go"
)

func main() {
	var sleep int64 = 3
	once := false

	for i := 1; i < len(os.Args); i++ {
		arg := string(os.Args[i])
		switch arg {
		default:
			{
				log.Printf("unknown option '%s'", arg)
			}
		case "--once":
			{
				once = true
			}
		}
	}

	token, tokenExists := os.LookupEnv("DISCORD_TOKEN")

	if !tokenExists {
		log.Printf("Please provide DISCORD_TOKEN")
		os.Exit(1)
	}

	guildId, guildIdExists := os.LookupEnv("GUILD_ID")

	if !guildIdExists {
		log.Printf("Please provide GUILD_ID")
		os.Exit(1)
	}

	channelId, channelIdExists := os.LookupEnv("CHANNEL_ID")

	if !channelIdExists {
		log.Printf("Please provide CHANNEL_ID")
		os.Exit(1)
	}

	sleepStr, sleepEnvExists := os.LookupEnv("SLEEP")

	if sleepEnvExists {
		sleep64, err := strconv.ParseInt(sleepStr, 10, 64)

		if err != nil {
			log.Printf("bad input for --sleep: %s\n", err)
			os.Exit(1)
		}

		sleep = sleep64
	}

	murult.Logger.Println("starting server...")
	server := murult.NewServer(token, guildId, channelId)
	defer server.CloseServer()

	for {
		server.Run()
		if once {
			return
		} else {
			time.Sleep(time.Duration(sleep * int64(time.Minute)))
		}
	}
}
