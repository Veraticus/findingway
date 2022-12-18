package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"

	murult "github.com/yuyuriyuri/murult"
)

var Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	os.Exit(run())
}

func run() int {
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
		log.Printf("Please provide a discord token")
		return 1
	}

	guildId, guildIdExists := os.LookupEnv("GUILD_ID")

	if !guildIdExists {
		log.Printf("Please provide a guild ID to get emojis from")
		return 1
	}

	channelId, channelIdExists := os.LookupEnv("CHANNEL_ID")

	if !channelIdExists {
		log.Printf("Please provide a channel ID to write to")
		return 1
	}

	world, worldExists := os.LookupEnv("WORLD")

	if !worldExists {
		log.Printf("Please provide a world to filter for")
		return 1
	}

	sleepStr, sleepEnvExists := os.LookupEnv("SLEEP")

	if sleepEnvExists {
		sleep64, err := strconv.ParseInt(sleepStr, 10, 64)

		if err != nil {
			log.Printf("bad input for --sleep: %s\n", err)
			return 1
		}

		sleep = sleep64
	}

	session, err := discordgo.New("Bot " + token)

	if err != nil {
		log.Printf("could not start Discord: %s\n", err)
		return 1
	}

	err = session.Open()

	if err != nil {
		log.Printf("could not open Discord session: %s\n", err)
		return 1
	}

	// Not sure how to check for the failure here :/
	defer session.Close()

	scraper := murult.NewScraper("https://xivpf.com/listings")
	server := &murult.Server{
		Token:     token,
		World:     world,
		GuildId:   guildId,
		ChannelId: channelId,
		Session:   session,
		Duties: []string{
			"The Weapon's Refrain (Ultimate)",
			"The Unending Coil of Bahamut (Ultimate)",
			"The Epic of Alexander (Ultimate)",
			"Dragonsong's Reprise (Ultimate)"},
	}

	first := true

	// TODO: Instead of just continuing, it should probably report the error somehow
	for {
		if !first {
			time.Sleep(time.Duration(sleep * int64(time.Minute)))
		}

		first = false
		err := scraper.Scrape()

		if err != nil {
			Logger.Printf("scraper error: %f\n", err)
			if once {
				return 0
			} else {
				continue
			}
		}

		err = server.CleanChannel()

		if err != nil {
			Logger.Printf("discord error cleaning channel: %f\n", err)
			if once {
				return 0
			} else {
				continue
			}
		}

		emojis, err := session.GuildEmojis(server.GuildId)

		if err != nil {
			log.Printf("could not get server emojis: %s\n", err)
			if once {
				return 0
			} else {
				continue
			}
		}

		err = server.PostListings(scraper.Listings, emojis)

		if err != nil {
			Logger.Printf("discord error updating messagea: %f\n", err)
			if once {
				return 0
			} else {
				continue
			}
		}
	}
}
