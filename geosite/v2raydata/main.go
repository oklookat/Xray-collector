package v2raydata

import (
	"fmt"

	"github.com/oklookat/xray-collector/geosite/cache"
	"github.com/oklookat/xray-collector/geosite/v2raydata/compiler"
)

type ResolveInputResult struct {
	// Paths to all lists and includes.
	Paths []string
	// [listName]Paths indexes of list and his includes.
	ListPaths map[string][]int
}

func ResolveInput(cfg *InputConfig, cache *cache.Geosite) (*ResolveInputResult, error) {
	isLocal, err := isLocalDirectory(cfg.URI)
	if err != nil {
		return nil, fmt.Errorf("isLocalDirectory: %w", err)
	}
	if isLocal {
		return newInputLocal(cfg).Resolve()
	}

	remote, err := NewInputRemote(cfg, cache)
	if err != nil {
		return nil, fmt.Errorf("NewInputRemote: %w", err)
	}
	return remote.Resolve()
}

func ResolveOutput(cfg *OutputConfig, inputResult *ResolveInputResult) error {
	var paths []string

	if len(cfg.Lists) > 0 {
		for _, wantedList := range cfg.Lists {
			listPaths, ok := inputResult.ListPaths[wantedList]
			if !ok {
				return fmt.Errorf("list '%s' not found in ResolveInputResult", wantedList)
			}
			for _, idx := range listPaths {
				paths = append(paths, inputResult.Paths[idx])
			}
		}
	} else {
		paths = inputResult.Paths
	}

	return compiler.Compile(paths, cfg.Path)
}
