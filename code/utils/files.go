package utils

import "os"

func ListFolderFiles(path string) ([]string, error) {
	folder, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer folder.Close()

	filesNames, err := folder.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	return filesNames, nil
}
