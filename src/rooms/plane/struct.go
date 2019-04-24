package plane

import (
	"errors"
	"session"
	"time"
	"user"
)

type FightData struct {
	session.Session
	MatchRequest
	LastActionTime int64 `json:"-"`
	Action         bool  `json:"-"` //是否轮到操作
}

//飞机游戏数据
type PlaneData struct {
	Rotate   int   `json:"rotate"`
	CenterId int   `json:"centerid"`
	HeadId   int   `json:"head"`
	BodyPos  []int `json:"body"`
	IsDied   bool  `json:"-"`
}

//房间属性
type Room struct {
	// Name         string `json:"name"`
	RoomId       int64
	DisabledJoin bool
	Owner        FightData
	Guest        FightData
	CraeteTime   int64
	Tm           *time.Timer
}

type UserRoom struct {
	Role     string
	InRoomId int64
}
type PlayerHistoryData struct {
	WinNum int `json:"winNum"`
	AllNum int `json:"allNum"`
}
type FirstHit struct {
	UserId string `json:"userId"`
}

type OverGameNotice struct {
	Players []OverPlayer `josn:"players"`
}

type OverPlayer struct {
	user.User `json:"User"`
	GetScore  int `json:"getscore"`
	IsWin     int `json:"iswin"`
}

//matchStatus 进入匹配队列 等待匹配 2匹配成功 加入房间
type MatchResponse struct {
	MatchStatus int `json:"matchStatus"`
}

//action match 摆放的飞机数据
type MatchRequest struct {
	Plane1 PlaneData `json:"plane1"`
	Plane2 PlaneData `json:"plane2"`
	Plane3 PlaneData `json:"plane3"`
}

//匹配进入房间通知
type MatchEnterNotice struct {
	FirstHit          `json:"firstHit"`
	user.User         `json:"player"`
	MatchRequest      `json:"playerData"`
	PlayerHistoryData `json:"playerHistoryData"`
}

//action 为hit
type HitRequest struct {
	HitId int `json:"hitId"`
}

type HitResponse struct {
	HitResult         int `json:"hitResult"`
	DestroyPlaneIndex int `json:"destroyPlaneIndex"`
}

type HitNotice struct {
	HitId             int `json:"hitId"`
	HitResult         int `json:"hitResult"`
	DestroyPlaneIndex int `json:"destroyPlaneIndex"`
}

type SuccessResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	Code int `json:"code"`
}

var (
	PlayingRooms  = make(map[int64]*Room)      //玩家房间
	MatchingRooms = make(map[int64]*Room)      //匹配房间
	UserRooms     = make(map[string]*UserRoom) //用户所在房间及角色

	ErrorRoomNoExist = errors.New("room not exist")
	ErrorParamSyntax = errors.New("room not exist")
)

func NewSuccessResponse(data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Code: 1001,
		Data: data,
	}
}

func NewErrorResponse(code int) *ErrorResponse {
	return &ErrorResponse{code}
}
