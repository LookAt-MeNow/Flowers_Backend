package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

//小写
func Md5Code(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	tempStr := h.Sum(nil)
	return hex.EncodeToString(tempStr)
}

//大写
func MD5Code(data string) string {
	return strings.ToUpper(Md5Code(data))
}

//加密
func MakePassword(plainpwd,salt string) string {
	return Md5Code(plainpwd + salt)
}

//解密
func ValidPassword(plainpwd,salt string,password string) bool {
	return Md5Code(plainpwd + salt) == password
}