package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v2"

	murult "github.com/yuyuriyuri/murult"
)

var Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	os.Exit(run())
}

func run() int {
	configPath := "./config.yaml"
	once := false
	var sleep int64 = 3

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
		case "--config":
			{
				if i+1 >= len(os.Args) {
					log.Println("please specify the token string")
					os.Exit(1)
				}
				i++
				configPath = string(os.Args[i])
			}
		case "--sleep":
			{
				if i+1 >= len(os.Args) {
					log.Println("please specify the token string")
					os.Exit(1)
				}
				i++
				sleepStr := string(os.Args[i])
				sleep64, err := strconv.ParseInt(sleepStr, 10, 64)

				if err != nil {
					log.Printf("bad input for --sleep: %s\n", err)
					return 1
				}

				sleep = sleep64
			}
		}
	}

	configBytes, err := os.ReadFile(configPath)

	if err != nil {
		log.Printf("could not read config.yaml: %s\n", err)
		return 1
	}

	var config Config
	err = yaml.Unmarshal(configBytes, &config)

	if err != nil {
		log.Printf("failed to open config file: %s\n", err)
		return 1
	}

	session, err := discordgo.New("Bot " + config.Token)

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
		Token:     config.Token,
		World:     config.World,
		GuildId:   config.GuildId,
		ChannelId: config.ChannelId,
		Session:   session,
		Duties:    config.Duties,
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

		if once {
			return 0
		}
	}
}

type Config struct {
	Token     string            `yaml:"token"`
	World     string            `yaml:"world"`
	GuildId   string            `yaml:"guildId"`
	ChannelId string            `yaml:"channelId"`
	Duties    []string          `yaml:"duties"`
	EmojiDb   map[string]string `yaml:"emojis"`
}
