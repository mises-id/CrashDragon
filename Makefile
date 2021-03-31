SASSC   ?= sassc
GO      ?= go
INSTALL ?= install

prefix = /usr/local
exec_prefix = $(prefix)

bindir = $(exec_prefix)/bin
sysconfdir = $(prefix)/etc
datarootdir = $(prefix)/share
datadir = $(datarootdir)
docdir = $(datarootdir)/doc/crashdragon

INSTALL_PROGRAM = $(INSTALL) -c
INSTALL_DATA = $(INSTALL) -c -m 644
INSTALL_SCRIPT = $(INSTALL) -c

HTML_TEMPLATES = \
	web/templates/admin_index.html \
	web/templates/admin_product.html \
	web/templates/admin_products.html \
	web/templates/admin_symfiles.html \
	web/templates/admin_user.html \
	web/templates/admin_users.html \
	web/templates/admin_version.html \
	web/templates/admin_versions.html \
	web/templates/crashes.html \
	web/templates/crash.html \
	web/templates/foot.html \
	web/templates/head.html \
	web/templates/index.html \
	web/templates/report.html \
	web/templates/reports.html \
	web/templates/symfiles.html

ASSETS_FONTS = \
	web/assets/fonts/bootstrap/glyphicons-halflings-regular.eot \
	web/assets/fonts/bootstrap/glyphicons-halflings-regular.svg \
	web/assets/fonts/bootstrap/glyphicons-halflings-regular.ttf \
	web/assets/fonts/bootstrap/glyphicons-halflings-regular.woff \
	web/assets/fonts/bootstrap/glyphicons-halflings-regular.woff2

ASSETS_JS = \
	web/assets/javascripts/app.js \
	web/assets/javascripts/bootstrap.js \
	web/assets/javascripts/bootstrap.min.js \
	web/assets/javascripts/jquery.min.js \
	web/assets/javascripts/bootstrap-sprockets.js \
	web/assets/javascripts/Chart.bundle.min.js

GO_SRC   = ./cmd/crashdragon/$(wildcard *.go)

SASSCFLAGS ?= -t compressed

all: build/bin/crashdragon build/bin/minidump_stackwalk

build/bin/crashdragon: $(GO_SRC) web/assets/stylesheets/app.css
	$(GO) build -o build/bin/crashdragon $(GO_SRC)

web/assets/stylesheets/app.css:
	$(SASSC) $(SASSCFLAGS) $(@D)/app.scss > $@.tmp && mv $@.tmp $@

build/bin/minidump_stackwalk:
	cd third_party/breakpad && ./autogen.sh 
	cd third_party/breakpad && ./configure --libdir="$(CURDIR)/build/lib" --prefix="$(CURDIR)/build" CXXFLAGS="-Wno-error" CFLAGS="-Wno-error"
	cd third_party/breakpad && $(MAKE) install

clean:
	cd third_party/breakpad && $(MAKE) uninstall
	rm -f build/bin/crashdragon
	rm -f web/assets/stylesheets/app.css.tmp
	rm -f web/assets/stylesheets/app.css
	rm -rf build/lib build/bin build/include build/share
	cd third_party/breakpad && $(MAKE) distclean

install: all
	$(INSTALL) -d $(DESTDIR)$(bindir)
	$(INSTALL_PROGRAM) build/bin/crashdragon $(DESTDIR)$(bindir)
	$(INSTALL_PROGRAM) build/bin/minidump_stackwalk $(DESTDIR)$(bindir)
	$(INSTALL) -d $(DESTDIR)$(datadir)/crashdragon/assets/stylesheets
	$(INSTALL_DATA) web/assets/stylesheets/app.css $(DESTDIR)$(datadir)/crashdragon/assets/stylesheets
	$(INSTALL) -d $(DESTDIR)$(datadir)/crashdragon/assets/javascripts
	$(INSTALL_DATA) $(ASSETS_JS) $(DESTDIR)$(datadir)/crashdragon/assets/javascripts
	$(INSTALL) -d $(DESTDIR)$(datadir)/crashdragon/assets/fonts/bootstrap
	$(INSTALL_DATA) $(ASSETS_FONTS) $(DESTDIR)$(datadir)/crashdragon/assets/fonts/bootstrap
	$(INSTALL) -d $(DESTDIR)$(datadir)/crashdragon/templates
	$(INSTALL_DATA) $(HTML_TEMPLATES) $(DESTDIR)$(datadir)/crashdragon/templates
	$(INSTALL) -d $(DESTDIR)$(sysconfdir)

uninstall:
	rm $(DESTDIR)$(bindir)/crashdragon
	rm $(DESTDIR)$(bindir)/minidump_stackwalk
	rm $(DESTDIR)$(datadir)/crashdragon/assets/stylesheets/app.css
	rm $(addprefix $(DESTDIR)$(datadir)/crashdragon/assets/javascripts/,$(notdir $(ASSETS_JS)))
	rm $(addprefix $(DESTDIR)$(datadir)/crashdragon/assets/fonts/bootstrap/,$(notdir $(ASSETS_FONTS)))
	rm $(addprefix $(DESTDIR)$(datadir)/crashdragon/templates/,$(notdir $(HTML_TEMPLATES)))
	rmdir $(DESTDIR)$(datadir)/crashdragon/assets/fonts/bootstrap/
	rmdir $(DESTDIR)$(datadir)/crashdragon/assets/fonts/
	rmdir $(DESTDIR)$(datadir)/crashdragon/assets/javascripts/
	rmdir $(DESTDIR)$(datadir)/crashdragon/assets/stylesheets/
	rmdir $(DESTDIR)$(datadir)/crashdragon/assets/
	rmdir $(DESTDIR)$(datadir)/crashdragon/templates/
	rmdir $(DESTDIR)$(datadir)/crashdragon/

.PHONY: uninstall install clean all
