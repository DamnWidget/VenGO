
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

default: build

.PHONY: default

clean:
	go clean

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

build:
	go build -v -x -o vengo ./application

.SILENT: clean build test cache_test env_test commands_test
