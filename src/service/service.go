package service

import (
	"log"
	"net"
	"session"
	"strings"
	"sync"
	"time"
)

type Service struct {
	port      string
	listener  net.Listener
	exiting   bool
	forceExit bool
}

var (
	stopChannel = make([]chan bool, 0)
)

func Create(port string) *Service {
	s := new(Service)
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	s.port = port
	return s
}

func (s *Service) Stop(force bool) {
	s.forceExit = force
	s.exiting = true
}

func (s *Service) Start(manager *SessionManager) error {
	var (
		wg  sync.WaitGroup
		err error
	)
	if s.listener, err = net.Listen("tcp", s.port); err != nil {
		log.Println("Start Failed:", err.Error())
		panic(err.Error())
	}

	defer s.listener.Close()

	log.Println("Listened on", s.port)

	for {
		if s.exiting {
			break
		}
		if conn, err := s.listener.Accept(); err != nil {
			log.Println("Accept Failed:", err.Error())
			continue
		} else {
			go func(c net.Conn) {
				wg.Add(1)
				defer wg.Done()
				defer func() {
					if err := c.Close(); err != nil {
						log.Println("Close failed:", err.Error())
					}
				}()
				stop := make(chan bool)
				stopChannel = append(stopChannel, stop)
				if err = manager.Process(session.Create(c), stop); err != nil {
					log.Println("manager start failed:", err.Error())
				}
			}(conn)
		}
	}
	if s.forceExit {
		log.Println("Force Exit")
		return nil
	}
	log.Println("Service Exiting at", time.Now(), "...")
	go func() {
		wg.Add(1)
		defer wg.Done()
		for _, s := range stopChannel {
			s <- true
		}
	}()
	wg.Wait()
	return nil
}
