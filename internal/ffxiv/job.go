package ffxiv

type Job int

const (
	GNB Job = iota
	PLD
	GLD
	DRK
	WAR
	MRD
	SCH
	ACN // Arcanist
	SGE
	AST
	WHM
	CNJ
	SAM
	DRG
	NIN
	MNK
	RPR
	VPR
	BRD
	MCH
	DNC
	BLM
	BLU
	SMN
	PCT
	RDM
	LNC
	PUG
	ROG
	THM
	ARC // Archer
	Unknown
)

func JobFromAbbreviation(abbreviation string) Job {
	switch abbreviation {
	case "GNB":
		return GNB
	case "PLD":
		return PLD
	case "GLD":
		return GLD
	case "DRK":
		return DRK
	case "WAR":
		return WAR
	case "MRD":
		return MRD
	case "SCH":
		return SCH
	case "ACN":
		return ACN
	case "SGE":
		return SGE
	case "AST":
		return AST
	case "WHM":
		return WHM
	case "CNJ":
		return CNJ
	case "SAM":
		return SAM
	case "DRG":
		return DRG
	case "NIN":
		return NIN
	case "MNK":
		return MNK
	case "RPR":
		return RPR
	case "VPR":
		return VPR
	case "BRD":
		return BRD
	case "MCH":
		return MCH
	case "DNC":
		return DNC
	case "BLM":
		return BLM
	case "BLU":
		return BLU
	case "SMN":
		return SMN
	case "PCT":
		return PCT
	case "RDM":
		return RDM
	case "LNC":
		return LNC
	case "PUG":
		return PUG
	case "ROG":
		return ROG
	case "THM":
		return THM
	case "ARC":
		return ARC
	}
	return Unknown
}

func (j Job) Emoji() string {
	switch j {
	case GNB:
		return "<:Gunbreaker:1240636120249794620> "
	case PLD:
		return "<:Paladin:1240636121650827274>"
	case GLD:
		return "<:gld:1277675881334046761>"
	case DRK:
		return "<:DarkKnight:1240636123965816862>"
	case WAR:
		return "<:Warrior:1240636122711986257>"
	case MRD:
		return "<:mrd:1277675899629600768>"
	case SCH:
		return "<:Scholar:1240636166542331944>"
	case ACN:
		return "<:acn:1277675927374925824>"
	case SGE:
		return "<:Sage:1240636165057548338>"
	case AST:
		return "<:Astrologian:1240636168685617154>"
	case WHM:
		return "<:WhiteMage:1240636167511080960>"
	case CNJ:
		return "<:cnj:1277675946358345788>"
	case SAM:
		return "<:Samurai:1240636206237089852>"
	case DRG:
		return "<:Dragoon:1240636213216411689>"
	case NIN:
		return "<:Ninja:1240636202193784922>"
	case MNK:
		return "<:Monk:1240636200285372498>"
	case RPR:
		return "<:Reaper:1240636203527573545>"
	case VPR:
		return "<:Viper:1243593551254651000>"
	case BRD:
		return "<:Bard:1240636208758128691>"
	case MCH:
		return "<:Machinist:1240636198960234526>"
	case DNC:
		return "<:Dancer:1240636251858665542>"
	case BLM:
		return "<:BlackMage:1240636210133602407>"
	case BLU:
		return "<:blu:1277665256730001520>"
	case SMN:
		return "<:Summoner:1240636207717810247>"
	case RDM:
		return "<:RedMage:1240636205025067038>"
	case PCT:
		return "<:Pictomancer:1243593616920936598>"
	case LNC:
		return "<:lnc:1277675968701661217>"
	case PUG:
		return "<:pgl:1277675984862052392>"
	case ROG:
		return "<:rog:1277676011386962064>"
	case THM:
		return "<:thm:1277676028113719296>"
	case ARC:
		return "<:arc:1277676046992408636>"
	}
	return "<:CUL:991374283836227604>"
}
