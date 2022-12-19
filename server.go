package murult

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type Server struct {
	Token     string
	GuildId   string
	ChannelId string
	// TODO: Add ability for admins to trigger updates on emojis
	Emojis  []*discordgo.Emoji
	Session *discordgo.Session
	Scraper *Scraper
	Duties  []string
}

func NewServer(token, guildId, channelId string) *Server {
	session, err := discordgo.New("Bot " + token)

	if err != nil {
		log.Printf("could not start Discord: '%s'\n", err)
		return nil
	}

	err = session.Open()

	if err != nil {
		log.Printf("could not open Discord session: '%s'\n", err)
		return nil
	}

	emojis, err := session.GuildEmojis(guildId)

	if err != nil {
		log.Printf("could not get server emojis: %s\n", err)
	}

	scraper := NewScraper("https://xivpf.com/listings")

	return &Server{
		Token:     token,
		GuildId:   guildId,
		ChannelId: channelId,
		Emojis:    emojis,
		Session:   session,
		Scraper:   scraper,
		Duties: []string{
			"The Weapon's Refrain (Ultimate)",
			"The Unending Coil of Bahamut (Ultimate)",
			"The Epic of Alexander (Ultimate)",
			"Dragonsong's Reprise (Ultimate)"},
	}
}

func (s *Server) CloseServer() {
	err := s.Session.Close()

	if err != nil {
		log.Printf("unable to close Discord websocket session: '%s'\n", err)
	}
}

func (s *Server) Run() {
	err := s.Scraper.Scrape()

	if err != nil {
		Logger.Printf("scraper error: '%s'\n", err)
		return
	}

	err = s.CleanChannel()

	if err != nil {
		Logger.Printf("Discord error cleaning channel: %f\n", err)
		return
	}

	err = s.PostListings(s.Scraper.Listings, s.Emojis)

	if err != nil {
		Logger.Printf("Discord error updating message: %f\n", err)
		return
	}

	log.Println("updated listing")
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

func (s *Server) PostListings(pf *PfState, emojis []*discordgo.Emoji) error {
	pf.FilterForUltimatesInMateria(s.Duties)

	for _, listing := range pf.Listings {
		messageSend := &discordgo.MessageSend{
			Content: listing.PartyDisplay(emojis),
		}

		_, err := s.Session.ChannelMessageSendComplex(s.ChannelId, messageSend)

		if err != nil {
			return err
		}
	}

	return nil
}
