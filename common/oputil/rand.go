package oputil

import (
	"crypto/rand"
	"io"
)

// GenerateNonce 生成指定长度的随机数
func GenerateNonce(len int) string {
	randBytes := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

	nonce := make([]byte, len)
	io.ReadFull(rand.Reader, nonce)

	for i, v := range nonce {
		nonce[i] = randBytes[v%10]
	}

	return string(nonce)
}
