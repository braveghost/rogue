package file

import (
	"os"
)

// 判断所给路径文件/文件夹是否存在
func PathExist(pt string) bool {
	_, err := os.Stat(pt) //os.Stat获取文件信息
	if err != nil {
		if !os.IsExist(err) {
			return false
		}
	}
	return true
}
