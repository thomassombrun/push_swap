package stack

import (
	"reflect"
	"testing"
)

// Chaque instruction est testée sur un état de départ, avec l'état attendu.
func TestExec(t *testing.T) {
	tests := []struct {
		instr string
		a, b  Stack // état de départ
		wa    Stack // a attendu
		wb    Stack // b attendu
	}{
		{"pa", Stack{1}, Stack{2, 3}, Stack{2, 1}, Stack{3}},
		{"pb", Stack{1, 2}, Stack{3}, Stack{2}, Stack{1, 3}},
		{"sa", Stack{1, 2, 3}, Stack{4, 5}, Stack{2, 1, 3}, Stack{4, 5}},
		{"sb", Stack{1, 2, 3}, Stack{4, 5}, Stack{1, 2, 3}, Stack{5, 4}},
		{"ss", Stack{1, 2, 3}, Stack{4, 5, 6}, Stack{2, 1, 3}, Stack{5, 4, 6}},
		{"ra", Stack{1, 2, 3}, Stack{4, 5, 6}, Stack{2, 3, 1}, Stack{4, 5, 6}},
		{"rb", Stack{1, 2, 3}, Stack{4, 5, 6}, Stack{1, 2, 3}, Stack{5, 6, 4}},
		{"rr", Stack{1, 2, 3}, Stack{4, 5, 6}, Stack{2, 3, 1}, Stack{5, 6, 4}},
		{"rra", Stack{1, 2, 3}, Stack{4, 5, 6}, Stack{3, 1, 2}, Stack{4, 5, 6}},
		{"rrb", Stack{1, 2, 3}, Stack{4, 5, 6}, Stack{1, 2, 3}, Stack{6, 4, 5}},
		{"rrr", Stack{1, 2, 3}, Stack{4, 5, 6}, Stack{3, 1, 2}, Stack{6, 4, 5}},
	}

	for _, tt := range tests {
		t.Run(tt.instr, func(t *testing.T) {
			st := State{A: append(Stack{}, tt.a...), B: append(Stack{}, tt.b...)}
			if !st.Exec(tt.instr) {
				t.Fatalf("Exec(%q) a renvoyé false", tt.instr)
			}
			if !reflect.DeepEqual(st.A, tt.wa) || !reflect.DeepEqual(st.B, tt.wb) {
				t.Errorf("après %s : a=%v b=%v, attendu a=%v b=%v", tt.instr, st.A, st.B, tt.wa, tt.wb)
			}
		})
	}
}

func TestExecInstructionInconnue(t *testing.T) {
	invalides := []string{"", "xx", "PA", "pa ", " pa", "pc", "sa\r", "rrx", "ra rb", "1"}
	for _, instr := range invalides {
		t.Run("["+instr+"]", func(t *testing.T) {
			st := State{A: Stack{1, 2, 3}, B: Stack{4}}
			if st.Exec(instr) {
				t.Errorf("Exec(%q) a renvoyé true, attendu false", instr)
			}
			if !reflect.DeepEqual(st.A, Stack{1, 2, 3}) || !reflect.DeepEqual(st.B, Stack{4}) {
				t.Errorf("l'état a changé après une instruction invalide : a=%v b=%v", st.A, st.B)
			}
		})
	}
}

// Sur des piles vides ou trop petites, aucune instruction ne doit paniquer.
func TestExecSurPilesVides(t *testing.T) {
	instrs := []string{"pa", "pb", "sa", "sb", "ss", "ra", "rb", "rr", "rra", "rrb", "rrr"}
	for _, instr := range instrs {
		t.Run(instr, func(t *testing.T) {
			st := State{}
			if !st.Exec(instr) {
				t.Errorf("Exec(%q) sur piles vides a renvoyé false", instr)
			}
			if len(st.A) != 0 || len(st.B) != 0 {
				t.Errorf("piles vides modifiées par %s : a=%v b=%v", instr, st.A, st.B)
			}
		})
	}
}

func TestPaPbSurPileSourceVide(t *testing.T) {
	st := State{A: Stack{1, 2}}
	st.Exec("pa") // B est vide : sans effet
	if !reflect.DeepEqual(st.A, Stack{1, 2}) || len(st.B) != 0 {
		t.Errorf("pa avec b vide a modifié les piles : a=%v b=%v", st.A, st.B)
	}

	st = State{B: Stack{1, 2}}
	st.Exec("pb") // A est vide : sans effet
	if len(st.A) != 0 || !reflect.DeepEqual(st.B, Stack{1, 2}) {
		t.Errorf("pb avec a vide a modifié les piles : a=%v b=%v", st.A, st.B)
	}
}

// Les instructions combinées doivent équivaloir aux deux instructions séparées.
func TestInstructionsCombinees(t *testing.T) {
	combos := map[string][2]string{
		"ss":  {"sa", "sb"},
		"rr":  {"ra", "rb"},
		"rrr": {"rra", "rrb"},
	}
	for combo, parts := range combos {
		t.Run(combo, func(t *testing.T) {
			s1 := State{A: Stack{1, 2, 3, 4}, B: Stack{5, 6, 7}}
			s2 := State{A: Stack{1, 2, 3, 4}, B: Stack{5, 6, 7}}
			s1.Exec(combo)
			s2.Exec(parts[0])
			s2.Exec(parts[1])
			if !reflect.DeepEqual(s1, s2) {
				t.Errorf("%s = %v, mais %s+%s = %v", combo, s1, parts[0], parts[1], s2)
			}
		})
	}
}

// Exemple du sujet : "2 1 3 6 8 5" avec sa, pb pb pb, rb, rrr, pa pa pa.
func TestExempleDuSujet(t *testing.T) {
	st := State{A: Stack{2, 1, 3, 6, 8, 5}}
	for _, instr := range []string{"sa", "pb", "pb", "pb", "rb", "rrr", "pa", "pa", "pa"} {
		if !st.Exec(instr) {
			t.Fatalf("instruction %q refusée", instr)
		}
	}
	if want := (Stack{1, 2, 3, 5, 6, 8}); !reflect.DeepEqual(st.A, want) {
		t.Errorf("a = %v, attendu %v", st.A, want)
	}
	if len(st.B) != 0 {
		t.Errorf("b = %v, attendu vide", st.B)
	}
	if !st.IsSorted() {
		t.Error("IsSorted() = false, attendu true")
	}
}

// Exemple du checker du sujet : "3 2 1 0" avec rra pb sa rra pa donne OK.
func TestExempleCheckerOK(t *testing.T) {
	st := State{A: Stack{3, 2, 1, 0}}
	for _, instr := range []string{"rra", "pb", "sa", "rra", "pa"} {
		st.Exec(instr)
	}
	if !st.IsSorted() {
		t.Errorf("attendu OK, mais a=%v b=%v", st.A, st.B)
	}
}

// Exemple du checker du sujet : "3 2 1 0" avec sa rra pb donne KO.
func TestExempleCheckerKO(t *testing.T) {
	st := State{A: Stack{3, 2, 1, 0}}
	for _, instr := range []string{"sa", "rra", "pb"} {
		st.Exec(instr)
	}
	if st.IsSorted() {
		t.Errorf("attendu KO, mais a=%v b=%v", st.A, st.B)
	}
}

func TestIsSorted(t *testing.T) {
	tests := []struct {
		name string
		st   State
		want bool
	}{
		{"a triée, b vide", State{A: Stack{1, 2, 3}}, true},
		{"piles vides", State{}, true},
		{"a triée mais b non vide", State{A: Stack{1, 2, 3}, B: Stack{4}}, false},
		{"a non triée, b vide", State{A: Stack{2, 1, 3}}, false},
		{"a vide mais b non vide", State{B: Stack{1}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.st.IsSorted(); got != tt.want {
				t.Errorf("IsSorted() = %v, attendu %v", got, tt.want)
			}
		})
	}
}
