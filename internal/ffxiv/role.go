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
		return "<:dps:1277660881538383967>"
	}
	if reflect.DeepEqual(rs.Roles, []Role{Healer}) {
		return "<:healer:1277660862966009966>"
	}
	if reflect.DeepEqual(rs.Roles, []Role{Tank}) {
		return "<:tank:1277660836382376028>"
	}
	if reflect.DeepEqual(rs.Roles, []Role{DPS, Healer}) {
		return "<:healerdps:1277660987981565982>"
	}
	if reflect.DeepEqual(rs.Roles, []Role{DPS, Tank}) {
		return "<:tankdps:1277660932935516200>"
	}
	if reflect.DeepEqual(rs.Roles, []Role{Healer, Tank}) {
		return "<:tankhealer:1277661018184749249>"
	}

	if reflect.DeepEqual(rs.Roles, []Role{Healer, Tank, DPS}) {
		return "<:tankhealerdps:1277661062661017620>"
	}

	return "<:DOH:1278745659079524352>"
}
