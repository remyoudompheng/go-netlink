package netlink

import (
	"os"
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

type CtrlOp struct {
	ID    uint32 `netlink:"1" type:"fixed"` // CTRL_ATTR_OP_ID
	Flags uint32 `netlink:"2" type:"fixed"` // CTRL_ATTR_OP_FLAGS
}

type CtrlMcastGroup struct {
	Name string `netlink:"1" type:"string"` // CTRL_ATTR_MCAST_GRP_NAME
	ID   uint32 `netlink:"2" type:"fixed"`  // CTRL_ATTR_MCAST_GRP_ID
}

type GenlCtrlMessage struct {
	Header      syscall.NlMsghdr // 16 bytes
	GenHeader   GenlMsghdr       // 4 bytes
	FamilyID    uint16           `netlink:"1" type:"fixed"`      // CTRL_ATTR_FAMILY_ID
	FamilyName  string           `netlink:"2" type:"string"`     // CTRL_ATTR_FAMILY_NAME
	Version     uint32           `netlink:"3" type:"fixed"`      // CTRL_ATTR_VERSION
	HdrSize     uint32           `netlink:"4" type:"fixed"`      // CTRL_ATTR_HDR_SIZE
	MaxAttr     uint32           `netlink:"5" type:"fixed"`      // CTRL_ATTR_MAXATTR
	Ops         []CtrlOp         `netlink:"6" type:"nestedlist"` // CTRL_ATTR_OPS
	McastGroups []CtrlMcastGroup `netlink:"7" type:"nestedlist"` // CTRL_ATTR_MCAST_GROUPS
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
	er := readManyAttributes(buf, m)
	return m, er
}
