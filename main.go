package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/Veraticus/trappingway/internal/discord"
	"github.com/Veraticus/trappingway/internal/scraper"

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

	d := &discord.Discord{
		Token: discordToken,
	}

	config, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		panic(fmt.Errorf("Could not read config.yaml: %w", err))
	}
	yaml.Unmarshal(config, &d)

	err = d.Start()
	defer d.Session.Close()
	if err != nil {
		panic(fmt.Errorf("Could not instantiate Discord: %f", err))
	}

	scraper := scraper.New("https://xivpf.com")

	fmt.Printf("Starting Trappingway...\n")
	for {
		fmt.Printf("Scraping source...\n")
		err := scraper.Scrape()
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
				err = d.PostListings(c.ID, scraper.Listings, c.Duty, c.DataCentres)
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
