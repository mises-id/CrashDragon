SRC=$(wildcard *.go)
.PHONY: clean all

all: CrashDragon upload_syms minidump-stackwalk/stackwalker

CrashDragon: $(SRC) assets/stylesheets/app.css
	go build -o bin/CrashDragon $(SRC)

upload_syms: upload_syms/main.go
	go build -o bin/upload_syms upload_syms/main.go

assets/stylesheets/app.css:
	sassc -t compressed $(@D)/app.scss > $@.tmp && mv $@.tmp $@

minidump-stackwalk/stackwalker:
	cd minidump-stackwalk && make

clean:
	rm -rf bin/
	rm -f assets/stylesheets/app.css
	cd minidump-stackwalk && make distclean
