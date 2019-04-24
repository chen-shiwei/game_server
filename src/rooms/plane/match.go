package plane

import (
	"encoding/json"
	"log"
	"protocol"
	"session"
	"sort"
)

func MatchAction(sess *session.Session, pkt *protocol.Packet) (interface{}, error) {
	log.Print(sess.User, "进入匹配")
	var posData MatchRequest
	if err := json.Unmarshal(pkt.Data(), &posData); err != nil {
		return NewErrorResponse(601), nil
	}
	log.Print(posData)
	log.Print(MatchingRooms, "匹配房间")

	if len(MatchingRooms) < 1 {
		log.Print(sess.User, "创建房间")
		CreateMatchRoom(sess, posData)
		log.Print(UserRooms, "房间用户")
		return NewSuccessResponse(MatchResponse{MatchStatus: 1}), nil
	}
	//get time is wait longer room
	var keys []int64
	for key := range MatchingRooms {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i int, j int) bool {
		return keys[i] < keys[j]
	})
	log.Print("房间号数据", keys)
	log.Print("房间号数据1", MatchingRooms[keys[0]])
	r, err := GetRoom(MatchingRooms[keys[0]].RoomId)
	if err != nil {
		return NewErrorResponse(600), nil
	}
	log.Print("获得的房间", r)
	if r.Join(sess, posData) {
		log.Print("开始计时器63")
		r.NewTimer("guest", 63)
		return NewSuccessResponse(MatchResponse{MatchStatus: 2}), nil
	}
	return NewErrorResponse(600), nil
}

func MatchActionCancel(sess *session.Session, pkt *protocol.Packet) (interface{}, error) {
	log.Print(sess.User, "取消匹配操作")
	log.Print(UserRooms, MatchingRooms)
	if ur, exist := UserRooms[sess.User.UserId]; exist {

		if room, exist := MatchingRooms[ur.InRoomId]; exist {
			room.DestoryRoom()
			log.Print(sess.User, "取消匹配成功")
			return NewSuccessResponse(nil), nil
		}
	}
	log.Print(sess.User, "取消匹配失败")

	return NewErrorResponse(1301), nil
}
