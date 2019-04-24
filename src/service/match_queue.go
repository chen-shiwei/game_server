package service

import (
	"session"
	"sync"
)

var (
	matchQueueList     = make([]*MatchQueue, 0)
	matchQueueListLock = sync.Mutex{}
)

type MatchQueue struct {
	Num   int
	Queue chan *session.Session
}

func CreateMatchQueue(num int) *MatchQueue {
	mq := &MatchQueue{Num: num, Queue: make(chan *session.Session, num)}
	matchQueueList = append(matchQueueList, mq)
	return mq
}

func (s *MatchQueue) Push(sess *session.Session) {
	s.Queue <- sess
	s.Num--
	if s.Num <= 0 {
		close(s.Queue)
		matchQueueList = matchQueueList[1:]
	}
}

func AttachMatchQueue(sess *session.Session) {
	if len(matchQueueList) == 0 {
		return
	}
	matchQueueListLock.Lock()
	defer matchQueueListLock.Unlock()
	matchQueueList[0].Push(sess)
}
