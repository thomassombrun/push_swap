package stack

import (
	"reflect"
	"testing"
)

func TestPush(t *testing.T) {
	s := Stack{2, 3}
	s.Push(1)
	if want := (Stack{1, 2, 3}); !reflect.DeepEqual(s, want) {
		t.Errorf("Push = %v, attendu %v", s, want)
	}

	var vide Stack
	vide.Push(5)
	if want := (Stack{5}); !reflect.DeepEqual(vide, want) {
		t.Errorf("Push sur pile vide = %v, attendu %v", vide, want)
	}
}

func TestPop(t *testing.T) {
	s := Stack{1, 2, 3}
	v, ok := s.Pop()
	if !ok || v != 1 {
		t.Errorf("Pop = (%d, %v), attendu (1, true)", v, ok)
	}
	if want := (Stack{2, 3}); !reflect.DeepEqual(s, want) {
		t.Errorf("pile après Pop = %v, attendu %v", s, want)
	}

	vide := Stack{}
	if v, ok := vide.Pop(); ok || v != 0 {
		t.Errorf("Pop sur pile vide = (%d, %v), attendu (0, false)", v, ok)
	}
}

func TestSwapRotateReverseRotate(t *testing.T) {
	tests := []struct {
		name string
		op   func(*Stack)
		in   Stack
		want Stack
	}{
		{"Swap normal", (*Stack).Swap, Stack{1, 2, 3}, Stack{2, 1, 3}},
		{"Swap deux éléments", (*Stack).Swap, Stack{1, 2}, Stack{2, 1}},
		{"Swap un élément", (*Stack).Swap, Stack{1}, Stack{1}},
		{"Swap pile vide", (*Stack).Swap, Stack{}, Stack{}},

		{"Rotate normal", (*Stack).Rotate, Stack{1, 2, 3}, Stack{2, 3, 1}},
		{"Rotate deux éléments", (*Stack).Rotate, Stack{1, 2}, Stack{2, 1}},
		{"Rotate un élément", (*Stack).Rotate, Stack{1}, Stack{1}},
		{"Rotate pile vide", (*Stack).Rotate, Stack{}, Stack{}},

		{"ReverseRotate normal", (*Stack).ReverseRotate, Stack{1, 2, 3}, Stack{3, 1, 2}},
		{"ReverseRotate deux éléments", (*Stack).ReverseRotate, Stack{1, 2}, Stack{2, 1}},
		{"ReverseRotate un élément", (*Stack).ReverseRotate, Stack{1}, Stack{1}},
		{"ReverseRotate pile vide", (*Stack).ReverseRotate, Stack{}, Stack{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := append(Stack{}, tt.in...)
			tt.op(&s)
			if !reflect.DeepEqual(s, tt.want) {
				t.Errorf("%v -> %v, attendu %v", tt.in, s, tt.want)
			}
		})
	}
}

// Rotate puis ReverseRotate (et inversement) doivent se compenser.
func TestRotateEtReverseRotateSontInverses(t *testing.T) {
	orig := Stack{4, 8, 15, 16, 23, 42}

	s := append(Stack{}, orig...)
	s.Rotate()
	s.ReverseRotate()
	if !reflect.DeepEqual(s, orig) {
		t.Errorf("Rotate puis ReverseRotate = %v, attendu %v", s, orig)
	}

	s = append(Stack{}, orig...)
	s.ReverseRotate()
	s.Rotate()
	if !reflect.DeepEqual(s, orig) {
		t.Errorf("ReverseRotate puis Rotate = %v, attendu %v", s, orig)
	}
}

func TestIsSortedAsc(t *testing.T) {
	tests := []struct {
		name string
		in   Stack
		want bool
	}{
		{"pile vide", Stack{}, true},
		{"un élément", Stack{7}, true},
		{"triée", Stack{1, 2, 3, 4}, true},
		{"triée avec négatifs", Stack{-3, -1, 0, 5}, true},
		{"non triée", Stack{1, 3, 2}, false},
		{"décroissante", Stack{3, 2, 1}, false},
		{"dernier élément trop petit", Stack{1, 2, 3, 0}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.in.IsSortedAsc(); got != tt.want {
				t.Errorf("IsSortedAsc(%v) = %v, attendu %v", tt.in, got, tt.want)
			}
		})
	}
}
