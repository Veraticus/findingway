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
		return "<:gunbreaker:985322473337782384>"
	case PLD:
		return "<:paladin:985322479318892584>"
	case GLD:
		return "<:gladiator:985322472079491152>"
	case DRK:
		return "<:darkknight:985322469873303624>"
	case WAR:
		return "<:warrior:985322493143318578>"
	case MRD:
		return "<:marauder:985322476986826782>"
	case SCH:
		return "<:scholar:985322486231089212>"
	case ACN:
		return "<:arcanist:985322461866369094>"
	case SGE:
		return "<:sage:985322483823566908>"
	case AST:
		return "<:astrologian:985322464127107093>"
	case WHM:
		return "<:whitemage:985322493919244328>"
	case CNJ:
		return "<:conjurer:985322468308811886>"
	case SAM:
		return "<:samurai:985322484842758235>"
	case DRG:
		return "<:dragoon:985322471232245860>"
	case NIN:
		return "<:ffxivninja:985322478521966612>"
	case MNK:
		return "<:monk:985322477683089418>"
	case RPR:
		return "<:reaper:985322481025966150>"
	case BRD:
		return "<:bard:985322465733533736>"
	case MCH:
		return "<:machinist:985322476244443246>"
	case DNC:
		return "<:ffxivdancer:985322469172850728>"
	case BLM:
		return "<:blackmage:985322466723377202>"
	case BLU:
		return "<:bluemage:985322467599974421>"
	case SMN:
		return "<:summoner:985322487191584839>"
	case RDM:
		return "<:redmage:985322481889996890>"
	case LNC:
		return "<:lancer:985322475225219084>"
	case PUG:
		return "<:pugilist:985322480203862056>"
	case ROG:
		return "<:rogue:985322482879848458>"
	case THM:
		return "<:thaumaturge:985322492258295818>"
	case ARC:
		return "<:archer:985322463552495616>"
	}
	return "❓"
}
