package landlord

import (
	"encoding/json"
	"errors"
	"log"
	"protocol"
	"room"
	"rooms/landlord/game"
	"session"
)

var (
	ErrorCommandSupport = errors.New("does not support")
	ErrorInvalidSess    = errors.New("sess id invalid")
)

type LandLordPlayer struct {
	mingpai bool
	online  bool
	deposit bool
	sess    *session.Session
}

type LandLordRoom struct {
	mingpai      bool
	playerIdList []string
	player       map[string]*LandLordPlayer
	gameTable    *game.LandLordGame
}

func init() {
	err := room.RegisterGameEntry("1001", "1001", &LandLordRoom{}, 0)
	if err != nil {
		log.Print(err.Error())
	}
}

func (s *LandLordRoom) New() room.GameHandler {
	r := new(LandLordRoom)
	r.playerIdList = make([]string, 0)
	r.player = make(map[string]*LandLordPlayer)
	return r
}

func (s *LandLordRoom) Join(sess *session.Session, back bool) (interface{}, error) {
	return nil, nil
}

func (s *LandLordRoom) Leave(sess *session.Session, quit bool) (interface{}, error) {
	return nil, nil
}

// 房间主逻辑仅处理请求发起方的响应，向其它玩家发送的信息请自行处理。
// data 为 []byte，需要使用JSON解码
func (s *LandLordRoom) Request(sess *session.Session, data []byte) (interface{}, error) {
	if s.gameTable == nil {
		s.gameTable = game.New(s, s.playerIdList...)
	}
	cmd := new(RequestCommand)
	if err := json.Unmarshal(data, cmd); err != nil {
		log.Println("command parse failed,", err.Error())
		return nil, err
	}
	var (
		ret *game.GameResponse
		err error
	)
	switch cmd.Cmd {
	case game.CMD_READY:
		ret, err = s.gameTable.Ready(sess.HashCode(), data)
	case game.CMD_START:
		ret, err = s.gameTable.Start(sess.HashCode(), data)
	case game.CMD_MINGPAI:
		ret, err = s.gameTable.Mingpai(sess.HashCode(), data)
	case game.CMD_CALLLORD:
		ret, err = s.gameTable.CallLoad(sess.HashCode(), data)
	case game.CMD_DOUBLE:
		ret, err = s.gameTable.Double(sess.HashCode(), data)
	case game.CMD_OUTPOKE:
		ret, err = s.gameTable.OutPoke(sess.HashCode(), data)
	default:
		return nil, ErrorCommandSupport
	}
	return ret, err
}

// 当有匹配行为时请处理
func (s *LandLordRoom) Match(sess *session.Session, data []byte) bool {
	jsonReq := new(RequestMatchPlayer)
	if err := json.Unmarshal(data, jsonReq); err != nil {
		log.Println("Match Player Failed, Error Json Syntax,", err.Error())
		panic(room.ErrorJsonDecode)
	}
	if _, ok := s.player[sess.HashCode()]; ok {
		return false
	}
	// TODO: check user level
	s.player[sess.HashCode()] = &LandLordPlayer{sess: sess, mingpai: s.mingpai}
	return true
}

func (s *LandLordRoom) Create(sess *session.Session, data []byte) error {
	jsonReq := new(RequestMatchPlayer)
	if err := json.Unmarshal(data, jsonReq); err != nil {
		log.Println("Create Room Failed,", err.Error())
		return err
	}
	s.mingpai = jsonReq.IsMingpai == 1
	s.playerIdList = append(s.playerIdList, sess.HashCode())
	s.player[sess.HashCode()] = &LandLordPlayer{mingpai: jsonReq.IsMingpai == 1, sess: sess}
	return nil
}

func (s *LandLordRoom) Broadcast(msg *game.GameResponse) error {
	pkt := protocol.NewPacket()
	pkt.SetCommand(protocol.CMD_REQUEST_ROOM)
	d, _ := json.Marshal(msg.Data)
	pkt.SetType(protocol.PACKET_JSON)
	pkt.SetData(d)
	for _, p := range s.player {
		p.sess.Send(pkt)
	}
	return nil
}

func (s *LandLordRoom) Notice(sessId string, msg *game.GameResponse) error {
	if p, ok := s.player[sessId]; !ok {
		return ErrorInvalidSess
	} else {
		pkt := protocol.NewPacket()
		pkt.SetCommand(protocol.CMD_REQUEST_ROOM)
		d, _ := json.Marshal(msg.Data)
		pkt.SetType(protocol.PACKET_JSON)
		pkt.SetData(d)
		p.sess.Send(pkt)
	}
	return nil
}
