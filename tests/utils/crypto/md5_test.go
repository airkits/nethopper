package crypto_test

import (
	"testing"

	"github.com/airkits/nethopper/utils/crypto/md5"
)

func TestMd5Encrypt(t *testing.T) {
	str := "gametest"
	newStr := md5.Encrypt(str)
	if newStr != "b77fcd439f9818b58154a167d0395156" {
		t.Error("encrypt md5 failed")
	}

}
func TestHex2Int(t *testing.T) {
	result := md5.Hex2Int("aaaaaaaaaaaa")
	if result != 0 {
		t.Error("hex2int convert error")
	}
	result = md5.Hex2Int("aaaa")
	if result != 43690 {
		t.Error("hex2int convert error")
	}
}
func TestMd5HashMod(t *testing.T) {
	str := "gametest"
	result := md5.HashMod(str, 8)
	if result != 6 {
		t.Error("HashValue failed")
	}

}
