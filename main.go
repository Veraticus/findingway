package main

import (
	"fmt"
	"io/ioutil"
	"os"
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

	discord := &discord.Discord{
		Token: discordToken,
	}

	config, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		panic(fmt.Errorf("Could not read config.yaml: %w", err))
	}
	yaml.Unmarshal(config, &discord)

	err = discord.Start()
	defer discord.Session.Close()
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

		for _, channel := range discord.Channels {
			fmt.Printf("Cleaning Discord for %v...\n", channel.Duty)
			err = discord.CleanChannel(channel.ID)
			if err != nil {
				fmt.Printf("Discord error cleaning channel: %f\n", err)
			}

			fmt.Printf("Updating Discord for %v...\n", channel.Duty)
			err = discord.PostListings(channel.ID, scraper.Listings, channel.Duty, channel.DataCentres)
			if err != nil {
				fmt.Printf("Discord error updating messagea: %f\n", err)
			}

			if once != "false" {
				os.Exit(0)
			}
		}
		time.Sleep(3 * time.Minute)
	}

}
