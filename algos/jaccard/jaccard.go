package jaccard

import (
	"github.com/spiderdev86/trendee-api/data"
	"sort"
)

type Classifier interface {
	Attributes() []string
	SetCoef(float64)
	Coef() float64
}

type Attributor interface {
	Attributes() []string
}

type ByCoef []*data.Article

func (a ByCoef) Len() int {
	return len(a)
}

func (a ByCoef) Swap(i, j int) {
	tmp := a[j]
	a[j] = a[i]
	a[i] = tmp
}

func (a ByCoef) Less(i, j int) bool {
	return a[i].Coef() > a[j].Coef()
}

func BuildComparaisonBag(a1 Attributor, a2 Attributor) []string {
	var fullbag []string
	keywordsA := a1.Attributes()
	keywordsB := a2.Attributes()
	for _, keyword := range keywordsA {
		fullbag = append(fullbag, keyword)
	}
	for _, keyword := range keywordsB {
		inSet := false
		for _, present := range fullbag {
			if keyword == present {
				inSet = true
			}
		}
		if !inSet {
			fullbag = append(fullbag, keyword)
		}
	}
	return fullbag
}

func JaccardIndex(comparee Attributor, with Attributor) float64 {
	fullBag := BuildComparaisonBag(comparee, with)
	compareelMap := BuildBooleanArrtibuteMap(comparee.Attributes(), fullBag)
	withMap := BuildBooleanArrtibuteMap(with.Attributes(), fullBag)
	m00 := CalucalteM00(withMap, compareelMap)
	m01 := CalucalteM01(withMap, compareelMap)
	m10 := CalucalteM10(withMap, compareelMap)
	m11 := CalucalteM11(withMap, compareelMap)
	jaccard := CaluclateJ(m00, m01, m10, m11)
	return jaccard
}

func JaccardOrder(comparee Attributor, with []data.Article) []*data.Article {
	var result []*data.Article
	for i := 0; i < len(with); i++ {
		art := &with[i]
		jaccardIndex := JaccardIndex(comparee, art)
		art.SetCoef(jaccardIndex)
		result = append(result, art)
	}
	sort.Sort(ByCoef(result))
	return result
}

func JaccardOrderSelfie(comparee data.Selfie, with []data.Article) []*data.Article {
	var result []*data.Article
	for i := 0; i < len(with); i++ {
		art := &with[i]
		jaccardIndex := JaccardIndex(&comparee, art)
		art.SetCoef(jaccardIndex)
		result = append(result, art)
	}
	sort.Sort(ByCoef(result))
	return result
}

func BuildBooleanArrtibuteMap(comparee []string, fullBag []string) map[string]bool {
	boolMap := make(map[string]bool)
	// Create the attribue map
	for _, t := range fullBag {
		boolMap[t] = false
	}
	for _, t := range comparee {
		boolMap[t] = true
	}
	return boolMap
}

func CaluclateJ(m00 int, m01 int, m10 int, m11 int) float64 {
	var result float64 = float64(m11) / float64(m01+m10+m11)
	return result
}

func CalucalteM00(a1 map[string]bool, a2 map[string]bool) int {
	m := 0
	for key, value := range a1 {
		if !value && !a2[key] {
			m++
		}
	}
	return m
}

func CalucalteM01(a1 map[string]bool, a2 map[string]bool) int {
	m := 0
	for key, value := range a1 {
		if !value && a2[key] {
			m++
		}
	}
	return m
}

func CalucalteM10(a1 map[string]bool, a2 map[string]bool) int {
	m := 0
	for key, value := range a1 {
		if value && !a2[key] {
			m++
		}
	}
	return m
}

func CalucalteM11(a1 map[string]bool, a2 map[string]bool) int {
	m := 0
	for key, value := range a1 {
		if value && a2[key] {
			m++
		}
	}
	return m
}
