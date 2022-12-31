package murult

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Post struct {
	ChannelId   string
	MessageId   string
	DataCentre  string
	Duty        string
	Tags        string
	Description string
	Creator     string
	Expires     time.Time
	Updated     time.Time
	Slots       []string
}

func NewPostFromRawPost(p *RawPost, channelId, messageId string) *Post {
	return &Post{
		ChannelId:   channelId,
		MessageId:   messageId,
		DataCentre:  p.DataCentre,
		Duty:        p.Duty,
		Tags:        p.Tags,
		Description: p.Description,
		Creator:     p.Creator,
		Expires:     p.ExpiresAt(),
		Updated:     p.UpdatedAt(),
		Slots:       p.Slots,
	}
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
	result.WriteString(" Expires At: ")
	result.WriteString(fmt.Sprintf("<t:%d>", p.Expires.Unix()))
	result.WriteByte('\n')

	// Last activity
	result.WriteString(EmojiFromStr("stopwatch", emojis))
	result.WriteString(" Last updated: ")
	result.WriteString(fmt.Sprintf("<t:%d>", p.Updated.Unix()))
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

func (p *Post) GetUpdated(emojis []*discordgo.Emoji) string {
	return fmt.Sprintf("%s <t:%d>", EmojiFromStr("stopwatch", emojis), p.Updated)
}
