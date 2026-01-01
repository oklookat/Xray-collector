package v2raydata

import (
	"bufio"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type includeResult struct {
	// Absolute path to included file.
	Path string
	// File name used in include.
	ListName string
	// Whether file existed at parse time.
	Found bool
}

type includeParser struct {
	rootFilePath, rootDir string
}

// Создаёт парсер для rootFile
func newIncludeParser(rootFile string) (*includeParser, error) {
	rootAbs, err := filepath.Abs(rootFile)
	if err != nil {
		return nil, err
	}
	rootAbs = filepath.Clean(rootAbs)
	return &includeParser{
		rootFilePath: rootAbs,
		rootDir:      filepath.Dir(rootAbs),
	}, nil
}

// Parse возвращает все includes reachable from rootFilePath, без дубликатов
func (p *includeParser) Parse() ([]includeResult, error) {
	visited := make(map[string]bool)
	var result []includeResult

	if err := p.parseFile(p.rootFilePath, visited, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (p *includeParser) parseFile(
	filePath string,
	visited map[string]bool,
	result *[]includeResult,
) error {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}
	absPath = filepath.Clean(absPath)

	// ограничение: include не выходит за пределы rootDir
	rel, err := filepath.Rel(p.rootDir, absPath)
	if err != nil {
		return err
	}
	if strings.HasPrefix(rel, ".."+string(filepath.Separator)) || rel == ".." {
		// файл за пределами rootDir → игнорируем
		return nil
	}

	if visited[absPath] {
		return nil
	}
	visited[absPath] = true

	listName := filepath.Base(absPath)

	file, err := os.Open(absPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			*result = append(*result, includeResult{
				Path:     absPath,
				ListName: listName,
				Found:    false,
			})
			return nil
		}
		return err
	}
	defer file.Close()

	*result = append(*result, includeResult{
		Path:     absPath,
		ListName: listName,
		Found:    true,
	})

	dir := filepath.Dir(absPath)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if !strings.HasPrefix(line, "include:") {
			continue
		}

		includeName := strings.TrimSpace(strings.TrimPrefix(line, "include:"))
		if includeName == "" {
			continue
		}

		includePath := filepath.Join(dir, includeName)
		includePath, err := filepath.Abs(includePath)
		if err != nil {
			continue
		}
		includePath = filepath.Clean(includePath)

		// рекурсивный вызов
		if err := p.parseFile(includePath, visited, result); err != nil {
			return err
		}
	}

	return scanner.Err()
}
