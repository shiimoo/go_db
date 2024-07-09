package base

type List []any

func (s *List) Add(ele any) {
	*s = append(*s, ele)
}
func (s *List) Pop() any {
	if len(*s) == 0 {
		return nil
	}
	ele := (*s)[0]
	*s = (*s)[1:]
	return ele
}

func NewList() List {
	return make(List, 0)
}
