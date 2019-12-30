package file

// 判断所给路径是否为文件
func IsFile(pt string) bool {
	return !IsDir(pt)
}
