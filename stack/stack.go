package stack

type Stack []int

// La fonction Push ajoute un élément au sommet de la pile
func (s *Stack) Push(v int) {
	*s = append(Stack{v}, *s...)
}

// La fonction Pop retire et retourne l'élément au sommet de la pile
func (s *Stack) Pop() (int, bool) {
	if len(*s) == 0 {
		return 0, false
	}
	v := (*s)[0]
	*s = (*s)[1:]
	return v, true
}

// La fonction Swap échange les deux éléments au sommet de la pile
func (s *Stack) Swap() {
	if len(*s) < 2 {
		return
	}
	(*s)[0], (*s)[1] = (*s)[1], (*s)[0]
}

// La fonction Rotate déplace l'élément au sommet de la pile vers le bas
func (s *Stack) Rotate() {
	if len(*s) < 2 {
		return
	}
	v := (*s)[0]
	*s = append((*s)[1:], v)
}

// La fonction ReverseRotate déplace l'élément en bas de la pile vers le sommet
func (s *Stack) ReverseRotate() {
	if len(*s) < 2 {
		return
	}
	v := (*s)[len(*s)-1]
	*s = append(Stack{v}, (*s)[:len(*s)-1]...)
}
