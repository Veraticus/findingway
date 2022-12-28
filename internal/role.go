package murult

import "github.com/bwmarrin/discordgo"

type Role struct {
	Tank   bool
	Healer bool
	Dps    bool
	Empty  bool
}

func (rs Role) Emoji(emojis []*discordgo.Emoji) string {
	if rs.Empty || (rs.Tank && rs.Healer && rs.Dps) {
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
