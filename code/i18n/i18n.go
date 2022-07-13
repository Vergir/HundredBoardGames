package i18n

import (
	"bufio"
	"hundred-board-games/code/utils"
	"os"
	"path/filepath"
	"strings"
)

type locale string

const (
	LOCALE_EN_GB locale = "en-gb"
	LOCALE_RU_RU locale = "ru-ru"
)

//[section_foo] => {key_abc => lang_abc, key_def => lang_def, ...}
//[section_bar] => {key_xyz => lang_xyz, key_shk => lang_shk, ...}
var keys map[string]map[string]string

func Init() error {
	locale := GetCurrentLocale()
	//todo: check if already inited to same locald and skip

	err := LoadLocale(locale)
	if err != nil {
		return err
	}

	return nil
}

func LoadLocale(locale locale) error {
	keys = make(map[string]map[string]string)

	sectionsFilenamesPath := filepath.Join("langs", string(locale))
	sectionsFilenames, err := utils.ListFolderFiles(sectionsFilenamesPath)
	if err != nil {
		return err
	}

	for _, sectionFilename := range sectionsFilenames {
		sectionName, _, _ := strings.Cut(sectionFilename, ".")
		keys[sectionName] = make(map[string]string)
		sectionPath := filepath.Join(sectionsFilenamesPath, sectionFilename)
		file, err := os.Open(sectionPath)
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			key, value, foundSomething := strings.Cut(scanner.Text(), ":")
			if !foundSomething {
				continue
			}
			keys[sectionName][key] = strings.ReplaceAll(strings.TrimSpace(value), "\\n", "\n")
		}
		if err := scanner.Err(); err != nil {
			return err
		}
	}

	return nil
}

func GetCurrentLocale() locale {
	//todo: cookie or guess
	return LOCALE_EN_GB
}

func Get(section string, key string) string {
	fallback := key

	if len(keys[section]) == 0 {
		return fallback
	}

	lang := keys[section][key]
	if lang == "" {
		return fallback
	}

	return lang
}

func Expand(label string, tokensExpansions map[string]string) string {
	newLabel := label
	for token, expansion := range tokensExpansions {
		newLabel = strings.ReplaceAll(newLabel, "{"+token+"}", expansion)
	}

	return newLabel
}

func GetSection(section string) map[string]string {
	return keys[section]
}
