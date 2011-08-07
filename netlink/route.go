package netlink

import (
	"os"
	"net"
	"encoding/binary"
	"syscall"
	"bytes"
	"fmt"
)


func MakeRouteMessage(proto, family int) (msg RawNetlinkMessage) {
	msg.Header.Type = uint16(proto)
	msg.Header.Flags = syscall.NLM_F_DUMP | syscall.NLM_F_REQUEST
	msg.Data = make([]byte, syscall.SizeofRtGenmsg)
	msg.Data[0] = uint8(family)
	return msg
}

type IfMap struct {
	MemStart uint64
	MemEnd   uint64
	BaseAddr uint64
	IRQ      uint16
	DMA      uint8
	Port     uint8
}

type RouteLinkMessage struct {
	Header syscall.NlMsghdr
	IfInfo syscall.IfInfomsg
	// attributes
	Address   net.HardwareAddr
	BcastAddr net.HardwareAddr
	Ifname    string
	MTU       uint32
	LinkType  int32
	QDisc     []byte
	Stats     [23]uint32
	Stats64   [23]uint64
	Map       *IfMap
	ProtInfo  []byte
	Weight    uint32
	Master    uint32
	TxQLen    uint32
	OperState uint8
	LinkMode  uint8
	Ifalias   string
	NumVF     uint32
	Group     uint32
}

const (
	IFLA_NUM_VF  = 0x15
	IFLA_STATS64 = 0x17
	IFLA_GROUP   = 0x1b
)

func ParseRouteLinkMessage(msg syscall.NetlinkMessage) (ParsedNetlinkMessage, os.Error) {
	m := new(RouteLinkMessage)
	m.Header = msg.Header
	buf := bytes.NewBuffer(msg.Data)
	binary.Read(buf, systemEndianness, &m.IfInfo)
	// read link attributes
	for {
		var attr syscall.RtAttr
		er := binary.Read(buf, systemEndianness, &attr)
		dataLen := int(attr.Len) - syscall.SizeofRtAttr
		if er != nil || dataLen > buf.Len() {
			fmt.Println(attr)
			break
		}
		switch attr.Type {
		case syscall.IFLA_ADDRESS: // 1
			readAlignedFromSlice(buf, &m.Address, dataLen)
		case syscall.IFLA_BROADCAST: // 2
			readAlignedFromSlice(buf, &m.BcastAddr, dataLen)
		case syscall.IFLA_IFNAME: // 3
			readAlignedFromSlice(buf, &m.Ifname, dataLen)
		case syscall.IFLA_MTU: // 4
			readAlignedFromSlice(buf, &m.MTU, dataLen)
		case syscall.IFLA_LINK: // 5
			readAlignedFromSlice(buf, &m.LinkType, dataLen)
		case syscall.IFLA_QDISC: // 6
			readAlignedFromSlice(buf, &m.QDisc, dataLen)
		case syscall.IFLA_STATS: // 7
			readAlignedFromSlice(buf, &m.Stats, dataLen)
		case syscall.IFLA_PROTINFO: // 12
			readAlignedFromSlice(buf, &m.ProtInfo, dataLen)
		case syscall.IFLA_TXQLEN: // 13
			readAlignedFromSlice(buf, &m.TxQLen, dataLen)
		case syscall.IFLA_MAP: // 14
			m.Map = new(IfMap)
			readAlignedFromSlice(buf, m.Map, dataLen)
		case syscall.IFLA_LINKMODE: // 17
			readAlignedFromSlice(buf, &m.LinkMode, dataLen)
		case syscall.IFLA_OPERSTATE: // 18
			readAlignedFromSlice(buf, &m.OperState, dataLen)
		case syscall.IFLA_IFALIAS: // 20
			readAlignedFromSlice(buf, &m.Ifalias, dataLen)
		case IFLA_NUM_VF: // 21
			readAlignedFromSlice(buf, &m.NumVF, dataLen)
		case IFLA_STATS64: // 23
			readAlignedFromSlice(buf, &m.Stats64, dataLen)
		case IFLA_GROUP: // 27
			readAlignedFromSlice(buf, &m.Group, dataLen)
		default:
			fmt.Println(attr)
			skipAlignedFromSlice(buf, dataLen)
		}
	}
	return m, nil
}

// address messages
type RouteAddrMessage struct {
	Header syscall.NlMsghdr
	IfAddr syscall.IfAddrmsg
	// attributes
	Address   net.IP
	Local     net.IP
	Label     string
	Broadcast net.IP
	Anycast   net.IP
	Cacheinfo *IfAddrCacheInfo
}

type IfAddrCacheInfo struct {
	Preferred uint32
	Valid     uint32
	CStamp    uint32
	TStamp    uint32
}

func ParseRouteAddrMessage(msg syscall.NetlinkMessage) (ParsedNetlinkMessage, os.Error) {
	m := new(RouteAddrMessage)
	m.Header = msg.Header
	buf := bytes.NewBuffer(msg.Data)

	binary.Read(buf, systemEndianness, &m.IfAddr)
	// read Address attributes
	for {
		var attr syscall.RtAttr
		er := binary.Read(buf, systemEndianness, &attr)
		dataLen := int(attr.Len) - syscall.SizeofRtAttr
		if er != nil || dataLen > buf.Len() {
			break
		}
		switch attr.Type {
		case syscall.IFA_ADDRESS:
			readAlignedFromSlice(buf, &m.Address, dataLen)
		case syscall.IFA_LOCAL:
			readAlignedFromSlice(buf, &m.Local, dataLen)
		case syscall.IFA_LABEL:
			readAlignedFromSlice(buf, &m.Label, dataLen)
		case syscall.IFA_BROADCAST:
			readAlignedFromSlice(buf, &m.Broadcast, dataLen)
		case syscall.IFA_ANYCAST:
			readAlignedFromSlice(buf, &m.Anycast, dataLen)
		case syscall.IFA_CACHEINFO:
			m.Cacheinfo = new(IfAddrCacheInfo)
			readAlignedFromSlice(buf, m.Cacheinfo, dataLen)
		default:
			fmt.Println(attr)
			skipAlignedFromSlice(buf, dataLen)
		}
	}
	return m, nil
}

func ParseRouteMessage(msg syscall.NetlinkMessage) (ParsedNetlinkMessage, os.Error) {
	switch msg.Header.Type {
	case syscall.RTM_NEWADDR, syscall.RTM_GETADDR, syscall.RTM_DELADDR:
		return ParseRouteAddrMessage(msg)
	case syscall.RTM_NEWLINK, syscall.RTM_GETLINK, syscall.RTM_DELLINK:
		return ParseRouteLinkMessage(msg)
	}
	return nil, fmt.Errorf("Invalid Route netlink message type: %d", msg.Header.Type)
}
