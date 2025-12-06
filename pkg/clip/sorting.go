package clip

import (
	"os"
	"sort"
	"strings"
)

// SortByName implements sort.Interface for sorting os.DirEntry by name.
type SortByName []os.DirEntry

func (s SortByName) Len() int {
	return len(s)
}

func (s SortByName) Less(i, j int) bool {
	return strings.Compare(s[i].Name(), s[j].Name()) == -1
}

func (s SortByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

var _ sort.Interface = (*SortByName)(nil)
