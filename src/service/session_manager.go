package service

import (
	"controller"
	"encoding/json"
	"fmt"
	"log"
	"protocol"
	"session"
	"time"
)

type SessionManager struct {
	onlineNumber uint64
	sessions     map[uint64]*session.Session
}

var (
	manager *SessionManager
)

func NewSessionManager() *SessionManager {
	if manager == nil {
		manager = &SessionManager{sessions: make(map[uint64]*session.Session)}
	}
	return manager
}

func GetSessionManager() *SessionManager {
	return NewSessionManager()
}

func (s *SessionManager) Process(sess *session.Session, stop <-chan bool) error {
	s.onlineNumber++
	defer func(sid uint64) {
		s.onlineNumber--
		delete(s.sessions, sid)
	}(s.onlineNumber)

	sess.Start()

	sess.RegisterEvent("login", func(s *session.Session) error {
		AttachMatchQueue(s)
		return nil
	})

	for {
		select {
		case <-sess.WillClosed:
			break
		case <-stop:
			sess.Close(true)
			log.Println("get a STOP sig")
			break
		default:
			var pkt *protocol.Packet
			var err error
			packet := sess.Receive()
			pkt, err = controller.Do(sess, packet)
			if err != nil {
				pkt = new(protocol.Packet)
				log.Println("Controller process failed:", err.Error())
				pkt.SetAck(packet.Ack())
				pkt.SetCommand(packet.Command())
				pkt.SetData([]byte(err.Error()))
			}
			sess.Send(pkt)
		}
	}
	log.Println("session destroy")
	return nil
}

func (s *SessionManager) Match(num uint16, matcher session.SessionMatcher, finishCallback func(result []*session.Session) error) []*session.Session {
	if num < 1 {
		return nil
	}
	sessionList := make([]*session.Session, 0)
	for _, sess := range s.sessions {
		fmt.Println("Begin Match:", sess)
		if matched, err := matcher.Match(sess); err == nil {
			fmt.Println("Matched", sess)
			sessionList = append(sessionList, sess)
			pkt := protocol.NewPacket()
			pkt.SetCommand(protocol.CMD_RESPONSE_MATCH)
			pkt.SetType(protocol.PACKET_JSON)
			d, _ := json.Marshal(matched)
			pkt.SetData(d)
			sess.Send(pkt)
			num--
		}
		if num <= 0 {
			break
		}
	}

	if len(sessionList) < int(num) {
		timer := time.NewTimer(60 * time.Second)
		mq := CreateMatchQueue(int(num) - len(sessionList))

		for {
			select {
			case <-timer.C:
				// timeout error
				fmt.Println("match timeout")
				timer.Stop()
				goto MATCH_FINISH
			case ns, ok := <-mq.Queue:
				if !ok {
					goto MATCH_FINISH
				}
				fmt.Println("Matched:", ns)
				sessionList = append(sessionList, ns)
				if len(sessionList) >= int(num) {
					goto MATCH_FINISH
				}
			}
		}
	}
MATCH_FINISH:
	if finishCallback != nil {
		finishCallback(sessionList)
	}
	return sessionList
}
