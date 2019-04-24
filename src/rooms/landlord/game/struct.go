package game

const (
	CMD_READY       int = 1
	CMD_START       int = 2
	CMD_TIME        int = 3
	CMD_MINGPAI     int = 11
	CMD_CALLLORD    int = 12
	CMD_DOUBLE      int = 13
	CMD_OUTPOKE     int = 14
	CMD_CONFIRMLORD int = 15
)

type RequestReady struct {
	Cmd   int `json:"cmd"`
	Ready int `json:"ready"`
}

type RequestStart struct {
	Cmd   int `json:"cmd"`
	Start int `json:"start"`
}

type ResponseError struct {
	Cmd  int    `json:"cmd"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// 计时器剩余时间
type ResponseTimer struct {
	Cmd      int `json:"cmd"`
	Code     int `json:"code"`
	Timeleft int `json:"timeleft"`
}

// 在发牌时自己点击明牌，或者地主第一次可以点击明牌
type RequestMingpai struct {
	Cmd           int `json:"cmd"`
	MingpaiNumber int `json:"mingpainum"` // 点击明牌时显示的明牌倍数
}

type ResponseMingpai struct {
	Cmd  int `json:"cmd"`
	Code int `json:"code" // 判定是否可以明牌`
}

// 我叫地主的结果传给服务器，便于服务器推送
type RequestMineCallLandlord struct {
	Cmd            int `json:"cmd"`
	IsCallLandLord int `json:"iscalllandlord"` // 是否叫地主
}

type ResponseMineCallLandlord struct {
	Cmd  int `json:"cmd"`
	Code int `json:"code"` //
}

// 自己出牌
type RequestMineOutPoke struct {
	Cmd       int   `json:"cmd"`
	IsOutPoke int   `json:"isoutpoke"` // 出不出另说
	OutPokeId []int `json:"outpokeid"` // 出牌的ID
}

type ResponseMineOutPoke struct {
	Cmd  int `json:"cmd"`
	Code int `json:"code"` // code判断出牌是否成功
}

// 查看结算倍数相信信息
type ResponseMineResultDataInfo struct {
	Cmd         int `json:"cmd"`
	Code        int `json:"code"`
	BaseTime    int `json:"basetime"`    // 初始倍数
	MingpaiTime int `json:"mingpaitime"` // 明牌倍数
	RobLandLord int `json:"roblandlord"` // 抢地主的倍数
	DipaiTime   int `json:"dipaitime"`   // 底牌倍数
	BombTime    int `json:"bombtime"`    // 炸弹倍数
	SpringTime  int `json:"sprinttime"`  // 春天倍数
	AllTime     int `json:"alltime"`     // 总倍数
}

// 加倍
type RequestMineDoubleTime struct {
	Cmd     int `json:"cmd"`
	TimeNum int `json:"timenum"` // 加倍的倍数    暂时默认是2倍。由客户端给服务器传
}

type ResponseMineDoubleTime struct {
	Cmd  int `json:"cmd"`
	Code int `json:"code"` // 判断出牌是否成功
}

// 对局中单个玩家座次信息
type PlayerInfo struct {
	Pos       int    `json:"pos"` // 玩家座次信息
	PlayerId  string `json:"playerid"`
	Nickname  string `json:"nickname"`  // 玩家昵称
	Icon      string `json:"headicon"`  // 玩家头像
	Sex       int    `json:"sex"`       // 玩家性别
	GoldNum   int    `json:"goldnum"`   // 金币数
	IsReady   int    `json:"isready"`   // 是否准备
	IsMingpai int    `json:"ismingpai"` // 玩家是否明牌开始
}

// 玩家匹配成功
type ResponseMatchPlayerSuccess struct {
	Cmd        int          `json:"cmd"`
	Code       int          `json:"code"`
	CurTimeNum int          `json:"curtimenum"` // 当前对局倍数
	PlayerInfo []PlayerInfo `json:"playerinfo"` // 对局中玩家座次以及相关信息  --数组
}

// 发牌
type PokeInfo struct {
	PlayerId string `json:"playerid"` // 玩家座次信息
	CardId   []int  `json:"cardid"`   // 对应玩家的手牌的ID  --如果玩家没有明牌并且不是自己 则ID全部传默认值0
}

// 发牌
type ResponseDealLandlordPoke struct {
	Cmd        int        `json:"cmd"`
	Code       int        `json:"code"`
	PlayerPoke []PokeInfo `json:"playerpoke"` //
}

// 在发牌时玩家点击明牌
type ResponseOtherMingpai struct {
	Cmd         int    `json:"cmd"`
	Code        int    `json:"code"`
	PlayerId    string `json:"playerid"`    // 广播玩家叫地主结果
	MingpaiTime int    `json:"mingpaitime"` // 明牌的倍数
}

// 别人叫地主的广播
type ResponseOtherCallLandlord struct {
	Cmd            int    `json:"cmd"`
	Code           int    `json:"code"`
	PlayerId       string `json:"playerid"`       // 广播玩家叫地主结果
	IsCallLandlord int    `json:"iscalllandlord"` // 是否叫地主或者抢地主   1---叫或者抢   2---不叫
}

// 换人叫地主
type ResponsePushCallLandlord struct {
	Cmd      int    `json:"cmd"`
	Code     int    `json:"code"`
	PlayerId string `json:"playerid"` // 广播玩家叫地主结果
	IsFirst  int    `json:"isfirst"`  // 我是不是第一个叫地主的  1---是第一个/前面没有人叫地主(客户端进行显示叫地主按钮)  2---不是第一个&前面有人叫地主(客户端显示抢地主按钮)
}

// 服务器确定谁是地主
type ResponseConfirmLandlord struct {
	Cmd          int    `json:"cmd"` //
	Code         int    `json:"code"`
	PlayerId     string `json:"playerid"`     // 谁是地主
	LandlordPoke []int  `json:"landlordpoke"` // 地主的三张牌
}

// 发牌
type ResponsePushDoubleTime struct {
	Cmd      int    `json:"cmd"`
	Code     int    `json:"code"`
	PlayerId string `json:"playerid"`
	IsDouble int    `json:"isdouble"` // 是否加倍
}

// 出牌
type ResponseOtherOutPoke struct {
	Cmd       int    `json:"cmd"`
	Code      int    `json:"code"`
	PlayerId  string `json:"playerid"`  // 出牌的人
	IsOutPoke int    `json:"isoutpoke"` // 是否出牌，是要不起还是怎么滴
	PokeId    []int  `json:"pokeid"`    // 出牌的ID
	PokeType  int    `json:"poketype"`  // 出牌的类型   --单张，对子，三带一，三带一对，四带二，四带两个对，顺子，连对，飞机(不带)，飞机(单张)，飞机(对子)，炸弹，王炸
}

//
type ResultInfo struct {
	PlayerId        string `json:"playerid"`        // 玩家ID
	NickName        string `json:"nickname"`        // 玩家昵称
	IsWin           int    `json:"iswin"`           // 是否胜利  1---胜利  2---失败
	IsFinishVictory int    `json:"isfinishvictory"` //是否连胜中断 1---中断  2--没中断
	BaseScore       int    `json:"basescore"`       // 对局底分
	GoldNum         int    `json:"goldnum"`         // 赢或者输的分数
	TimeNum         int    `json:"timenum"`         // 对局中总倍数
	WinsNum         int    `json:"winsnum"`         // 连胜数
	IsLandlord      int    `json:"islandlord"`      // 是否是地主  1---地主  2---农民
	IsBankRuptcy    int    `json:"isbankruptcy"`    // 是否破产  1---破产，2---没破产
	IsCapping       int    `json:"iscapping"`       // 是否封顶 1---封顶 2---没封顶
}

// 服务器推送结算消息
type ResponseResultData struct {
	Cmd        int          `json:"cmd"`
	Code       int          `json:"code"`
	ResultInfo []ResultInfo `json:"resultinfo"` // 结算信息面板
}

// 更新对局中倍数的显示
type ResponsePushUpdateCurTime struct {
	Cmd        int `json:"cmd"`
	Code       int `json:"code"`
	CurTimeNum int `json:"curtimenum"` // 当前的倍数
}
