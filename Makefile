
# Copyright (C) 2014  Oscar Campos <oscar.campos@member.fsf.org>

# This program is free software; you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation; either version 2 of the License, or
# (at your option) any later version.

# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.

# You should have received a copy of the GNU General Public License along
# with this program; if not, write to the Free Software Foundation, Inc.,
# 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

# See LICENSE file for more details.

RM = rm -Rf
CREATE_DIR = -D
INSTALL = install
ACTIVATE_TPL = env/tpl/activate
INSTALL_DATA = $(INSTALL) -m 644 -p
TARGET = $(HOME)/.VenGO
BINDIR = bin
VERSION = VERSION
BUILD = go build -v -x -o

default: build

.PHONY: default

clean:
	go clean
	$(RM) $(BINDIR)/list
	$(RM) $(BINDIR)/lsenvs
	$(RM) $(BINDIR)/install
	$(RM) $(BINDIR)/uninstall
	$(RM) $(BINDIR)/mkenv
	$(RM) $(BINDIR)/rmenv

test: cache_test env_test

cache_test:
	cd cache && ginkgo -r --randomizeAllSpecs --failOnPending --randomizeSuites --trace --race
.PHONY: cache_test

env_test:
	cd env && ginkgo -r --randomizeAllSpecs --failOnPending --randomizeSuites --trace --race
.PHONY: env_test

commands_test:
	cd commands && ginkgo -r --randomizeAllSpecs --failOnPending --randomizeSuites --trace --race
.PHONY: commands_test

build: clean test
	$(BUILD) bin/list ./applications/list
	$(BUILD) bin/lsenvs ./applications/lsenvs
	$(BUILD) bin/install ./applications/install
	$(BUILD) bin/uninstall ./applications/uninstall
	$(BUILD) bin/mkenv ./applications/mkenv
	$(BUILD) bin/rmenv ./applications/rmenv

fast_build:
	$(BUILD) bin/list ./applications/list
	$(BUILD) bin/lsenvs ./applications/lsenvs
	$(BUILD) bin/install ./applications/install
	$(BUILD) bin/uninstall ./applications/uninstall
	$(BUILD) bin/mkenv ./applications/mkenv
	$(BUILD) bin/rmenv ./applications/rmenv

install: fast_build installdirs
	$(INSTALL_DATA) $(CREATE_DIR) $(ACTIVATE_TPL) $(TARGET)/scripts/tpl/activate
	$(INSTALL_DATA) $(VERSION) $(TARGET)/version
	echo ""
	echo "\033[32mVenGO is now installed in your system\033[0m"
	echo "add 'source $(HOME)/.VenGO/bin/vengo' to your .bashrc or .profile to activate it"
	echo "you can also do '. $(HOME)/.VenGO/bin/vengo' to start using it right now"
.PHONY: install

installdirs:
	install -d $(TARGET)
	(tar -cf - bin) | (cd $(TARGET) && tar -xf -)
.PHONY: installdirs

uninstall:
	$(RM) $(TARGET)
	go build -tags clean -o cleaner ./application/cleaner
	./cleaner
	$(RM) ./cleaner
.PHONY: uninstall

.SILENT: clean build fast_build test cache_test env_test commands_test install installdirs
