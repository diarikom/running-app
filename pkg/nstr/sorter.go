package nstr

import (
	"sort"
)

type StringSorter struct {
	slice sort.StringSlice
}

func NewStringSorter(coll []string) StringSorter {
	// Convert to StringSlice
	slice := sort.StringSlice(coll)

	// Sort string
	slice.Sort()

	return StringSorter{
		slice: slice,
	}
}

func (s StringSorter) AtIndex(needle string) int {
	// Get index
	idx := s.slice.Search(needle)
	// If found return index
	if s.slice[idx] == needle {
		return idx
	}
	// If not found return -1
	return -1
}

func (s StringSorter) Contains(needle string) bool {
	// Get search
	idx := s.slice.Search(needle)

	// Check index
	if idx > len(s.slice)-1 {
		return false
	}

	// Compare
	return s.slice[idx] == needle
}
