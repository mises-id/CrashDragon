SRC=$(wildcard *.go)
.PHONY: clean all

all: CrashDragon upload_syms minidump-stackwalk/stackwalker

CrashDragon: $(SRC)
	go build -o bin/CrashDragon $(SRC)

upload_syms: $(SRC)
	go build -o bin/upload_syms upload_syms/main.go

minidump-stackwalk/stackwalker:
	cd minidump-stackwalk && make

clean:
	rm -rf bin/
	cd minidump-stackwalk && make distclean
