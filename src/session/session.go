package session

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"protocol"
	"user"
)

type Session struct {
	seq        uint32
	conn       net.Conn
	hashCode   string
	LoginEvent chan bool
	User       *user.User
	WillClosed chan bool
	request    chan *protocol.Packet
	response   chan *protocol.Packet
}

func Create(conn net.Conn) *Session {
	return &Session{
		conn:       conn,
		WillClosed: make(chan bool, 1),
		LoginEvent: make(chan bool, 1),
		request:    make(chan *protocol.Packet),
		response:   make(chan *protocol.Packet)}
}

func (s *Session) SetUser(u *user.User) {
	if u == nil {
		return
	}
	s.User = u
	s.LoginEvent <- true
	close(s.LoginEvent)
}

func (s *Session) Start() {
	go s.Response()
	go s.Request()
}

func (s *Session) RegisterEvent(event string, eventCallback func(*Session) error) {
	if eventCallback == nil {
		return
	}
	switch event {
	case "login":
		go func() {
			<-s.LoginEvent
			eventCallback(s)
		}()
	}
}

func (s *Session) HashCode() string {
	if s.hashCode == "" {
		h := md5.New()
		h.Write([]byte(s.conn.RemoteAddr().String()))
		s.hashCode = hex.EncodeToString(h.Sum(nil))
	}
	return s.hashCode
}

func (s *Session) Send(pkt *protocol.Packet) {
	s.response <- pkt
}

func (s *Session) Receive() *protocol.Packet {
	return <-s.request
}

func (s *Session) ReadHeader() (*protocol.Packet, error) {
	var (
		buf       = make([]byte, 16)
		n   int   = 0
		l   int   = 0
		err error = nil
	)
	for {
		l, err = s.conn.Read(buf[n:])
		if err != nil {
			log.Println("Read Header Failed:", err.Error())
			return nil, err
		}
		n = n + l
		if n == 16 {
			break
		}
	}
	return protocol.BuildPacketFromHeader(buf), nil
}

func (s *Session) ReadContent(length int, pkt *protocol.Packet) error {
	defer func() {
		if errmsg := recover(); errmsg != nil {
			log.Println(errmsg)
		}
	}()
	log.Println("Will Read Content, Length:", length)
	var (
		buf       = make([]byte, length)
		n   int   = 0
		l   int   = 0
		err error = nil
	)
	for {
		l, err = s.conn.Read(buf)
		if err != nil {
			log.Println("Read Content Failed:", err.Error())
		}
		n = n + l
		if n >= length {
			break
		}
	}
	pkt.SetData(buf)
	return nil
}

func (s *Session) Request() {
	for {
		pkt, err := s.ReadHeader()
		if err != nil {
			log.Println("ReadHeader Failed")
			break
		}
		err = s.ReadContent(int(pkt.Length()), pkt)
		if err != nil {
			log.Println("ReadContent Failed")
			break
		}
		s.request <- pkt
	}
	s.Close(false)
}

func (s *Session) Response() {
	for {
		select {
		case packet := <-s.response:
			s.seq++
			packet.SetSequence(s.seq)
			if n, err := s.conn.Write(packet.Pack()); err != nil {
				log.Println("Send Failed:", err.Error())

			} else {
				s.conn.Write(packet.Pack()[n:])
			}
		case <-s.WillClosed:
			break
		}
	}
}

func (s *Session) Close(f bool) {
	s.WillClosed <- true
	if !f {
		s.WillClosed <- true
	}
	close(s.WillClosed)
}

func (s *Session) String() string {
	return fmt.Sprintf("%s %s %v", s.conn.RemoteAddr().String(), s.HashCode(), s.User)
}
