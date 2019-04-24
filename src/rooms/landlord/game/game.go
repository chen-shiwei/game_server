package game

import (
	"encoding/json"
	"errors"
	"math/rand"
	"time"
)

var (
	ErrorTickerRunning = errors.New("ticker running")
)

type Messager interface {
	Broadcast(msg *GameResponse, except ...string) error
	Notice(sessId string, msg *GameResponse) error
}

type GameResponse struct {
	OpType int
	Data   interface{}
}

type LandLordGame struct {
	times          int              // 加倍
	messager       Messager         // 房间通知接口
	readyList      map[string]int   // 玩家准备列表
	stopTickerList map[string]bool  // 停止计时列表
	startTime      time.Time        // start time
	timeWait       int              // 当前出牌位等待时间
	lord           string           // 地主
	pos            map[string]int   // 当前玩家对应座位
	players        map[int]string   // 桌面位置对应玩家
	pokesOwn       map[string][]int // 当前玩家持有牌
	pokesOut       map[string][]int // 当前桌面牌
	lordPoke       []int            // 地主牌
	bombCount      int              // 炸弹计数
	curPoke        []int            // 上一次出牌
	curIdx         int              // 上一个出版人
	ticker         *time.Ticker     // 计时器
}

func New(messager Messager, players ...string) *LandLordGame {
	game := &LandLordGame{
		messager:       messager,
		readyList:      make(map[string]int),
		stopTickerList: make(map[string]bool),
		pos:            make(map[string]int),
		players:        make(map[int]string),
		pokesOwn:       make(map[string][]int),
		pokesOut:       make(map[string][]int),
		lordPoke:       make([]int, 0)}
	for i, p := range players {
		game.pokesOwn[p] = make([]int, 0)
		game.pokesOut[p] = make([]int, 0)
		game.readyList[p] = 1 // 默认准备
		game.pos[p] = i + 1
		game.players[i+1] = p
	}
	return game
}

func (s *LandLordGame) Ready(sid string, d []byte) (ret *GameResponse, err error) {
	return
}

func (s *LandLordGame) Start(sid string, d []byte) (ret *GameResponse, err error) {
	s.freshPoke()
	ret = &GameResponse{}
	s.startTime = time.Now()

	s.startTicker(30, func() {
		msg := &GameResponse{Data: &ResponseTimer{Cmd: CMD_TIME, Timeleft: s.timeWait}}
		sendTimes := 0
		for sid, _ := range s.pos {
			if _, ok := s.stopTickerList[sid]; ok {
				continue
			}
			s.messager.Notice(sid, msg)
			sendTimes++
		}
		if sendTimes == 0 {
			s.stopTicker()
		}
	}, func() {
		//
	})
	ret.Data = &struct {
		Code int `json:"code"`
	}{1}
	return
}

func (s *LandLordGame) Mingpai(sid string, d []byte) (ret *GameResponse, err error) {
	return
}

func (s *LandLordGame) CallLoad(sid string, d []byte) (ret *GameResponse, err error) {
	defer s.nextPlayerTicker()
	s.stopTicker()
	if _, ok := s.stopTickerList[sid]; ok {
		ret = &GameResponse{Data: &ResponseError{Code: -1, Msg: "denied"}}
		return
	}
	jsonReq := new(RequestMineCallLandlord)
	if err = json.Unmarshal(d, jsonReq); err != nil {
		ret = &GameResponse{Data: &ResponseError{Code: -1, Msg: err.Error()}}
		return
	}
	s.stopTickerList[sid] = true
	if jsonReq.IsCallLandLord == 1 {
		s.lord = sid
	}
	if len(s.stopTickerList) == 3 {
		msg := &GameResponse{Data: &ResponseConfirmLandlord{
			PlayerId:     s.lord,
			LandlordPoke: s.lordPoke}}
		s.messager.Broadcast(msg, sid)
	} else {
		msg := &GameResponse{Data: &ResponseOtherCallLandlord{
			PlayerId:       sid,
			IsCallLandlord: jsonReq.IsCallLandLord}}
		s.messager.Broadcast(msg, sid)
	}
	ret = &GameResponse{Data: &ResponseMineCallLandlord{Code: 0}}
	return
}

func (s *LandLordGame) Double(sid string, d []byte) (ret *GameResponse, err error) {
	return
}

func (s *LandLordGame) OutPoke(sid string, d []byte) (ret *GameResponse, err error) {
	return
}

func (s *LandLordGame) freshPoke() {
	// 3-10 JQKA2 Joker(0,1)
	pokerList := make([]int, 0)
	for i := 3; i <= 15; i++ {
		pokerList = append(pokerList, i*10+1, i*10+2, i*10+3, i*10+4)
	}
	pokerList = append(pokerList, 16*10+0, 16*10+1)
	for len(pokerList) > 3 {
		for i, _ := range s.pokesOwn {
			rand.Seed(time.Now().UnixNano())
			k := rand.Intn(len(pokerList))
			s.pokesOwn[i] = append(s.pokesOwn[i], pokerList[k])
			pokerList = append(pokerList[:k], pokerList[k+1:]...)
		}
	}
	s.lordPoke = append(s.lordPoke, pokerList...)
}

func (s *LandLordGame) startTicker(seconds int, tick func(), complete func()) error {
	if s.ticker != nil {
		return ErrorTickerRunning
	}
	s.timeWait = seconds
	s.ticker = time.NewTicker(1 * time.Second)
	go func() {
		for {
			<-s.ticker.C
			if s.timeWait > 0 {
				tick()
			} else if s.timeWait < 0 {
				s.timeWait = 0
				s.ticker = nil
				return
			} else {
				complete()
				break
			}
			s.timeWait--
		}
	}()
	return nil
}

func (s *LandLordGame) stopTicker() {
	if s.ticker == nil {
		s.timeWait = 0
	} else {
		s.timeWait = -1
		s.ticker.Stop()
	}
}

func (s *LandLordGame) nextPlayerTicker() {
	s.curIdx++
	if s.curIdx > 3 {
		s.curIdx = 1
	}
	sid := s.players[s.curIdx]
	s.startTicker(30, func() {
		msg := &GameResponse{Data: &ResponseTimer{Cmd: CMD_TIME, Timeleft: s.timeWait}}
		s.messager.Notice(sid, msg)
	}, func() {
		s.nextPlayerTicker()
	})
}
