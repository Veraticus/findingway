package murult

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var permission int64 = discordgo.PermissionManageServer

var Commands []*discordgo.ApplicationCommand = []*discordgo.ApplicationCommand{{
	Name:        "add-duty",
	Description: "Add the specified duty to the watch list (Must be the exact and not the abbreviations of it)",
	Options: []*discordgo.ApplicationCommandOption{{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "duty-name",
		Description: "Name of duty to add to the list",
		Required:    true,
	}},
	DefaultMemberPermissions: &permission,
}, {
	Name:        "remove-duty",
	Description: "Remove the specified duty from the watch list (Must be the exact and not the abbreviations of it)",
	Options: []*discordgo.ApplicationCommandOption{{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "duty-name",
		Description: "Name of duty to remove from the list",
		Required:    true,
	}},
	DefaultMemberPermissions: &permission,
}, {
	Name:                     "list-duties",
	Description:              "List all the duties that we are currently watching",
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
	// Duty management
	"add-duty": func(
		s *Server,
		d *discordgo.Session,
		i *discordgo.InteractionCreate) {
		dutyName, exists := getDutyName(i.ApplicationCommandData())

		if !exists {
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Bad argument",
				},
			})
		} else {
			s.AddDuty(i.GuildID, i.ChannelID, dutyName)
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Added `%s` to this channel", dutyName),
				},
			})
		}
	},
	"remove-duty": func(
		s *Server,
		d *discordgo.Session,
		i *discordgo.InteractionCreate) {
		dutyName, exists := getDutyName(i.ApplicationCommandData())

		if !exists {
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Bad argument",
				},
			})
		} else {
			s.RemoveDuty(i.ChannelID, dutyName)
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Removed `%s` from this channel", dutyName),
				},
			})
		}
	},
	"list-duties": func(
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

		var respond strings.Builder
		respond.WriteString("We are currently tracking the following duties:\n")

		for i, d := range duties {
			respond.WriteString(fmt.Sprintf("%d. `%s`\n", i+1, d))
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
