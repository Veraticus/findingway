package murult

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Server struct {
	lock    sync.RWMutex
	Token   string
	GuildId string
	// TODO: Add ability for admins to trigger updates on emojis
	Emojis   map[string]*discordgo.Emoji
	Session  *discordgo.Session
	Scraper  *Scraper
	channels map[string]*Channel
}

func NewServer(token string) *Server {
	session, err := discordgo.New("Bot " + token)

	if err != nil {
		Logger.Printf("could not start Discord: '%s'\n", err)
		return nil
	}

	err = session.Open()

	if err != nil {
		Logger.Printf("could not open Discord session: '%s'\n", err)
		return nil
	}

	scraper := NewScraper("https://xivpf.com/listings")

	if scraper == nil {
		Logger.Printf("unable to initialize scraper")
		return nil
	}

	server := &Server{
		Token:    token,
		Emojis:   make(map[string]*discordgo.Emoji),
		Session:  session,
		Scraper:  scraper,
		channels: make(map[string]*Channel),
	}

	server.clearCommands()
	server.registerCommands()

	return server
}

func (s *Server) CloseServer() {
	err := s.Session.Close()

	if err != nil {
		Logger.Printf("unable to close Discord websocket session: '%s'\n", err)
	}
}

func (s *Server) Run(sleep int64) {
	for {
		s.lock.Lock()
		err := s.Scraper.Scrape()

		if err != nil {
			Logger.Printf("scraper error: '%s'\n", err)
			return
		}

		for cid := range s.channels {
			err = s.CleanChannel(cid)

			if err != nil {
				Logger.Printf("Discord error cleaning channel: %f\n", err)
				return
			}

			err = s.PostListingsToChannel(cid)

			if err != nil {
				Logger.Printf("Discord error updating message: %f\n", err)
				return
			}

			Logger.Printf("updated listing for channel %s\n", cid)

		}

		s.lock.Unlock()
		time.Sleep(time.Duration(sleep * int64(time.Minute)))
	}
}

// CleanChannel cleans "all" the messages in a given channel.
// The current implementation is quite bad, it cannot be relied on to delete all the
// messages nor it is efficient in doing so.
func (s *Server) CleanChannel(channelId string) error {
	messages, err := s.Session.ChannelMessages(channelId, 100, "", "", "")

	if err != nil {
		return fmt.Errorf("could not list messages: %f", err)
	}

	for _, message := range messages {
		err := s.Session.ChannelMessageDelete(channelId, message.ID)
		if err != nil {
			return fmt.Errorf("could not delete message %+v: %f", message, err)
		}
	}

	return nil
}

func (s *Server) PostListingsToChannel(channelId string) error {
	channel, exists := s.channels[channelId]

	if !exists {
		// TODO: Should we delete the channel here?
		Logger.Println("Attempting to get listings for a channel that we are not tracking?")
		return nil
	}

	duties := channel.Duties()

	listings := s.Scraper.Listings.GetListings(duties)

	for _, l := range listings {
		message := &discordgo.MessageSend{
			Content: l.Stringify(channel.Emojis()),
		}

		_, err := s.Session.ChannelMessageSendComplex(channelId, message)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) clearCommands() {
	registeredCommands, err := s.Session.ApplicationCommands(s.Session.State.User.ID, s.GuildId)

	if err != nil {
		log.Fatalf("Could not fetch registered commands: %v", err)
	}

	for _, v := range registeredCommands {
		err := s.Session.ApplicationCommandDelete(s.Session.State.User.ID, s.GuildId, v.ID)

		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}

		fmt.Printf("Deleted command '%s'\n", v.Name)
	}
}

func (s *Server) registerCommands() {
	s.Session.AddHandler(func(d *discordgo.Session, i *discordgo.InteractionCreate) {
		if cmd, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
			cmd(s, d, i)
		}
	})

	for _, v := range Commands {
		cmd, err := s.Session.ApplicationCommandCreate(s.Session.State.User.ID, s.GuildId, v)

		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}

		fmt.Printf("Created command '%s'\n", cmd.Name)
	}
}

func (s *Server) Channel(channelId string) (*Channel, bool) {
	v, e := s.channels[channelId]
	return v, e
}

func (s *Server) Channels() []string {
	result := make([]string, 0, len(s.channels))
	for k := range s.channels {
		result = append(result, k)
	}
	return result
}

func (s *Server) AddDuty(channelId, duty string) bool {
	config, exists := s.channels[channelId]
	if !exists {
		return false
	} else {
		return config.AddDuty(duty)
	}
}

func (s *Server) RemoveDuty(channelId, duty string) bool {
	config, exists := s.channels[channelId]
	if !exists {
		return false
	} else {
		return config.RemoveDuty(duty)
	}
}

func (s *Server) AddChannel(channelId string, emojis []*discordgo.Emoji) bool {
	lenBefore := len(s.channels)
	channel := NewChannel()
	channel.UpdateEmojis(emojis)
	s.channels[channelId] = channel
	return len(s.channels) != lenBefore
}

func (s *Server) RemoveChannel(channelId string) bool {
	_, exists := s.channels[channelId]
	delete(s.channels, channelId)
	return exists
}
