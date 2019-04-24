package plane

import (
	"encoding/json"
	"log"
	"protocol"
	"session"
	"time"
)

// create match room
func CreateMatchRoom(sess *session.Session, planePos MatchRequest) {
	r := new(Room)
	r.RoomId = time.Now().UnixNano()
	r.Owner = FightData{Session: *sess, MatchRequest: planePos, Action: true}
	r.Guest = FightData{}
	r.CraeteTime = r.RoomId
	MatchingRooms[r.RoomId] = r
	UserRooms[sess.User.UserId] = &UserRoom{
		InRoomId: r.RoomId,
		Role:     "owner"}
	return
}

func (r *Room) Join(sess *session.Session, planePos MatchRequest) bool {
	if r.DisabledJoin {
		return false
	}
	r.Guest = FightData{Session: *sess, MatchRequest: planePos, Action: false}
	r.DisabledJoin = true
	PlayingRooms[r.RoomId] = r
	UserRooms[sess.User.UserId] = &UserRoom{
		InRoomId: r.RoomId,
		Role:     "guest"}
	delete(MatchingRooms, r.RoomId)
	log.Print("加入房间", PlayingRooms)
	log.Print("匹配房间", MatchingRooms)
	go r.notice()
	return true
}

func (r *Room) DestoryRoom() {
	if _, exist := PlayingRooms[r.RoomId]; exist {
		delete(PlayingRooms, r.RoomId)
		return
	}
	if _, exist := MatchingRooms[r.RoomId]; exist {
		delete(PlayingRooms, r.RoomId)
		return
	}
	delete(UserRooms, r.Guest.User.UserId)
	delete(UserRooms, r.Owner.User.UserId)
	return
}

func GetRoom(rid int64) (*Room, error) {
	if room, exist := PlayingRooms[rid]; exist {
		return room, nil
	}
	if room, exist := MatchingRooms[rid]; exist {
		return room, nil
	}
	return nil, ErrorRoomNoExist
}

func (r *Room) notice() error {
	time.Sleep(time.Second * 1)
	//notice owner
	ownerBuffer, err := json.Marshal(NewSuccessResponse(MatchEnterNotice{
		FirstHit:          FirstHit{UserId: r.Owner.User.UserId},
		User:              *r.Guest.Session.User,
		MatchRequest:      r.Guest.MatchRequest,
		PlayerHistoryData: PlayerHistoryData{}}))
	if err != nil {
		return ErrorRoomNoExist
	}

	noticeOwner := protocol.BuildPacket(0, 0, protocol.CMD_PLAN_MATCH_NOTICE,
		protocol.PACKET_JSON)
	noticeOwner.SetData(ownerBuffer)
	r.Owner.Session.Send(noticeOwner)

	//notice guest
	guestBuffer, err := json.Marshal(NewSuccessResponse(MatchEnterNotice{
		FirstHit{UserId: r.Owner.User.UserId},
		*r.Owner.Session.User,
		r.Owner.MatchRequest,
		PlayerHistoryData{}}))
	if err != nil {
		return ErrorRoomNoExist
	}
	noticeGuest := protocol.BuildPacket(0, 0, protocol.CMD_PLAN_MATCH_NOTICE,
		protocol.PACKET_JSON)
	noticeGuest.SetData(guestBuffer)
	r.Guest.Session.Send(noticeGuest)
	return nil
}

func (r *Room) NewTimer(role string, seconds time.Duration) {
	var winer string
	if role == "owner" {
		r.Owner.Action = false
		r.Guest.Action = true
		winer = "guest"
	} else {
		r.Owner.Action = true
		r.Guest.Action = false
		winer = "owner"
	}
	r.Tm = time.NewTimer(time.Second * seconds)
	go func(winer string) {
		<-r.Tm.C

		OverGame(r, winer)
	}(winer)
}

func GetUserRole(sess *session.Session) (*Room, string) {
	userRoom := UserRooms[sess.User.UserId]
	r := PlayingRooms[userRoom.InRoomId]
	return r, userRoom.Role
}
