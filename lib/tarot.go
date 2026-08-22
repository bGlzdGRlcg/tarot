package lib

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// 随机获取一张塔罗牌
func Get_tarot() (Card, int, error) {
	iRand, err := rand.Int(rand.Reader, big.NewInt(int64(len(Cards))))
	if err != nil {
		return Card{}, 0, fmt.Errorf("generate tarot card index: %w", err)
	}

	isDown, err := rand.Int(rand.Reader, big.NewInt(2))
	if err != nil {
		return Card{}, 0, fmt.Errorf("generate tarot orientation: %w", err)
	}

	return Cards[iRand.Int64()], int(isDown.Int64()), nil
}
