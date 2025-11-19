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

// 随机获取一张塔罗牌，返回他的图片链接和描述
func Get_tarot_inline(url string) string {
	card, is_down := Get_tarot()
	var desc string

	if is_down == 0 {
		desc = url + card.Card_file + ".jpg\n" + card.Card_name + "「正位」\n" + card.Card_up
	} else {
		desc = url + "_" + card.Card_file + ".jpg\n" + card.Card_name + "「逆位」\n" + card.Card_down
	}
	return desc
}
