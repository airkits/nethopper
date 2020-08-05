package md5

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
)

//Encrypt 计算hash值
func Encrypt(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	md5Hash := hex.EncodeToString(h.Sum(nil))
	return md5Hash
}

//Hex2Int HEX string convert to int
func Hex2Int(s string) uint32 {
	n, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return 0
	}
	return uint32(n)
}

// HashMod 获取hash取模
func HashMod(str string, mod uint32) uint32 {
	result := Encrypt(str)
	rs := []rune(result)
	l := len(rs)
	hashStr := string(rs[l-7 : l])
	return Hex2Int(hashStr) % mod
}
