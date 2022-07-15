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

var i18nTokenReplacer = strings.NewReplacer("\\n", "\n", "&nbsp;", "\u00A0")

//[section_foo] => {key_abc => lang_abc, key_def => lang_def, ...}
//[section_bar] => {key_xyz => lang_xyz, key_shk => lang_shk, ...}
var keys map[string]map[string]string

func Init() error {
	locale := GetCurrentLocale()
	//todo: check if already inited to same locale and skip

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
		sectionPath := filepath.Join(sectionsFilenamesPath, sectionFilename)

		keys[sectionName], err = loadSection(sectionPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func loadSection(sectionFilepath string) (map[string]string, error) {
	sectionKeys := make(map[string]string)

	file, err := os.Open(sectionFilepath)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key, value, foundKey := strings.Cut(scanner.Text(), ":")
		if !foundKey {
			continue
		}
		sectionKeys[key] = i18nTokenReplacer.Replace(strings.TrimSpace(value))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if err = file.Close(); err != nil {
		return nil, err
	}

	return sectionKeys, nil
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
