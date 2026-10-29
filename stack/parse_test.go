package stack

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   Stack
		wantOK bool
	}{
		{"liste normale", "3 2 1 0", Stack{3, 2, 1, 0}, true},
		{"un seul entier", "42", Stack{42}, true},
		{"nombres négatifs", "-1 -5 3", Stack{-1, -5, 3}, true},
		{"signe plus", "+4 2", Stack{4, 2}, true},
		{"espaces multiples", "  1   2\t3  ", Stack{1, 2, 3}, true},
		{"chaîne vide", "", Stack{}, true},
		{"que des espaces", "   ", Stack{}, true},
		{"exemple du sujet", "2 1 3 6 5 8", Stack{2, 1, 3, 6, 5, 8}, true},
		{"mot non entier", "0 one 2 3", nil, false},
		{"décimal", "1 2.5 3", nil, false},
		{"lettre collée", "12a 3", nil, false},
		{"dépassement d'int", "99999999999999999999", nil, false},
		{"doublon simple", "1 2 1", nil, false},
		{"doublon adjacent", "5 5", nil, false},
		{"doublon avec zéros", "1 01", nil, false},
		{"doublon zéro et -0", "0 -0", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := Parse(tt.input)
			if ok != tt.wantOK {
				t.Fatalf("Parse(%q) ok = %v, attendu %v", tt.input, ok, tt.wantOK)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse(%q) = %v, attendu %v", tt.input, got, tt.want)
			}
		})
	}
}

// Le premier entier doit se trouver au sommet de la pile (index 0).
func TestParseSommetEstLePremier(t *testing.T) {
	st, ok := Parse("7 8 9")
	if !ok {
		t.Fatal("Parse a échoué")
	}
	if v, _ := st.Pop(); v != 7 {
		t.Errorf("sommet = %d, attendu 7", v)
	}
}
