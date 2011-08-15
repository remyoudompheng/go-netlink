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

func TestRouteLink(r *bufio.Reader, w *bufio.Writer) {
	msg := netlink.MakeRouteMessage(syscall.RTM_GETLINK, syscall.AF_UNSPEC)
	netlink.WriteMessage(w, msg)
	for {
		resp, _ := netlink.ReadMessage(r)
		parsedmsg, er := netlink.ParseRouteMessage(resp)
		if parsedmsg == nil {
			break
		}
		msg_s, _ := json.MarshalIndent(parsedmsg, "", "  ")
		fmt.Printf("Errmsg: %#v\nLinkmsg = %s\n", er, msg_s)
	}
}

func TestRouteAddr(r *bufio.Reader, w *bufio.Writer) {
	msg := netlink.MakeRouteMessage(syscall.RTM_GETADDR, syscall.AF_UNSPEC)
	netlink.WriteMessage(w, msg)
	for {
		resp, _ := netlink.ReadMessage(r)
		parsedmsg, er := netlink.ParseRouteMessage(resp)
		if parsedmsg == nil {
			break
		}
		if er != nil {
			fmt.Println(resp)
		}
		msg_s, _ := json.MarshalIndent(parsedmsg, "", "  ")
		fmt.Printf("Errmsg: %#v\nAddrMsg = %s\n", er, msg_s)
	}
}

func TestGenericFamily(r *bufio.Reader, w *bufio.Writer) {
	msg := genl.MakeGenCtrlCmd(genl.CTRL_CMD_GETFAMILY)
	netlink.WriteMessage(w, &msg)
	b := bytes.NewBuffer([]byte{})
	buf := bufio.NewWriter(b)
	netlink.WriteMessage(buf, &msg)
	fmt.Println(b.Bytes())
	fmt.Printf("%#v\n", msg)

	for {
		resp, _ := netlink.ReadMessage(r)
		parsedmsg, _ := genl.ParseGenlFamilyMessage(resp)
		switch m := parsedmsg.(type) {
		case nil:
			return
		case netlink.ErrorMessage:
			msg_s, _ := json.MarshalIndent(m, "", "  ")
			fmt.Printf("ErrorMsg = %s\n%s\n", msg_s, os.NewSyscallError("netlink", int(-m.Errno)))
			break
		default:
			msg_s, _ := json.MarshalIndent(m, "", "  ")
			fmt.Printf("GenlFamily = %s\n", msg_s)
		}
	}
}

func main() {
	s, _ := netlink.DialNetlink("route", 0)
	r := bufio.NewReader(s)
	w := bufio.NewWriter(s)
	TestRouteLink(r, w)
	TestRouteAddr(r, w)

	// NETLINK_GENERIC tests
	s, _ = netlink.DialNetlink("generic", 0)
	r = bufio.NewReader(s)
	w = bufio.NewWriter(s)
	fmt.Println("Testing generic family messages")
	TestGenericFamily(r, w)
}
