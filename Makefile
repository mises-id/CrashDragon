SRC=$(wildcard *.go)

all: CrashDragon minidump-stackwalk/stackwalker

CrashDragon: $(SRC)
	go build -o CrashDragon $(SRC)

minidump-stackwalk/stackwalker:
	cd minidump-stackwalk && make
