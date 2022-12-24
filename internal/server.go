package murult

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	lock     sync.RWMutex
	token    string
	session  *discordgo.Session
	scraper  *Scraper
	channels map[string]*Channel
	emojis   map[string][]*discordgo.Emoji
	db       *Db
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

	db := NewDb()

	db.CreateChannelsTable()
	db.CreateDutiesTable()
	db.CreatePostsTable()

	channels, ok := db.SelectAllChannels()

	if !ok {
		return nil
	}

	server := &Server{
		token:    token,
		session:  session,
		scraper:  scraper,
		channels: channels,
		emojis:   make(map[string][]*discordgo.Emoji),
		db:       db,
	}

	server.clearCommands()
	server.registerCommands()

	return server
}

func (s *Server) CloseServer() {
	s.db.Close()
	s.session.Close()
}

func (s *Server) Run(sleep int64) {
	for {
		s.lock.Lock()
		pfState, err := s.scraper.Scrape()

		if err != nil {
			Logger.Printf("scraper error: '%s'\n", err)
			time.Sleep(time.Duration(sleep * int64(time.Minute)))
			continue
		}

		for cid, channel := range s.channels {
			removedPosts, updatedPosts, newPosts := channel.UpdatePosts(pfState)

			for _, r := range removedPosts {
				if r.MessageId != "" {
					err := s.session.ChannelMessageDelete(cid, r.MessageId)

					if err != nil {
						Logger.Printf("Discord error cleaning channel: %f\n", err)
					}
				}
			}

			for _, p := range updatedPosts {
				if p.MessageId != "" {
					message, err := s.session.ChannelMessageEdit(cid, p.MessageId, p.Stringify(s.Emojis(channel.guildId)))

					if err != nil {
						Logger.Printf("Discord error updating message: %f\n", err)
						continue
					}

					p.MessageId = message.ID
				}
			}

			for _, p := range newPosts {
				message, err := s.session.ChannelMessageSendComplex(cid, &discordgo.MessageSend{
					Content: p.Stringify(s.Emojis(channel.guildId)),
				})

				if err != nil {
					Logger.Printf("Discord error creating message: %f\n", err)
					continue
				}

				p.MessageId = message.ID
				s.db.InsertPost(cid, message.ID, p.Creator)
			}

			Logger.Printf("updated channel `%s`\n", cid)
		}

		s.lock.Unlock()
		time.Sleep(time.Duration(sleep * int64(time.Minute)))
	}
}

func (s *Server) clearCommands() {
	registeredCommands, err := s.session.ApplicationCommands(s.session.State.User.ID, "")

	if err != nil {
		log.Fatalf("Could not fetch registered commands: %v", err)
	}

	for _, v := range registeredCommands {
		err := s.session.ApplicationCommandDelete(s.session.State.User.ID, "", v.ID)

		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}

		fmt.Printf("Deleted command `%s`\n", v.Name)
	}
}

func (s *Server) registerCommands() {
	s.session.AddHandler(func(d *discordgo.Session, i *discordgo.InteractionCreate) {
		if cmd, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
			cmd(s, d, i)
		}
	})

	for _, v := range Commands {
		cmd, err := s.session.ApplicationCommandCreate(s.session.State.User.ID, "", v)

		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}

		fmt.Printf("Created command `%s`\n", cmd.Name)
	}
}

func (s *Server) AddDuty(guildId, channelId, duty string) {
	channel, exists := s.channels[channelId]

	if exists {
		channel.duties[duty] = struct{}{}

		s.db.InsertDuty(channelId, duty)
	} else {
		s.channels[channelId] = NewChannel(
			guildId,
			map[string]struct{}{
				duty: {},
			},
			map[string]*Post{})

		s.db.InsertChannel(guildId, channelId)
		s.db.InsertDuty(channelId, duty)
	}
}

func (s *Server) RemoveDuty(channelId, duty string) {
	channel, exists := s.channels[channelId]

	if exists {
		delete(channel.duties, duty)

		s.db.RemoveDuty(channelId, duty)

		if len(channel.duties) == 0 {
			s.db.RemoveChannel(channel.guildId, channelId)
		}
	}
}

func (s *Server) Duties(channelId string) map[string]struct{} {
	channel, exists := s.channels[channelId]
	if exists {
		return channel.duties
	} else {
		return map[string]struct{}{}
	}
}

func (c *Server) UpdateEmojis(guildId string) {
	emojis, err := c.session.GuildEmojis(guildId)

	if err != nil {
		Logger.Printf("Unable to update emojis for '%s'\n", guildId)
	}

	c.emojis[guildId] = emojis
}

func (c *Server) Emojis(guildId string) []*discordgo.Emoji {
	_, exists := c.emojis[guildId]

	if !exists {
		c.UpdateEmojis(guildId)
	}

	return c.emojis[guildId]
}
