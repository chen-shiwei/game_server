package plane

import (
	"encoding/json"
	"log"
	"protocol"
	"session"
)

func HitAction(sess *session.Session, pkt *protocol.Packet) (interface{}, error) {
	var hit HitRequest
	if err := json.Unmarshal(pkt.Data(), &hit); err != nil {
		return NewErrorResponse(601), nil
	}
	r, role := GetUserRole(sess)
	var (
		destroyPlaneIndex int
		hitResult         int
	)
	noticeHit := protocol.NewPacket()
	noticeHit.SetCommand(protocol.CMD_PLAN_HIT_NOTICE)
	noticeHit.SetType(protocol.PACKET_JSON)

	if role == "owner" {
		if !r.Owner.Action {
			return NewErrorResponse(1302), nil
		}
		destroyPlaneIndex, hitResult = r.Guest.MatchRequest.isHited(hit.HitId)
		if r.Guest.MatchRequest.Plane1.IsDied &&
			r.Guest.MatchRequest.Plane2.IsDied &&
			r.Guest.MatchRequest.Plane3.IsDied {
			OverGame(r, "owner")
		} else {
			stop := r.Tm.Stop()
			r.NewTimer(role, 63)
			log.Print("停止计时器结果：", stop)
		}
		notice, _ := json.Marshal(NewSuccessResponse(HitNotice{
			HitId:             hit.HitId,
			DestroyPlaneIndex: destroyPlaneIndex,
			HitResult:         hitResult,
		}))
		noticeHit.SetData(notice)
		r.Guest.Send(noticeHit)
	} else {
		if !r.Guest.Action {
			return NewErrorResponse(1302), nil
		}
		destroyPlaneIndex, hitResult = r.Owner.MatchRequest.isHited(hit.HitId)
		if r.Owner.MatchRequest.Plane1.IsDied &&
			r.Owner.MatchRequest.Plane2.IsDied &&
			r.Owner.MatchRequest.Plane3.IsDied {
			OverGame(r, "guest")
		} else {
			stop := r.Tm.Stop()
			r.NewTimer(role, 63)
			log.Print("停止计时器结果：", stop)
		}
		notice, _ := json.Marshal(NewSuccessResponse(HitNotice{
			HitId:             hit.HitId,
			DestroyPlaneIndex: destroyPlaneIndex,
			HitResult:         hitResult,
		}))
		noticeHit.SetData(notice)
		r.Owner.Send(noticeHit)

	}

	return NewSuccessResponse(HitResponse{
		HitResult:         hitResult,
		DestroyPlaneIndex: destroyPlaneIndex,
	}), nil
}

/**
planNum 打中的飞机index
		  1 plane1
		  2 plane2
		  3 plane3
hitStatus 1 击中头部
		  2 击中身体
		  3 没有击中
*/
func (oppo *MatchRequest) isHited(hitId int) (destroyPlaneIndex int, hitResult int) {
	for _, v := range oppo.Plane1.BodyPos {
		if hitId == v {
			return 0, 2
		}
	}
	for _, v := range oppo.Plane2.BodyPos {
		if hitId == v {
			return 0, 2
		}
	}
	for _, v := range oppo.Plane3.BodyPos {
		if hitId == v {
			return 0, 2
		}
	}
	switch hitId {
	case oppo.Plane1.HeadId:
		if oppo.Plane1.IsDied {
			return 0, 3
		}
		oppo.Plane1.IsDied = true
		return 1, 1
	case oppo.Plane2.HeadId:
		if oppo.Plane2.IsDied {
			return 0, 3
		}
		oppo.Plane2.IsDied = true
		return 2, 1
	case oppo.Plane3.HeadId:
		if oppo.Plane3.IsDied {
			return 0, 3
		}
		oppo.Plane3.IsDied = true
		return 3, 1
	default:
		return 0, 3
	}

}
