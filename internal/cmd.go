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
	Name:                     "register-channel",
	Description:              "Tells the bot to post to this channel",
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
		options := i.ApplicationCommandData().Options
		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

		for _, opt := range options {
			optionMap[opt.Name] = opt
		}

		maybeDuty, ok := optionMap["duty-name"]

		if !ok {
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Missing required field 'duty-name'",
				},
			})
			return
		}

		if maybeDuty.Type != discordgo.ApplicationCommandOptionString {
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("'duty-name' should be a string but we received a value of type '%s'\n", maybeDuty.Type.String()),
				},
			})
			return
		}

		channel, exists := s.Channel(i.ChannelID)

		if !exists {
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "This channel is not registered for this bot",
				},
			})
			return
		}

		dutyName := maybeDuty.StringValue()

		if channel.AddDuty(dutyName) {
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Added `%s` to the list of tracked duties for this channel", dutyName),
				},
			})
		} else {
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("We are already tacking a duty called `%s` for this channel", dutyName),
				},
			})
		}
	},
	"remove-duty": func(
		s *Server,
		d *discordgo.Session,
		i *discordgo.InteractionCreate) {
		options := i.ApplicationCommandData().Options
		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

		for _, opt := range options {
			optionMap[opt.Name] = opt
		}

		maybeDuty, ok := optionMap["duty-name"]

		if !ok {
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Missing required field 'duty-name'",
				},
			})
			return
		}

		if maybeDuty.Type != discordgo.ApplicationCommandOptionString {
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("'duty-name' should be a string but we received a value of type '%s'\n", maybeDuty.Type.String()),
				},
			})
			return
		}

		channel, exists := s.Channel(i.ChannelID)

		if !exists {
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "This channel is not registered for this bot",
				},
			})
			return
		}

		dutyName := maybeDuty.StringValue()

		if channel.RemoveDuty(dutyName) {
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Removed `%s` from the list of tracked duties for this channel", dutyName),
				},
			})
		} else {
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("This channel is not tracking a duty called `%s` for this channel", dutyName),
				},
			})
		}
	},
	"list-duties": func(
		s *Server,
		d *discordgo.Session,
		i *discordgo.InteractionCreate) {
		channel, exists := s.Channel(i.ChannelID)

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
	"register-channel": func(
		s *Server,
		d *discordgo.Session,
		i *discordgo.InteractionCreate) {
		emojis, err := s.Session.GuildEmojis(i.GuildID)

		if err != nil {
			emojis = make([]*discordgo.Emoji, 0)
		}

		if s.AddChannel(i.ChannelID, emojis) {
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Registration successful. You should use `/add-duty <duty-name>` command to add duties you want this channel to track",
				},
			})
		} else {
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "This channel is already registered",
				},
			})
		}

	},
	"remove-channel": func(
		s *Server,
		d *discordgo.Session,
		i *discordgo.InteractionCreate) {

		if s.RemoveChannel(i.ChannelID) {
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Removal successful. We will no long post listings to this channel",
				},
			})
		} else {
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "This channel is not even registered yet",
				},
			})
		}
	},
	"update-emojis": func(
		s *Server,
		d *discordgo.Session,
		i *discordgo.InteractionCreate) {
		emojis, err := s.Session.GuildEmojis(i.GuildID)

		if err != nil {
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Failed to update emojis database because '%s'", err.Error()),
				},
			})
			return
		}

		channel, exists := s.Channel(i.ChannelID)

		if !exists {
			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "This channel is not registered for this bot",
				},
			})
			return
		}

		channel.UpdateEmojis(emojis)

		d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Updated the emoji database for this guild",
			},
		})
	}}
