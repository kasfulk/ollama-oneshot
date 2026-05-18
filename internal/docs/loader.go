package docs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Load(files []string) (string, error) {
	var entries []string
	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: cannot read %s: %v\n", f, err)
			continue
		}
		entries = append(entries, FormatDoc(f, string(content)))
	}
	return WrapDocumentation(entries), nil
}

func LoadGlob(patterns []string) (string, error) {
	var entries []string
	for _, pattern := range patterns {
		matches, err := filepath.Glob(strings.TrimSpace(pattern))
		if err != nil {
			continue
		}
		for _, match := range matches {
			content, err := os.ReadFile(match)
			if err != nil {
				continue
			}
			entries = append(entries, FormatDoc(match, string(content)))
		}
	}
	return WrapDocumentation(entries), nil
}