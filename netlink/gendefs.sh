#!/bin/bash

godefs -g netlink rtnetlink.c | gofmt > rtnetlink_defs.go

