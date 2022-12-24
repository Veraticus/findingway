package murult

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
func (pf *PfState) GetPosts(duties []string) map[string]*Post {
	list := make(map[string]*Post, 0)

	for _, l := range pf.posts {
		if l.DataCentre == "Materia" {
			for _, d := range duties {
				if l.Duty == d {
					list[l.Creator] = l
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
