package murult

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type RawPost struct {
	DataCentre  string
	PfCategory  string
	Duty        string `selector:".left .duty"`
	Tags        string `selector:".left .description span"`
	TagsColor   string `selector:".left .description span" attr:"class"`
	Description string
	MinIL       string `selector:".middle .stat .value"`
	Creator     string `selector:".right .creator .text"`
	World       string `selector:".right .world .text"`
	Expires     string `selector:".right .expires .text"`
	Updated     string `selector:".right .updated .text"`
	Slots       []string
}

func NewRawPost() *RawPost {
	return &RawPost{
		Slots: make([]string, 0),
	}
}

func (p *RawPost) Stringify(emojis []*discordgo.Emoji) string {
	var result strings.Builder
	result.Grow(100)

	// Title
	result.WriteString("\n***")
	result.WriteString(strings.ToUpper(p.Duty))
	result.WriteString("***\n")

	// Creator
	result.WriteString("Created by: ")
	result.WriteString(p.Creator)
	result.WriteByte('\n')

	// Creation time
	result.WriteString(EmojiFromStr("hourglass", emojis))
	result.WriteString(" Expires At: ")
	result.WriteString(fmt.Sprintf("<t:%d>", p.ExpiresAt().Unix()))
	result.WriteByte('\n')

	// Last activity
	result.WriteString(EmojiFromStr("stopwatch", emojis))
	result.WriteString(" Last updated: ")
	result.WriteString(fmt.Sprintf("<t:%d>", p.UpdatedAt().Unix()))
	result.WriteByte('\n')

	// Description
	result.WriteString("-----------\n")
	result.WriteString(p.Description)
	result.WriteString("\n-----------")
	result.WriteByte('\n')

	// Roster
	result.WriteString("Roster: ")
	for _, slot := range p.Slots {
		result.WriteString(JobEmojiFromStr(slot, emojis))
	}
	result.WriteByte('\n')

	// Tags
	result.WriteString("Tags: ")
	result.WriteString(p.Tags)

	return result.String()
}

var expiresSecondsRegexp = regexp.MustCompile(`in (\d+) seconds`)
var expiresMinutesRegexp = regexp.MustCompile(`in (\d+) minutes`)
var expiresHoursRegexp = regexp.MustCompile(`in (\d+) hours`)

func (l *RawPost) ExpiresAt() time.Time {
	now := time.Now()

	if l.Expires == "" {
		return now
	}

	if l.Expires == "now" {
		return now
	}

	if l.Expires == "in a second" {
		return now.Add(time.Duration(1) * time.Second)
	}

	if l.Expires == "in a minute" {
		return now.Add(time.Duration(1) * time.Minute)
	}

	if l.Expires == "in an hour" {
		return now.Add(time.Duration(1) * time.Hour)
	}

	match := expiresSecondsRegexp.FindStringSubmatch(l.Expires)
	if len(match) != 0 {
		seconds, err := strconv.Atoi(match[1])
		if err != nil {
			Logger.Printf("could not parse time %v: %s", l.Expires, err)
			return time.Time{}
		}
		return now.Add(time.Duration(seconds) * time.Second)
	}

	match = expiresMinutesRegexp.FindStringSubmatch(l.Expires)
	if len(match) != 0 {
		minutes, err := strconv.Atoi(match[1])
		if err != nil {
			Logger.Printf("could not parse time %v: %s", l.Expires, err)
			return time.Time{}
		}
		return now.Add(time.Duration(minutes) * time.Minute)
	}

	match = expiresHoursRegexp.FindStringSubmatch(l.Expires)
	if len(match) != 0 {
		hours, err := strconv.Atoi(match[1])
		if err != nil {
			Logger.Printf("could not parse time %v: %s", l.Expires, err)
			return time.Time{}
		}
		return now.Add(time.Duration(hours) * time.Hour)
	}

	Logger.Printf("failed to parse time %v", l.Expires)
	return time.Time{}
}

var updatedSecondsRegexp = regexp.MustCompile(`(\d+) seconds ago`)
var updatedMinutesRegexp = regexp.MustCompile(`(\d+) minutes ago`)
var updatedHoursRegexp = regexp.MustCompile(`(\d+) hours ago`)

func (l *RawPost) UpdatedAt() time.Time {
	now := time.Now()

	if l.Updated == "" {
		return now
	}

	if l.Updated == "now" {
		return now
	}

	if l.Updated == "a second ago" {
		return now.Add(time.Duration(-1) * time.Second)
	}

	if l.Updated == "a minute ago" {
		return now.Add(time.Duration(-1) * time.Minute)
	}

	if l.Updated == "an hour ago" {
		return now.Add(time.Duration(-1) * time.Hour)
	}

	match := updatedSecondsRegexp.FindStringSubmatch(l.Updated)
	if len(match) != 0 {
		seconds, err := strconv.Atoi(match[1])
		if err != nil {
			Logger.Printf("could not parse time %v: %s", l.Updated, err)
			return time.Time{}
		}
		return now.Add(time.Duration(-seconds) * time.Second)
	}

	match = updatedMinutesRegexp.FindStringSubmatch(l.Updated)
	if len(match) != 0 {
		minutes, err := strconv.Atoi(match[1])
		if err != nil {
			Logger.Printf("could not parse time %v: %s", l.Updated, err)
			return time.Time{}
		}
		return now.Add(time.Duration(-minutes) * time.Minute)
	}

	match = updatedHoursRegexp.FindStringSubmatch(l.Updated)
	if len(match) != 0 {
		hours, err := strconv.Atoi(match[1])
		if err != nil {
			Logger.Printf("could not parse time %v: %s", l.Updated, err)
			return time.Time{}
		}
		return now.Add(time.Duration(-hours) * time.Hour)
	}

	Logger.Printf("failed to parse time %v", l.Updated)
	return time.Time{}
}
