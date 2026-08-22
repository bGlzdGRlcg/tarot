package main

import (
	"errors"
	"fmt"
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

func normalizeAssetURL(rawURL string) string {
	if rawURL == "" {
		rawURL = "https://tarot.listder.xyz/"
	}
	return strings.TrimRight(rawURL, "/") + "/"
}

func senderName(c tele.Context) string {
	if message := c.Message(); message != nil && message.GuestUser != nil && message.GuestUser.FirstName != "" {
		return message.GuestUser.FirstName
	}
	if sender := c.Sender(); sender != nil && sender.FirstName != "" {
		return sender.FirstName
	}
	return "你"
}

func tarotPhoto(assetURL string, card tarot.Card, isDown int, prefix string) *tele.Photo {
	position := "正位"
	meaning := card.Card_up
	file := card.Card_file
	if isDown == 1 {
		position = "逆位"
		meaning = card.Card_down
		file = "_" + file
	}

	return &tele.Photo{
		File:    tele.FromURL(assetURL + file + ".jpg"),
		Caption: prefix + card.Card_name + " 「" + position + "」\n" + meaning,
	}
}

func tarotPhotoResult(assetURL string, card tarot.Card, isDown int, prefix string) *tele.PhotoResult {
	photo := tarotPhoto(assetURL, card, isDown, prefix)
	photoURL := photo.File.FileURL
	result := &tele.PhotoResult{
		URL:         photoURL,
		Title:       "塔罗牌",
		Description: card.Card_name,
		Caption:     photo.Caption,
		ThumbURL:    photoURL,
	}
	result.SetResultID(fmt.Sprintf("tarot-%d-%d", card.Card_id, isDown))
	return result
}

func drawUniqueTarot(used map[int]struct{}) (tarot.Card, int, error) {
	for len(used) < len(tarot.Cards) {
		card, isDown, err := tarot.Get_tarot()
		if err != nil {
			return tarot.Card{}, 0, err
		}
		if _, exists := used[card.Card_id]; exists {
			continue
		}
		used[card.Card_id] = struct{}{}
		return card, isDown, nil
	}
	return tarot.Card{}, 0, fmt.Errorf("no tarot cards left to draw")
}

func init() {
	godotenv.Load()
}

func main() {
	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	assetURL := normalizeAssetURL(os.Getenv("URL"))

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
		if strings.Contains(c.Text(), "何意味") {
			return c.Reply("意味何")
		}
		if strings.Contains(c.Text(), "意味何") {
			return c.Reply("何意味")
		}
		return nil
	})

	b.Handle("/tarot", func(c tele.Context) error {
		card, isDown, err := tarot.Get_tarot()
		if err != nil {
			return err
		}
		t := tarotPhoto(assetURL, card, isDown, "看看 "+senderName(c)+" 抽到了什么：\n")
		return c.SendAlbum(tele.Album{t})
	})

	b.Handle("/formation", func(c tele.Context) error {
		chat := c.Chat()
		if chat == nil || chat.Type != tele.ChatPrivate {
			return c.Send("请在私聊中使用此命令。")
		}

		id, err := strconv.Atoi(c.Message().Payload)
		if err != nil || id < 0 || id >= 10 {
			return c.Send(formation_reply)
		}
		f := tarot.Formations[id]
		if err := c.Send("启用" + f.Fname + "，少女祈祷中..."); err != nil {
			return err
		}
		used := make(map[int]struct{}, f.Fnum)
		for i := 0; i < f.Fnum; i++ {
			card, isDown, err := drawUniqueTarot(used)
			if err != nil {
				return err
			}
			num := strconv.Itoa(i + 1)
			frep := f.Frep[i] + "\n"
			t := tarotPhoto(assetURL, card, isDown, frep+"第"+num+"张牌：")
			if err := c.SendAlbum(tele.Album{t}); err != nil {
				return err
			}
		}
		if f.IsCut {
			card, isDown, err := drawUniqueTarot(used)
			if err != nil {
				return err
			}
			t := tarotPhoto(assetURL, card, isDown, "切牌：")
			if err := c.SendAlbum(tele.Album{t}); err != nil {
				return err
			}
		}
		return nil
	})

	b.Handle(tele.OnQuery, func(c tele.Context) error {
		bg := assetURL + "Extra/BG.jpg"

		results := make(tele.Results, 1)

		results[0] = &tele.PhotoResult{
			URL:         bg,
			Title:       "塔罗牌",
			Description: "喵～",
			Caption:     "点击发送即可占卜",
			ThumbURL:    bg,
		}

		results[0].SetResultID("0")

		return c.Answer(&tele.QueryResponse{
			Results:   results,
			CacheTime: 1,
		})
	})

	b.Handle(tele.OnGuestMessage, func(c tele.Context) error {
		message := c.Message()
		if message == nil || message.GuestQueryID == "" {
			return fmt.Errorf("guest message has no guest query id")
		}

		card, isDown, err := tarot.Get_tarot()
		if err != nil {
			return err
		}
		result := tarotPhotoResult(assetURL, card, isDown, "看看 "+senderName(c)+" 抽到了什么：\n")
		return c.AnswerGuest(result)
	})

	b.Handle(tele.OnInlineResult, func(c tele.Context) error {
		card, isDown, err := tarot.Get_tarot()
		if err != nil {
			return err
		}
		t := tarotPhoto(assetURL, card, isDown, "看看 "+senderName(c)+" 抽到了什么：\n")
		if err := c.Edit(t); err != nil && !errors.Is(err, tele.ErrTrueResult) {
			return err
		}
		return nil
	})

	b.Start()
}
