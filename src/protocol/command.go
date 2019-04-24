package protocol

const (
	CMD_ECHO           uint8 = 255
	CMD_KEEPALIVE      uint8 = 0
	CMD_USER_LOGIN     uint8 = 1
	CMD_USER_LOGOUT    uint8 = 2
	CMD_USER_BACK      uint8 = 3
	CMD_CREATE_ROOM    uint8 = 11
	CMD_JOIN_ROOM      uint8 = 12
	CMD_LEAVE_ROOM     uint8 = 13
	CMD_REQUEST_ROOM   uint8 = 14
	CMD_RESPONSE_MATCH uint8 = 21

	CMD_PLAN_MATCH        uint8 = 50
	CMD_PLAN_HIT          uint8 = 51
	CMD_PLAN_MATCH_NOTICE uint8 = 52
	CMD_PLAN_CANCEL_MATCH uint8 = 53
	CMD_PLAN_HIT_NOTICE   uint8 = 54
	CMD_PLAN_OVER         uint8 = 55
)
