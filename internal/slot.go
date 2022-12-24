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
