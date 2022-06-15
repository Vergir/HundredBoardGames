package utils

import (
	"strings"
)

func FormFullFilename(entityName string, extensionString string) string {
	extensionIndex := strings.LastIndex(extensionString, ".") + len(".")

	return entityName + "." + extensionString[extensionIndex:]
}

func StaticJs(filename string) string {
	return "static/js/" + filename
}

func AnyStringHasSubstring(ss []string, substr string) bool {
	for _, s := range ss {
		if strings.Contains(s, substr) {
			return true
		}
	}

	return false
}
