SASSC   ?= sassc
GO      ?= go

GO_SRC   = ./server/$(wildcard *.go)

SASSCFLAGS ?= -t compressed

all: crashdragon upload_syms minidump-stackwalk/stackwalker

crashdragon: $(GO_SRC) assets/stylesheets/app.css
	$(GO) build -o bin/crashdragon $(GO_SRC)

upload_syms: upload_syms/main.go
	$(GO) build -o bin/upload_syms upload_syms/main.go

assets/stylesheets/app.css:
	$(SASSC) $(SASSCFLAGS) $(@D)/app.scss > $@.tmp && mv $@.tmp $@

minidump-stackwalk/stackwalker:
	cd minidump-stackwalk && $(MAKE)

clean:
	rm -f bin/crashdragon
	rm -f bin/upload_syms
	rm -f assets/stylesheets/app.css.tmp
	rm -f assets/stylesheets/app.css
	cd minidump-stackwalk && $(MAKE) distclean

.PHONY: clean all
