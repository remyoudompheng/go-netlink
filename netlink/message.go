package netlink

import (
	"os"
	"bufio"
	"encoding/binary"
	"syscall"
	"bytes"
)

var (
	systemEndianness = binary.LittleEndian
	globalSeq        = uint32(0)
)

type NetlinkMsg interface {
	toRawMsg() syscall.NetlinkMessage
}

type RawNetlinkMessage syscall.NetlinkMessage

func (m RawNetlinkMessage) toRawMsg() syscall.NetlinkMessage {
	return syscall.NetlinkMessage(m)
}

// Higher level implementation: let's suppose we're on a little-endian platform

func WriteMessage(s *bufio.Writer, m NetlinkMsg) os.Error {
	msg := m.toRawMsg()
	msg.Header.Len = uint32(syscall.NLMSG_HDRLEN + len(msg.Data))
	msg.Header.Seq = globalSeq
	globalSeq++
	binary.Write(s, systemEndianness, msg.Header.Len)   // 4 bytes
	binary.Write(s, systemEndianness, msg.Header.Type)  // 2 bytes
	binary.Write(s, systemEndianness, msg.Header.Flags) // 2 bytes
	binary.Write(s, systemEndianness, msg.Header.Seq)   // 4 bytes
	binary.Write(s, systemEndianness, msg.Header.Pid)   // 4 bytes
	_, er := s.Write(msg.Data)
	s.Flush()
	return er
}

func ReadMessage(s *bufio.Reader) (msg syscall.NetlinkMessage, er os.Error) {
	binary.Read(s, systemEndianness, &msg.Header)
	msg.Data = make([]byte, msg.Header.Len-syscall.NLMSG_HDRLEN)
	_, er = s.Read(msg.Data)
	return msg, er
}

type Attr struct {
	Len  uint16
	Type uint16
}

type ParsedNetlinkMessage interface{}

type ErrorMessage struct {
	Header      syscall.NlMsghdr
	Errno       int32
	WrongHeader syscall.NlMsghdr
}

func ParseErrorMessage(msg syscall.NetlinkMessage) ParsedNetlinkMessage {
	var parsed ErrorMessage
	parsed.Header = msg.Header
	buf := bytes.NewBuffer(msg.Data)
	binary.Read(buf, systemEndianness, &parsed.Errno)
	binary.Read(buf, systemEndianness, &parsed.WrongHeader)
	return parsed
}
