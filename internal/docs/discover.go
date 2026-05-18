package docs

import "os"

var autoDiscoverFiles = []string{
	"AGENT.md",
	"PRD.md",
	"DATABASE.md",
	"ARCHITECTURE.md",
	"TASKS.md",
}

func Discover(dir string) ([]string, error) {
	var found []string
	for _, name := range autoDiscoverFiles {
		path := dir + "/" + name
		if _, err := os.Stat(path); err == nil {
			found = append(found, path)
		}
	}
	return found, nil
}