package util

import (
	"sort"

	"github.com/golang/geo/s2"
)

type s2List []s2.CellID

func (s2l s2List) Len() int {
	return len(s2l)
}

func (s2l s2List) Less(i int, j int) bool {
	return s2l[i].Level() < s2l[j].Level()
}

func (s2l s2List) Swap(i int, j int) {
	s2l[i], s2l[j] = s2l[j], s2l[i]
}

func toS2List(ints []uint64) s2List {
	lst := s2List{}
	for _, id := range ints {
		lst = append(lst, s2.CellID(id))
	}
	return lst
}

func ContainsOverlappingS2IDs(ids []uint64) bool {
	lst := toS2List(ids)

	sort.Sort(lst)

	length := lst.Len()

	for i := 0; i < length; i++ {
		higher := lst[i]
		for j := i + 1; j < length; j++ {
			lower := lst[j]
			if higher.Level() < lower.Level() && higher.Contains(lower) {
				return true
			}
		}
	}
	return false
}
