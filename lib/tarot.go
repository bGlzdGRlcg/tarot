package lib

import (
	"crypto/rand"
	"math/big"
)

// 随机获取一张塔罗牌
func Get_tarot() (Card, int) {
	irand, _ := rand.Int(rand.Reader, big.NewInt(78))
	is_down, _ := rand.Int(rand.Reader, big.NewInt(2))
	return Cards[irand.Int64()], int(is_down.Int64())
}
