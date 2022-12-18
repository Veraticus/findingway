package murult

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func EmojiFromStr(str string, emojis []*discordgo.Emoji) string {
	for _, e := range emojis {
		if e.Name == str {
			return fmt.Sprintf("<:%s>", e.APIName())
		}
	}
	return fmt.Sprintf(":question: (%s)", str)
}

func JobEmojiFromStr(str string, emojis []*discordgo.Emoji) string {
	for _, e := range emojis {
		if e.Name == str {
			return fmt.Sprintf("<:%s>", e.APIName())
		}
	}
	return fmt.Sprintf(":clown: (%s)", str)
}

type Roles struct {
	Tank   bool
	Healer bool
	Dps    bool
	Empty  bool
}

func (rs Roles) Emoji(emojis []*discordgo.Emoji) string {
	if rs.Empty {
		return EmojiFromStr("any_slot", emojis)
	} else if rs.Tank && rs.Healer && rs.Dps {
		return EmojiFromStr("tank_healer_dps_slot", emojis)
	}
	if rs.Tank && rs.Healer && !rs.Dps {
		return EmojiFromStr("tank_healer_slot", emojis)
	}
	if rs.Tank && !rs.Healer && rs.Dps {
		return EmojiFromStr("tank_dps_slot", emojis)
	}
	if rs.Tank && !rs.Healer && !rs.Dps {
		return EmojiFromStr("tank_slot", emojis)
	}
	if !rs.Tank && rs.Healer && rs.Dps {
		return EmojiFromStr("healer_dps_slot", emojis)
	}
	if !rs.Tank && rs.Healer && !rs.Dps {
		return EmojiFromStr("healer_slot", emojis)
	}
	if !rs.Tank && !rs.Healer && rs.Dps {
		return EmojiFromStr("dps_slot", emojis)
	}

	return ":question:"
}
