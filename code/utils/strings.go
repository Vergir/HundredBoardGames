package utils

import (
	"strings"
)

func FormComplexFilename(fileName string, stringWithExtension string) string {
	extensionIndex := strings.LastIndex(stringWithExtension, ".") + len(".")

	return fileName + "." + stringWithExtension[extensionIndex:]
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
