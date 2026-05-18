package docs

import "fmt"

func FormatDoc(filePath, content string) string {
	return fmt.Sprintf("[FILE: %s]\n%s", filePath, content)
}

func WrapDocumentation(entries []string) string {
	if len(entries) == 0 {
		return ""
	}
	result := "<PROJECT_DOCUMENTATION>\n\n"
	for _, entry := range entries {
		result += entry + "\n\n"
	}
	result += "</PROJECT_DOCUMENTATION>"
	return result
}