package data

const LANG_FR = "fr"
const LANG_EN = "en"

// All data is comming in english
type Translator interface {
	TranslateIn(lang string) string
}

type Gender string

var genderTranslationMap map[Gender]map[string]string = map[Gender]map[string]string{
	"woman": {
		LANG_FR: "femme",
	},
	"men": {
		LANG_FR: "homme",
	},
}

func (g Gender) TranslateIn(lang string) string {
	return genderTranslationMap[g][lang]
}

var categoryTranslationMap map[Category]map[string]string = map[Category]map[string]string{
	"top": {
		LANG_FR: "tops",
	},
	"dress": {
		LANG_FR: "robes",
	},
	"skirts": {
		LANG_FR: "jupes",
	},
	"jackets": {
		LANG_FR: "vestes",
	},
	"outwear": {
		LANG_FR: "manteaux",
	},
	"pants": {
		LANG_FR: "pantalons",
	},
	"shorts": {
		LANG_FR: "shorts",
	},
	"sweaters": {
		LANG_FR: "pull & gilets",
	},
	"shoes": {
		LANG_FR: "chaussures",
	},
	"beach wear": {
		LANG_FR: "beach wear",
	},
	"bags": {
		LANG_FR: "sacs",
	},
	"jewelery": {
		LANG_FR: "bijoux",
	},
	"accessories": {
		LANG_FR: "accessoires",
	},
	"beauty": {
		LANG_FR: "beauté",
	},
	"shirt": {
		LANG_FR: "chemise",
	},
	"suit": {
		LANG_FR: "costume",
	},
}

type Category string

func (cat Category) TranslateIn(lang string) string {
	return categoryTranslationMap[cat][lang]
}

type SubCategory string

var subCategoryTranslationMap map[SubCategory]map[string]string = map[SubCategory]map[string]string{
	"overcoat": {
		LANG_FR: "pardessus",
	},
	"shortsleeve": {
		LANG_FR: "manches courtes",
	},
	"t-shirt": {
		LANG_FR: "t-shirt",
	},
	"sleeveless": {
		LANG_FR: "sans manches",
	},
	"longsleeve": {
		LANG_FR: "manches longues",
	},
	"tank": {
		LANG_FR: "débardeur",
	},
	"printed": {
		LANG_FR: "imprimée",
	},
	"tunics": {
		LANG_FR: "tunique",
	},
	"polo": {
		LANG_FR: "polo",
	},
	"shirt": {
		LANG_FR: "chemise",
	},
	"cashmere": {
		LANG_FR: "cashmere",
	},
	"camisole": {
		LANG_FR: "caraco",
	},
	"halter": {
		LANG_FR: "dos nu",
	},
	"midi": {
		LANG_FR: "midi",
	},
	"maxi": {
		LANG_FR: "maxi",
	},
	"evening": {
		LANG_FR: "soirée",
	},
	"cocktail": {
		LANG_FR: "cocktail",
	},
	"bride": {
		LANG_FR: "mariée",
	},
	"bridesmaid": {
		LANG_FR: "demoiselle d'honneur",
	},
	"mini": {
		LANG_FR: "mini",
	},
	"mid-length": {
		LANG_FR: "mi-longue",
	},
	"long": {
		LANG_FR: "longue",
	},
	"denim": {
		LANG_FR: "jeans",
	},
	"leather": {
		LANG_FR: "cuir",
	},
	"blazers": {
		LANG_FR: "blazers",
	},
	"faux fur": {
		LANG_FR: "fourrure",
	},
	"suede": {
		LANG_FR: "daim",
	},
	"velvet": {
		LANG_FR: "velour",
	},
	"coat": {
		LANG_FR: "manteau de pluie",
	},
	"fur & shearling": {
		LANG_FR: "fourrure",
	},
	"leather & suede": {
		LANG_FR: "cuir",
	},
	"puffers": {
		LANG_FR: "doudoune",
	},
	"trenchcoat": {
		LANG_FR: "trench",
	},
	"wool": {
		LANG_FR: "laine",
	},
	"cropped": {
		LANG_FR: "cropped",
	},
	"legging": {
		LANG_FR: "legging",
	},
	"skinny": {
		LANG_FR: "stretch",
	},
	"wide leg": {
		LANG_FR: "ample",
	},
	"culotte": {
		LANG_FR: "jupe-culotte",
	},
	"classic": {
		LANG_FR: "fluide",
	},
	"bernuda": {
		LANG_FR: "bernuda",
	},
	"cardigan": {
		LANG_FR: "gilet",
	},
	"crewneck & scoopneck": {
		LANG_FR: "col rond",
	},
	"turtleneck": {
		LANG_FR: "col roulé",
	},
	"v-neck": {
		LANG_FR: "col V",
	},
	"hoodies": {
		LANG_FR: "sweat à capuche",
	},
	"athletic": {
		LANG_FR: "running",
	},
	"boots": {
		LANG_FR: "bottines",
	},
	"leather boots": {
		LANG_FR: "bottes en cuir",
	},
	"over the knee boots": {
		LANG_FR: "bottes hautes",
	},
	"flat": {
		LANG_FR: "plate",
	},
	"mules & clogs": {
		LANG_FR: "sandales à talons",
	},
	"platforms": {
		LANG_FR: "talons hauts",
	},
	"pumps": {
		LANG_FR: "escarpins",
	},
	"sandals": {
		LANG_FR: "sandales",
	},
	"sneakers": {
		LANG_FR: "baskets et tennis",
	},
	"wedges": {
		LANG_FR: "semelles compensées",
	},
	"one-piece": {
		LANG_FR: "monokini",
	},
	"two-piece": {
		LANG_FR: "bikini",
	},
	"cover-up": {
		LANG_FR: "paréo",
	},
	"backpack": {
		LANG_FR: "sac à dos",
	},
	"clutch": {
		LANG_FR: "pochettes",
	},
	"hobo": {
		LANG_FR: "besace",
	},
	"satchel": {
		LANG_FR: "cartable",
	},
	"shoulder": {
		LANG_FR: "sacs à main",
	},
	"wallet": {
		LANG_FR: "porte-monnaie",
	},
	"duffel & tote": {
		LANG_FR: "caba",
	},
	"bracelet": {
		LANG_FR: "bracelet",
	},
	"charm": {
		LANG_FR: "breloque",
	},
	"earrings": {
		LANG_FR: "boucles d'oreilles",
	},
	"necklace": {
		LANG_FR: "collier",
	},
	"pins": {
		LANG_FR: "pins",
	},
	"ring": {
		LANG_FR: "bague",
	},
	"diamond ring": {
		LANG_FR: "bague en diamant",
	},
	"watch": {
		LANG_FR: "montre",
	},
	"gloves": {
		LANG_FR: "gants",
	},
	"hat": {
		LANG_FR: "chapeau",
	},
	"scarve & wrap": {
		LANG_FR: "écharpe & châle",
	},
	"sunglasses": {
		LANG_FR: "lunettes de soleil",
	},
	"eyeglasses": {
		LANG_FR: "lunettes de vue",
	},
	"key chains": {
		LANG_FR: "porte clés",
	},
	"umbrella": {
		LANG_FR: "parapluie",
	},
	"headband": {
		LANG_FR: "serre-tête",
	},
	"barrette": {
		LANG_FR: "barrette",
	},
	"ear muffs": {
		LANG_FR: "cache-oreilles",
	},
	"hair accessories": {
		LANG_FR: "autres accessoires pour cheveux",
	},
	"belt": {
		LANG_FR: "ceinture",
	},
	"bath & body": {
		LANG_FR: "savons & crème",
	},
	"parfume": {
		LANG_FR: "parfum",
	},
	"haircare": {
		LANG_FR: "soin des cheveux",
	},
	"makeup": {
		LANG_FR: "maquillage",
	},
	"skincare": {
		LANG_FR: "soin de la peau",
	},
	"dress": {
		LANG_FR: "habillé",
	},
	"gingham": {
		LANG_FR: "à carreaux",
	},
	"mao collar": {
		LANG_FR: "à col mao",
	},
	"fur": {
		LANG_FR: "fourrure",
	},
	"chinos & khakis": {
		LANG_FR: "chinos",
	},
	"casual": {
		LANG_FR: "décontracté",
	},
	"cargo": {
		LANG_FR: "treillis",
	},
	"jogger": {
		LANG_FR: "jogging",
	},
	"classic  bermuda": {
		LANG_FR: "bermuda classic",
	},
	"zip up cardigan": {
		LANG_FR: "gilet à zip",
	},
	"lace-up": {
		LANG_FR: "chaussures à lacets",
	},
	"slip-ons & loafers": {
		LANG_FR: "mocassin",
	},
	"classic swim short": {
		LANG_FR: "short de bain classique",
	},
	"printed swim short": {
		LANG_FR: "short de bain imprimé",
	},
	"brief": {
		LANG_FR: "slip de bain",
	},
	"printed brief": {
		LANG_FR: "slip de bain imprimé",
	},
	"sport watch": {
		LANG_FR: "montre de sport",
	},
	"classic watch": {
		LANG_FR: "montre classique",
	},
	"silver watch": {
		LANG_FR: "montre en argent",
	},
	"gold watch": {
		LANG_FR: "montre en or",
	},
	"cuff links": {
		LANG_FR: "bouton de manchette",
	},
	"cross body bag": {
		LANG_FR: "sacs à bandoulière",
	},
	"borsalino": {
		LANG_FR: "chapeau",
	},
	"luggage": {
		LANG_FR: "valise",
	},
	"ties": {
		LANG_FR: "cravate",
	},
	"rompers": {
		LANG_FR: "barboteuse",
	},
	"jumpsuit": {
		LANG_FR: "combinaison",
	},
	"jumpsuit short": {
		LANG_FR: "combishort",
	},
	"top jumpsuit": {
		LANG_FR: "combinaison",
	},
	"short jumpsuit": {
		LANG_FR: "combinaison short",
	},
	"playsuit": {
		LANG_FR: "combishort",
	},
	"tech": {
		LANG_FR: "coque de téléphone",
	},
	"dressed": {
		LANG_FR: "habillé",
	},
}

func (subCat SubCategory) TranslateIn(lang string) string {
	return subCategoryTranslationMap[subCat][lang]
}

type Color string

var colorTranslationMap map[Color]map[string]string = map[Color]map[string]string{
	"brown": {
		LANG_FR: "marron",
	},
	"orange": {
		LANG_FR: "orange",
	},
	"yellow": {
		LANG_FR: "jaune",
	},
	"red": {
		LANG_FR: "rouge",
	},
	"purple": {
		LANG_FR: "violet",
	},
	"blue": {
		LANG_FR: "bleu",
	},
	"green": {
		LANG_FR: "vert",
	},
	"gray": {
		LANG_FR: "gris",
	},
	"white": {
		LANG_FR: "blanc",
	},
	"black": {
		LANG_FR: "noir",
	},
	"pink": {
		LANG_FR: "rose",
	},
	"gold": {
		LANG_FR: "doré",
	},
	"silver": {
		LANG_FR: "argent",
	},
	"beige": {
		LANG_FR: "beige",
	},
}

func (color Color) TranslateIn(lang string) string {
	return colorTranslationMap[color][lang]
}
