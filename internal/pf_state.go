package murult

type PfState struct {
	posts []*RawPost
}

func NewPfState(posts []*RawPost) *PfState {
	return &PfState{
		posts: posts,
	}
}

// GetPosts returns an array of posts that is in the PF.
// Can be filtered based on the arguments.
// TODO: Improve this so that we don't have to loop the entire posts for every channel
func (pf *PfState) GetPosts(duties []string, regions []string) map[string]*RawPost {
	list := make(map[string]*RawPost, 0)

	for _, post := range pf.posts {
		for _, region := range regions {
			for _, dc := range DcsFromRegion(region) {
				for _, duty := range duties {
					if post.Duty == duty && post.DataCentre == dc {
						dupPost, exists := list[post.Creator]

						if exists {
							dupTime := dupPost.ExpiresAt()
							newTime := post.ExpiresAt()

							if newTime.After(dupTime) {
								list[post.Creator] = post
							} else {
								list[post.Creator] = dupPost
							}
						} else {
							list[post.Creator] = post
						}
					}
				}
			}
		}
	}

	return list
}
