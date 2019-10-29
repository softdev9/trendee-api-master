package jaccard

import (
	"testing"
)

// Test interface
//	SetCoef(float64)
//	Coef() float64

type TestAttributor struct {
	Keywords []string
	Coefi    float64
}

func (t TestAttributor) Attributes() []string {
	return t.Keywords
}

func (t TestAttributor) Coef() float64 {
	return t.Coefi
}

func (t TestAttributor) SetCoef(c float64) {
	t.Coefi = c
}

func TestJaccard(t *testing.T) {
	a1 := TestAttributor{Keywords: []string{"yellow"}}
	a2 := TestAttributor{Keywords: []string{"yellow"}}
	a3 := TestAttributor{Keywords: []string{"blue"}}
	testWith := []TestAttributor{a2, a3}
	// We compare a1 with the article a2 and a3
	var interfaceSlice []interface{} = make([]interface{}, len(testWith))
	for i, d := range testWith {
		interfaceSlice[i] = d
	}
	result := JaccardIndex(a1, a2)
	if result != 1 {
		t.Error("Expected was 1 but got ", result)
	}
}
