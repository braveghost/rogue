package file

import (
	"errors"
	"os"
	"path"
)

// 判断所给路径是否为文件夹
func IsDir(pt string) bool {
	s, err := os.Stat(pt)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func DirNotExistCreate(pt string) bool {
	if !PathExist(pt) {
		err := os.Mkdir(pt, os.ModePerm)
		if err == nil {
			return true
		}
	}
	if IsDir(pt) {
		return true
	}

	return false
}

// 验证目录是否存在
func DirVerify(dirName string) (string, error) {
	currentDir, _ := os.Getwd()
	tmpPath := path.Join(currentDir, dirName)
	tmpFileInfo, err := os.Stat(tmpPath)
	if err == nil {
		if ! tmpFileInfo.IsDir() {
			err = errors.New("File of the same name")
		}
	} else {
		err = os.Mkdir(tmpPath, 0755)
	}
	return tmpPath, err
}
