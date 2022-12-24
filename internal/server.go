package murult

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Server struct {
	lock     sync.RWMutex
	Token    string
	Session  *discordgo.Session
	Scraper  *Scraper
	channels map[string]*Channel
	emojis   map[string][]*discordgo.Emoji
}

func InitServer(token string) *Server {
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
		Session:  session,
		Scraper:  scraper,
		channels: make(map[string]*Channel),
		emojis:   make(map[string][]*discordgo.Emoji),
	}

	server.clearCommands()
	server.registerCommands()

	return server
}

func (s *Server) Run(sleep int64) {
	for {
		s.lock.Lock()
		pfState, err := s.Scraper.Scrape()

		if err != nil {
			Logger.Printf("scraper error: '%s'\n", err)
			time.Sleep(time.Duration(sleep * int64(time.Minute)))
			continue
		}

		for cid, channel := range s.channels {
			removedPosts, updatedPosts, newPosts := channel.UpdatePosts(pfState)

			for _, r := range removedPosts {
				if r.MessageId != "" {
					err := s.Session.ChannelMessageDelete(cid, r.MessageId)

					if err != nil {
						Logger.Printf("Discord error cleaning channel: %f\n", err)
					}
				}
			}

			for _, p := range updatedPosts {
				message, err := s.Session.ChannelMessageEdit(cid, p.MessageId, p.Stringify(s.Emojis(channel.guildId)))

				if err != nil {
					Logger.Printf("Discord error updating message: %f\n", err)
					continue
				}

				p.MessageId = message.ID
			}

			for _, p := range newPosts {
				message, err := s.Session.ChannelMessageSendComplex(cid, &discordgo.MessageSend{
					Content: p.Stringify(s.Emojis(channel.guildId)),
				})

				if err != nil {
					Logger.Printf("Discord error creating message: %f\n", err)
					continue
				}

				p.MessageId = message.ID
			}

			Logger.Printf("updated channel `%s`\n", cid)
		}

		s.lock.Unlock()
		time.Sleep(time.Duration(sleep * int64(time.Minute)))
	}
}

func (s *Server) clearCommands() {
	registeredCommands, err := s.Session.ApplicationCommands(s.Session.State.User.ID, "")

	if err != nil {
		log.Fatalf("Could not fetch registered commands: %v", err)
	}

	for _, v := range registeredCommands {
		err := s.Session.ApplicationCommandDelete(s.Session.State.User.ID, "", v.ID)

		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}

		fmt.Printf("Deleted command `%s`\n", v.Name)
	}
}

func (s *Server) registerCommands() {
	s.Session.AddHandler(func(d *discordgo.Session, i *discordgo.InteractionCreate) {
		if cmd, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
			cmd(s, d, i)
		}
	})

	for _, v := range Commands {
		cmd, err := s.Session.ApplicationCommandCreate(s.Session.State.User.ID, "", v)

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
	} else {
		s.channels[channelId] = NewChannel(
			guildId,
			map[string]struct{}{
				duty: {},
			},
			map[string]*Post{})
	}
}

func (s *Server) RemoveDuty(channelId, duty string) {
	channel, exists := s.channels[channelId]
	if exists {
		delete(channel.duties, duty)
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
	emojis, err := c.Session.GuildEmojis(guildId)

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
