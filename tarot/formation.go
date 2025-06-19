package tarot

// Formation 定义了塔罗牌阵
type Formation struct {
	Fid   int      `json:"id"`              // 牌阵ID
	Fname string   `json:"name"`            // 牌阵名称
	Fnum  int      `json:"formations_num"`  // 牌阵需要的牌数
	IsCut bool     `json:"is_cut"`          // 是否需要切牌
	Frep  []string `json:"representations"` // 牌阵的含义
}

// Formations 定义了所有的牌阵
// 每个牌阵都包含牌阵ID、名称、需要的牌数、是否需要切牌以及每张牌的含义
// 共有10种牌阵
var Formations [10]Formation = [10]Formation{{
	Fid:   0,
	Fname: "圣三角牌阵",
	Fnum:  3,
	IsCut: false,
	Frep:  []string{"处境", "行动", "结果"},
}, {
	Fid:   1,
	Fname: "圣三角牌阵",
	Fnum:  3,
	IsCut: false,
	Frep:  []string{"现状", "愿望", "行动"},
}, {
	Fid:   2,
	Fname: "时间之流牌阵",
	Fnum:  3,
	IsCut: true,
	Frep:  []string{"过去", "现在", "未来"},
}, {
	Fid:   3,
	Fname: "四要素牌阵",
	Fnum:  4,
	IsCut: false,
	Frep:  []string{"火，象征行动，行动上的建议", "气，象征言语，言语上的对策", "水，象征感情，感情上的态度", "土，象征物质，物质上的准备"},
}, {
	Fid:   4,
	Fname: "五牌阵",
	Fnum:  5,
	IsCut: true,
	Frep:  []string{"现在或主要问题", "过去的影响", "未来", "主要原因", "行动可能带来的结果"},
}, {
	Fid:   5,
	Fname: "吉普赛十字阵",
	Fnum:  5,
	IsCut: false,
	Frep:  []string{"对方的想法", "你的想法", "相处中存在的问题", "二人目前的环境", "关系发展的结果"},
}, {
	Fid:   6,
	Fname: "马蹄牌阵",
	Fnum:  6,
	IsCut: true,
	Frep:  []string{"现状", "可预知的情况", "不可预知的情况", "即将发生的", "结果", "问卜者的主观想法"},
}, {
	Fid:   7,
	Fname: "六芒星牌阵",
	Fnum:  7,
	IsCut: true,
	Frep:  []string{"过去", "现在", "未来", "对策", "环境", "态度", "预测结果"},
}, {
	Fid:   8,
	Fname: "平安扇牌阵",
	Fnum:  4,
	IsCut: false,
	Frep:  []string{"人际关系现状", "与对方结识的因缘", "双方关系的发展", "双方关系的结论"},
}, {
	Fid:   9,
	Fname: "沙迪若之星牌阵",
	Fnum:  6,
	IsCut: true,
	Frep:  []string{"问卜者的感受", "问卜者的问题", "问题下的影响因素", "将问卜者与问题纠缠在一起的往事", "需要注意/考虑的", "可能的结果"},
}}
