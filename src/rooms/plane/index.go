package plane

import (
	"controller"
	"log"
	"protocol"
	"room"
	"session"
)

type PlanRoom struct {
	sessionList *session.Session
}

func init() {
	//注册游戏入口
	err := room.RegisterGameEntry("plan", "1002", &PlanRoom{}, 10)
	if err != nil {
		log.Print(err.Error())
	}
	controller.Register(protocol.CMD_PLAN_MATCH, MatchAction)
	controller.Register(protocol.CMD_PLAN_HIT, HitAction)
	controller.Register(protocol.CMD_PLAN_CANCEL_MATCH, MatchActionCancel)
}

func (p *PlanRoom) New() room.GameHandler {
	return new(PlanRoom)
}

func (p *PlanRoom) Join(sess *session.Session, back bool) (interface{}, error) {
	return nil, nil
}

func (p *PlanRoom) Leave(sess *session.Session, quit bool) (interface{}, error) {
	return nil, nil
}

func (p *PlanRoom) Request(sess *session.Session, data []byte) (interface{}, error) {
	return nil, nil
}

func (s *PlanRoom) Match(sess *session.Session, data []byte) bool {
	return true
}

func (s *PlanRoom) Create(sess *session.Session, data []byte) error {
	return nil
}
