package lib

// Card 结构体
type Card struct {
	Card_id   int    `json:"id"`   // ID 用来唯一标识一张塔罗牌
	Card_name string `json:"name"` // Name 是塔罗牌的名称
	Card_type string `json:"type"` // Type 是塔罗牌的类型
	Card_up   string `json:"up"`   // Up 是正位的含义
	Card_down string `json:"down"` // Down 是逆位的含义
	Card_file string `json:"file"` // File 是塔罗牌对应的文件路径
}

// Cards 是一个包含所有塔罗牌的数组
// 每张塔罗牌都包含 ID、名称、类型、正位含义、逆位含义和文件路径
// 共有 78 张塔罗牌
var Cards [1]Card = [1]Card{{
	Card_id:   0,
	Card_name: "愚者",
	Card_type: "MajorArcana",
	Card_up:   "请扫码支付 ￥10.00 查看结果",
	Card_down: "请扫码支付 ￥10.00 查看结果",
	Card_file: "MajorArcana/QRCODE"}
}
