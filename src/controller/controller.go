package controller

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"protocol"
	"session"
	"user"
)

type Controller func(*session.Session, *protocol.Packet) (interface{}, error)

var (
	controller = make(map[uint8]Controller)

	ErrorControllerRegister = errors.New("does not support")
	ErrorControllerExists   = errors.New("controller exists")
)

func init() {
	Register(protocol.CMD_ECHO, func(_ *session.Session, pkt *protocol.Packet) (interface{}, error) {
		return pkt, nil
	})
	Register(protocol.CMD_USER_LOGIN, func(sess *session.Session, pkt *protocol.Packet) (interface{}, error) {
		var lr = new(user.LoginRequest)
		if err := json.Unmarshal(pkt.Data(), lr); err != nil {
			return user.ErrorParamSyntax, nil
		}
		UserInfo, err := user.GetUserLass(lr.Token)
		if err != nil {
			return user.ErrorParamSyntax, nil
		}
		sess.SetUser(UserInfo)
		return UserInfo, nil
	})

}

func Do(sess *session.Session, packet *protocol.Packet) (*protocol.Packet, error) {
	if r, ok := controller[packet.Command()]; !ok {
		return nil, ErrorControllerRegister
	} else {
		result, err := r(sess, packet)
		if err != nil {
			pkt := protocol.BuildPacket(0, packet.Sequence(), 0, protocol.PACKET_STRING)
			pkt.SetData([]byte(err.Error()))
			return pkt, err
		} else {
			pkt := BuildPacket(result, packet)
			return pkt, err
		}
	}
}

func Register(command uint8, fn Controller) error {
	if _, ok := controller[command]; ok {
		return ErrorControllerExists
	}
	controller[command] = fn
	return nil
}

func BuildPacket(data interface{}, packet *protocol.Packet) *protocol.Packet {
	if data, ok := data.(*protocol.Packet); ok {
		data.SetAck(packet.Ack())
		data.SetCommand(packet.Command())
		return data
	}
	pkt := protocol.NewPacket()
	pkt.SetAck(packet.Ack())
	pkt.SetCommand(packet.Command())
	switch data.(type) {
	case string:
		pkt.SetData([]byte(data.(string)))
		pkt.SetType(protocol.PACKET_STRING)
	case int, int64, int32, int16:
		tmpNum64 := make([]byte, 8)
		binary.BigEndian.PutUint64(tmpNum64, data.(uint64))
		pkt.SetData(tmpNum64)
		pkt.SetType(protocol.PACKET_NUMBER)
	case uint, uint64, uint32, uint16:
		tmpNum64 := make([]byte, 8)
		binary.BigEndian.PutUint64(tmpNum64, data.(uint64))
		pkt.SetData(tmpNum64)
		pkt.SetType(protocol.PACKET_UNSIGNED_NUMBER)
	case bool:
		if data.(bool) {
			pkt.SetData([]byte{1})
		} else {
			pkt.SetData([]byte{0})
		}
		pkt.SetType(protocol.PACKET_BOOL)
	default:
		if d, err := json.Marshal(data); err != nil {
			pkt.SetData([]byte(err.Error()))
			pkt.SetType(protocol.PACKET_STRING)
		} else {
			pkt.SetData(d)
			pkt.SetType(protocol.PACKET_JSON)
		}
	}
	return pkt

}
