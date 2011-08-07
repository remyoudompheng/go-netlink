package main

import (
	"os"
	"bufio"
	"bytes"
	"fmt"
	"netlink"
	"syscall"
)

func TestRouteLink(r *bufio.Reader, w *bufio.Writer) {
	msg := netlink.MakeRouteMessage(syscall.RTM_GETLINK, syscall.AF_UNSPEC)
	netlink.WriteMessage(w, msg)
	for {
		resp, _ := netlink.ReadMessage(r)
		parsedmsg, _ := netlink.ParseRouteMessage(resp)
		if parsedmsg == nil {
			break
		}
		fmt.Printf("%#v\n", parsedmsg)
	}
}

func TestRouteAddr(r *bufio.Reader, w *bufio.Writer) {
	msg := netlink.MakeRouteMessage(syscall.RTM_GETADDR, syscall.AF_UNSPEC)
	netlink.WriteMessage(w, msg)
	for {
		resp, _ := netlink.ReadMessage(r)
		parsedmsg, _ := netlink.ParseRouteMessage(resp)
		if parsedmsg == nil {
			break
		}
		fmt.Printf("%#v\n", parsedmsg)
	}
}

func TestGenericFamily(r *bufio.Reader, w *bufio.Writer) {
	msg := netlink.MakeGenCtrlCmd(netlink.CTRL_CMD_GETFAMILY)
	netlink.WriteMessage(w, &msg)
	b := bytes.NewBuffer([]byte{})
	buf := bufio.NewWriter(b)
	netlink.WriteMessage(buf, &msg)
	fmt.Println(b.Bytes())
	fmt.Printf("%#v\n", msg)

	for {
		resp, _ := netlink.ReadMessage(r)
		parsedmsg, _ := netlink.ParseGenlFamilyMessage(resp)
		switch m := parsedmsg.(type) {
		case netlink.ErrorMessage:
			fmt.Printf("%#v, %s\n", m, os.NewSyscallError("netlink", int(-m.Errno)))
			break
		default:
			fmt.Printf("%#v\n", m)
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
