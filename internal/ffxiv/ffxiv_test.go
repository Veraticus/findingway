package ffxiv

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExpiresAt(t *testing.T) {
	now := time.Now()

	tests := []struct {
		listing *Listing
		wants   time.Time
	}{
		{
			listing: &Listing{Expires: "in a second"},
			wants:   now.Add(time.Duration(1) * time.Second),
		},
		{
			listing: &Listing{Expires: "in a minute"},
			wants:   now.Add(time.Duration(1) * time.Minute),
		},
		{
			listing: &Listing{Expires: "in an hour"},
			wants:   now.Add(time.Duration(1) * time.Hour),
		},
		{
			listing: &Listing{Expires: "in 37 seconds"},
			wants:   now.Add(time.Duration(37) * time.Second),
		},
		{
			listing: &Listing{Expires: "in 22 minutes"},
			wants:   now.Add(time.Duration(22) * time.Minute),
		},
		{
			listing: &Listing{Expires: "in 2 hours"},
			wants:   now.Add(time.Duration(2) * time.Hour),
		},
	}

	for _, tt := range tests {
		expiresAt, err := tt.listing.ExpiresAt()
		assert.Nil(t, err)
		assert.WithinDuration(t, tt.wants, expiresAt, 10*time.Millisecond)
	}
}

func TestUpdatedAt(t *testing.T) {
	now := time.Now()

	tests := []struct {
		listing *Listing
		wants   time.Time
	}{
		{
			listing: &Listing{Updated: "a second ago"},
			wants:   now.Add(time.Duration(-1) * time.Second),
		},
		{
			listing: &Listing{Updated: "a minute ago"},
			wants:   now.Add(time.Duration(-1) * time.Minute),
		},
		{
			listing: &Listing{Updated: "an hour ago"},
			wants:   now.Add(time.Duration(-1) * time.Hour),
		},
		{
			listing: &Listing{Updated: "37 seconds ago"},
			wants:   now.Add(time.Duration(-37) * time.Second),
		},
		{
			listing: &Listing{Updated: "22 minutes ago"},
			wants:   now.Add(time.Duration(-22) * time.Minute),
		},
		{
			listing: &Listing{Updated: "2 hours ago"},
			wants:   now.Add(time.Duration(-2) * time.Hour),
		},
	}

	for _, tt := range tests {
		updatedAt, err := tt.listing.UpdatedAt()
		assert.Nil(t, err)
		assert.WithinDuration(t, tt.wants, updatedAt, 10*time.Millisecond)
	}
}
