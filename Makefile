rootdir = $(shell printenv PWD)

all:
	cd src; go build -o $(rootdir)/outputs/orange-cat-server
