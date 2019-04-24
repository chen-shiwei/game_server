package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

var (
	ErrorPacketHeader = errors.New("packet header does not complete")
	ErrorPacketData   = errors.New("packet data does not complete")
)

type Packet struct {
	seq     uint32 // 3 bytes
	ack     uint32 // 3 bytes
	command uint8  // 1 byte
	typ     uint8  // 1 byte
	length  uint64 // 8 bytes too large
	data    []byte
}

func NewPacket() *Packet {
	return &Packet{}
}

func BuildPacketFromHeader(header []byte) *Packet {
	s := &Packet{}
	s.UnpackHeader(header[:16])
	return s
}

func BuildPacket(seq uint32, ack uint32, command uint8, typ uint8) *Packet {
	return &Packet{seq: seq, ack: ack, command: command, typ: typ}
}

func (s *Packet) UnpackHeader(data []byte) error {
	if len(data) < 16 {
		return ErrorPacketHeader
	}
	s.seq = binary.BigEndian.Uint32([]byte{0, data[0], data[1], data[2]})
	s.ack = binary.BigEndian.Uint32([]byte{0, data[3], data[4], data[5]})
	s.command = uint8(data[6])
	s.typ = uint8(data[7])
	s.length = binary.BigEndian.Uint64(data[8:16])
	return nil
}

func (s *Packet) Unpack(data []byte) error {
	if len(data) < 16 {
		return ErrorPacketHeader
	}
	s.UnpackHeader(data[:16])
	if uint64(len(data[16:])) != s.length {
		return ErrorPacketData
	}
	s.data = data[16:]
	return nil
}

func (s *Packet) Pack() []byte {
	buf := bytes.NewBuffer(nil)
	tmpNum32 := make([]byte, 4)
	binary.BigEndian.PutUint32(tmpNum32, s.seq)
	buf.Write(tmpNum32[1:])
	binary.BigEndian.PutUint32(tmpNum32, s.ack)
	buf.Write(tmpNum32[1:])
	buf.WriteByte(byte(s.command))
	buf.WriteByte(byte(s.typ))
	tmpNum64 := make([]byte, 8)
	binary.BigEndian.PutUint64(tmpNum64, s.length)
	buf.Write(tmpNum64)
	header := buf.Bytes()
	return append(header, s.data...)
}

func (s *Packet) String() string {
	return fmt.Sprintf("Seq: %d, Ack: %d, Command: %d, Type: %d, Length: %d",
		s.seq, s.ack, s.command, s.typ, s.length)
}

func (s *Packet) SetSequence(seq uint32) {
	s.seq = seq
}

func (s *Packet) Sequence() uint32 {
	return s.seq
}

func (s *Packet) SetAck(ack uint32) {
	s.ack = ack
}

func (s *Packet) Ack() uint32 {
	return s.ack
}

func (s *Packet) SetCommand(cmd uint8) {
	s.command = cmd
}

func (s *Packet) Command() uint8 {
	return s.command
}

func (s *Packet) SetType(typ uint8) {
	s.typ = typ
}

func (s *Packet) Type() uint8 {
	return s.typ
}

func (s *Packet) Length() uint64 {
	return s.length
}

func (s *Packet) SetData(data []byte) {
	s.length = uint64(len(data))
	if s.length == 0 {
		return
	}
	s.data = data
}

func (s *Packet) Data() []byte {
	return s.data
}
