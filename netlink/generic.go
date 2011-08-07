package netlink

import (
  "bytes"
  "encoding/binary"
 "syscall"
)

type GenlMsghdr struct {
	Command  uint8
	Version  uint8
	Reserved uint16
}

// Netlink messages are aligned to 4 bytes multiples
type GenericNetlinkMessage struct {
	Header    syscall.NlMsghdr // 12 bytes
	GenHeader GenlMsghdr       // 4 bytes
  Data      []byte
}

func (msg *GenericNetlinkMessage)toRawMsg() (rawmsg syscall.NetlinkMessage) {
  rawmsg.Header = msg.Header
  rawmsg.Data = make([]byte, 4 + len(msg.Data))
  w := bytes.NewBuffer(rawmsg.Data)
  binary.Write(w, systemEndianness, msg.GenHeader)
  w.Write(msg.Data)
  return rawmsg
}


