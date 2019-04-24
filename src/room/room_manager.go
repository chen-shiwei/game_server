package room

import (
	"controller"
	"encoding/json"
	"errors"
	"log"
	"protocol"
	"session"
)

type RoomManager struct {
	rooms     map[string]*Room
	usersRoom map[string]string
}

var (
	roomManager            *RoomManager
	ErrorUserInRoom        = errors.New("user already in a room")
	ErrorRoomExists        = errors.New("room exists")
	ErrorRoomNotExists     = errors.New("room does not exists")
	ErrorUserDoesNotInRoom = errors.New("user does not in room")
	ErrorUnknowGame        = errors.New("unknown game")
	ErrorUnsupportRoom     = errors.New("unsupported room")
	ErrorParamSyntax       = errors.New("params syntax error")
)

func init() {
	roomManager = &RoomManager{rooms: make(map[string]*Room), usersRoom: make(map[string]string)}
	controller.Register(protocol.CMD_CREATE_ROOM, CreateRoom)
	controller.Register(protocol.CMD_JOIN_ROOM, JoinRoom)
	controller.Register(protocol.CMD_LEAVE_ROOM, LeaveRoom)
	controller.Register(protocol.CMD_REQUEST_ROOM, RequestRoom)

}

func (s *RoomManager) AddRoom(sess *session.Session, r *Room) error {
	if _, ok := s.rooms[r.HashCode()]; !ok {
		s.rooms[r.HashCode()] = r
		s.usersRoom[sess.User.UserId] = r.HashCode()
		return nil
	} else {
		return ErrorRoomExists
	}
}

func (s *RoomManager) GetRoom(r *Room) (*Room, error) {
	if room, ok := s.rooms[r.HashCode()]; ok {
		return room, nil
	} else {
		return nil, ErrorRoomNotExists
	}
}

func (s *RoomManager) GetRoomByRoomId(roomId string) (*Room, error) {
	if room, ok := s.rooms[roomId]; ok {
		return room, nil
	} else {
		return nil, ErrorRoomNotExists
	}
}

func (s *RoomManager) GetRoomBySession(sess *session.Session) (*Room, error) {
	if roomId, ok := s.usersRoom["1"]; ok {
		if room, ok := s.rooms[roomId]; ok {
			return room, nil
		} else {
			return nil, ErrorRoomNotExists
		}
	} else {
		return nil, ErrorUserDoesNotInRoom
	}
}

func (s *RoomManager) JoinRoom(sess *session.Session, roomId string, password string) (*ResponseJoinRoom, error) {
	if _, err := s.GetRoomBySession(sess); err == nil {
		return nil, ErrorUserInRoom
	}
	if r, err := s.GetRoomByRoomId(roomId); err != nil {
		return nil, ErrorParamSyntax
	} else {
		if err := r.Join(sess, password); err != nil {
			return nil, err
		} else {
			s.usersRoom[sess.User.UserId] = r.HashCode()
			resp := &ResponseJoinRoom{RoomId: roomId}
			resp.PlayerList = make([]*Player, 0)
			for _, u := range r.Members() {
				var manager uint8
				if r.owner.User.UserId == u.UserId {
					manager = 1
				}
				resp.PlayerList = append(resp.PlayerList, &Player{
					Name: u.Name, Id: u.UserId, Manager: manager})
			}
			return resp, nil
		}
	}
}

func (s *RoomManager) LeaveRoom(sess *session.Session) error {
	if r, err := s.GetRoomBySession(sess); err != nil {
		return ErrorUserDoesNotInRoom
	} else {
		if err = r.Leave(sess, true); err != nil {
			return err
		}
		delete(s.usersRoom, sess.User.UserId)
		if r.MemberCount() == 0 {
			delete(s.rooms, r.HashCode())
		}
		return nil
	}
}

func RequestRoom(sess *session.Session, pack *protocol.Packet) (interface{}, error) {
	if room, err := roomManager.GetRoomBySession(sess); err != nil {
		return nil, err
	} else {
		return room.Request(sess, pack)
	}
}

func LeaveRoom(sess *session.Session, data *protocol.Packet) (interface{}, error) {
	resp := &ResponseLeaveRoom{}
	if err := roomManager.LeaveRoom(sess); err != nil {
		log.Println("LeaveRoom Error:", err.Error())
		resp.ResultCode = -1
		resp.Reason = err.Error()
		return resp, err
	}
	return resp, nil
}

func JoinRoom(sess *session.Session, data *protocol.Packet) (interface{}, error) {
	jsonReq := new(RequestJoinRoom)
	if err := json.Unmarshal(data.Data(), jsonReq); err != nil {
		log.Println("JoinRoom Failed:", err.Error())
		return nil, ErrorParamSyntax
	}
	if resp, err := roomManager.JoinRoom(sess, jsonReq.RoomId, jsonReq.Password); err != nil {
		log.Println("JoinRoom Failed:", err.Error())
		return nil, err
	} else {
		return resp, nil
	}
}

func CreateRoom(sess *session.Session, data *protocol.Packet) (interface{}, error) {
	if _, err := roomManager.GetRoomBySession(sess); err == nil {
		return nil, ErrorUserInRoom
	}
	jsonReq := new(RequstCreateRoom)
	if err := json.Unmarshal(data.Data(), jsonReq); err != nil {
		log.Println("CreateRoom Failed:", err.Error())
		return nil, ErrorParamSyntax
	}
	if jsonReq.GameId == "" {
		return nil, ErrorUnknowGame
	}
	var (
		r     *Room
		entry *GameRoomEntry
		err   error
	)
	if jsonReq.Type == TYPE_CHATROOM {
		jsonReq.GameId = "chat_room"
	}
	entry, err = GetGameRoomEntry(jsonReq.GameId)
	if err != nil {
		data.SetType(protocol.PACKET_STRING)
		data.SetData([]byte(err.Error()))
		return nil, err
	}
	switch jsonReq.Type {
	case TYPE_CHATROOM:
		r = CreateChatRoom(sess, entry.Handler, data.Data(), jsonReq.Password)
	case TYPE_GAMEROOM:
		// pass the request to create a room for custom parameters
		r, err = CreateGameRoom(sess, entry.MemberLimit, entry.Handler, data.Data(), jsonReq.Password, string(jsonReq.Match))
	default:
		log.Println("room type does not supported")
		return nil, ErrorUnsupportRoom
	}
	if r == nil {
		return nil, errors.New("create room failed")
	}
	res := &ResponseCreateRoom{Name: r.HashCode(), RoomId: r.HashCode()}
	return res, nil
}
