
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
INSTALL_DATA = $(INSTALL) -m 644 -p
TARGET = $(HOME)/.VenGO
TARGET_BINDIR = $(TARGET)/bin/vengo
SCRIPTS = commands
VERSION = VERSION
VENGO = bin/vengo
BUILD = go build -v -x -o

default: build

.PHONY: default

clean:
	go clean
	$(RM) $(SCRIPTS)/scripts/list
	$(RM) $(SCRIPTS)/scripts/lsenvs
	$(RM) $(SCRIPTS)/scripts/install
	$(RM) $(SCRIPTS)/scripts/uninstall
	$(RM) $(SCRIPTS)/scripts/mkenv
	$(RM) $(SCRIPTS)/scripts/rmenv

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
	$(BUILD) commands/scripts/list ./applications/list
	$(BUILD) commands/scripts/lsenvs ./applications/lsenvs
	$(BUILD) commands/scripts/install ./applications/install
	$(BUILD) commands/scripts/uninstall ./applications/uninstall
	$(BUILD) commands/scripts/mkenv ./applications/mkenv
	$(BUILD) commands/scripts/rmenv ./applications/rmenv

fast_build:
	$(BUILD) commands/scripts/list ./applications/list
	$(BUILD) commands/scripts/lsenvs ./applications/lsenvs
	$(BUILD) commands/scripts/install ./applications/install
	$(BUILD) commands/scripts/uninstall ./applications/uninstall
	$(BUILD) commands/scripts/mkenv ./applications/mkenv
	$(BUILD) commands/scripts/rmenv ./applications/rmenv

install: build installdirs
	$(INSTALL) $(CREATE_DIR) $(VENGO) $(TARGET_BINDIR)
	$(INSTALL_DATA) $(VERSION) $(TARGET)/version
	$(INSTALL_DATA) -D env/tpl/activate $(TARGET)/scripts/tpl/activate
.PHONY: install

installdirs:
	(cd $(SCRIPTS) && tar -cf - scripts) | (cd $(TARGET) && tar -xf -)
.PHONY: installdirs

uninstall:
	$(RM) $(TARGET)
	go build -o cleaner ./application/cleaner
	./cleaner
	$(RM) ./cleaner
.PHONY: uninstall

.SILENT: clean build test cache_test env_test commands_test
