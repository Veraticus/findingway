package murult

import (
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
	pfState  *PfState
}

func NewServer(token, path string) *Server {
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

	db := NewDb(path)

	db.CreateChannelsTable()
	db.CreateRegionsTable()
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

	server.RegisterCommands()

	session.AddHandler(func(d *discordgo.Session, i *discordgo.ChannelDelete) {
		Logger.Println("received channel deletion event")
		server.db.RemoveChannel(i.GuildID, i.ID)
		server.lock.Lock()
		defer server.lock.Unlock()
		delete(server.channels, i.ID)
	})
	session.AddHandler(func(d *discordgo.Session, i *discordgo.InteractionCreate) {
		Logger.Printf("received interaction event of type '%s'\n", i.Type.String())
		if i.Type == discordgo.InteractionApplicationCommand {
			if cmd, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
				cmd(server, d, i)
			}
		}
	})

	return server
}

func (s *Server) CloseServer() {
	s.lock.Lock()
	defer s.lock.Unlock()

	err := s.db.Close()

	if err != nil {
		Logger.Printf("Unable to close SQLITE db connection because '%s'\n", err)
	}

	err = s.session.Close()

	if err != nil {
		Logger.Printf("Unable to close Discord session connection because '%s'\n", err)
	}
}

func (s *Server) Run(sleep int64) {
	for {
		s.lock.Lock()
		pfState, err := s.scraper.Scrape()
		s.pfState = pfState
		s.lock.Unlock()

		if err != nil {
			Logger.Printf("Unable to scrape website because '%s'\n", err)
			Logger.Printf("Sleeping for %d minutes\n", sleep)
			time.Sleep(time.Duration(sleep * int64(time.Minute)))
			continue
		}

		s.sendUpdates()

		Logger.Printf("Sleeping for %d minutes\n", sleep)
		time.Sleep(time.Duration(sleep * int64(time.Minute)))
	}
}

func (s *Server) RegisterCommands() bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	cmds, err := s.session.ApplicationCommandBulkOverwrite(s.session.State.User.ID, "", Commands)

	if err != nil {
		Logger.Printf("Cannot create bulk commands because %s\n", err)
		return false
	}

	for _, cmd := range cmds {
		Logger.Printf("Created command `%s`\n", cmd.Name)
	}

	return true
}

func (s *Server) AddRegion(guildId, channelId string, region Region) {
	s.lock.Lock()
	defer s.lock.Unlock()

	channel, exists := s.channels[channelId]

	if exists {
		channel.regions[region] = struct{}{}
		s.db.InsertRegion(channelId, region)
	} else {
		s.channels[channelId] = NewChannel(
			guildId,
			channelId,
			map[Region]struct{}{
				region: {},
			},
			map[string]struct{}{},
			map[string]*Post{})

		s.db.InsertChannel(guildId, channelId)
		s.db.InsertRegion(channelId, region)
	}
}

func (s *Server) RemoveRegion(channelId string, region Region) {
	s.lock.Lock()
	defer s.lock.Unlock()

	channel, exists := s.channels[channelId]

	if exists {
		delete(channel.regions, region)

		s.db.RemoveRegion(channelId, region)

		if len(channel.regions) == 0 {
			s.db.RemoveChannel(channel.guildId, channelId)
		}
	}
}

func (s *Server) AddDuty(guildId, channelId, duty string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	channel, exists := s.channels[channelId]

	if exists {
		channel.duties[duty] = struct{}{}

		s.db.InsertDuty(channelId, duty)
	} else {
		s.channels[channelId] = NewChannel(
			guildId,
			channelId,
			map[Region]struct{}{},
			map[string]struct{}{
				duty: {},
			},
			map[string]*Post{})

		s.db.InsertChannel(guildId, channelId)
		s.db.InsertDuty(channelId, duty)
	}
}

func (s *Server) RemoveDuty(channelId, duty string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	channel, exists := s.channels[channelId]

	if exists {
		delete(channel.duties, duty)

		s.db.RemoveDuty(channelId, duty)

		if len(channel.duties) == 0 {
			s.db.RemoveChannel(channel.guildId, channelId)
		}
	}
}

func (s *Server) Duties(channelId string) []string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	channel, exists := s.channels[channelId]
	if exists {
		duties := make([]string, 0, len(channel.duties))

		for d := range channel.duties {
			duties = append(duties, d)
		}

		return duties
	} else {
		return []string{}
	}
}

func (s *Server) Regions(channelId string) []string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	channel, exists := s.channels[channelId]
	if exists {
		regions := make([]string, 0, len(channel.regions))

		for d := range channel.regions {
			regions = append(regions, d)
		}

		return regions
	} else {
		return []string{}
	}
}

func (s *Server) UpdateEmojis(guildId string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.updateEmojis(guildId)
}

func (s *Server) updateEmojis(guildId string) {
	emojis, err := s.session.GuildEmojis(guildId)

	if err != nil {
		Logger.Printf("Unable to update emojis for '%s'\n", guildId)
		return
	}

	s.emojis[guildId] = emojis
}

func (s *Server) Emojis(guildId string) []*discordgo.Emoji {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.getEmojis(guildId)
}

func (s *Server) getEmojis(guildId string) []*discordgo.Emoji {
	_, exists := s.emojis[guildId]

	if !exists {
		s.updateEmojis(guildId)
	}

	return s.emojis[guildId]
}

func (s *Server) sendUpdates() {
	for _, channel := range s.channels {
		s.lock.Lock()
		sleepTime := s.sendUpdateToChannel(channel)
		s.lock.Unlock()
		time.Sleep(time.Second * time.Duration(sleepTime) / 2)
	}
}

func (s *Server) sendUpdateToChannel(channel *Channel) int {
	count := 0
	Logger.Printf("starting to send updates to channelId=%s", channel.channelId)
	removedPosts, updatedPosts, newPosts := channel.UpdatePosts(s.pfState)
	channel.posts = make(map[string]*Post, len(updatedPosts)+len(newPosts))

	for _, p := range removedPosts {
		err := s.session.ChannelMessageDelete(p.ChannelId, p.MessageId)
		count++

		if err != nil {
			Logger.Printf("Discord error cleaning message '%s' in channel '%s' because '%s'\n", p.MessageId, channel.channelId, err)
		}

		s.db.RemovePost(p)
	}

	for _, p := range updatedPosts {
		_, err := s.session.ChannelMessageEdit(channel.channelId, p.MessageId, p.Stringify(s.getEmojis(channel.guildId)))
		count++

		if err != nil {
			Logger.Printf("Discord error updating message '%s' in channel '%s' because '%s'\n", p.MessageId, channel.channelId, err)
			s.db.RemovePost(p)
			continue
		}

		channel.posts[p.Creator] = p
		s.db.InsertPost(p)
	}

	for _, rp := range newPosts {
		message, err := s.session.ChannelMessageSendComplex(channel.channelId, &discordgo.MessageSend{
			Content: rp.Stringify(s.getEmojis(channel.guildId)),
		})
		count++

		if err != nil {
			Logger.Printf("Discord error creating message in channel '%s' because '%s'\n", channel.channelId, err)
			continue
		}

		p := NewPostFromRawPost(rp, channel.channelId, message.ID)
		channel.posts[p.Creator] = p
		s.db.InsertPost(p)
	}

	Logger.Printf("finished sending updates to channelId=%s", channel.channelId)
	return count
}
