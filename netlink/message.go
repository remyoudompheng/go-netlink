package netlink

import (
	"os"
  "bufio"
	"encoding/binary"
	"syscall"
)

var (
	systemEndianness = binary.LittleEndian
)

type NetlinkMsg interface {
  toRawMsg() syscall.NetlinkMessage
}

type RawNetlinkMessage syscall.NetlinkMessage

func (m RawNetlinkMessage)toRawMsg() syscall.NetlinkMessage {
  return syscall.NetlinkMessage(m)
}

// Higher level implementation: let's suppose we're on a little-endian platform

func WriteMessage(s *bufio.Writer, m NetlinkMsg) os.Error {
  msg := m.toRawMsg()
	binary.Write(s, systemEndianness, msg.Header.Len)   // 4 bytes
	binary.Write(s, systemEndianness, msg.Header.Type)  // 2 bytes
	binary.Write(s, systemEndianness, msg.Header.Flags) // 2 bytes
	binary.Write(s, systemEndianness, msg.Header.Seq)   // 4 bytes
	binary.Write(s, systemEndianness, msg.Header.Pid)   // 4 bytes
	_, er := s.Write(msg.Data[:msg.Header.Len-syscall.NLMSG_HDRLEN])
  s.Flush()
	return er
}

func ReadMessage(s *bufio.Reader) (msg syscall.NetlinkMessage, er os.Error) {
	binary.Read(s, systemEndianness, &msg.Header)
	msg.Data = make([]byte, msg.Header.Len-syscall.NLMSG_HDRLEN)
	_, er = s.Read(msg.Data)
	return msg, er
}

