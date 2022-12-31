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

func (c *Channel) UpdatePosts(pf *PfState) (map[string]*Post, map[string]*Post, map[string]*RawPost) {
	currentPosts := pf.GetPosts(c.Duties(), c.Regions())
	removedPosts := make(map[string]*Post, 0)
	updatedPosts := make(map[string]*Post, 0)
	newPosts := make(map[string]*RawPost, 0)

	for creator, newPost := range currentPosts {
		oldPost, exists := c.posts[creator]

		if !exists {
			newPosts[creator] = newPost
		} else {
			updatedPosts[creator] = NewPostFromRawPost(newPost, c.channelId, oldPost.MessageId)
		}
	}

	for creator, oldPost := range c.posts {
		_, exists := currentPosts[creator]
		if !exists {
			removedPosts[creator] = oldPost
		}
	}

	return removedPosts, updatedPosts, newPosts
}
