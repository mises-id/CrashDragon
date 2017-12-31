SASSC   ?= sassc
GO      ?= go

GO_SRC   = ./server/$(wildcard *.go)

SASSCFLAGS ?= -t compressed

all: crashdragon upload_syms build/bin/minidump_stackwalk

crashdragon: $(GO_SRC) assets/stylesheets/app.css
	$(GO) build -o bin/crashdragon $(GO_SRC)

upload_syms: upload_syms/main.go
	$(GO) build -o bin/upload_syms upload_syms/main.go

assets/stylesheets/app.css:
	$(SASSC) $(SASSCFLAGS) $(@D)/app.scss > $@.tmp && mv $@.tmp $@

build/bin/minidump_stackwalk:
	cd breakpad && ./autogen.sh && ./configure --prefix="$(CURDIR)/build" && $(MAKE) install

clean:
	rm -f bin/crashdragon
	rm -f bin/upload_syms
	rm -f assets/stylesheets/app.css.tmp
	rm -f assets/stylesheets/app.css
	rm -rf build/
	cd breakpad && $(MAKE) distclean

.PHONY: clean all
