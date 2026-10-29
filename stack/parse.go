package stack

import (
	"strconv"
	"strings"
)

// La fonction Parse transforme un texte comme "3 2 1 0" en pile.
// Elle renvoie false si un mot n'est pas un entier ou s'il y a un doublon.
func Parse(s string) (Stack, bool) {
	st := Stack{}
	seen := map[int]bool{}
	for _, word := range strings.Fields(s) {
		n, err := strconv.Atoi(word)
		if err != nil || seen[n] {
			return nil, false
		}
		seen[n] = true
		st = append(st, n)
	}
	return st, true
}
