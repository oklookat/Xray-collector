package v2raydata

import (
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
)

func MergeResolveInputResults(results []*ResolveInputResult) (*ResolveInputResult, error) {
	merged := &ResolveInputResult{
		Paths:     make([]string, 0),
		ListPaths: make(map[string][]int),
	}

	globalIndex := make(map[string]int) // path -> global index
	globalListUnique := make(map[string]map[int]struct{})

	for _, res := range results {
		for listName, indexes := range res.ListPaths {
			if _, ok := merged.ListPaths[listName]; !ok {
				merged.ListPaths[listName] = []int{}
				globalListUnique[listName] = map[int]struct{}{}
			}

			for _, localIdx := range indexes {
				if localIdx < 0 || localIdx >= len(res.Paths) {
					return nil, fmt.Errorf(
						"invalid ResolveInputResult: list %q contains index %d out of bounds (paths=%d)",
						listName, localIdx, len(res.Paths),
					)
				}

				path := res.Paths[localIdx]

				globalIdx, ok := globalIndex[path]
				if !ok {
					merged.Paths = append(merged.Paths, path)
					globalIdx = len(merged.Paths) - 1
					globalIndex[path] = globalIdx
				}

				if _, ok := globalListUnique[listName][globalIdx]; !ok {
					merged.ListPaths[listName] = append(
						merged.ListPaths[listName],
						globalIdx,
					)
					globalListUnique[listName][globalIdx] = struct{}{}
				}
			}
		}
	}

	return merged, nil
}

func isLocalDirectory(path string) (bool, error) {
	u, err := url.Parse(path)
	if err == nil && u.Scheme != "" {
		return false, nil
	}

	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	return info.IsDir(), nil
}
