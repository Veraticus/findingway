package murult

type Slot struct {
	Roles  Roles
	Job    string
	Filled bool
}

func NewSlot() *Slot {
	return &Slot{
		Roles: Roles{
			Tank:   false,
			Healer: false,
			Dps:    false,
			Empty:  false,
		},
	}
}
