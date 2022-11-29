package utils

import (
	"crypto/md5"
	"fmt"
)

func Hash(args ...string) string {
	var str string
	for _, s := range args {
		str += s
	}
	value := md5.Sum([]byte(str))
	rs := []rune(fmt.Sprintf("%x", value))
	return string(rs)
}
