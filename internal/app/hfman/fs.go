package hfman

import (
	"os"
	"path"
)

func CreateDirs(root string, folderPerm os.FileMode) (err error) {
	if err = os.MkdirAll(root, folderPerm); err != nil {
		return err
	}
	if err = os.MkdirAll(path.Join(root, "download"), folderPerm); err != nil {
		return err
	}
	if err = os.MkdirAll(path.Join(root, "upload"), folderPerm); err != nil {
		return err
	}
	return nil
}
