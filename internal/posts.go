package murult

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Post struct {
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
	Party       []*Slot
	MessageId   string
}

func (p *Post) Stringify(emojis []*discordgo.Emoji) string {
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
	result.WriteString(" Time left: ")
	result.WriteString(p.Expires)
	result.WriteByte('\n')

	// Last activity
	result.WriteString(EmojiFromStr("stopwatch", emojis))
	result.WriteString(" Last updated: ")
	result.WriteString(p.Updated)
	result.WriteByte('\n')

	// Description
	result.WriteString("-----------\n")
	result.WriteString(p.Description)
	result.WriteString("\n-----------")
	result.WriteByte('\n')

	// Roster
	result.WriteString("Roster: ")
	for _, slot := range p.Party {
		if slot.Filled {
			result.WriteString(JobEmojiFromStr(slot.Job, emojis))
		} else {
			result.WriteString(slot.Roles.Emoji(emojis) + " ")
		}
	}
	result.WriteByte('\n')

	// Tags
	result.WriteString("Tags: ")
	tags := parseTags(p.Tags)
	for _, t := range tags {
		result.WriteString(t)
		result.WriteByte(' ')
	}

	return result.String()

}

func (p *Post) GetUpdated(emojis []*discordgo.Emoji) string {
	return fmt.Sprintf("%s %s", EmojiFromStr("stopwatch", emojis), p.Updated)
}

func parseTags(tags string) []string {
	result := make([]string, 0)
	raw := []byte(tags)
	start := 0

	for i, v := range raw {
		if v == '[' {
			start = i
		} else if v == ']' {
			result = append(result, string(raw[start:i+1]))
		}
	}

	return result
}

var expiresSecondsRegexp = regexp.MustCompile(`in (\d+) seconds`)
var expiresMinutesRegexp = regexp.MustCompile(`in (\d+) minutes`)
var expiresHoursRegexp = regexp.MustCompile(`in (\d+) hours`)

func (p *Post) ExpiresAt() (time.Time, error) {
	now := time.Now()

	if p.Expires == "" {
		return now, nil
	}

	if p.Expires == "now" {
		return now, nil
	}

	if p.Expires == "in a second" {
		return now.Add(time.Duration(1) * time.Second), nil
	}

	if p.Expires == "in a minute" {
		return now.Add(time.Duration(1) * time.Minute), nil
	}

	if p.Expires == "in an hour" {
		return now.Add(time.Duration(1) * time.Hour), nil
	}

	match := expiresSecondsRegexp.FindStringSubmatch(p.Expires)
	if len(match) != 0 {
		seconds, err := strconv.Atoi(match[1])
		if err != nil {
			return now, fmt.Errorf("could not parse time %v: %w", p.Expires, err)
		}
		return now.Add(time.Duration(seconds) * time.Second), nil
	}

	match = expiresMinutesRegexp.FindStringSubmatch(p.Expires)
	if len(match) != 0 {
		minutes, err := strconv.Atoi(match[1])
		if err != nil {
			return now, fmt.Errorf("could not parse time %v: %w", p.Expires, err)
		}
		return now.Add(time.Duration(minutes) * time.Minute), nil
	}

	match = expiresHoursRegexp.FindStringSubmatch(p.Expires)
	if len(match) != 0 {
		hours, err := strconv.Atoi(match[1])
		if err != nil {
			return now, fmt.Errorf("could not parse time %v: %w", p.Expires, err)
		}
		return now.Add(time.Duration(hours) * time.Hour), nil
	}

	return now, fmt.Errorf("failed to parse time %v", p.Expires)
}

var updatedSecondsRegexp = regexp.MustCompile(`(\d+) seconds ago`)
var updatedMinutesRegexp = regexp.MustCompile(`(\d+) minutes ago`)
var updatedHoursRegexp = regexp.MustCompile(`(\d+) hours ago`)

func (p *Post) UpdatedAt() (time.Time, error) {
	now := time.Now()

	if p.Updated == "" {
		return now, nil
	}

	if p.Updated == "now" {
		return now, nil
	}

	if p.Updated == "a second ago" {
		return now.Add(time.Duration(-1) * time.Second), nil
	}

	if p.Updated == "a minute ago" {
		return now.Add(time.Duration(-1) * time.Minute), nil
	}

	if p.Updated == "an hour ago" {
		return now.Add(time.Duration(-1) * time.Hour), nil
	}

	match := updatedSecondsRegexp.FindStringSubmatch(p.Updated)
	if len(match) != 0 {
		seconds, err := strconv.Atoi(match[1])
		if err != nil {
			return now, fmt.Errorf("could not parse time %v: %w", p.Updated, err)
		}
		return now.Add(time.Duration(-seconds) * time.Second), nil
	}

	match = updatedMinutesRegexp.FindStringSubmatch(p.Updated)
	if len(match) != 0 {
		minutes, err := strconv.Atoi(match[1])
		if err != nil {
			return now, fmt.Errorf("could not parse time %v: %w", p.Updated, err)
		}
		return now.Add(time.Duration(-minutes) * time.Minute), nil
	}

	match = updatedHoursRegexp.FindStringSubmatch(p.Updated)
	if len(match) != 0 {
		hours, err := strconv.Atoi(match[1])
		if err != nil {
			return now, fmt.Errorf("could not parse time %v: %w", p.Updated, err)
		}
		return now.Add(time.Duration(-hours) * time.Hour), nil
	}

	return now, fmt.Errorf("failed to parse time %v", p.Updated)
}
