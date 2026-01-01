package v2raydata

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"time"

	"github.com/oklookat/xray-collector/geosite/cache"
)

func newInputRemoteList(cache *cache.Geosite, rootUrl *url.URL) *inputRemoteList {
	return &inputRemoteList{
		cache:   cache,
		rootUrl: rootUrl,
		client: newHttpClient(
			"Xray-collector (github.com/oklookat/Xray-collector)",
			1.0,
			10*time.Second,
		),
	}
}

type inputRemoteList struct {
	cache   *cache.Geosite
	rootUrl *url.URL
	client  *httpClient
}

// Resolve downloads the root list and all transitive includes
// and returns a slice of cached paths (normalized, without duplicates).
func (r *inputRemoteList) Resolve(rootListName string) ([]string, error) {
	// Use absolute normalized paths as map keys to avoid duplicates
	resolved := make(map[string]string) // absPath -> cachedPath

	rootCachedPath, err := r.resolveSingle(rootListName, r.buildListURL(rootListName))
	if err != nil {
		return nil, err
	}

	rootCachedPath, _ = filepath.Abs(rootCachedPath)
	rootCachedPath = filepath.Clean(rootCachedPath)
	resolved[rootCachedPath] = rootCachedPath

	parseQueue := []string{rootCachedPath}

	for len(parseQueue) > 0 {
		currentPath := parseQueue[0]
		parseQueue = parseQueue[1:]

		parser, err := newIncludeParser(currentPath)
		if err != nil {
			return nil, fmt.Errorf("newIncludeParser: %w", err)
		}

		includes, err := parser.Parse()
		if err != nil {
			return nil, fmt.Errorf("include parse '%s': %w", currentPath, err)
		}

		for _, inc := range includes {
			incAbs, err := filepath.Abs(inc.Path)
			if err != nil {
				continue
			}
			incAbs = filepath.Clean(incAbs)

			// skip if already resolved
			if _, ok := resolved[incAbs]; ok {
				continue
			}

			listURL := r.buildListURL(inc.ListName)
			cachedPath, err := r.resolveSingle(inc.ListName, listURL)
			if err != nil {
				return nil, err
			}

			cachedPath, _ = filepath.Abs(cachedPath)
			cachedPath = filepath.Clean(cachedPath)
			resolved[cachedPath] = cachedPath

			// queue for parsing includes
			parseQueue = append(parseQueue, cachedPath)
		}
	}

	// convert map to slice
	paths := make([]string, 0, len(resolved))
	for _, p := range resolved {
		paths = append(paths, p)
	}

	return paths, nil
}

// Downloads a single list and updates the cache.
func (r *inputRemoteList) resolveSingle(listName, listURL string) (string, error) {
	slog.Info("Resolving list", "name", listName, "url", listURL)

	cachedPath, etag, err := r.cache.IsCached(listName)
	if err != nil {
		return "", fmt.Errorf("is cached '%s': %w", listName, err)
	}

	resp, err := r.client.Get(listURL, etag)
	if err != nil {
		return "", fmt.Errorf("GET '%s': %w", listURL, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotModified:
		if cachedPath != "" {
			slog.Info("Not modified", "list", listName)
			return cachedPath, nil
		}

		// force download if cache is empty
		resp.Body.Close()
		resp, err = r.client.Get(listURL, "")
		if err != nil {
			return "", fmt.Errorf("forced GET '%s': %w", listURL, err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("forced GET '%s': unexpected status %d", listURL, resp.StatusCode)
		}
		path, err := r.cache.Cache(listName, resp.Header.Get("ETag"), resp.Body)
		if err != nil {
			return "", fmt.Errorf("cache '%s': %w", listName, err)
		}
		return path, nil

	case http.StatusOK:
		path, err := r.cache.Cache(listName, resp.Header.Get("ETag"), resp.Body)
		if err != nil {
			return "", fmt.Errorf("cache '%s': %w", listName, err)
		}
		return path, nil

	default:
		return "", fmt.Errorf("GET '%s': unexpected status %d", listURL, resp.StatusCode)
	}
}

// Safely build URL for list, escaping the path
func (r *inputRemoteList) buildListURL(listName string) string {
	if r.rootUrl == nil {
		return ""
	}
	u := *r.rootUrl
	u.Path = path.Join(u.Path, url.PathEscape(listName))
	return u.String()
}
