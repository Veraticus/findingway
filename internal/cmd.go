package murult

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var permission int64 = discordgo.PermissionManageServer
var gmtOffsetMin float64 = -12
var gmtOffsetMax float64 = 12

var Commands []*discordgo.ApplicationCommand = []*discordgo.ApplicationCommand{{
	Name:        "add-duty",
	Description: "Add the specified duty to the watch list",
	Type:        discordgo.ChatApplicationCommand,
	Options: []*discordgo.ApplicationCommandOption{{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "duty-name",
		Description: "Name of duty to add to the list",
		Required:    true,
		Choices:     CreateDiscordDutyChoices(),
	}},
	DefaultMemberPermissions: &permission,
}, {
	Name:        "remove-duty",
	Description: "Remove the specified duty from the watch list",
	Type:        discordgo.ChatApplicationCommand,
	Options: []*discordgo.ApplicationCommandOption{{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "duty-name",
		Description: "Name of duty to remove from the list",
		Required:    true,
		Choices:     CreateDiscordDutyChoices(),
	}},
	DefaultMemberPermissions: &permission,
}, {
	Name:        "add-region",
	Description: "Add the specified region",
	Type:        discordgo.ChatApplicationCommand,
	Options: []*discordgo.ApplicationCommandOption{{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "region-name",
		Description: "Name of region to track",
		Required:    true,
		Choices:     CreateDiscordRegionChoices(),
	}},
	DefaultMemberPermissions: &permission,
}, {
	Name:        "remove-region",
	Description: "Remove the specified region",
	Type:        discordgo.ChatApplicationCommand,
	Options: []*discordgo.ApplicationCommandOption{{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "region-name",
		Description: "Name of region to track",
		Required:    true,
		Choices:     CreateDiscordRegionChoices(),
	}},
	DefaultMemberPermissions: &permission,
}, {
	Name:                     "info",
	Description:              "Show relevant information about this channel",
	Type:                     discordgo.ChatApplicationCommand,
	DefaultMemberPermissions: &permission,
}, {
	Name:                     "update-emojis",
	Description:              "Update the emoji database for this channel",
	Type:                     discordgo.ChatApplicationCommand,
	DefaultMemberPermissions: &permission,
}, {
	Name:                     "schedule-duty",
	Description:              "Generate timestamp",
	Type:                     discordgo.ChatApplicationCommand,
	DefaultMemberPermissions: &permission,
	Options: []*discordgo.ApplicationCommandOption{{
		Type:        discordgo.ApplicationCommandOptionInteger,
		Name:        "year",
		Description: "Year",
		Required:    true,
	}, {
		Type:        discordgo.ApplicationCommandOptionInteger,
		Name:        "month",
		Description: "Month",
		Required:    true,
	}, {
		Type:        discordgo.ApplicationCommandOptionInteger,
		Name:        "day",
		Description: "Day",
		Required:    true,
	}, {
		Type:        discordgo.ApplicationCommandOptionInteger,
		Name:        "hour",
		Description: "Hour",
		Required:    true,
	}, {
		Type:        discordgo.ApplicationCommandOptionInteger,
		Name:        "minute",
		Description: "Minute",
		Required:    true,
	}, {
		Type:        discordgo.ApplicationCommandOptionInteger,
		Name:        "offset",
		Description: "Timezone offset from UTC",
		Required:    true,
		MinValue:    &gmtOffsetMin,
		MaxValue:    gmtOffsetMax,
	}, {
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "region-name",
		Description: "Name of region",
		Required:    true,
		Choices:     CreateDiscordRegionChoices(),
	}, {
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "duty-name",
		Description: "Name of duty",
		Required:    true,
		Choices:     CreateDiscordDutyChoices(),
	}},
}}

type CommandHandler = func(
	server *Server,
	d *discordgo.Session,
	i *discordgo.InteractionCreate)

var CommandHandlers map[string]CommandHandler = map[string]CommandHandler{
	"add-duty": func(
		s *Server,
		d *discordgo.Session,
		i *discordgo.InteractionCreate) {
		duty, exists := getStringValue("duty-name", i.ApplicationCommandData())

		if !exists {
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Bad argument",
				},
			})
		} else {
			s.AddDuty(i.GuildID, i.ChannelID, duty)
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Added duty `%s` to this channel", duty),
				},
			})
		}
	},
	"remove-duty": func(
		s *Server,
		d *discordgo.Session,
		i *discordgo.InteractionCreate) {
		duty, exists := getStringValue("duty-name", i.ApplicationCommandData())

		if !exists {
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Bad argument",
				},
			})
		} else {
			s.RemoveDuty(i.ChannelID, duty)
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Removed duty `%s` from this channel", duty),
				},
			})
		}
	},
	"add-region": func(
		s *Server,
		d *discordgo.Session,
		i *discordgo.InteractionCreate) {
		region, exists := getStringValue("region-name", i.ApplicationCommandData())

		if !exists {
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Bad argument",
				},
			})
		} else {
			s.AddRegion(i.GuildID, i.ChannelID, region)
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Added region `%s` to this channel", region),
				},
			})
		}
	},
	"remove-region": func(
		s *Server,
		d *discordgo.Session,
		i *discordgo.InteractionCreate) {
		region, exists := getStringValue("region-name", i.ApplicationCommandData())

		if !exists {
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Bad argument",
				},
			})
		} else {
			s.RemoveRegion(i.ChannelID, region)
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Removed region `%s` from this channel", region),
				},
			})
		}
	},
	"info": func(
		s *Server,
		d *discordgo.Session,
		i *discordgo.InteractionCreate) {
		channel, exists := s.channels[i.ChannelID]

		if !exists {
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "This channel is not registered for this bot",
				},
			})
			return
		}

		duties := channel.Duties()

		if len(duties) == 0 {
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "This channel is not tracking any duties",
				},
			})
			return
		}

		regions := channel.Regions()

		if len(duties) == 0 {
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "This channel is not tracking any regions",
				},
			})
			return
		}

		var respond strings.Builder
		respond.WriteString("We are currently tracking the following duties:\n")

		for i, d := range duties {
			respond.WriteString(fmt.Sprintf("%d. `%s`\n", i+1, d))
		}

		respond.WriteString("We are currently tracking the following regions:\n")

		for i, r := range regions {
			respond.WriteString(fmt.Sprintf("%d. `%s`\n", i+1, r))
		}

		d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: respond.String(),
			},
		})
	},
	"update-emojis": func(
		s *Server,
		d *discordgo.Session,
		i *discordgo.InteractionCreate) {
		s.UpdateEmojis(i.GuildID)
		d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Updated the emoji database for this guild",
			},
		})
	},
	"schedule-duty": func(
		s *Server,
		d *discordgo.Session,
		i *discordgo.InteractionCreate) {
		year, yearE := getIntValue("year", i.ApplicationCommandData())
		month, monthE := getIntValue("month", i.ApplicationCommandData())
		day, dayE := getIntValue("day", i.ApplicationCommandData())
		hour, hourE := getIntValue("hour", i.ApplicationCommandData())
		minute, minuteE := getIntValue("minute", i.ApplicationCommandData())
		offset, offsetE := getIntValue("offset", i.ApplicationCommandData())
		duty, dutyE := getStringValue("duty-name", i.ApplicationCommandData())
		region, regionE := getStringValue("region-name", i.ApplicationCommandData())

		if !yearE || !monthE || !dayE || !hourE || !minuteE || !offsetE || !dutyE || !regionE {
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Missing arguments",
				},
			})
		} else {
			date := time.Date(year, time.Month(month), day, hour-offset, minute, 0, 0, time.UTC)
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("%s plans on creating a PF for %s in %s at <t:%d> your time.\n Please react with your roles if you are interested.", i.Member.Mention(), duty, region, date.Unix()),
				},
			})
		}
	}}

func getStringValue(key string, options discordgo.ApplicationCommandInteractionData) (string, bool) {
	for _, opt := range options.Options {
		if opt.Name == key && opt.Type == discordgo.ApplicationCommandOptionString {
			return opt.StringValue(), true
		}
	}
	return "", false
}

func getIntValue(key string, options discordgo.ApplicationCommandInteractionData) (int, bool) {
	for _, opt := range options.Options {
		if opt.Name == key && opt.Type == discordgo.ApplicationCommandOptionInteger {
			return int(opt.IntValue()), true
		}
	}
	return 0, false
}
