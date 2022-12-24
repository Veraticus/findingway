package murult

type Channel struct {
	guildId string
	duties  map[string]struct{}
	posts   map[string]*Post
}

func NewChannel(
	guildId string,
	duties map[string]struct{},
	posts map[string]*Post) *Channel {
	return &Channel{
		guildId: guildId,
		duties:  duties,
		posts:   posts,
	}
}

func (s *Channel) Duties() []string {
	result := make([]string, 0, len(s.duties))

	for k := range s.duties {
		result = append(result, k)
	}

	return result
}

func (c *Channel) UpdatePosts(pf *PfState) (map[string]*Post, map[string]*Post, map[string]*Post) {
	currentPosts := pf.GetPosts(c.Duties())
	removedPosts := make(map[string]*Post, 0)
	updatedPosts := make(map[string]*Post, 0)
	newPosts := make(map[string]*Post, 0)

	for creator, newPost := range currentPosts {
		oldPost, exists := c.posts[creator]

		if !exists {
			newPosts[creator] = newPost
		} else {
			updatedPosts[creator] = newPost
			newPost.MessageId = oldPost.MessageId
		}
	}

	for creator, oldPost := range c.posts {
		_, exists := currentPosts[creator]
		if !exists {
			removedPosts[creator] = oldPost
		}
	}

	c.posts = currentPosts

	return removedPosts, updatedPosts, newPosts
}
