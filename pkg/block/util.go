package block

import "os"

func CreatePath(persistentStoragePath string) error {
	if _, err := os.Stat(persistentStoragePath); os.IsNotExist(err) {
		if err := os.MkdirAll(persistentStoragePath, os.FileMode(0755)); err != nil {
			return err
		}
	} else {
	}
	return nil
}
