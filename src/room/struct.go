package room

type RequstCreateRoom struct {
	Type     uint8  `json:"type"`
	GameId   string `json:"gameId"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Match    uint8  `json:"match"`
}

type ResponseCreateRoom struct {
	Name   string `json:"name"`
	RoomId string `json:"roomId"`
}

type RequestJoinRoom struct {
	RoomId      string `json:"roomId"`
	MemberCount uint16 `json:"count"`
	Password    string `json:"password"`
}

type Player struct {
	Name    string `json:"name"`
	Id      string `json:"id"`
	Icon    string `json:"icon"`
	Manager uint8  `json:"manager"`
}

type ResponseJoinRoom struct {
	RoomId     string    `json:"roomId"`
	PlayerList []*Player `json:"playerList"`
}

type RequestLeaveRoom struct {
}

type ResponseLeaveRoom struct {
	ResultCode int    `json:"code"`
	Reason     string `json:"reason"`
}
