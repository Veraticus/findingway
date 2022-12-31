package murult

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

var Logger = log.New(os.Stdout, "murult: ", log.Ldate|log.Ltime|log.Lshortfile)

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

type Region = string

const (
	RegionNA  Region = "NA"
	RegionEU  Region = "EU"
	RegionOCE Region = "OCE"
	RegionJP  Region = "JP"
)

func DcsFromRegion(region Region) []string {
	if region == RegionNA {
		return []string{"Aether", "Crystal", "Dynamis", "Primal"}
	} else if region == RegionEU {
		return []string{"Chaos", "Light"}
	} else if region == RegionOCE {
		return []string{"Materia"}
	} else if region == RegionJP {
		return []string{"Elemental", "Gaia", "Mana", "Meteor"}
	} else {
		return []string{}
	}
}

func CreateDiscordRegionChoices() []*discordgo.ApplicationCommandOptionChoice {
	return []*discordgo.ApplicationCommandOptionChoice{{
		Name:  "OCE",
		Value: RegionOCE,
	}, {
		Name:  "NA",
		Value: RegionNA,
	}, {
		Name:  "EU",
		Value: RegionEU,
	}, {
		Name:  "JP",
		Value: RegionJP,
	}}
}

type Duty = string

const (
	DutyUWU  Duty = "The Weapon's Refrain (Ultimate)"
	DutyTEA  Duty = "The Epic of Alexander (Ultimate)"
	DutyUCOB Duty = "The Unending Coil of Bahamut (Ultimate)"
	DutyDSR  Duty = "Dragonsong's Reprise (Ultimate)"
	DutyP5S  Duty = "Abyssos: The Fifth Circle (Savage)"
	DutyP6S  Duty = "Abyssos: The Sixth Circle (Savage)"
	DutyP7S  Duty = "Abyssos: The Seventh Circle (Savage)"
	DutyP8S  Duty = "Abyssos: The Eight Circle (Savage)"
)

func CreateDiscordDutyChoices() []*discordgo.ApplicationCommandOptionChoice {
	return []*discordgo.ApplicationCommandOptionChoice{{
		Name:  "UWU",
		Value: DutyUWU,
	}, {
		Name:  "TEA",
		Value: DutyTEA,
	}, {
		Name:  "UCOB",
		Value: DutyUCOB,
	}, {
		Name:  "DSR",
		Value: DutyDSR,
	}, {
		Name:  "P5S",
		Value: DutyP5S,
	}, {
		Name:  "P6S",
		Value: DutyP6S,
	}, {
		Name:  "P7S",
		Value: DutyP7S,
	}, {
		Name:  "P8S",
		Value: DutyP8S,
	}}
}
