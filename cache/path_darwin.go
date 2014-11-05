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

package cache

import (
	"fmt"
	"log"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
)

// Return the CacheDirctory for OS X, ~/Library/Caches/VenGO
func CacheDirectory() string {
	return path.Join(ExpandUser("~"), "Library", "Caches", "VenGO")
}

// return back the binary string version for downloads in OS X
func GetBinaryVersion(version string) string {
	cmd := exec.Command("sw_vers", "-productVersion")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	major_ver := "10.6"
	ver := strings.TrimRight(string(out), "\n")
	numeric_ver, _ := strconv.ParseInt(ver[3:], 10, 64)
	if numeric_ver >= int64(8) {
		major_ver = "10.8"
	}
	return fmt.Sprintf("%s.darwin-%s-osx%s", version, runtime.GOARCH, major_ver)
}
