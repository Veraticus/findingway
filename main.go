package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/Veraticus/trappingway/internal/discord"
	"github.com/Veraticus/trappingway/internal/scraper"

	"gopkg.in/yaml.v2"
)

func main() {
	os.Exit(run())
}

func run() int {
	var token string
	once := false

	for i := 1; i < len(os.Args); i++ {
		arg := string(os.Args[i])
		switch arg {
		default:
			{
				log.Printf("Unknown option '%s'", arg)
			}
		case "--once":
			{
				once = true
			}
		case "--token":
			{
				if i+1 >= len(os.Args) {
					log.Println("please specify the token string")
					os.Exit(1)
				}
				i++
				token = string(os.Args[i])
			}
		}
	}

	server := &discord.Discord{
		Token: token,
	}

	config, err := os.ReadFile("./config.yaml")

	if err != nil {
		log.Printf("Could not read config.yaml: %s\n", err)
		return 1
	}

	err = yaml.Unmarshal(config, &server)

	if err != nil {
		log.Printf("Bad config.yaml: %s\n", err)
		return 1
	}

	err = server.Start()
	defer server.Close()

	if err != nil {
		log.Printf("Could not instantiate Discord: %s\n", err)
		return 1
	}

	scraper := scraper.New("https://xivpf.com/listings")

	log.Printf("Starting Trappingway...\n")
	for {
		log.Printf("Scraping source...\n")

		err := scraper.Scrape()

		if err != nil {
			log.Printf("Scraper error: %f\n", err)
			return 1
		}

		var wg sync.WaitGroup
		for _, channel := range server.Channels {
			wg.Add(1)

			go func(c *discord.Channel) {
				log.Printf("Cleaning Discord for %v...\n", c.Duty)
				err = server.CleanChannel(c.ID)

				if err != nil {
					log.Printf("Error cleaning channel: %#v\n", err)
				}

				log.Printf("Updating Discord for %v...\n", c.Duty)
				err = server.PostListings(c.ID, scraper.Listings, c.Duty, c.DataCentres)

				if err != nil {
					log.Printf("Discord error updating message: %#v\n", err)
				}

				wg.Done()
			}(channel)
		}
		wg.Wait()

		if once {
			return 0
		}

		time.Sleep(3 * time.Minute)
	}
}
