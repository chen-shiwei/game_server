package chat

import (
	"room"
	"session"
)

type ChatRoom struct {
	sessionList *session.Session
}

var (
	chatRoom *ChatRoom
)

func (s *ChatRoom) New() room.GameHandler {
	chatRoom = new(ChatRoom)
}

// 新用户加入，sess 当前加入的新用户会话，back 是否为异常断线后返回用户
func (s *ChatRoom) Join(sess *session.Session, back bool) (interface{}, error) {
	return nil, nil
}

// 用户退出，sess 退出的用户会话，quit 是否主动要求退出
func (s *ChatRoom) Leave(sess *session.Session, quit bool) (interface{}, error) {
	return nil, nil
}

// 业务过程中的请求，sess 发出指令的用户会话，data 发送的请求参数
func (s *ChatRoom) Request(sess *session.Session, data []byte) (interface{}, error) {
	return nil, nil
}

func (s *ChatRoom) Match(sess *session.Session, data []byte) bool {
	return false
}

func (s *ChatRoom) Create(sess *session.Session, data []byte) error {
	return nil
}
