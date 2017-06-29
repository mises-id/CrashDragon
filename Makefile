SRC=$(wildcard *.go)
.PHONY: clean all

all: CrashDragon minidump-stackwalk/stackwalker

CrashDragon: $(SRC)
	go build -o bin/CrashDragon $(SRC)

minidump-stackwalk/stackwalker:
	cd minidump-stackwalk && make

clean:
	rm -rf bin/
	cd minidump-stackwalk && make distclean
