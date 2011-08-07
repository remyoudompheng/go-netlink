package netlink

import (
	"os"
	"fmt"
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
	Header    syscall.NlMsghdr // 16 bytes
	GenHeader GenlMsghdr       // 4 bytes
	Data      []byte
}

func (msg *GenericNetlinkMessage) toRawMsg() (rawmsg syscall.NetlinkMessage) {
	rawmsg.Header = msg.Header
	w := bytes.NewBuffer([]byte{})
	binary.Write(w, systemEndianness, msg.GenHeader)
	w.Write(msg.Data)
	rawmsg.Data = w.Bytes()
	return rawmsg
}

func ParseGenlMessage(msg syscall.NetlinkMessage) (genmsg GenericNetlinkMessage, er os.Error) {
	genmsg.Header = msg.Header
	buf := bytes.NewBuffer(msg.Data)
	binary.Read(buf, systemEndianness, &genmsg.GenHeader)
	genmsg.Data = buf.Bytes()
	return genmsg, nil
}

// Control messages for Generic Netlink interface

type GenlCtrlMessage struct {
	Header      syscall.NlMsghdr // 16 bytes
	GenHeader   GenlMsghdr       // 4 bytes
	FamilyID    uint16
	FamilyName  string
	Version     uint32
	HdrSize     uint32
	MaxAttr     uint32
	Ops         []byte
	McastGroups []byte
}

const (
	CTRL_CMD_UNSPEC = iota
	CTRL_CMD_NEWFAMILY
	CTRL_CMD_DELFAMILY
	CTRL_CMD_GETFAMILY
	CTRL_CMD_NEWOPS
	CTRL_CMD_DELOPS
	CTRL_CMD_GETOPS
	CTRL_CMD_NEWMCAST_GRP
	CTRL_CMD_DELMCAST_GRP
	CTRL_CMD_GETMCAST_GRP
)

const (
	CTRL_ATTR_UNSPEC = iota
	CTRL_ATTR_FAMILY_ID
	CTRL_ATTR_FAMILY_NAME
	CTRL_ATTR_VERSION
	CTRL_ATTR_HDRSIZE
	CTRL_ATTR_MAXATTR
	CTRL_ATTR_OPS
	CTRL_ATTR_MCAST_GROUPS
)

const (
	GENL_ID_CTRL      = syscall.NLMSG_MIN_TYPE
	GENL_VERSION_CTRL = 0x1
)

func MakeGenCtrlCmd(cmd uint8) (msg GenericNetlinkMessage) {
	msg.Header.Type = GENL_ID_CTRL
	msg.Header.Flags = syscall.NLM_F_REQUEST | syscall.NLM_F_DUMP
	msg.GenHeader.Command = cmd
	msg.GenHeader.Version = GENL_VERSION_CTRL
	return msg
}

func ParseGenlFamilyMessage(msg syscall.NetlinkMessage) (ParsedNetlinkMessage, os.Error) {
	m := new(GenlCtrlMessage)
	m.Header = msg.Header
	switch m.Header.Type {
	case syscall.NLMSG_DONE:
		return nil, nil
	case syscall.NLMSG_ERROR:
		return ParseErrorMessage(msg), nil
	}
	buf := bytes.NewBuffer(msg.Data)
	binary.Read(buf, systemEndianness, &m.GenHeader)

	// read Family attributes
	for {
		var attr syscall.RtAttr
		er := binary.Read(buf, systemEndianness, &attr)
		dataLen := int(attr.Len) - syscall.SizeofRtAttr
		if er != nil || dataLen > buf.Len() {
			break
		}
		switch attr.Type {
		case CTRL_ATTR_FAMILY_ID:
			readAlignedFromSlice(buf, &m.FamilyID, dataLen)
		case CTRL_ATTR_FAMILY_NAME:
			readAlignedFromSlice(buf, &m.FamilyName, dataLen)
		case CTRL_ATTR_VERSION:
			readAlignedFromSlice(buf, &m.Version, dataLen)
		case CTRL_ATTR_HDRSIZE:
			readAlignedFromSlice(buf, &m.HdrSize, dataLen)
		case CTRL_ATTR_MAXATTR:
			readAlignedFromSlice(buf, &m.MaxAttr, dataLen)
		case CTRL_ATTR_OPS:
			readAlignedFromSlice(buf, &m.Ops, dataLen)
		default:
			fmt.Println(attr)
			skipAlignedFromSlice(buf, dataLen)
		}
	}
	return m, nil
}
