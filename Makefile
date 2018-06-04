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
	templates/admin_index.html
	templates/admin_product.html
	templates/admin_products.html
	templates/admin_symfiles.html
	templates/admin_user.html
	templates/admin_users.html
	templates/admin_version.html
	templates/admin_versions.html
	templates/crashes.html
	templates/crash.html
	templates/foot.html
	templates/head.html
	templates/index.html
	templates/report.html
	templates/reports.html
	templates/symfiles.html

ASSETS_FONTS = \
	assets/fonts/bootstrap/glyphicons-halflings-regular.eot \
	assets/fonts/bootstrap/glyphicons-halflings-regular.svg \
	assets/fonts/bootstrap/glyphicons-halflings-regular.ttf \
	assets/fonts/bootstrap/glyphicons-halflings-regular.woff \
	assets/fonts/bootstrap/glyphicons-halflings-regular.woff2

ASSETS_JS = \
	assets/javascripts/app.js \
	assets/javascripts/bootstrap.js \
	assets/javascripts/bootstrap.min.js \
	assets/javascripts/jquery.min.js \
	assets/javascripts/bootstrap-sprockets.js \
	assets/javascripts/Chart.bundle.min.js

GO_SRC   = ./server/$(wildcard *.go)

SASSCFLAGS ?= -t compressed

all: crashdragon build/bin/minidump_stackwalk

crashdragon: $(GO_SRC) assets/stylesheets/app.css
	$(GO) build -o bin/crashdragon $(GO_SRC)

assets/stylesheets/app.css:
	$(SASSC) $(SASSCFLAGS) $(@D)/app.scss > $@.tmp && mv $@.tmp $@

build/bin/minidump_stackwalk:
	cd breakpad && ./autogen.sh && ./configure --prefix="$(CURDIR)/build" && $(MAKE) install

clean:
	rm -f bin/crashdragon
	rm -f assets/stylesheets/app.css.tmp
	rm -f assets/stylesheets/app.css
	rm -rf build/
	cd breakpad && $(MAKE) distclean

install:
	$(INSTALL_PROGRAM) bin/crashdragon $(DESTDIR)$(bindir)
	$(INSTALL_PROGRAM) build/bin/minidump_stackwalk $(DESTDIR)$(bindir)
	$(INSTALL_DATA) -d assets/stylesheets/app.css $(DESTDIR)$(datadir)/crashdragon/assets
	$(INSTALL_DATA) -d $(ASSETS_JS) $(DESTDIR)$(datadir)/crashdragon/assets/javascripts
	$(INSTALL_DATA) -d $(ASSETS_FONTS) $(DESTDIR)$(datadir)/crashdragon/assets/fonts/bootstrap
	$(INSTALL_DATA) -d $(HTML_TEMPLATES) $(DESTDIR)$(datadir)/crashdragon/templates


.PHONY: install clean all
