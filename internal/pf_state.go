package murult

import "time"

type PfState struct {
	posts map[string]*Post
}

func NewPfState() *PfState {
	return &PfState{
		posts: make(map[string]*Post),
	}
}

// GetPosts returns an array of posts that is in the PF.
// Can be filtered based on the arguments.
func (pf *PfState) GetPosts(duties []string, regions []Region) map[string]*Post {
	list := make(map[string]*Post, 0)

	for _, post := range pf.posts {
		for _, region := range regions {
			worlds := WorldsFromRegion(region)
			for _, world := range worlds {
				for _, duty := range duties {
					if post.Duty == duty && post.DataCentre == world {
						oldPost, exists := list[post.Creator]

						if exists {
							oldTime, err := oldPost.ExpiresAt()

							if err != nil {
								oldTime = time.Time{}
							}

							newTime, err := post.ExpiresAt()

							if err != nil {
								oldTime = time.Time{}
							}

							if newTime.After(oldTime) {
								list[post.Creator] = post
							} else {
								list[post.Creator] = oldPost
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

func (pf *PfState) Add(l *Post) {
	// TODO: Check if we already have this creator
	// If we do, check which one is the latest one.
	pf.posts[l.Creator] = l
}
