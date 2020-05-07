package md5

import (
	"crypto/md5"
	"encoding/hex"
)

//Encrypt 计算hash值
func Encrypt(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	md5Hash := hex.EncodeToString(h.Sum(nil))
	return md5Hash
}
