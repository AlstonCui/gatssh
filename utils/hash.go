package utils

import (
	"crypto/md5"
	"fmt"
)

func Md5Sum(str string) string {
	tmpHash := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", tmpHash)
}