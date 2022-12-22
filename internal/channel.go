package murult

import (
	"github.com/bwmarrin/discordgo"
)

type Channel struct {
	duties map[string]struct{}
	emojis []*discordgo.Emoji
}

func NewChannel() *Channel {
	return &Channel{
		duties: map[string]struct{}{
			"The Weapon's Refrain (Ultimate)":         {},
			"The Unending Coil of Bahamut (Ultimate)": {},
			"The Epic of Alexander (Ultimate)":        {},
			"Dragonsong's Reprise (Ultimate)":         {}},
		emojis: make([]*discordgo.Emoji, 0),
	}
}

func (c *Channel) AddDuty(duty string) bool {
	_, exists := c.duties[duty]
	c.duties[duty] = struct{}{}
	return !exists
}

func (c *Channel) RemoveDuty(duty string) bool {
	_, exists := c.duties[duty]
	delete(c.duties, duty)
	return exists
}

func (s *Channel) Duties() []string {
	result := make([]string, 0, len(s.duties))

	for k := range s.duties {
		result = append(result, k)
	}

	return result
}

func (c *Channel) UpdateEmojis(emojis []*discordgo.Emoji) {
	c.emojis = emojis
}

func (c *Channel) Emojis() []*discordgo.Emoji {
	return c.emojis
}
