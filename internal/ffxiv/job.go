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
	BRD
	MCH
	DNC
	BLM
	BLU
	SMN
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
		return "<:GNB:991374319190016080>"
	case PLD:
		return "<:PLD:991374318175015012>"
	case GLD:
		return "<:GLD:991374316904136775>"
	case DRK:
		return "<:DRK:991374315536793660>"
	case WAR:
		return "<:WAR:991374314454659103>"
	case MRD:
		return "<:MRD:991374313108291624>"
	case SCH:
		return "<:SCH:1053569722609840139>"
	case ACN:
		return "<:ACN:991374307328544879>"
	case SGE:
		return "<:SGE:991374305839554640>"
	case AST:
		return "<:AST:991374326823661659>"
	case WHM:
		return "<:WHM:991374304753221642>"
	case CNJ:
		return "<:CNJ:991374303138427000>"
	case SAM:
		return "<:SAM:991374325317906534>"
	case DRG:
		return "<:DRG:991374301636866048>"
	case NIN:
		return "<:NIN:991374324374188053>"
	case MNK:
		return "<:MNK:991374300370190346>"
	case RPR:
		return "<:RPR:991374299350978620>"
	case BRD:
		return "<:BRD:991374298168168528>"
	case MCH:
		return "<:MCH:991374296813412442>"
	case DNC:
		return "<:DNC:991374295475437588>"
	case BLM:
		return "<:BLM:991374293978071081>"
	case BLU:
		return "<:BLU:991374292929495151>"
	case SMN:
		return "<:SMN:991374323245924382>"
	case RDM:
		return "<:RDM:991374291776065576>"
	case LNC:
		return "<:LNC:991374290714898493>"
	case PUG:
		return "<:PUG:991374289032986666>"
	case ROG:
		return "<:ROG:991374288215081040>"
	case THM:
		return "<:THM:991374286604488744>"
	case ARC:
		return "<:ARC:991374285161631804>"
	}
	return "<:CUL:991374283836227604>"
}
