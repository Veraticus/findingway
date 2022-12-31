package murult

type Role struct {
	Tank   bool
	Healer bool
	Dps    bool
	Empty  bool
}

func (rs *Role) Emoji() string {
	if rs.Empty || (rs.Tank && rs.Healer && rs.Dps) {
		return RoleEmojiTankHealerDps
	}
	if rs.Tank && rs.Healer && !rs.Dps {
		return RoleEmojiTankHealer
	}
	if rs.Tank && !rs.Healer && rs.Dps {
		return RoleEmojiTankDps
	}
	if rs.Tank && !rs.Healer && !rs.Dps {
		return RoleEmojiTank
	}
	if !rs.Tank && rs.Healer && rs.Dps {
		return RoleEmojiHealerDps
	}
	if !rs.Tank && rs.Healer && !rs.Dps {
		return RoleEmojiHealer
	}
	if !rs.Tank && !rs.Healer && rs.Dps {
		return RoleEmojiDps
	}

	return ":question:"
}
