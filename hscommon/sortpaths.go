package hscommon

import (
	"sort"
	"strings"
)

// SortPaths sorts paths
func SortPaths(rootPath *PathEntry) {
	sort.Slice(rootPath.Children, func(i, j int) bool {
		if ((len(rootPath.Children[i].Children) == 0) && (len(rootPath.Children[j].Children) == 0)) ||
			((len(rootPath.Children[i].Children) != 0) && (len(rootPath.Children[j].Children) != 0)) {
			return strings.ToLower(rootPath.Children[i].Name) < strings.ToLower(rootPath.Children[j].Name)
		}

		return len(rootPath.Children[i].Children) > len(rootPath.Children[j].Children)
	})
}
