package stack

// State contient les deux piles a et b
type State struct {
	A Stack
	B Stack
}

// La fonction Exec applique l'instruction donnée sur l'état des piles
// Si l'instruction n'est pas reconnue, la fonction retourne false
func (st *State) Exec(instr string) bool {
	switch instr {
	case "pa":
		if v, ok := st.B.Pop(); ok {
			st.A.Push(v)
		}
	case "pb":
		if v, ok := st.A.Pop(); ok {
			st.B.Push(v)
		}
	case "sa":
		st.A.Swap()
	case "sb":
		st.B.Swap()
	case "ss":
		st.A.Swap()
		st.B.Swap()
	case "ra":
		st.A.Rotate()
	case "rb":
		st.B.Rotate()
	case "rr":
		st.A.Rotate()
		st.B.Rotate()
	case "rra":
		st.A.ReverseRotate()
	case "rrb":
		st.B.ReverseRotate()
	case "rrr":
		st.A.ReverseRotate()
		st.B.ReverseRotate()
	default:
		return false
	}
	return true
}

// La fonction IsSorted vérifie si la pile A est triée et si la pile B est vide.
func (st *State) IsSorted() bool {
	return len(st.B) == 0 && st.A.IsSortedAsc()
}
