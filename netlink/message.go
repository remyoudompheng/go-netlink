package netlink

import (
	"os"
	"bufio"
	"encoding/binary"
	"syscall"
	"bytes"
)

var (
	SystemEndianness = binary.LittleEndian
	globalSeq        = uint32(0)
)

type NetlinkMsg interface {
	toRawMsg() syscall.NetlinkMessage
}

// A netlink message with unparsed data
type RawNetlinkMessage syscall.NetlinkMessage

func (m RawNetlinkMessage) toRawMsg() syscall.NetlinkMessage {
	return syscall.NetlinkMessage(m)
}

// Higher level implementation: let's suppose we're on a little-endian platform

// Write a netlink message to a socket
func WriteMessage(s *bufio.Writer, m NetlinkMsg) os.Error {
	msg := m.toRawMsg()
	msg.Header.Len = uint32(syscall.NLMSG_HDRLEN + len(msg.Data))
	msg.Header.Seq = globalSeq
	globalSeq++
	binary.Write(s, SystemEndianness, msg.Header)   // 16 bytes
	_, er := s.Write(msg.Data)
	s.Flush()
	return er
}

// Reads a netlink message from a socket
func ReadMessage(s *bufio.Reader) (msg syscall.NetlinkMessage, er os.Error) {
	binary.Read(s, SystemEndianness, &msg.Header)
	msg.Data = make([]byte, msg.Header.Len-syscall.NLMSG_HDRLEN)
	_, er = s.Read(msg.Data)
	return msg, er
}

type Attr struct {
	Len  uint16
	Type uint16
}

type ParsedNetlinkMessage interface{}

// The structure of netlink error messages
type ErrorMessage struct {
	Header      syscall.NlMsghdr
	Errno       int32
	WrongHeader syscall.NlMsghdr
}

// Parses a netlink error message
func ParseErrorMessage(msg syscall.NetlinkMessage) ErrorMessage {
	var parsed ErrorMessage
	parsed.Header = msg.Header
	buf := bytes.NewBuffer(msg.Data)
	binary.Read(buf, SystemEndianness, &parsed.Errno)
	binary.Read(buf, SystemEndianness, &parsed.WrongHeader)
	return parsed
}
