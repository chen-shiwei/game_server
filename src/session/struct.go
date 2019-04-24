package session

import (
	"user"
)

type ResponseMatch struct {
	RoomName  string       `json:"name"`
	RoomId    string       `json:"roomId"`
	RoomLimit uint16       `json:"limit"`
	Users     []*user.User `json:"users"`
}
