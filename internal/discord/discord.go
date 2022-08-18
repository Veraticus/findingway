package discord

import (
	"fmt"
	"strings"
	"time"

	"github.com/Veraticus/trappingway/internal/ffxiv"

	"github.com/bwmarrin/discordgo"
)

type Discord struct {
	Token string

	Session  *discordgo.Session
	Channels []*Channel `yaml:"channels"`
}

type Channel struct {
	ID          string   `yaml:"id"`
	Duty        string   `yaml:"duty"`
	DataCentres []string `yaml:"dataCentres"`
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

func (d *Discord) CleanChannel(channelId string) error {
	messages, err := d.Session.ChannelMessages(channelId, 100, "", "", "")
	if err != nil {
		return fmt.Errorf("Could not list messages: %f", err)
	}
	for _, message := range messages {
		err := d.Session.ChannelMessageDelete(channelId, message.ID)
		if err != nil {
			return fmt.Errorf("Could not delete message %+v: %f", message, err)
		}
	}

	return nil
}

func (d *Discord) PostListings(channelId string, listings *ffxiv.Listings, duty string, dataCentres []string) error {
	scopedListings := listings.ForDutyAndDataCentres(duty, dataCentres)

	headerEmbed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s PFs", duty),
		Type:        discordgo.EmbedTypeRich,
		Color:       0x6600ff,
		Description: fmt.Sprintf("Found %v listings %v", len(scopedListings.Listings), fmt.Sprintf("<t:%v:R>", time.Now().Unix())),
		Footer: &discordgo.MessageEmbedFooter{
			Text: strings.Repeat("\u3000", 20),
		},
	}
	headerMessageSend := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{headerEmbed},
	}
	_, err := d.Session.ChannelMessageSendComplex(channelId, headerMessageSend)
	if err != nil {
		return fmt.Errorf("Could not send header: %f", err)
	}

	fields := []*discordgo.MessageEmbedField{}
	for i, listing := range scopedListings.Listings {
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
		if (i+1)%5 == 0 {
			err = d.sendMessage(channelId, fields)
			if err != nil {
				return fmt.Errorf("Could not send message: %f", err)
			}
			fields = []*discordgo.MessageEmbedField{}
		}
	}

	// Ensure we send any remaining messages
	if len(fields) != 0 {
		err = d.sendMessage(channelId, fields)
		if err != nil {
			return fmt.Errorf("Could not send message: %f", err)
		}
	}

	return nil
}

func (d *Discord) sendMessage(channelId string, fields []*discordgo.MessageEmbedField) error {
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
	_, err := d.Session.ChannelMessageSendComplex(channelId, messageSend)
	if err != nil {
		return err
	}

	return nil
}
