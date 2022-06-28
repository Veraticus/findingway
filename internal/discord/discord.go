package discord

import (
	"fmt"
	"strings"
	"time"

	"github.com/Veraticus/trappingway/internal/ffxiv"

	"github.com/bwmarrin/discordgo"
)

type Discord struct {
	Token     string
	ChannelId string

	Session *discordgo.Session
}

func (d *Discord) Start() error {
	s, err := discordgo.New("Bot " + d.Token)
	if err != nil {
		return fmt.Errorf("Could not start Discord: %f", err)
	}

	err = s.Open()
	if err != nil {
		return fmt.Errorf("Could not open Discord session: %f", err)
	}

	d.Session = s
	return nil
}

func (d *Discord) CleanChannel() error {
	messages, err := d.Session.ChannelMessages(d.ChannelId, 100, "", "", "")
	if err != nil {
		return fmt.Errorf("Could not list messages: %f", err)
	}
	for _, message := range messages {
		err := d.Session.ChannelMessageDelete(d.ChannelId, message.ID)
		if err != nil {
			return fmt.Errorf("Could not delete message %+v: %f", message, err)
		}
	}

	return nil
}

func (d *Discord) PostListings(listings *ffxiv.Listings, datacentre, duty string) error {
	scopedListings := listings.ForDataCentreAndDuty(datacentre, duty)

	headerEmbed := &discordgo.MessageEmbed{
		Title:       "Dragonsong's Reprise (Ultimate) PFs",
		Type:        discordgo.EmbedTypeRich,
		Color:       0x6600ff,
		Description: fmt.Sprintf("Found %v listings %v", len(scopedListings), fmt.Sprintf("<t:%v:R>", time.Now().Unix())),
		Footer: &discordgo.MessageEmbedFooter{
			Text: strings.Repeat("\u3000", 20),
		},
	}
	headerMessageSend := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{headerEmbed},
	}
	_, err := d.Session.ChannelMessageSendComplex(d.ChannelId, headerMessageSend)
	if err != nil {
		return fmt.Errorf("Could not send header: %f", err)
	}

	fields := []*discordgo.MessageEmbedField{}
	for i, listing := range scopedListings {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   listing.Creator,
			Value:  listing.PartyDisplay(),
			Inline: true,
		})
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   listing.GetTags(),
			Value:  listing.GetDescription(),
			Inline: true,
		})
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   listing.GetExpires(),
			Value:  listing.GetUpdated(),
			Inline: true,
		})

		// Send a message every 5 listings
		if i%5 == 0 && i != 0 {
			err = d.sendMessage(fields)
			if err != nil {
				return fmt.Errorf("Could not send message: %f", err)
			}
			fields = []*discordgo.MessageEmbedField{}
		}
	}

	// Ensure we send any remaining messages
	if len(fields) != 0 {
		err = d.sendMessage(fields)
		if err != nil {
			return fmt.Errorf("Could not send message: %f", err)
		}
	}

	return nil
}

func (d *Discord) sendMessage(fields []*discordgo.MessageEmbedField) error {
	embed := &discordgo.MessageEmbed{
		Type:   discordgo.EmbedTypeRich,
		Color:  0x6600ff,
		Fields: fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: strings.Repeat("\u3000", 20),
		},
	}
	messageSend := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{embed},
	}
	_, err := d.Session.ChannelMessageSendComplex(d.ChannelId, messageSend)
	if err != nil {
		return err
	}

	return nil
}
