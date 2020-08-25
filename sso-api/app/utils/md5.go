package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

func GetMd5(password string) string {
	h := md5.New()
	h.Write([]byte("heng" + password))
	return fmt.Sprintf("%x", hex.EncodeToString(h.Sum(nil)))
}

//生成Guid字串
func GetUniqueToken() string {
	b := make([]byte, 64)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5(base64.URLEncoding.EncodeToString(b))
}
