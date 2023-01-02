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
	taskChan chan func()
	die      bool
	dieChan  chan bool
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
		die:      false,
		dieChan:  make(chan bool),
		taskChan: make(chan func()),
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
	Logger.Println("Received signal to die. Cancelling all future jobs")

	s.lock.Lock()
	s.die = true
	s.lock.Unlock()

	<-s.dieChan
	<-s.dieChan

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

func (s *Server) StartScrapeJob(sleep int64) {
	for {
		Logger.Println("Starting scraping job")
		pfState, err := s.scraper.Scrape()
		s.lock.Lock()
		s.pfState = pfState
		s.lock.Unlock()

		if err != nil {
			Logger.Printf("Unable to scrape website because '%s'\n", err)
			Logger.Printf("Sleeping for %d minutes\n", sleep)
			time.Sleep(time.Duration(sleep * int64(time.Minute)))
			continue
		}

		s.lock.Lock()
		jobs := make([]func(*Server), 0)
		for channelId := range s.channels {
			channel := s.channels[channelId]
			removedPosts, updatedPosts, newPosts := channel.UpdatePosts(s.pfState)

			for creator := range removedPosts {
				p := removedPosts[creator]
				jobs = append(jobs, RemovePostHandler(p))
				Logger.Printf("Created job to remove duty for '%s' by '%s' in '%s'  in messageId='%s' in channelId='%s'", p.Duty, p.Creator, p.DataCentre, p.MessageId, p.ChannelId)
			}

			for creator := range updatedPosts {
				p := updatedPosts[creator]
				jobs = append(jobs, UpdatePostHandler(p))
				Logger.Printf("Created job to update duty for '%s' by '%s' in '%s'  in messageId=%s in channelId=%s", p.Duty, p.Creator, p.DataCentre, p.MessageId, p.ChannelId)
			}

			for creator := range newPosts {
				rp := newPosts[creator]
				jobs = append(jobs, CreatePostHandler(channelId, rp))
				Logger.Printf("Created job to create duty for '%s' by '%s' in '%s' in channelId=%s", rp.Duty, rp.Creator, rp.DataCentre, channel.channelId)
			}

			Logger.Printf("Created %d jobs for channelId=%s", len(jobs), channel.channelId)
		}
		s.lock.Unlock()

		for _, f := range jobs {
			s.lock.Lock()
			f(s)
			s.lock.Unlock()
		}

		Logger.Printf("Sleeping for %d minutes\n", sleep)
		time.Sleep(time.Duration(sleep * int64(time.Minute)))

		s.lock.Lock()
		if s.die {
			s.lock.Unlock()
			s.dieChan <- true
			return
		}
		s.lock.Unlock()
	}
}

func (s *Server) StartUpdateJob(sleep int64) {
	for task := range s.taskChan {
		s.lock.Lock()
		task()
		if s.die {
			s.lock.Unlock()
			s.dieChan <- true
			return
		}
		s.lock.Unlock()
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
			delete(s.channels, channelId)
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
			delete(s.channels, channelId)
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

func RemovePostHandler(p *Post) func(*Server) {
	return func(s *Server) {
		Logger.Printf("Executing job to remove duty for '%s' by '%s' in '%s'  in messageId='%s' in channelId='%s'\n", p.Duty, p.Creator, p.DataCentre, p.MessageId, p.ChannelId)

		channel, exists := s.channels[p.ChannelId]
		if !exists {
			Logger.Printf("Attempting to delete a post from channelId='%s' but the channel is already deleted\n", p.ChannelId)
			return
		}

		_, exists = channel.posts[p.Creator]

		if !exists {
			Logger.Printf("Attempting to delete a post from channelId='%s' but the post is already deleted\n", p.ChannelId)
			return
		}

		if exists {
			err := s.session.ChannelMessageDelete(p.ChannelId, p.MessageId)

			if err != nil {
				Logger.Printf("Discord error cleaning message '%s' in channelId='%s' because '%s'\n", p.MessageId, channel.channelId, err)
			}

			delete(channel.posts, p.Creator)
			s.db.RemovePost(p)
			Logger.Println("Job done")
		}
	}
}

func UpdatePostHandler(p *Post) func(*Server) {
	return func(s *Server) {
		Logger.Printf("Executing job to update duty for '%s' by '%s' in '%s'  in messageId='%s' in channelId='%s'", p.Duty, p.Creator, p.DataCentre, p.MessageId, p.ChannelId)

		channel, exists := s.channels[p.ChannelId]
		if !exists {
			Logger.Printf("Attempting to update a post from channelId='%s' but the channel is already deleted\n", p.ChannelId)
			return
		}

		_, exists = channel.posts[p.Creator]

		if !exists {
			Logger.Printf("Attempting to update a post from channelId='%s' but the post is already deleted\n", p.ChannelId)
			return
		}

		_, err := s.session.ChannelMessageEdit(channel.channelId, p.MessageId, p.Stringify(s.getEmojis(channel.guildId)))

		if err != nil {
			Logger.Printf("Discord error updating messageId='%s' in channelId='%s' because '%s'\n", p.MessageId, channel.channelId, err)
			s.db.RemovePost(p)
			return
		}

		channel.posts[p.Creator] = p
		s.db.InsertPost(p)
		Logger.Println("Job done")
	}
}

func CreatePostHandler(channelId string, rp *RawPost) func(*Server) {
	return func(s *Server) {
		Logger.Printf("Executing job to create duty for '%s' by '%s' in '%s'  in channelId=%s", rp.Duty, rp.Creator, rp.DataCentre, channelId)

		channel, exists := s.channels[channelId]

		if !exists {
			Logger.Printf("Attempting to create a post from channelId='%s' but the channel is already deleted\n", channelId)
			return
		}

		p, exists := channel.posts[rp.Creator]

		if exists {
			Logger.Printf("Attempting to create a post from channelId='%s' but creator='%s' already have a post. So we will update the old post instead\n", channel.channelId, rp.Creator)

			_, err := s.session.ChannelMessageEdit(channel.channelId, p.MessageId, p.Stringify(s.getEmojis(channel.guildId)))

			if err != nil {
				Logger.Printf("Discord error updating messageId='%s' in channelId='%s' because '%s'\n", p.MessageId, channel.channelId, err)
				s.db.RemovePost(p)
				return
			}

			channel.posts[p.Creator] = p
			s.db.InsertPost(p)
		} else {
			message, err := s.session.ChannelMessageSendComplex(channel.channelId, &discordgo.MessageSend{
				Content: rp.Stringify(s.getEmojis(channel.guildId)),
			})

			if err != nil {
				Logger.Printf("Discord error creating message in channelId='%s' because '%s'\n", channel.channelId, err)
				s.db.RemovePost(p)
				return
			}

			p := NewPostFromRawPost(rp, channel.channelId, message.ID)
			channel.posts[p.Creator] = p
			s.db.InsertPost(p)
		}
		Logger.Println("Job done")
	}
}
