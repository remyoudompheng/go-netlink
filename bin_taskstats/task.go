package main

import (
	"os"
	"bufio"
	"bytes"
	"fmt"
	"netlink"
	"netlink/genl"
	"syscall"
	"json"
)

func MakeCmdMessage() (msg netlink.GenericNetlinkMessage) {
  msg.Header.Type = 23
  msg.Header.Flags = syscall.NLM_F_REQUEST | syscall.NLM_F_DUMP
  msg.GenHeader.Command = genl.TASKSTATS_CMD_GET
  msg.GenHeader.Version = genl.TASKSTATS_GENL_VERSION
  buf := bytes.NewWriter([]byte{})
  binary.Write(buf, binary.LittleEndian, uint32(os.Getpid()))
  msg.Data = buf.Bytes
  return msg
}

func TestTaskStats(r *bufio.Reader, w *bufio.Writer) {
	msg := MakeCmdMessage()
	netlink.WriteMessage(w, &msg)
	b := bytes.NewBuffer([]byte{})
	buf := bufio.NewWriter(b)
	netlink.WriteMessage(buf, &msg)
	fmt.Println(b.Bytes())
	fmt.Printf("%#v\n", msg)

	for {
		resp, _ := netlink.ReadMessage(r)
		fmt.Printf("%#v\n", resp)
		parsedmsg, _ := genl.ParseGenlTaskstatsMsg(resp)
		switch m := parsedmsg.(type) {
		case nil:
			return
		case netlink.ErrorMessage:
			msg_s, _ := json.MarshalIndent(m, "", "  ")
			fmt.Printf("ErrorMsg = %s\n%s\n", msg_s, os.NewSyscallError("netlink", int(-m.Errno)))
			break
		default:
			msg_s, _ := json.MarshalIndent(m, "", "  ")
			fmt.Printf("Taskstats = %s\n", msg_s)
		}
	}
}

func main() {
	s, _ := netlink.DialNetlink("generic", 0)
	r := bufio.NewReader(s)
	w := bufio.NewWriter(s)
	TestTaskStats(r, w)
}
