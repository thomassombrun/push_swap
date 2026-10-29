package stack

import (
	"math/rand"
	"reflect"
	"testing"
)

// appliquer exécute les instructions sur une copie de a et renvoie l'état final.
func appliquer(t *testing.T, a Stack, instrs []string) State {
	t.Helper()
	st := State{A: append(Stack{}, a...)}
	for _, instr := range instrs {
		if !st.Exec(instr) {
			t.Fatalf("instruction invalide produite par Sort : %q", instr)
		}
	}
	return st
}

// verifierTri contrôle que les instructions trient a et vident b.
func verifierTri(t *testing.T, a Stack) []string {
	t.Helper()
	instrs := Sort(a)
	st := appliquer(t, a, instrs)
	if !st.IsSorted() {
		t.Fatalf("Sort(%v) = %v ne trie pas : a=%v b=%v", a, instrs, st.A, st.B)
	}
	if len(st.A) != len(a) {
		t.Fatalf("Sort(%v) a perdu des éléments : a=%v", a, st.A)
	}
	return instrs
}

func TestSortCasLimites(t *testing.T) {
	tests := []struct {
		name string
		in   Stack
	}{
		{"pile vide", Stack{}},
		{"un seul élément", Stack{1}},
		{"déjà triée (2)", Stack{1, 2}},
		{"déjà triée (3)", Stack{1, 2, 3}},
		{"déjà triée (10)", Stack{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{"triée avec négatifs", Stack{-5, -1, 0, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sort(tt.in); len(got) != 0 {
				t.Errorf("Sort(%v) = %v, attendu aucune instruction", tt.in, got)
			}
		})
	}
}

func TestSortDeuxElements(t *testing.T) {
	got := Sort(Stack{2, 1})
	if want := []string{"sa"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Sort({2,1}) = %v, attendu %v", got, want)
	}
}

func TestSortExemplesDuSujet(t *testing.T) {
	exemples := []Stack{
		{2, 1, 3, 6, 5, 8},
		{4, 67, 3, 87, 23},
		{2, 1, 3, 6, 8, 5},
	}
	for _, a := range exemples {
		verifierTri(t, a)
	}
}

func TestSortNeModifiePasLEntree(t *testing.T) {
	a := Stack{5, 3, 8, 1, 9, 2}
	copie := append(Stack{}, a...)
	Sort(a)
	if !reflect.DeepEqual(a, copie) {
		t.Errorf("Sort a modifié sa pile d'entrée : %v, attendu %v", a, copie)
	}
}

// permutations génère toutes les permutations de v (algorithme de Heap).
func permutations(v []int) [][]int {
	var res [][]int
	var generer func(k int)
	tab := append([]int{}, v...)
	generer = func(k int) {
		if k == 1 {
			res = append(res, append([]int{}, tab...))
			return
		}
		generer(k - 1)
		for i := 0; i < k-1; i++ {
			if k%2 == 0 {
				tab[i], tab[k-1] = tab[k-1], tab[i]
			} else {
				tab[0], tab[k-1] = tab[k-1], tab[0]
			}
			generer(k - 1)
		}
	}
	generer(len(tab))
	return res
}

// Toutes les permutations de 1 à 7 éléments doivent être triées correctement.
func TestSortToutesLesPermutations(t *testing.T) {
	for n := 1; n <= 7; n++ {
		base := make([]int, n)
		for i := range base {
			base[i] = i + 1
		}
		perms := permutations(base)
		for _, p := range perms {
			verifierTri(t, Stack(p))
		}
		t.Logf("n=%d : %d permutations vérifiées", n, len(perms))
	}
}

// Des valeurs négatives, zéro et grandes doivent être gérées.
func TestSortValeursVariees(t *testing.T) {
	cas := []Stack{
		{0, -1, 1},
		{-10, 50, -3, 0, 7},
		{1000000, -1000000, 0, 42},
		{3, 2, 1, 0, -1, -2, -3},
	}
	for _, a := range cas {
		verifierTri(t, a)
	}
}

func TestSortListesAleatoires(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for _, n := range []int{8, 10, 20, 50, 100, 200, 500} {
		for essai := 0; essai < 20; essai++ {
			a := Stack(rng.Perm(n))
			verifierTri(t, a)
		}
	}
}

// Sort ne doit produire que des instructions connues, sans espace ni ligne vide.
func TestSortInstructionsBienFormees(t *testing.T) {
	valides := map[string]bool{
		"pa": true, "pb": true, "sa": true, "sb": true, "ss": true,
		"ra": true, "rb": true, "rr": true, "rra": true, "rrb": true, "rrr": true,
	}
	rng := rand.New(rand.NewSource(7))
	for i := 0; i < 50; i++ {
		for _, instr := range Sort(Stack(rng.Perm(30))) {
			if !valides[instr] {
				t.Fatalf("instruction mal formée : %q", instr)
			}
		}
	}
}

// Pas d'instructions manifestement inutiles : jamais pa suivi de pb (ou l'inverse), ni ra suivi de rra.
func TestSortSansOperationsRedondantes(t *testing.T) {
	annulations := map[string]string{
		"pa": "pb", "pb": "pa",
		"ra": "rra", "rra": "ra",
		"rb": "rrb", "rrb": "rb",
		"rr": "rrr", "rrr": "rr",
	}
	rng := rand.New(rand.NewSource(11))
	for i := 0; i < 50; i++ {
		instrs := Sort(Stack(rng.Perm(25)))
		for j := 1; j < len(instrs); j++ {
			if annulations[instrs[j-1]] == instrs[j] {
				t.Fatalf("paire d'instructions qui s'annulent (%s, %s) en position %d", instrs[j-1], instrs[j], j)
			}
		}
	}
}

// Contrôle de la taille de la sortie : ces seuils sont larges et ne doivent pas être dépassés
// si l'algorithme est optimisé (à ajuster selon ton algorithme).
func TestSortNombreDInstructions(t *testing.T) {
	rng := rand.New(rand.NewSource(3))

	// Le sujet annonce 6 instructions pour "4 67 3 87 23".
	n := len(Sort(Stack{4, 67, 3, 87, 23}))
	t.Logf("4 67 3 87 23 -> %d instructions (sujet : 6)", n)

	for _, taille := range []int{5, 100, 500} {
		total := 0
		const essais = 20
		for i := 0; i < essais; i++ {
			total += len(Sort(Stack(rng.Perm(taille))))
		}
		t.Logf("n=%d : moyenne de %d instructions", taille, total/essais)
	}
}
