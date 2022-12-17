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

var Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

func main() {

	os.Exit(run())
}

func run() int {
	discordToken, ok := os.LookupEnv("DISCORD_TOKEN")
	if !ok {
		Logger.Println("You must supply a DISCORD_TOKEN to start!")
		return 1
	}
	once, ok := os.LookupEnv("ONCE")
	if !ok {
		once = "false"
	}

	d := &discord.Discord{
		Token: discordToken,
	}

	config, err := os.ReadFile("./config.yaml")
	if err != nil {
		Logger.Printf("Could not read config.yaml: %s\n", err)
		return 1
	}
	yaml.Unmarshal(config, &d)

	err = d.Start()
	defer d.Session.Close()
	if err != nil {
		log.Printf("Could not instantiate Discord: %s\n", err)
		return 1
	}

	scraper := scraper.New("https://xivpf.com")

	Logger.Printf("Starting Trappingway...\n")
	for {
		Logger.Printf("Scraping source...\n")
		err := scraper.Scrape()
		if err != nil {
			Logger.Printf("Scraper error: %f\n", err)
			continue
		}

		var wg sync.WaitGroup
		for _, channel := range d.Channels {
			wg.Add(1)
			go func(c *discord.Channel) {
				Logger.Printf("Cleaning Discord for %v...\n", c.Duty)
				err = d.CleanChannel(c.ID)
				if err != nil {
					Logger.Printf("Discord error cleaning channel: %f\n", err)
				}

				Logger.Printf("Updating Discord for %v...\n", c.Duty)
				err = d.PostListings(c.ID, scraper.Listings, c.Duty, c.DataCentres)
				if err != nil {
					Logger.Printf("Discord error updating messagea: %f\n", err)
				}
				wg.Done()
			}(channel)
		}
		wg.Wait()
		if once != "false" {
			return 0
		}

		time.Sleep(3 * time.Minute)
	}
}
