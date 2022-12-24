package murult

type Channel struct {
	guildId   string
	channelId string
	regions   map[Region]struct{}
	duties    map[string]struct{}
	posts     map[string]*Post
}

func NewChannel(
	guildId string,
	channelId string,
	regions map[Region]struct{},
	duties map[string]struct{},
	posts map[string]*Post) *Channel {
	return &Channel{
		guildId:   guildId,
		channelId: channelId,
		regions:   regions,
		duties:    duties,
		posts:     posts,
	}
}

func (s *Channel) Duties() []string {
	result := make([]string, 0, len(s.duties))

	for k := range s.duties {
		result = append(result, k)
	}

	return result
}

func (s *Channel) Regions() []Region {
	result := make([]Region, 0, len(s.regions))

	for k := range s.regions {
		result = append(result, k)
	}

	return result
}

func (c *Channel) UpdatePosts(pf *PfState) map[string][]*Post {
	c.posts = pf.GetPosts(c.Duties(), c.Regions())
	result := map[string][]*Post{}

	for _, p := range c.posts {
		l, exists := result[p.Duty]

		if !exists {
			l = []*Post{}
		}

		result[p.Duty] = append(l, p)
	}

	return result
}
