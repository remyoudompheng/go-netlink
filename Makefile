include $(GOROOT)/src/Make.inc

TARG=netlink
GOFILES=\
	attributes.go\
	message.go\
	socket.go\
	family_conn.go\
	family_route.go\
	family_generic.go\
	rtnetlink_defs.go\

CLEANFILES+=rtnetlink_defs.go

include $(GOROOT)/src/Make.pkg

rtnetlink_defs.go: rtnetlink.c
	godefs -g netlink rtnetlink.c | gofmt > rtnetlink_defs.go

clean: extra-clean

extra-clean:
	$(MAKE) -C genl clean
	$(MAKE) -C nl80211 clean

all: install
	$(MAKE) -C genl install
	$(MAKE) -C nl80211 install
