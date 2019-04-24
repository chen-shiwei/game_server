package plane

import (
	"encoding/json"
	"log"
	"protocol"
)

func OverGame(r *Room, winer string) {
	log.Print("结算", winer)
	var ogn OverGameNotice
	if winer == "owner" {
		ogn.Players = append(ogn.Players, OverPlayer{*r.Owner.User, countScore(&r.Guest), 1})
		ogn.Players = append(ogn.Players, OverPlayer{*r.Guest.User, countScore(&r.Owner), 2})
	} else {
		ogn.Players = append(ogn.Players, OverPlayer{*r.Owner.User, countScore(&r.Guest), 2})
		ogn.Players = append(ogn.Players, OverPlayer{*r.Guest.User, countScore(&r.Owner), 1})
	}
	data, _ := json.Marshal(NewSuccessResponse(ogn))
	packet := protocol.BuildPacket(0, 0, protocol.CMD_PLAN_OVER, protocol.PACKET_JSON)
	packet.SetData(data)
	r.Guest.Send(packet)
	r.Owner.Send(packet)
	stop := r.Tm.Stop()
	log.Print("关闭最后一个定时:", stop)
	r.DestoryRoom()
	return
}

func countScore(data *FightData) (getScore int) {
	if data.Plane1.IsDied {
		getScore++
	}
	if data.Plane2.IsDied {
		getScore++
	}
	if data.Plane2.IsDied {
		getScore++
	}
	return
}
