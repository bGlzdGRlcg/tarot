package lib

// User 定义了用户信息
type User struct {
	UserId       int64 `json:"user_id"`        // 用户ID
	UpdateTime   int64 `json:"update_time"`    // 更新时间
	Frequency    int   `json:"frequency"`      // 使用频率
	ResultCard   Card  `json:"result_card"`    // 结果卡牌
	IsResultDown int   `json:"is_result_down"` // 结果卡牌是否正位
}
