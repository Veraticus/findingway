package murult

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var permission int64 = discordgo.PermissionManageServer

var Commands []*discordgo.ApplicationCommand = []*discordgo.ApplicationCommand{{
	Name:        "add-duty",
	Description: "Add the specified duty to the watch list",
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
	DefaultMemberPermissions: &permission,
}, {
	Name:                     "update-emojis",
	Description:              "Update the emoji database for this channel",
	DefaultMemberPermissions: &permission,
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
		duty, exists := getDutyName(i.ApplicationCommandData())

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
		duty, exists := getDutyName(i.ApplicationCommandData())

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
		region, exists := getRegionName(i.ApplicationCommandData())

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
		region, exists := getRegionName(i.ApplicationCommandData())

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
	}}

func getDutyName(options discordgo.ApplicationCommandInteractionData) (string, bool) {
	for _, opt := range options.Options {
		if opt.Name == "duty-name" && opt.Type == discordgo.ApplicationCommandOptionString {
			return opt.StringValue(), true
		}
	}
	return "", false
}

func getRegionName(options discordgo.ApplicationCommandInteractionData) (Region, bool) {
	for _, opt := range options.Options {
		if opt.Name == "region-name" && opt.Type == discordgo.ApplicationCommandOptionString {
			return opt.StringValue(), true
		}
	}
	return "", false
}
