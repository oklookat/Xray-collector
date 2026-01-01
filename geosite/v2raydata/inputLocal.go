package v2raydata

import (
	"fmt"
	"path/filepath"
)

func newInputLocal(cfg *InputConfig) *inputLocal {
	return &inputLocal{
		cfg: cfg,
	}
}

type inputLocal struct {
	cfg *InputConfig
}

func normalizePath(p string) (string, error) {
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	return filepath.Clean(abs), nil
}

func (i inputLocal) Resolve() (*ResolveInputResult, error) {
	result := &ResolveInputResult{
		Paths:     []string{},
		ListPaths: map[string][]int{},
	}

	unique := make(map[string]int)
	listUnique := map[string]map[int]struct{}{}

	for _, listName := range i.cfg.Lists {
		listUnique[listName] = map[int]struct{}{}

		rootRaw := filepath.Join(i.cfg.URI, listName)
		rootPath, err := normalizePath(rootRaw)
		if err != nil {
			return nil, err
		}

		idx, ok := unique[rootPath]
		if !ok {
			result.Paths = append(result.Paths, rootPath)
			idx = len(result.Paths) - 1
			unique[rootPath] = idx
		}

		if _, ok := listUnique[listName][idx]; !ok {
			result.ListPaths[listName] = append(result.ListPaths[listName], idx)
			listUnique[listName][idx] = struct{}{}
		}

		parser, err := newIncludeParser(rootPath)
		if err != nil {
			return nil, fmt.Errorf("newIncludeParser: %w", err)
		}
		paths, err := parser.Parse()
		if err != nil {
			return nil, err
		}

		for _, p := range paths {
			if !p.Found || p.Path == "" {
				continue
			}

			normPath, err := normalizePath(p.Path)
			if err != nil {
				return nil, err
			}

			idx, ok := unique[normPath]
			if !ok {
				result.Paths = append(result.Paths, normPath)
				idx = len(result.Paths) - 1
				unique[normPath] = idx
			}

			if _, ok := listUnique[listName][idx]; !ok {
				result.ListPaths[listName] = append(result.ListPaths[listName], idx)
				listUnique[listName][idx] = struct{}{}
			}
		}
	}

	return result, nil
}
