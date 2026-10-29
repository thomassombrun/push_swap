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

// compile construit l'exécutable checker dans un dossier temporaire.
func compile(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "checker")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	if out, err := exec.Command("go", "build", "-o", bin, ".").CombinedOutput(); err != nil {
		t.Fatalf("échec de la compilation : %v\n%s", err, out)
	}
	return bin
}

// lancer exécute le programme avec des arguments et une entrée standard.
func lancer(t *testing.T, bin, stdin string, args ...string) (stdout, stderr string, code int) {
	t.Helper()
	cmd := exec.Command(bin, args...)
	cmd.Stdin = strings.NewReader(stdin)
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

func TestCheckerSansArgument(t *testing.T) {
	bin := compile(t)
	stdout, stderr, _ := lancer(t, bin, "sa\npb\n")
	if stdout != "" || stderr != "" {
		t.Errorf("sans argument : stdout=%q stderr=%q, attendu aucune sortie", stdout, stderr)
	}
}

func TestCheckerOK(t *testing.T) {
	bin := compile(t)
	// Exemple du sujet.
	stdout, stderr, _ := lancer(t, bin, "rra\npb\nsa\nrra\npa\n", "3 2 1 0")
	if stdout != "OK\n" || stderr != "" {
		t.Errorf("stdout=%q stderr=%q, attendu \"OK\\n\" et stderr vide", stdout, stderr)
	}
}

func TestCheckerKO(t *testing.T) {
	bin := compile(t)
	// Exemple du sujet.
	stdout, stderr, _ := lancer(t, bin, "sa\nrra\npb\n", "3 2 1 0")
	if stdout != "KO\n" || stderr != "" {
		t.Errorf("stdout=%q stderr=%q, attendu \"KO\\n\" et stderr vide", stdout, stderr)
	}
}

func TestCheckerSansInstruction(t *testing.T) {
	bin := compile(t)
	if stdout, _, _ := lancer(t, bin, "", "1 2 3"); stdout != "OK\n" {
		t.Errorf("liste triée sans instruction : %q, attendu OK", stdout)
	}
	if stdout, _, _ := lancer(t, bin, "", "2 1 3"); stdout != "KO\n" {
		t.Errorf("liste non triée sans instruction : %q, attendu KO", stdout)
	}
}

// OK exige que b soit vide : pb seul sur une pile triée doit donner KO.
func TestCheckerKOSiBNonVide(t *testing.T) {
	bin := compile(t)
	if stdout, _, _ := lancer(t, bin, "pb\n", "1 2 3"); stdout != "KO\n" {
		t.Errorf("stdout=%q, attendu KO (b non vide)", stdout)
	}
}

func TestCheckerLignesVidesIgnorees(t *testing.T) {
	bin := compile(t)
	// Exemple du sujet : echo -e "rra\npb\nsa\nrra\npa\n" laisse une ligne vide.
	if stdout, _, _ := lancer(t, bin, "rra\npb\nsa\nrra\npa\n\n", "3 2 1 0"); stdout != "OK\n" {
		t.Errorf("stdout=%q, attendu OK", stdout)
	}
}

func TestCheckerErreurs(t *testing.T) {
	bin := compile(t)
	tests := []struct {
		name  string
		stdin string
		args  []string
	}{
		{"argument non entier", "rra\npb\nsa\n", []string{"3 2 one 0"}},
		{"doublon", "sa\n", []string{"1 2 2"}},
		{"instruction inconnue", "sa\nxx\n", []string{"3 2 1"}},
		{"instruction en majuscules", "SA\n", []string{"3 2 1"}},
		{"instruction avec espace", "sa \n", []string{"3 2 1"}},
		{"deux instructions sur une ligne", "sa pb\n", []string{"3 2 1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, code := lancer(t, bin, tt.stdin, tt.args...)
			if stderr != "Error\n" {
				t.Errorf("stderr=%q, attendu \"Error\\n\"", stderr)
			}
			if stdout != "" {
				t.Errorf("stdout=%q, attendu vide en cas d'erreur", stdout)
			}
			if code == 0 {
				t.Error("code de sortie 0, attendu non nul")
			}
		})
	}
}

// --- Cas limites ci-dessous : ils échouent tant que ces deux points ne sont pas corrigés dans main.go ---

// Une liste vide ("") doit se comporter comme l'absence d'argument : aucune sortie.
func TestCheckerArgumentVide(t *testing.T) {
	bin := compile(t)
	stdout, stderr, _ := lancer(t, bin, "", "")
	if stdout != "" || stderr != "" {
		t.Errorf("argument vide : stdout=%q stderr=%q, attendu aucune sortie", stdout, stderr)
	}
}

// Les entiers peuvent être passés en arguments séparés : ./checker 3 2 1 0
func TestCheckerPlusieursArguments(t *testing.T) {
	bin := compile(t)
	// "2 1" non triée sans instruction : KO attendu. Si seul le premier argument est lu,
	// la pile vaut {2}, ce qui est trié, et le programme affiche OK à tort.
	stdout, stderr, _ := lancer(t, bin, "", "2", "1")
	if stdout != "KO\n" || stderr != "" {
		t.Errorf("stdout=%q stderr=%q, attendu KO", stdout, stderr)
	}
	// Avec sa, la pile {2,1} devient triée : OK attendu.
	stdout, stderr, _ = lancer(t, bin, "sa\n", "2", "1")
	if stdout != "OK\n" || stderr != "" {
		t.Errorf("avec sa : stdout=%q stderr=%q, attendu OK", stdout, stderr)
	}
}
