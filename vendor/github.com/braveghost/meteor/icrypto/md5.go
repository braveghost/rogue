package icrypto

import (
	"crypto/md5"
	"fmt"
	"io"
)

func BytesToMd5(b []byte) string {
	has := md5.Sum(b)
	return fmt.Sprintf("%x", has)
}
func StrToMd5(str string) string {
	has := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", has)
}

func SliceStrToMd5(li []string) string {
	hash := md5.New()
	for _, l := range li {
		io.WriteString(hash, l)
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}
