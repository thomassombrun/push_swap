package main

import (
	"bytes"
	"errors"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// compile construit l'exécutable push-swap dans un dossier temporaire.
func compile(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "push-swap")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	if out, err := exec.Command("go", "build", "-o", bin, ".").CombinedOutput(); err != nil {
		t.Fatalf("échec de la compilation : %v\n%s", err, out)
	}
	return bin
}

func lancer(t *testing.T, bin string, args ...string) (stdout, stderr string, code int) {
	t.Helper()
	cmd := exec.Command(bin, args...)
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	err := cmd.Run()
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		code = exitErr.ExitCode()
	} else if err != nil {
		t.Fatalf("exécution impossible : %v", err)
	}
	return out.String(), errb.String(), code
}

// instructions découpe la sortie en lignes (une instruction par ligne).
func instructions(sortie string) []string {
	if sortie == "" {
		return nil
	}
	return strings.Split(strings.TrimSuffix(sortie, "\n"), "\n")
}

func TestPushSwapSansArgument(t *testing.T) {
	bin := compile(t)
	stdout, stderr, _ := lancer(t, bin)
	if stdout != "" || stderr != "" {
		t.Errorf("sans argument : stdout=%q stderr=%q, attendu aucune sortie", stdout, stderr)
	}
}

func TestPushSwapDejaTriee(t *testing.T) {
	bin := compile(t)
	for _, arg := range []string{"1", "1 2", "1 2 3 4 5"} {
		stdout, stderr, _ := lancer(t, bin, arg)
		if stdout != "" || stderr != "" {
			t.Errorf("%q : stdout=%q stderr=%q, attendu aucune sortie", arg, stdout, stderr)
		}
	}
}

func TestPushSwapErreurs(t *testing.T) {
	bin := compile(t)
	tests := []struct{ name, arg string }{
		{"non entier", "0 one 2 3"},
		{"décimal", "1 2.5"},
		{"doublon", "1 2 3 2"},
		{"dépassement d'int", "1 99999999999999999999"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, code := lancer(t, bin, tt.arg)
			if stderr != "Error\n" {
				t.Errorf("stderr=%q, attendu \"Error\\n\"", stderr)
			}
			if stdout != "" {
				t.Errorf("stdout=%q, attendu vide", stdout)
			}
			if code == 0 {
				t.Error("code de sortie 0, attendu non nul")
			}
		})
	}
}

// La sortie doit contenir uniquement des instructions valides, une par ligne, terminées par \n.
func TestPushSwapFormatDeSortie(t *testing.T) {
	bin := compile(t)
	stdout, stderr, code := lancer(t, bin, "2 1 3 6 5 8")
	if code != 0 || stderr != "" {
		t.Fatalf("code=%d stderr=%q", code, stderr)
	}
	if !strings.HasSuffix(stdout, "\n") {
		t.Errorf("la sortie ne se termine pas par un retour à la ligne : %q", stdout)
	}
	if strings.Contains(stdout, " ") || strings.Contains(stdout, "\r") || strings.Contains(stdout, "\n\n") {
		t.Errorf("la sortie contient des espaces, des \\r ou des lignes vides : %q", stdout)
	}
	valides := map[string]bool{
		"pa": true, "pb": true, "sa": true, "sb": true, "ss": true,
		"ra": true, "rb": true, "rr": true, "rra": true, "rrb": true, "rrr": true,
	}
	for _, instr := range instructions(stdout) {
		if !valides[instr] {
			t.Errorf("instruction invalide dans la sortie : %q", instr)
		}
	}
}

// Test de bout en bout : la sortie de push-swap, donnée à un simulateur, doit trier la pile.
// (Équivalent de ./push-swap "$ARG" | ./checker "$ARG" sans dépendre du binaire checker.)
func TestPushSwapTriePourDeVrai(t *testing.T) {
	bin := compile(t)
	cas := []string{
		"2 1 3 6 5 8",
		"4 67 3 87 23",
		"3 2 1",
		"2 1",
		"5 4 3 2 1",
		"-3 10 0 -7 5 2",
		"9 8 7 6 5 4 3 2 1 0",
	}
	for _, arg := range cas {
		t.Run(arg, func(t *testing.T) {
			stdout, stderr, code := lancer(t, bin, arg)
			if code != 0 || stderr != "" {
				t.Fatalf("code=%d stderr=%q", code, stderr)
			}
			if !simulerEtVerifier(t, arg, instructions(stdout)) {
				t.Errorf("les instructions %v ne trient pas %q", instructions(stdout), arg)
			}
		})
	}
}

// Le sujet annonce 6 instructions pour "4 67 3 87 23" (cible à atteindre pour C4).
// Ce test ne fait qu'informer : il n'échoue pas tant que la sortie reste correcte.
func TestPushSwapNombreDInstructionsExempleDuSujet(t *testing.T) {
	bin := compile(t)
	stdout, _, _ := lancer(t, bin, "4 67 3 87 23")
	n := len(instructions(stdout))
	t.Logf("4 67 3 87 23 -> %d instructions (sujet : 6)", n)
}

// --- Petit simulateur indépendant du package stack, pour ne pas tester le code avec lui-même ---

func simulerEtVerifier(t *testing.T, arg string, instrs []string) bool {
	t.Helper()
	var a, b []int
	for _, mot := range strings.Fields(arg) {
		n := 0
		neg := false
		s := mot
		if strings.HasPrefix(s, "-") {
			neg, s = true, s[1:]
		}
		for _, c := range s {
			n = n*10 + int(c-'0')
		}
		if neg {
			n = -n
		}
		a = append(a, n)
	}
	for _, in := range instrs {
		switch in {
		case "pa":
			if len(b) > 0 {
				a = append([]int{b[0]}, a...)
				b = b[1:]
			}
		case "pb":
			if len(a) > 0 {
				b = append([]int{a[0]}, b...)
				a = a[1:]
			}
		case "sa":
			swap(a)
		case "sb":
			swap(b)
		case "ss":
			swap(a)
			swap(b)
		case "ra":
			a = rotate(a)
		case "rb":
			b = rotate(b)
		case "rr":
			a, b = rotate(a), rotate(b)
		case "rra":
			a = reverse(a)
		case "rrb":
			b = reverse(b)
		case "rrr":
			a, b = reverse(a), reverse(b)
		default:
			t.Fatalf("instruction inconnue : %q", in)
		}
	}
	if len(b) != 0 {
		return false
	}
	for i := 1; i < len(a); i++ {
		if a[i-1] > a[i] {
			return false
		}
	}
	return true
}

func swap(s []int) {
	if len(s) >= 2 {
		s[0], s[1] = s[1], s[0]
	}
}

func rotate(s []int) []int {
	if len(s) < 2 {
		return s
	}
	return append(append([]int{}, s[1:]...), s[0])
}

func reverse(s []int) []int {
	if len(s) < 2 {
		return s
	}
	return append([]int{s[len(s)-1]}, s[:len(s)-1]...)
}
