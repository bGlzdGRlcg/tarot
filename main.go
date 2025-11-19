package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tarot "tarot/lib"

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

	url := os.Getenv("URL")

	if url == "" {
		url = "https://tarot.listder.xyz/"
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/ping", func(c tele.Context) error {
		return c.Send("pong!")
	})

	b.Handle(tele.OnText, func(c tele.Context) error {
		if strings.Contains(c.Text(), "Ciallo") || strings.Contains(c.Text(), "ciallo") {
			return c.Reply("柚子厨真恶心")
		}
		return nil
	})

	b.Handle("/tarot", func(c tele.Context) error {
		card, is_down := tarot.Get_tarot()
		t := &tele.Photo{}
		if is_down == 0 {
			t = &tele.Photo{File: tele.FromURL(url + card.Card_file + ".jpg"), Caption: "看看 " + c.Sender().FirstName + " 抽到了什么：\n" + card.Card_name + " 「正位」\n" + card.Card_up}
		}
		if is_down == 1 {
			t = &tele.Photo{File: tele.FromURL(url + "_" + card.Card_file + ".jpg"), Caption: "看看 " + c.Sender().FirstName + " 抽到了什么：\n" + card.Card_name + " 「逆位」\n" + card.Card_down}
		}
		return c.SendAlbum(tele.Album{t})
	})

	b.Handle("/formation", func(c tele.Context) error {
		if c.Chat().Type != tele.ChatPrivate {
			return c.Send("请在私聊中使用此命令。")
		}

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
				t = &tele.Photo{File: tele.FromURL(url + card.Card_file + ".jpg"), Caption: frep + "第" + num + "张牌：" + card.Card_name + " 「正位」\n" + card.Card_up}
			}
			if is_down == 1 {
				t = &tele.Photo{File: tele.FromURL(url + "_" + card.Card_file + ".jpg"), Caption: frep + "第" + num + "张牌：" + card.Card_name + " 「逆位」\n" + card.Card_down}
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
			t = &tele.Photo{File: tele.FromURL(url + card.Card_file + ".jpg"), Caption: rpy + card.Card_name + " 「正位」\n" + card.Card_up}
		}
		if is_down == 1 {
			t = &tele.Photo{File: tele.FromURL(url + "_" + card.Card_file + ".jpg"), Caption: rpy + card.Card_name + " 「逆位」\n" + card.Card_down}
		}
		return c.SendAlbum(tele.Album{t})
	})

	b.Handle(tele.OnQuery, func(c tele.Context) error {
		bg := "https://tarot.listder.xyz/Extra/BG.jpg"

		results := make(tele.Results, 1)

		results[0] = &tele.ArticleResult{
			Title:    "Tarot",
			Text:     "点击按钮即可占卜",
			URL:      url,
			ThumbURL: bg,
		}

		results[0].SetResultID("0")

		results[0].SetReplyMarkup(&tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{
				{
					tele.InlineButton{
						Text: "少女祈祷中...",
						Data: "tarot",
					},
				},
			},
		})

		return c.Answer(&tele.QueryResponse{
			Results:   results,
			CacheTime: 1,
		})
	})

	b.Handle(tele.OnInlineResult, func(c tele.Context) error {
		desc := tarot.Get_tarot_inline(url)
		return c.EditCaption(desc)
	})

	b.Start()
}
