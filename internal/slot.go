package murult

type Slot struct {
	Roles  Role
	Job    string
	Filled bool
}

func NewSlot() *Slot {
	return &Slot{
		Roles: Role{
			Tank:   false,
			Healer: false,
			Dps:    false,
			Empty:  false,
		},
	}
}

func (rs *Slot) Emoji() string {
	if rs.Filled {
		return rs.Job
	} else {
		return rs.Roles.Emoji()
	}
}
