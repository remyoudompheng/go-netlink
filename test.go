package main

import (
	"bufio"
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


func main() {
	s, _ := netlink.DialNetlink("route", 0)
	r := bufio.NewReader(s)
	w := bufio.NewWriter(s)
	TestRouteLink(r, w)
	TestRouteAddr(r, w)
}
