package main

type SZList struct {
	List []int
}

func (l *SZList) Push(item int) bool {
	if !l.Exists(item) {
		l.List = append(l.List, item)
		return true
	}
	return false
}

func (l SZList) Exists(item int) bool {
	for _, i := range l.List {
		if item == i {
			return true
		}
	}
	return false
}
