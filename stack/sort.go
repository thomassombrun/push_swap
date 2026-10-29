package stack

// La fonction IsSortedAsc vérifie si la pile est triée par ordre croissant.
func (s Stack) IsSortedAsc() bool {
	for i := 1; i < len(s); i++ {
		if s[i-1] > s[i] {
			return false
		}
	}
	return true
}

// La fonction Sort renvoie la liste d'instructions qui trie la pile a.
func Sort(a Stack) []string {
	if a.IsSortedAsc() {
		return []string{}
	}
	if len(a) == 2 {
		return []string{"sa"}
	}
	if len(a) == 3 {
		return sortThree(a)
	}
	if len(a) > 10 {
		return radixSort(a)
	}
	return sortMany(a)
}

// La fonction sortThree renvoie les instructions qui trient une pile de 3 nombres.
func sortThree(a Stack) []string {
	x, y, z := a[0], a[1], a[2]
	switch {
	case x > y && y < z && x < z:
		return []string{"sa"}
	case x > y && y > z:
		return []string{"sa", "rra"}
	case x > y && y < z && x > z:
		return []string{"ra"}
	case x < y && y > z && x < z:
		return []string{"sa", "ra"}
	default:
		return []string{"rra"}
	}
}

// La fonction indexMinimum renvoie la position du plus petit nombre de la pile.
func indexMinimum(s Stack) int {
	minIdx := 0
	for i := 1; i < len(s); i++ {
		if s[i] < s[minIdx] {
			minIdx = i
		}
	}
	return minIdx
}

// La fonction moveToTop renvoie les instructions qui amènent l'élément à la position idx au sommet.
func moveToTop(size, idx int) []string {
	var moves []string
	if idx <= size/2 {
		for i := 0; i < idx; i++ {
			moves = append(moves, "ra")
		}
	} else {
		for i := 0; i < size-idx; i++ {
			moves = append(moves, "rra")
		}
	}
	return moves
}

// La fonction sortMany trie une pile de plus de 3 nombres en déplaçant le plus petit élément au sommet et de le pousser dans la pile B, puis en triant les 3 éléments restants et en ramenant les éléments de B dans A.
func sortMany(a Stack) []string {
	st := State{A: append(Stack{}, a...)}
	var moves []string

	for len(st.A) > 3 {
		for _, m := range moveToTop(len(st.A), indexMinimum(st.A)) {
			st.Exec(m)
			moves = append(moves, m)
		}
		st.Exec("pb")
		moves = append(moves, "pb")
	}

	if !st.A.IsSortedAsc() {
		moves = append(moves, sortThree(st.A)...)
	}

	for range st.B {
		st.Exec("pa")
		moves = append(moves, "pa")
	}
	return simplify(moves)

}

// La fonction simplify supprime les paires d'instructions consécutives qui s'annulent (ex : pb puis pa).
func simplify(moves []string) []string {
	res := []string{}
	for _, m := range moves {
		if n := len(res); n > 0 && annule(res[n-1], m) {
			res = res[:n-1]
			continue
		}
		res = append(res, m)
	}
	return res
}

// La fonction annule indique si deux instructions qui se suivent se compensent.
func annule(a, b string) bool {
	paires := map[string]string{
		"pa": "pb", "pb": "pa",
		"ra": "rra", "rra": "ra",
		"rb": "rrb", "rrb": "rb",
	}
	return paires[a] == b
}

// La fonction ranks remplace chaque valeur par son rang (0 = le plus petit).
func ranks(a Stack) Stack {
	r := make(Stack, len(a))
	for i, v := range a {
		rank := 0
		for _, w := range a {
			if w < v {
				rank++
			}
		}
		r[i] = rank
	}
	return r
}

// radixSort renvoie les instructions qui trient a avec un tri radix sur les rangs.
func radixSort(a Stack) []string {
	st := State{A: ranks(a)}
	n := len(a)
	var moves []string

	for bit := 0; (n-1)>>bit != 0; bit++ {
		for i := 0; i < n; i++ {
			move := "ra"
			if (st.A[0]>>bit)&1 == 0 {
				move = "pb"
			}
			st.Exec(move)
			moves = addMove(moves, move)
		}
		for len(st.B) > 0 {
			st.Exec("pa")
			moves = addMove(moves, "pa")
		}
	}
	return moves
}

// La fonction addMove ajoute m à la liste, sauf si m annule le dernier mouvement (pa/pb).
func addMove(moves []string, m string) []string {
	if n := len(moves); n > 0 {
		last := moves[n-1]
		if (m == "pa" && last == "pb") || (m == "pb" && last == "pa") {
			return moves[:n-1]
		}
	}
	return append(moves, m)
}
