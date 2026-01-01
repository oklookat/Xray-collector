package v2raydata

import (
	"fmt"
	"net/url"

	"github.com/oklookat/xray-collector/geosite/cache"
)

func NewInputRemote(cfg *InputConfig, cache *cache.Geosite) (*InputRemote, error) {
	rem := &InputRemote{
		cache: cache,
		cfg:   cfg,
	}
	rootURL, err := rem.parseURL(cfg.URI)
	if err != nil {
		return nil, fmt.Errorf("parseURL: %w", err)
	}
	rem.rootURL = rootURL
	return rem, err
}

type InputRemote struct {
	cache   *cache.Geosite
	cfg     *InputConfig
	rootURL *url.URL
}

// Returns paths to all resolved lists.
func (r InputRemote) Resolve() (*ResolveInputResult, error) {
	result := &ResolveInputResult{
		Paths:     []string{},
		ListPaths: map[string][]int{},
	}

	// path -> index in result.Paths
	unique := make(map[string]int)
	// listName -> set of indices
	listUnique := make(map[string]map[int]struct{})

	for _, listName := range r.cfg.Lists {
		if _, ok := listUnique[listName]; !ok {
			listUnique[listName] = map[int]struct{}{}
		}

		remList := newInputRemoteList(r.cache, r.rootURL)

		paths, err := remList.Resolve(listName)
		if err != nil {
			return nil, fmt.Errorf("remote list resolve '%s': %w", listName, err)
		}

		for _, p := range paths {
			if p == "" {
				continue
			}

			normPath, err := normalizePath(p)
			if err != nil {
				return nil, fmt.Errorf("normalize path '%s': %w", p, err)
			}

			idx, exists := unique[normPath]
			if !exists {
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

func (r InputRemote) parseURL(path string) (*url.URL, error) {
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	switch u.Scheme {
	case "http", "https":
		return u, err
	default:
		return nil, fmt.Errorf("'%s' not an remote URL", u.String())
	}
}
