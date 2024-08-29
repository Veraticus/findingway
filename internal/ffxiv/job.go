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
		return "<:GBR:1277726613773094943>"
	case PLD:
		return "<:PLD:1277726628721463401>"
	case GLD:
		return "<:gld:1277675881334046761>"
	case DRK:
		return "<:DRK:1277726599860715663>"
	case WAR:
		return "<:WAR:1277726641287725088>"
	case MRD:
		return "<:mrd:1277675899629600768>"
	case SCH:
		return "<:SCH:1277726703279538227>"
	case ACN:
		return "<:acn:1277675927374925824>"
	case SGE:
		return "<:SGE:1277726715988152330>"
	case AST:
		return "<:AST:1277726729695137936>"
	case WHM:
		return "<:WHM:1277726691137163315>"
	case CNJ:
		return "<:cnj:1277675946358345788>"
	case SAM:
		return "<:SAM:1277726922419474482>"
	case DRG:
		return "<:DRG:1277726895324008532>"
	case NIN:
		return "<:NIN:1277726908234203196>"
	case MNK:
		return "<:MNK:1277726879331123221>"
	case RPR:
		return "<:RPR:1277726936457674763>"
	case VPR:
		return "<:VPR:1277726962076487680>"
	case BRD:
		return "<:BRD:1277727023715979358>"
	case MCH:
		return "<:MCH:1277727037540532385>"
	case DNC:
		return "<:DNC:1277727057131868261>"
	case BLM:
		return "<:BLM:1277727132159705200>"
	case BLU:
		return "<:BLU:1277727209003683850>"
	case SMN:
		return "<:SMN:1277727150598000742>"
	case RDM:
		return "<:RDM:1277727163516194907>"
	case PCT:
		return "<:PCT:1277727183497990288>"
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
	return "<:DOH:1278745659079524352>"
}
