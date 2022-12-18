package murult

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Server struct {
	Token     string
	World     string
	GuildId   string
	ChannelId string
	Session   *discordgo.Session
	Duties    []string
}

func (s *Server) CleanChannel() error {
	messages, err := s.Session.ChannelMessages(s.ChannelId, 100, "", "", "")
	if err != nil {
		return fmt.Errorf("could not list messages: %f", err)
	}
	for _, message := range messages {
		err := s.Session.ChannelMessageDelete(s.ChannelId, message.ID)
		if err != nil {
			return fmt.Errorf("could not delete message %+v: %f", message, err)
		}
	}

	return nil
}

func (s *Server) PostListings(listings *Listings, emojis []*discordgo.Emoji) error {
	scopedListings := listings.ForUltimatesInMateria(s.Duties, s.World)

	for _, listing := range scopedListings.Listings {
		field := &discordgo.MessageEmbedField{
			Name:   listing.Duty,
			Value:  listing.PartyDisplay(emojis),
			Inline: true,
		}
		embed := &discordgo.MessageEmbed{
			Type:   discordgo.EmbedTypeRich,
			Color:  0x6600ff,
			Fields: []*discordgo.MessageEmbedField{field},
			Footer: &discordgo.MessageEmbedFooter{
				Text: strings.Repeat("\u3000", 20),
			},
		}
		messageSend := &discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{embed},
		}

		_, err := s.Session.ChannelMessageSendComplex(s.ChannelId, messageSend)

		if err != nil {
			return err
		}
	}

	return nil
}
