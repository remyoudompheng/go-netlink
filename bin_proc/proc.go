package main

import (
	"os"
	"bufio"
	"fmt"
	"syscall"
	"netlink"
	"json"
	"bytes"
	"encoding/binary"
)

func MakeProcConnectorMsg() netlink.ConnMessage {
	var msg netlink.ConnMessage
	msg.Header.Type = syscall.NLMSG_DONE
	msg.Header.Flags = 0
	msg.Header.Pid = uint32(os.Getpid())
	msg.ConnHdr.Id = netlink.ConnMsgid{Idx: netlink.CN_IDX_PROC, Val: netlink.CN_VAL_PROC}
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, int32(netlink.PROC_CN_MCAST_LISTEN))
	msg.Data = buf.Bytes()
	return msg
}

func get_events(r *bufio.Reader, w *bufio.Writer) {
	msg := MakeProcConnectorMsg()
	netlink.WriteMessage(w, &msg)
	for {
		resp, er := netlink.ReadMessage(r)
		if er != nil {
			fmt.Println(er)
			break
		}

		cnmsg, er := netlink.ParseConnMessage(resp)
		if er != nil {
			fmt.Println(er)
			break
		}

		msg_s, er := json.MarshalIndent(cnmsg, "", "  ")
		fmt.Printf("%s\n", msg_s)
		if resp.Header.Type == syscall.NLMSG_DONE {
			return
		}
	}
}

func main() {
	s, er := netlink.DialNetlink("conn", netlink.CN_IDX_PROC)
	if er != nil {
		fmt.Println(er)
		return
	}
	r := bufio.NewReader(s)
	w := bufio.NewWriter(s)
	for {
		get_events(r, w)
	}
}
