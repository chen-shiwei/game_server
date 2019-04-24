package landlord

type RequestCommand struct {
	Cmd int `json:"cmd"`
}

// 匹配玩家
type RequestMatchPlayer struct {
	IsMingpai  int `json:"ismingpai"`  // 玩家是否明牌  5倍
	MatchLevel int `json:"matchlevel"` // 玩家匹配的场次等级
}

type ResponseMatchPlayer struct {
	Code int `json:"code"` // 确定成功回调
}
