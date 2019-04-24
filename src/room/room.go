package room

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"protocol"
	"service"
	"session"
	"time"
	"user"
)

var (
	ErrorHandlerInvalid = errors.New("invalid handle")
	ErrorPasswdInvalid  = errors.New("password invalid")
	ErrorSessionInvalid = errors.New("invalid session or user")
	ErrorNoPermissions  = errors.New("permission denied")
	ErrorRoomLimit      = errors.New("too many users")
	ErrorInRoom         = errors.New("in room")
	ErrorJsonDecode     = errors.New("error when decode json string")
	ErrorNotSupportData = errors.New("not supported data")
)

const (
	TYPE_CHATROOM uint8 = 0
	TYPE_GAMEROOM uint8 = 1
)

type GameHandler interface {
	New() GameHandler
	// init a room for game logic
	Create(*session.Session, []byte) error
	// A new member join to room
	Join(*session.Session, bool) (interface{}, error)
	// A request coming
	Request(*session.Session, []byte) (interface{}, error)
	// A member leave the room, the second param `broken` means connection error, user would be return back
	Leave(*session.Session, bool) (interface{}, error)
	// confirm a matching
	Match(*session.Session, []byte) bool
}

type Room struct {
	roomType        uint8
	createTime      int64
	hashCode        string
	memberLimit     uint16
	memberCount     uint16
	owner           *session.Session
	passwd          string
	sessions        map[string]*session.Session
	sessionMatching map[string]bool
	kickUsers       map[string]int64
	handle          GameHandler
	requestData     []byte
}

func CreateChatRoom(owner *session.Session, handler GameHandler, request []byte, options ...string) *Room {
	r := &Room{roomType: TYPE_CHATROOM, memberLimit: math.MaxUint16, owner: owner, handle: handler.New()}
	r.initRoom()
	r.requestData = request
	if owner != nil {
		r.sessions[owner.HashCode()] = owner
		r.memberCount = 1
	}
	if len(options) > 0 {
		r.passwd = options[0]
	}
	return r
}

func CreateGameRoom(owner *session.Session, memberLimit uint16, handle GameHandler, request []byte, options ...string) (*Room, error) {
	r := &Room{roomType: TYPE_GAMEROOM, memberLimit: memberLimit, owner: owner, memberCount: 1, handle: handle.New()}
	r.initRoom()
	r.requestData = request
	r.sessions[owner.HashCode()] = owner
	if len(options) > 0 {
		r.passwd = options[0]
	}
	defer func() {
		if e := recover(); e != nil {
			log.Println(e)
		}
	}()
	if len(options) > 1 && options[1] == "1" {
		go service.GetSessionManager().Match(r.memberLimit-1, r, func(rslt []*session.Session) error {
			pkt := protocol.NewPacket()
			pkt.SetCommand(protocol.CMD_RESPONSE_MATCH)
			pkt.SetType(protocol.PACKET_JSON)
			matched := r.matched()
			d, _ := json.Marshal(matched)
			pkt.SetData(d)
			for _, s := range r.sessions {
				s.Send(pkt)
			}
			return nil
		})
	}
	return r, nil
}

func (s *Room) initRoom() {
	s.createTime = time.Now().UnixNano()
	s.sessions = make(map[string]*session.Session)
	s.kickUsers = make(map[string]int64)
	s.sessionMatching = make(map[string]bool)
}

func (s *Room) Match(sess *session.Session) (*session.ResponseMatch, error) {
	if _, ok := s.sessionMatching[sess.HashCode()]; ok {
		return nil, ErrorInRoom
	}
	if !s.handle.Match(sess, s.requestData) {
		return nil, ErrorInRoom
	}
	s.sessionMatching[s.hashCode] = true
	s.Join(sess, s.passwd) // ignore error report
	matched := s.matched()
	return matched, nil
}

func (s *Room) matched() *session.ResponseMatch {
	matched := new(session.ResponseMatch)
	matched = &session.ResponseMatch{}
	matched.RoomId = s.HashCode()
	matched.RoomLimit = s.memberLimit
	matched.RoomName = s.HashCode()
	matched.Users = make([]*user.User, 0)
	for _, sess := range s.sessions {
		matched.Users = append(matched.Users, sess.User)
	}
	return matched
}

func (s *Room) MemberCount() uint16 {
	return s.memberCount
}

func (s *Room) Join(sess *session.Session, passwd string) error {
	if sess == nil {
		return ErrorSessionInvalid
	}
	if s.passwd != "" && passwd != s.passwd {
		return ErrorPasswdInvalid
	}
	if deniedTime, ok := s.kickUsers[sess.HashCode()]; ok && deniedTime > time.Now().Unix() {
		return ErrorNoPermissions
	}
	if s.memberCount >= s.memberLimit && s.memberLimit != 0 {
		return ErrorRoomLimit
	}
	var back bool
	// call handler method to check session status
	if s.handle != nil {
		if _, ok := s.sessions[sess.HashCode()]; ok {
			back = true
		} else {
			back = false
		}
		if _, err := s.handle.Join(sess, back); err != nil {
			return err
		}
	}
	s.sessions[sess.HashCode()] = sess
	if !back {
		s.memberCount++
	}
	return nil
}

func (s *Room) Members() []*user.User {
	players := make([]*user.User, 0)
	for _, u := range s.sessions {
		players = append(players, u.User)
	}
	return players
}

func (s *Room) Leave(sess *session.Session, quit bool) error {
	if sess == nil {
		return ErrorSessionInvalid
	}
	if _, ok := s.sessions[sess.HashCode()]; !ok {
		return ErrorSessionInvalid
	}
	if s.handle != nil {
		if _, err := s.handle.Leave(sess, quit); err != nil {
			return err
		}
	}
	if quit {
		s.memberCount--
		delete(s.sessions, sess.HashCode())
	}
	return nil
}

func (s *Room) Request(sess *session.Session, data interface{}) (interface{}, error) {
	if d, ok := data.(*protocol.Packet); ok {
		return s.handle.Request(sess, d.Data())
	} else if d, ok := data.([]byte); ok {
		return s.handle.Request(sess, d)
	} else {
		return nil, ErrorNotSupportData
	}
}

func (s *Room) HashCode() string {
	if s.hashCode == "" {
		h := md5.New()
		h.Write([]byte(fmt.Sprintf("%d", s.createTime)))
		s.hashCode = hex.EncodeToString(h.Sum(nil))
	}
	return s.hashCode
}
