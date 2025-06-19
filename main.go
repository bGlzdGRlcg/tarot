package main

import (
	"bytes"
	"image"
	"log"
	"os"
	"strconv"
	"time"

	"tarot/tarot"

	"github.com/disintegration/imaging"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
)

const formation_reply = `塔罗牌阵：
0. 圣三角牌阵
1. 圣三角牌阵v2
2. 时间之流牌阵
3. 四要素牌阵
4. 五牌阵
5. 吉普赛十字阵
6. 马蹄牌阵
7. 六芒星牌阵
8. 平安扇牌阵
9. 沙迪若之星牌阵
请发送 /formation <数字> 来使用对应的牌阵。
`

func init() {
	godotenv.Load()
}

func main() {
	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/ping", func(c tele.Context) error {
		return c.Send("pong!")
	})

	b.Handle("/tarot", func(c tele.Context) error {
		card, is_down := tarot.Get_tarot()
		t := &tele.Photo{}
		if is_down == 0 {
			img, _ := os.Open("assets/" + card.Card_file + ".png")
			defer img.Close()
			t = &tele.Photo{File: tele.FromReader(img), Caption: "看看 " + c.Sender().FirstName + " 抽到了什么：" + card.Card_name + " 「正位」\n" + card.Card_up}
		}
		if is_down == 1 {
			img, _ := os.Open("assets/" + card.Card_file + ".png")
			defer img.Close()
			i, _, _ := image.Decode(img)
			rotated := imaging.Rotate180(i)
			var buf bytes.Buffer
			imaging.Encode(&buf, rotated, imaging.PNG)
			t = &tele.Photo{File: tele.FromReader(&buf), Caption: "看看 " + c.Sender().FirstName + " 抽到了什么：" + card.Card_name + " 「逆位」\n" + card.Card_down}
		}
		return c.SendAlbum(tele.Album{t})
	})

	b.Handle("/formation", func(c tele.Context) error {
		id, err := strconv.Atoi(c.Message().Payload)
		if err != nil || id < 0 || id >= 10 {
			return c.Send(formation_reply)
		}
		f := tarot.Formations[id]
		c.Send("启用" + f.Fname + "，少女祈祷中...")
		for i := 0; i < f.Fnum-1; i++ {
			card, is_down := tarot.Get_tarot()
			t := &tele.Photo{}
			num := strconv.Itoa(i + 1)
			frep := f.Frep[i] + "\n"
			if is_down == 0 {
				img, _ := os.Open("assets/" + card.Card_file + ".png")
				defer img.Close()
				t = &tele.Photo{File: tele.FromReader(img), Caption: frep + "第" + num + "张牌：\n" + card.Card_name + " 「正位」\n" + card.Card_up}
			}
			if is_down == 1 {
				img, _ := os.Open("assets/" + card.Card_file + ".png")
				defer img.Close()
				i, _, _ := image.Decode(img)
				rotated := imaging.Rotate180(i)
				var buf bytes.Buffer
				imaging.Encode(&buf, rotated, imaging.PNG)
				t = &tele.Photo{File: tele.FromReader(&buf), Caption: frep + "第" + num + "张牌：" + card.Card_name + " 「逆位」\n" + card.Card_down}
			}
			c.SendAlbum(tele.Album{t})
		}
		rpy := f.Frep[f.Fnum-1] + "\n第" + strconv.Itoa(f.Fnum) + "张牌："
		if f.IsCut {
			rpy = f.Frep[f.Fnum-1] + "\n切牌："
		}
		card, is_down := tarot.Get_tarot()
		t := &tele.Photo{}
		if is_down == 0 {
			img, _ := os.Open("assets/" + card.Card_file + ".png")
			defer img.Close()
			t = &tele.Photo{File: tele.FromReader(img), Caption: rpy + card.Card_name + " 「正位」\n" + card.Card_up}
		}
		if is_down == 1 {
			img, _ := os.Open("assets/" + card.Card_file + ".png")
			defer img.Close()
			i, _, _ := image.Decode(img)
			rotated := imaging.Rotate180(i)
			var buf bytes.Buffer
			imaging.Encode(&buf, rotated, imaging.PNG)
			t = &tele.Photo{File: tele.FromReader(&buf), Caption: rpy + card.Card_name + " 「逆位」\n" + card.Card_down}
		}
		return c.SendAlbum(tele.Album{t})
	})

	b.Start()
}
