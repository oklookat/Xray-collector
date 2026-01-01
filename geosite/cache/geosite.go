package cache

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

const (
	_defaultCacheDirPerm  = 0755
	_defaultCacheFilePerm = 0644
)

// maxSizeBytes == 0: unlimited
func NewGeosite(cacheDir string, maxSizeBytes int64) (*Geosite, error) {
	if err := os.MkdirAll(cacheDir, _defaultCacheDirPerm); err != nil {
		return nil, fmt.Errorf("create cache dir: %w", err)
	}

	c := &Geosite{
		cacheDir:     cacheDir,
		maxSizeBytes: maxSizeBytes,
	}

	if err := c.cleanupIfOversize(); err != nil {
		return nil, err
	}

	return c, nil
}

type Geosite struct {
	cacheDir     string
	maxSizeBytes int64
	mu           sync.Mutex
}

func (c *Geosite) Directory() string {
	return c.cacheDir
}

func (c *Geosite) Cache(listName, etag string, data io.Reader) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	listPath := c.getListPath(listName)
	etagPath := c.getETagPath(listName)

	tmpList := listPath + ".tmp"
	listFile, err := os.OpenFile(tmpList, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, _defaultCacheFilePerm)
	if err != nil {
		return "", fmt.Errorf("open tmp list: %w", err)
	}

	if _, err := io.Copy(listFile, data); err != nil {
		listFile.Close()
		_ = os.Remove(tmpList)
		return "", fmt.Errorf("write tmp list: %w", err)
	}
	if err := listFile.Close(); err != nil {
		_ = os.Remove(tmpList)
		return "", fmt.Errorf("close tmp list: %w", err)
	}

	if err := os.Rename(tmpList, listPath); err != nil {
		_ = os.Remove(tmpList)
		return "", fmt.Errorf("rename tmp list: %w", err)
	}

	// etag
	tmpETag := etagPath + ".tmp"
	if err := os.WriteFile(tmpETag, []byte(etag), _defaultCacheFilePerm); err != nil {
		_ = os.Remove(tmpETag)
		return "", fmt.Errorf("write tmp etag: %w", err)
	}
	if err := os.Rename(tmpETag, etagPath); err != nil {
		_ = os.Remove(tmpETag)
		return "", fmt.Errorf("rename tmp etag: %w", err)
	}

	if err := c.cleanupIfOversize(); err != nil {
		return "", err
	}

	return listPath, nil
}

func (c *Geosite) IsCached(listName string) (string, string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	listPath := c.getListPath(listName)
	etagPath := c.getETagPath(listName)

	if _, err := os.Stat(listPath); err != nil {
		if os.IsNotExist(err) {
			return "", "", nil
		}
		return "", "", err
	}

	etagBytes, err := os.ReadFile(etagPath)
	if err != nil {
		if os.IsNotExist(err) {
			_ = c.removePair(listName)
			return "", "", nil
		}
		return "", "", err
	}

	return listPath, string(etagBytes), nil
}

func (c *Geosite) removePair(listName string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	listPath := c.getListPath(listName)
	etagPath := c.getETagPath(listName)

	var err1, err2 error
	if err1 = os.Remove(listPath); err1 != nil && !os.IsNotExist(err1) {
		err1 = fmt.Errorf("remove list: %w", err1)
	} else {
		err1 = nil
	}
	if err2 = os.Remove(etagPath); err2 != nil && !os.IsNotExist(err2) {
		err2 = fmt.Errorf("remove etag: %w", err2)
	} else {
		err2 = nil
	}

	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}

func (c *Geosite) cleanupIfOversize() error {
	if c.maxSizeBytes <= 0 {
		return nil
	}

	var totalSize int64
	files := []fs.FileInfo{}

	err := filepath.Walk(c.cacheDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		totalSize += info.Size()
		files = append(files, info)
		return nil
	})
	if err != nil {
		return err
	}

	if totalSize <= c.maxSizeBytes {
		return nil
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Before(files[j].ModTime())
	})

	for _, f := range files {
		path := filepath.Join(c.cacheDir, f.Name())
		_ = os.Remove(path)
		totalSize -= f.Size()
		if totalSize <= c.maxSizeBytes {
			break
		}
	}

	return nil
}

func (c *Geosite) getListPath(listName string) string {
	return filepath.Join(c.cacheDir, filepath.Base(listName))
}

func (c *Geosite) getETagPath(listName string) string {
	return filepath.Join(c.cacheDir, filepath.Base(listName)+".etag")
}
