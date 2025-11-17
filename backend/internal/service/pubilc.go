// 此处存放一些公共函数
package service

import (
	"crypto/rand"
	"log"
	"math/big"
)

func GetHash(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	buf := make([]byte, n)
	for i := range buf {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			// 出错就退回一个固定值，避免 panic
			log.Println("rand.Int err:", err)
			return "DefaultRandomSecretKey123"
		}
		buf[i] = letters[num.Int64()]
	}
	return string(buf)
}
