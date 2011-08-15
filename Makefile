godag: netlink/rtnetlink_defs.go netlink/genl/taskstats_defs.go
	gd .
	gd -I. -o proc bin_proc
	gd -I. -o taskstats bin_taskstats

.PHONY: godag

netlink/rtnetlink_defs.go:
	godefs -g netlink netlink/rtnetlink.c | gofmt > netlink/rtnetlink_defs.go

netlink/genl/taskstats_defs.go:
	godefs -g genl netlink/genl/taskstats.c | gofmt > netlink/genl/taskstats_defs.go

