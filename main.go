package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Veraticus/trappingway/internal/discord"
	"github.com/Veraticus/trappingway/internal/scraper"
)

func main() {
	discordToken, ok := os.LookupEnv("DISCORD_TOKEN")
	if !ok {
		panic("You must supply a DISCORD_TOKEN to start!")
	}

	discordChannelId, ok := os.LookupEnv("DISCORD_CHANNEL_ID")
	if !ok {
		panic("You must supply a DISCORD_CHANNEL_ID to start!")
	}

	dataCentre, ok := os.LookupEnv("DATA_CENTRE")
	if !ok {
		panic("You must supply a DATA_CENTRE to start!")
	}

	duty, ok := os.LookupEnv("DUTY")
	if !ok {
		panic("You must supply a DUTY to start!")
	}

	discordMessageId, ok := os.LookupEnv("DISCORD_MESSAGE_ID")

	discord := &discord.Discord{
		Token:     discordToken,
		ChannelId: discordChannelId,
		MessageId: discordMessageId,
	}
	err := discord.Start()
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

		fmt.Printf("Updating Discord...\n")
		err = discord.UpdateMessage(scraper.Listings, dataCentre, duty)
		if err != nil {
			fmt.Printf("Discord error: %f\n", err)
		}
		time.Sleep(1 * time.Minute)
	}

}
