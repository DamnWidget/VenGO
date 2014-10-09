/*
	Copyright (C) 2014  Oscar Campos <oscar.campos@member.fsf.org>

	This program is free software; you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation; either version 2 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License along
    with this program; if not, write to the Free Software Foundation, Inc.,
    51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

	See LICENSE file for more details.
*/

// +build linux dragonfly freebsd netbsd openbsd solaris

package cache

import (
	"os"
	"os/path"
)

// Return the CacheDirctory for not darwin Unix. By default it is
// ~/.cache/VenGO. On Linux, if the environemnt variable XDG_CACHE_HOME
// exists it will be XDG_CACHE_HOME/VenGO
func CacheDirectory() string {
	XDG_CACHE_HOME := os.Getenv("XDG_CACHE_HOME")
	if XDG_CACHE_HOME == "" {
		XDG_CACHE_HOME = ExpandUser("~/.cache")
	}
	return path.Join(XDG_CACHE_HOME, "VenGO")
}
