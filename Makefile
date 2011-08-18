DEFS=\
     netlink/rtnetlink_defs.go\
     netlink/genl/taskstats_defs.go\
     nl80211/nl80211_defs.go

godag: $(DEFS)
	gd .
	gd -I. -o proc bin_proc
	gd -I. -o taskstats bin_taskstats

.PHONY: godag $(DEFS)

netlink/rtnetlink_defs.go:
	godefs -g netlink netlink/rtnetlink.c | gofmt > netlink/rtnetlink_defs.go

netlink/genl/taskstats_defs.go:
	godefs -g genl netlink/genl/taskstats.c | gofmt > netlink/genl/taskstats_defs.go

nl80211/nl80211_defs.go:
	./gen_nl80211.sh | godefs -gnl80211 - > nl80211/nl80211_defs.go
