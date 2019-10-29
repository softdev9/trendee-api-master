package data

import (
	"testing"
)

func TestTranslateGenderFr(t *testing.T) {
	// Test if gender implement translator
	gFr := Gender("woman").TranslateIn(LANG_FR)
	if gFr != "femme" {
		t.Errorf("Expeced translation was femme but go \"%s\"", gFr)
	}
	gFr = Gender("men").TranslateIn(LANG_FR)
	if gFr != "homme" {
		t.Errorf("Expeced translation was homme but go \"%s\"", gFr)
	}
}

func TestTranslateColorFr(t *testing.T) {
	tests := []struct {
		Color Color
		Fr    string
	}{
		{
			Color("brown"),
			"marron",
		},
		{
			Color("orange"),
			"orange",
		},
		{
			Color("yellow"),
			"jaune",
		},
		{
			Color("red"),
			"rouge",
		},
		{
			Color("purple"),
			"violet",
		},
		{
			Color("blue"),
			"bleu",
		},
		{
			Color("green"),
			"vert",
		},
		{
			Color("gray"),
			"gris",
		},
		{
			Color("white"),
			"blanc",
		},
		{
			Color("black"),
			"noir",
		},
		{
			Color("pink"),
			"rose",
		},
		{
			Color("gold"),
			"doré",
		},
		{
			Color("silver"),
			"argent",
		},
		{
			Color("beige"),
			"beige",
		},
	}
	for _, test := range tests {
		result := test.Color.TranslateIn(LANG_FR)
		if result != test.Fr {
			t.Errorf(" Expected %s for color %s but got %v ", test.Fr, test.Color, result)
		}
	}
}

func TestTranslateCaterogyFr(t *testing.T) {
	tests := []struct {
		Cat Category
		Fr  string
	}{
		{
			Cat: Category("top"),
			Fr:  "tops",
		},
		{
			Cat: Category("dress"),
			Fr:  "robes",
		},
		{
			Cat: Category("skrits"),
			Fr:  "jupes",
		},
		{
			Cat: Category("jackets"),
			Fr:  "vestes",
		},
		{
			Cat: Category("outwear"),
			Fr:  "manteaux",
		},
		{
			Cat: Category("pants"),
			Fr:  "pantalons",
		},
		{
			Cat: Category("shorts"),
			Fr:  "shorts",
		},
		{
			Cat: Category("sweaters"),
			Fr:  "pull & gilets",
		},
		{
			Cat: Category("shoes"),
			Fr:  "chaussures",
		},
		{
			Cat: Category("beach wear"),
			Fr:  "beach wear",
		},
		{
			Cat: Category("bags"),
			Fr:  "sacs",
		},
		{
			Cat: Category("jewelery"),
			Fr:  "bijoux",
		},
		{
			Cat: Category("accessories"),
			Fr:  "accessoires",
		},
		{
			Cat: Category("beauty"),
			Fr:  "beauté",
		},
		{
			Cat: Category("shirt"),
			Fr:  "chemise",
		},
		{
			Cat: Category("suit"),
			Fr:  "costume",
		},
	}
	for _, test := range tests {
		result := test.Cat.TranslateIn(LANG_FR)
		if result != test.Fr {
			t.Errorf(" Expected %v for %s but got %v ", test.Fr, test.Cat, result)
		}
	}
}

func TestTranslateSubCaterogyFr(t *testing.T) {
	tests := []struct {
		Subcat SubCategory
		Fr     string
	}{
		{
			Subcat: SubCategory("overcoat"),
			Fr:     "pardessus",
		},
		{
			Subcat: SubCategory("top jumpsuit"),
			Fr:     "combinaison",
		},
		{
			Subcat: SubCategory("shortsleeve"),
			Fr:     "manches courtes",
		},
		{
			Subcat: SubCategory("t-shirt"),
			Fr:     "t-shirt",
		},
		{
			Subcat: SubCategory("sleeveless"),
			Fr:     "sans manches",
		},
		{
			Subcat: SubCategory("longsleeve"),
			Fr:     "manches longues",
		},
		{
			Subcat: SubCategory("tank"),
			Fr:     "débardeur",
		},
		{
			Subcat: SubCategory("printed"),
			Fr:     "imprimée",
		},
		{
			Subcat: SubCategory("tunics"),
			Fr:     "tunique",
		},
		{
			Subcat: SubCategory("polo"),
			Fr:     "polo",
		},
		{
			Subcat: SubCategory("shirt"),
			Fr:     "chemise",
		},
		{
			Subcat: SubCategory("cashmere"),
			Fr:     "cashmere",
		},
		{
			Subcat: SubCategory("camisole"),
			Fr:     "caraco",
		},
		{
			Subcat: SubCategory("halter"),
			Fr:     "dos nu",
		},
		{
			Subcat: SubCategory("midi"),
			Fr:     "midi",
		},
		{
			Subcat: SubCategory("maxi"),
			Fr:     "maxi",
		},
		{
			Subcat: SubCategory("evening"),
			Fr:     "soirée",
		},
		{
			Subcat: SubCategory("cocktail"),
			Fr:     "cocktail",
		},
		{
			Subcat: SubCategory("bride"),
			Fr:     "mariée",
		},
		{
			Subcat: SubCategory("bridesmaid"),
			Fr:     "demoiselle d'honneur",
		},
		{
			Subcat: SubCategory("mini"),
			Fr:     "mini",
		},
		{
			Subcat: SubCategory("mid-length"),
			Fr:     "mi-longue",
		},
		{
			Subcat: SubCategory("long"),
			Fr:     "longue",
		},
		{
			Subcat: SubCategory("printed"),
			Fr:     "imprimée",
		},
		{
			Subcat: SubCategory("denim"),
			Fr:     "jeans",
		},
		{
			Subcat: SubCategory("leather"),
			Fr:     "cuir",
		},
		{
			Subcat: SubCategory("blazers"),
			Fr:     "blazers",
		},
		{
			Subcat: SubCategory("faux fur"),
			Fr:     "fourrure",
		},
		{
			Subcat: SubCategory("suede"),
			Fr:     "daim",
		},
		{
			Subcat: SubCategory("velvet"),
			Fr:     "velour",
		},
		{
			Subcat: SubCategory("coat"),
			Fr:     "manteau de pluie",
		},
		{
			Subcat: SubCategory("fur & shearling"),
			Fr:     "fourrure",
		},
		{
			Subcat: SubCategory("leather & suede"),
			Fr:     "cuir",
		},
		{
			Subcat: SubCategory("puffers"),
			Fr:     "doudoune",
		},
		{
			Subcat: SubCategory("trenchcoat"),
			Fr:     "trench",
		},
		{
			Subcat: SubCategory("wool"),
			Fr:     "laine",
		},
		{
			Subcat: SubCategory("cropped"),
			Fr:     "cropped",
		},
		{
			Subcat: SubCategory("legging"),
			Fr:     "legging",
		},
		{
			Subcat: SubCategory("skinny"),
			Fr:     "stretch",
		},
		{
			Subcat: SubCategory("wide leg"),
			Fr:     "ample",
		},
		{
			Subcat: SubCategory("culotte"),
			Fr:     "jupe-culotte",
		},
		{
			Subcat: SubCategory("classic"),
			Fr:     "fluide",
		},
		{
			Subcat: SubCategory("bernuda"),
			Fr:     "bernuda",
		},
		{
			Subcat: SubCategory("cardigan"),
			Fr:     "gilet",
		},
		{
			Subcat: SubCategory("cashmere"),
			Fr:     "cashmere",
		},
		{
			Subcat: SubCategory("crewneck & scoopneck"),
			Fr:     "col rond",
		},
		{
			Subcat: SubCategory("turtleneck"),
			Fr:     "col roulé",
		},
		{
			Subcat: SubCategory("v-neck"),
			Fr:     "col V",
		},
		{
			Subcat: SubCategory("hoodies"),
			Fr:     "sweat à capuche",
		},
		{
			Subcat: SubCategory("athletic"),
			Fr:     "running",
		},
		{
			Subcat: SubCategory("boots"),
			Fr:     "bottines",
		},
		{
			Subcat: SubCategory("leather boots"),
			Fr:     "bottes en cuir",
		},
		{
			Subcat: SubCategory("over the knee boots"),
			Fr:     "bottes hautes",
		},
		{
			Subcat: SubCategory("flat"),
			Fr:     "plate",
		},
		{
			Subcat: SubCategory("mules & clogs"),
			Fr:     "sandales à talons",
		},
		{
			Subcat: SubCategory("platforms"),
			Fr:     "talons hauts",
		},
		{
			Subcat: SubCategory("pumps"),
			Fr:     "escarpins",
		},
		{
			Subcat: SubCategory("sandals"),
			Fr:     "sandales",
		},
		{
			Subcat: SubCategory("sneakers"),
			Fr:     "baskets et tennis",
		},
		{
			Subcat: SubCategory("wedges"),
			Fr:     "semelles compensées",
		},
		{
			Subcat: SubCategory("one-piece"),
			Fr:     "monokini",
		},
		{
			Subcat: SubCategory("two-piece"),
			Fr:     "bikini",
		},
		{
			Subcat: SubCategory("cover-up"),
			Fr:     "paréo",
		},
		{
			Subcat: SubCategory("backpack"),
			Fr:     "sac à dos",
		},
		{
			Subcat: SubCategory("clutch"),
			Fr:     "pochettes",
		},
		{
			Subcat: SubCategory("hobo"),
			Fr:     "besace",
		},
		{
			Subcat: SubCategory("satchel"),
			Fr:     "cartable",
		},
		{
			Subcat: SubCategory("shoulder"),
			Fr:     "sacs à main",
		},
		{
			Subcat: SubCategory("duffel & tote"),
			Fr:     "caba",
		},
		{
			Subcat: SubCategory("wallet"),
			Fr:     "porte-monnaie",
		},
		{
			Subcat: SubCategory("bracelet"),
			Fr:     "bracelet",
		},
		{
			Subcat: SubCategory("charm"),
			Fr:     "breloque",
		},
		{
			Subcat: SubCategory("earrings"),
			Fr:     "boucles d'oreilles",
		},
		{
			Subcat: SubCategory("necklace"),
			Fr:     "collier",
		},
		{
			Subcat: SubCategory("pins"),
			Fr:     "pins",
		},
		{
			Subcat: SubCategory("ring"),
			Fr:     "bague",
		},
		{
			Subcat: SubCategory("diamond ring"),
			Fr:     "bague en diamant",
		},
		{
			Subcat: SubCategory("watch"),
			Fr:     "montre",
		},
		{
			Subcat: SubCategory("bracelet"),
			Fr:     "bracelet",
		},
		{
			Subcat: SubCategory("charm"),
			Fr:     "breloque",
		},
		{
			Subcat: SubCategory("gloves"),
			Fr:     "gants",
		},
		{
			Subcat: SubCategory("hat"),
			Fr:     "chapeau",
		},
		{
			Subcat: SubCategory("scarve & wrap"),
			Fr:     "écharpe & châle",
		},
		{
			Subcat: SubCategory("sunglasses"),
			Fr:     "lunettes de soleil",
		},
		{
			Subcat: SubCategory("eyeglasses"),
			Fr:     "lunettes de vue",
		},
		{
			Subcat: SubCategory("key chains"),
			Fr:     "porte clés",
		},
		{
			Subcat: SubCategory("umbrella"),
			Fr:     "parapluie",
		},
		{
			Subcat: SubCategory("headband"),
			Fr:     "serre-tête",
		},
		{
			Subcat: SubCategory("barrette"),
			Fr:     "barrette",
		},
		{
			Subcat: SubCategory("ear muffs"),
			Fr:     "cache-oreilles",
		},
		{
			Subcat: SubCategory("hair accessories"),
			Fr:     "autres accessoires pour cheveux",
		},
		{
			Subcat: SubCategory("belt"),
			Fr:     "ceinture",
		},
		{
			Subcat: SubCategory("bath & body"),
			Fr:     "savons & crème",
		},
		{
			Subcat: SubCategory("parfume"),
			Fr:     "parfum",
		},
		{
			Subcat: SubCategory("haircare"),
			Fr:     "soin des cheveux",
		},
		{
			Subcat: SubCategory("makeup"),
			Fr:     "maquillage",
		},
		{
			Subcat: SubCategory("skincare"),
			Fr:     "soin de la peau",
		},
		{
			Subcat: SubCategory("gingham"),
			Fr:     "à carreaux",
		},
		{
			Subcat: SubCategory("gingham"),
			Fr:     "à carreaux",
		},
		{
			Subcat: SubCategory("dress"),
			Fr:     "habillé",
		},
		{
			Subcat: SubCategory("mao collar"),
			Fr:     "à col mao",
		},
		{
			Subcat: SubCategory("fur"),
			Fr:     "fourrure",
		},
		{
			Subcat: SubCategory("chinos & khakis"),
			Fr:     "chinos",
		},
		{
			Subcat: SubCategory("casual"),
			Fr:     "décontracté",
		},
		{
			Subcat: SubCategory("cargo"),
			Fr:     "treillis",
		},
		{
			Subcat: SubCategory("jogger"),
			Fr:     "jogging",
		},
		{
			Subcat: SubCategory("classic  bermuda"),
			Fr:     "bermuda classic",
		},
		{
			Subcat: SubCategory("cardigan"),
			Fr:     "gilet",
		},
		{
			Subcat: SubCategory("zip up cardigan"),
			Fr:     "gilet à zip",
		},
		{
			Subcat: SubCategory("cashmere"),
			Fr:     "cashmere",
		},
		{
			Subcat: SubCategory("lace-up"),
			Fr:     "chaussures à lacets",
		},
		{
			Subcat: SubCategory("slip-ons & loafers"),
			Fr:     "mocassin",
		},
		{
			Subcat: SubCategory("classic swim short"),
			Fr:     "short de bain classique",
		},
		{
			Subcat: SubCategory("printed swim short"),
			Fr:     "short de bain imprimé",
		},
		{
			Subcat: SubCategory("brief"),
			Fr:     "slip de bain",
		},
		{
			Subcat: SubCategory("printed brief"),
			Fr:     "slip de bain imprimé",
		},
		{
			Subcat: SubCategory("sport watch"),
			Fr:     "montre de sport",
		},
		{
			Subcat: SubCategory("classic watch"),
			Fr:     "montre classique",
		},
		{
			Subcat: SubCategory("silver watch"),
			Fr:     "montre en argent",
		},
		{
			Subcat: SubCategory("gold watch"),
			Fr:     "montre en or",
		},
		{
			Subcat: SubCategory("cuff links"),
			Fr:     "bouton de manchette",
		},
		{
			Subcat: SubCategory("cross body bag"),
			Fr:     "sacs à bandoulière",
		},
		{
			Subcat: SubCategory("borsalino"),
			Fr:     "chapeau",
		},
		{
			Subcat: SubCategory("luggage"),
			Fr:     "valise",
		},
		{
			Subcat: SubCategory("ties"),
			Fr:     "cravate",
		},
		{
			Subcat: SubCategory("short jumpsuit"),
			Fr:     "combinaison short",
		},
		{
			Subcat: SubCategory("rompers"),
			Fr:     "combishort",
		},
		{
			Subcat: SubCategory("playsuit"),
			Fr:     "combishort",
		},
	}
	for _, test := range tests {
		result := test.Subcat.TranslateIn(LANG_FR)
		if result != test.Fr {
			t.Errorf(" Expected %s for %s but got %s", test.Fr, test.Subcat, result)
		}
	}
}
