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

func TestMostRecentUpdateTime(t *testing.T) {
	expected := &Listing{Updated: "a second ago"}

	listings := &Listings{
		Listings: []*Listing{
			{Updated: "20 seconds ago"},
			{Updated: "43 minutes ago"},
			expected,
			{Updated: "an hour ago"},
		},
	}

	mostRecent, err := listings.MostRecentUpdated()
	assert.Nil(t, err)
	assert.Equal(t, expected, mostRecent)
}

func TestUpdatedWithinLast(t *testing.T) {
	expected1 := &Listing{Updated: "a second ago"}
	expected2 := &Listing{Updated: "40 seconds ago"}
	expected3 := &Listing{Updated: "a minute ago"}
	expected4 := &Listing{Updated: "2 minutes ago"}
	bad1 := &Listing{Updated: "5 minutes ago"}
	bad2 := &Listing{Updated: "22 minutes ago"}
	bad3 := &Listing{Updated: "an hour ago"}
	bad4 := &Listing{Updated: "3 hours ago"}

	listings := &Listings{
		Listings: []*Listing{
			bad4,
			expected1,
			bad3,
			expected4,
			expected2,
			bad2,
			bad1,
			expected3,
		},
	}

	mostRecentListings, err := listings.UpdatedWithinLast(3 * time.Minute)
	assert.Nil(t, err)
	assert.Equal(t, []*Listing{expected1, expected4, expected2, expected3}, mostRecentListings.Listings)
}
