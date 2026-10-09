package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	tarot "tarot/lib"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
)

const formationReply = `塔罗牌阵：
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

var (
	userData   map[int64]*tarot.User
	userDataMu sync.Mutex
)

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
	meaning := card.Up
	file := card.File
	if isDown == 1 {
		position = "逆位"
		meaning = card.Down
		file = "_" + file
	}

	return &tele.Photo{
		File:    tele.FromURL(assetURL + file + ".jpg"),
		Caption: prefix + card.Name + " 「" + position + "」\n" + meaning,
	}
}

func tarotPhotoResult(assetURL string, card tarot.Card, isDown int, prefix string) *tele.PhotoResult {
	photo := tarotPhoto(assetURL, card, isDown, prefix)
	photoURL := photo.File.FileURL
	result := &tele.PhotoResult{
		URL:         photoURL,
		Title:       "塔罗牌",
		Description: card.Name,
		Caption:     photo.Caption,
		ThumbURL:    photoURL,
	}
	result.SetResultID(fmt.Sprintf("tarot-%d-%d", card.Id, isDown))
	return result
}

func drawUniqueTarot(used map[int]struct{}) (tarot.Card, int, error) {
	for len(used) < len(tarot.Cards) {
		card, isDown, err := tarot.Get_tarot()
		if err != nil {
			return tarot.Card{}, 0, err
		}
		if _, exists := used[card.Id]; exists {
			continue
		}
		used[card.Id] = struct{}{}
		return card, isDown, nil
	}
	return tarot.Card{}, 0, fmt.Errorf("no tarot cards left to draw")
}

func IsToday(timestamp int64) bool {
	loc := time.FixedZone("UTC+8", 8*60*60)
	now := time.Now().In(loc)
	t := time.Unix(timestamp, 0).In(loc)

	return t.Format("2006-01-02") == now.Format("2006-01-02")
}

func init() {
	godotenv.Load()
}

func drawTarot(c tele.Context) (tarot.Card, int, bool, error) {
	var userID int64
	if message := c.Message(); message != nil && message.GuestUser != nil {
		userID = message.GuestUser.ID
	} else if sender := c.Sender(); sender != nil {
		userID = sender.ID
	} else {
		return tarot.Card{}, 0, false, fmt.Errorf("cannot resolve user from context")
	}

	userDataMu.Lock()
	defer userDataMu.Unlock()

	user := userData[userID]
	if user == nil {
		user = &tarot.User{UserId: userID}
		userData[userID] = user
	}
	if !IsToday(user.UpdateTime) {
		user.Frequency = 0
	}
	if user.Frequency >= 3 {
		return user.ResultCard, user.IsResultDown, true, nil
	}
	card, isDown, err := tarot.Get_tarot()
	if err != nil {
		return tarot.Card{}, 0, false, err
	}
	user.Frequency++
	user.UpdateTime = time.Now().Unix()
	user.ResultCard = card
	user.IsResultDown = isDown
	return card, isDown, false, nil
}

func drawPrefix(c tele.Context, limited bool) string {
	prefix := "看看 " + senderName(c) + " 抽到了什么：\n"
	if limited {
		prefix = "你今天已经抽过三次了 这是你的最终结果\n" + prefix
	}
	return prefix
}

func sendTarot(assetURL string, c tele.Context) (*tele.Photo, error) {
	card, isDown, limited, err := drawTarot(c)
	if err != nil {
		return nil, err
	}
	return tarotPhoto(assetURL, card, isDown, drawPrefix(c, limited)), nil
}

func main() {
	userData = make(map[int64]*tarot.User)

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
		t, err := sendTarot(assetURL, c)
		if err != nil {
			return err
		}
		return c.SendAlbum(tele.Album{t})
	})

	b.Handle("/formation", func(c tele.Context) error {
		chat := c.Chat()
		if chat == nil || chat.Type != tele.ChatPrivate {
			return c.Send("请在私聊中使用此命令。")
		}

		id, err := strconv.Atoi(c.Message().Payload)
		if err != nil || id < 0 || id >= 10 {
			return c.Send(formationReply)
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
			Caption:     "点击下方按钮开始占卜",
			ThumbURL:    bg,
		}

		results[0].SetResultID("0")
		markup := &tele.ReplyMarkup{}
		button := markup.Data("抽一张", "tarot_draw", strconv.FormatInt(c.Sender().ID, 10))
		markup.Inline(markup.Row(button))
		results[0].SetReplyMarkup(markup)

		return c.Answer(&tele.QueryResponse{
			Results:    results,
			CacheTime:  1,
			IsPersonal: true,
		})
	})

	b.Handle(tele.OnGuestMessage, func(c tele.Context) error {
		message := c.Message()
		if message == nil || message.GuestQueryID == "" {
			return fmt.Errorf("guest message has no guest query id")
		}

		if message.ReplyTo != nil && message.ReplyTo.Sender != nil && message.ReplyTo.Sender.IsBot {
			return nil
		}

		card, isDown, limited, err := drawTarot(c)
		if err != nil {
			return err
		}

		return c.AnswerGuest(tarotPhotoResult(assetURL, card, isDown, drawPrefix(c, limited)))
	})

	b.Handle(&tele.Btn{Unique: "tarot_draw"}, func(c tele.Context) error {
		if c.Callback() == nil || !c.Callback().IsInline() {
			return c.Respond(&tele.CallbackResponse{Text: "出错了 呜呜呜"})
		}
		if c.Sender() == nil || c.Data() != strconv.FormatInt(c.Sender().ID, 10) {
			return c.Respond(&tele.CallbackResponse{Text: "你不是主人喵～", ShowAlert: true})
		}
		if err := c.Respond(); err != nil {
			return err
		}
		card, isDown, limited, err := drawTarot(c)
		if err != nil {
			return err
		}
		t := tarotPhoto(assetURL, card, isDown, drawPrefix(c, limited))
		if err := c.Edit(t, &tele.ReplyMarkup{InlineKeyboard: [][]tele.InlineButton{}}); err != nil && !errors.Is(err, tele.ErrTrueResult) {
			return err
		}
		return nil
	})

	b.Start()
}
