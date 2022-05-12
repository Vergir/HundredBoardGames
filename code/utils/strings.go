package utils

import (
	"strconv"
	"strings"
)

func FormFullFilename(entityId int, extensionString string) string {
	extensionIndex := strings.LastIndex(extensionString, ".") + len(".")

	return strconv.Itoa(int(entityId)) + "." + extensionString[extensionIndex:]
}
