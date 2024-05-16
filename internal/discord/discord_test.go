package discord

import (
	"testing"
	"time"

	"os"

	"github.com/Veraticus/findingway/internal/ffxiv"

	"github.com/stretchr/testify/assert"
)

func TestStartDiscord(t *testing.T) {
	token, ok := os.LookupEnv("DISCORD_TOKEN")
	assert.Equal(t, ok, true)

	disc := &Discord{Token: token}
	err := disc.Start()

	assert.Equal(t, err, nil)
}

func TestPostListings(t *testing.T) {

	token, ok := os.LookupEnv("DISCORD_TOKEN")
	assert.NotEqual(t, ok, nil)

	disc := &Discord{Token: token}
	err := disc.Start()

	assert.Equal(t, err, nil)

	now := time.Now()
	listings := []struct {
		listing *ffxiv.Listing
		wants   time.Time
	}{
		{
			listing: &ffxiv.Listing{Updated: "a second ago"},
			wants:   now.Add(time.Duration(-1) * time.Second),
		},
		{
			listing: &ffxiv.Listing{Updated: "a minute ago"},
			wants:   now.Add(time.Duration(-1) * time.Minute),
		},
		{
			listing: &ffxiv.Listing{Updated: "an hour ago"},
			wants:   now.Add(time.Duration(-1) * time.Hour),
		},
		{
			listing: &ffxiv.Listing{Updated: "37 seconds ago"},
			wants:   now.Add(time.Duration(-37) * time.Second),
		},
		{
			listing: &ffxiv.Listing{Updated: "22 minutes ago"},
			wants:   now.Add(time.Duration(-22) * time.Minute),
		},
		{
			listing: &ffxiv.Listing{Updated: "2 hours ago"},
			wants:   now.Add(time.Duration(-2) * time.Hour),
		},
	}

	var ffxivListings ffxiv.Listings
	for _, item := range listings {
		ffxivListings.Listings = append(ffxivListings.Listings, item.listing)
	}

	// testing in #staff-actions
	listErr := disc.PostListings("1174350271304958032", &ffxivListings, "Dragonsong's Reprise (Ultimate)", "Aether")
	assert.Equal(t, listErr, nil)
}
