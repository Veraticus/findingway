package ffxiv

import (
	"reflect"
)

type Role int

const (
	DPS Role = iota
	Healer
	Tank
	Empty
)

type Roles struct {
	Roles []Role
}

func (rs Roles) Emoji() string {
	if reflect.DeepEqual(rs.Roles, []Role{DPS}) {
		return "<:dps_slot:991374282460508252>"
	}
	if reflect.DeepEqual(rs.Roles, []Role{Healer}) {
		return "<:healer_slot:991374281479041034>"
	}
	if reflect.DeepEqual(rs.Roles, []Role{Tank}) {
		return "<:tank_slot:991374280380121180>"
	}
	if reflect.DeepEqual(rs.Roles, []Role{DPS, Healer}) {
		return "<:dps_healer_slot:991374279121850418>"
	}
	if reflect.DeepEqual(rs.Roles, []Role{DPS, Tank}) {
		return "<:dps_tank_slot:991374278060691566>"
	}
	if reflect.DeepEqual(rs.Roles, []Role{Healer, Tank}) {
		return "<:healer_tank_slot:991374276844327042>"
	}

	if reflect.DeepEqual(rs.Roles, []Role{Healer, Tank, DPS}) {
		return "<:tankhealerdps:985322491398459482>"
	}

	return "<:any_slot:991374273975435384>"
}
