package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Veraticus/findingway/internal/discord"
	"github.com/Veraticus/findingway/internal/scraper"

	"gopkg.in/yaml.v2"
)

func main() {
	discordToken, ok := os.LookupEnv("DISCORD_TOKEN")
	if !ok {
		panic("You must supply a DISCORD_TOKEN to start!")
	}
	once, ok := os.LookupEnv("ONCE")
	if !ok {
		once = "false"
	}
	discordToken = strings.TrimSpace(discordToken)

	d := &discord.Discord{
		Token: discordToken,
	}

	config, err := os.ReadFile("./config.yaml")
	if err != nil {
		panic(fmt.Errorf("Could not read config.yaml: %w", err))
	}
	yaml.Unmarshal(config, &d)

	err = d.Start()
	defer d.Session.Close()
	if err != nil {
		panic(fmt.Errorf("Could not instantiate Discord: %f", err))
	}

	scraper := &scraper.Scraper{Url: "https://xivpf.com"}

	fmt.Printf("Starting findingway...\n")
	for {
		fmt.Printf("Scraping source...\n")
		listings, err := scraper.Scrape()
		if err != nil {
			fmt.Printf("Scraper error: %f\n", err)
			continue
		}

		var wg sync.WaitGroup
		for _, channel := range d.Channels {
			wg.Add(1)
			go func(c *discord.Channel) {
				fmt.Printf("Cleaning Discord for %v...\n", c.Duty)
				err = d.CleanChannel(c.ID)
				if err != nil {
					fmt.Printf("Discord error cleaning channel: %f\n", err)
				}

				fmt.Printf("Updating Discord for %v...\n", c.Duty)
				err = d.PostListings(c.ID, listings, c.Duty, c.DataCentres)
				if err != nil {
					fmt.Printf("Discord error updating messagea: %f\n", err)
				}
				wg.Done()
			}(channel)
		}
		wg.Wait()
		if once != "false" {
			os.Exit(0)
		}

		time.Sleep(3 * time.Minute)
	}

}
