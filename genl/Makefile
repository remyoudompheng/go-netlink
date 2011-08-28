include $(GOROOT)/src/Make.inc

TARG=netlink/genl
GOFILES=\
	ctrl.go\
	taskstats.go\
	taskstats_defs.go

include $(GOROOT)/src/Make.pkg

taskstats_defs.go: taskstats.go
	godefs -g genl taskstats.c | gofmt > taskstats_defs.go

clean: clean-godefs

clean-godefs:
	rm taskstats_defs.go

