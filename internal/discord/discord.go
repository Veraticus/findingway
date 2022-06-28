package discord

import (
	"fmt"
	"time"

	"github.com/Veraticus/trappingway/internal/ffxiv"

	"github.com/bwmarrin/discordgo"
)

type Discord struct {
	Token     string
	MessageId string
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

func (d *Discord) UpdateMessage(listings *ffxiv.Listings, datacentre, duty string) error {
	fields := []*discordgo.MessageEmbedField{}
	for _, listing := range listings.ForDataCentreAndDuty(datacentre, duty) {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   listing.Creator,
			Value:  listing.PartyDisplay(),
			Inline: true,
		})
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   listing.GetTags(),
			Value:  listing.Description,
			Inline: true,
		})
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   listing.GetExpires(),
			Value:  listing.GetUpdated(),
			Inline: true,
		})
	}

	fmt.Printf("Total fields: %v\n", len(fields))
	for _, field := range fields {
		fmt.Printf("Field is: %+v\n", field)
	}

	embed := &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Color:       0x6600ff,
		Description: "Last updated at " + time.Now().Format("15:04:05"),
		Fields:      fields,
	}

	if len(d.MessageId) == 0 {
		messageSend := &discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{embed},
		}
		message, err := d.Session.ChannelMessageSendComplex(d.ChannelId, messageSend)
		if err != nil {
			return fmt.Errorf("Could not send message: %f", err)
		}
		fmt.Printf("Created new message! ID is: %v\n", message.ID)
		d.MessageId = message.ID
	} else {
		messageEdit := &discordgo.MessageEdit{
			Embeds:  []*discordgo.MessageEmbed{embed},
			ID:      d.MessageId,
			Channel: d.ChannelId,
		}
		message, err := d.Session.ChannelMessageEditComplex(messageEdit)
		if err != nil {
			return fmt.Errorf("Could not update message: %f", err)
		}
		fmt.Printf("Returned message is: %+v\n", message)
	}

	return nil
}
