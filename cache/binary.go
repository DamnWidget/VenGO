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
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mcuadros/go-version"
)

// Download an specific version of Golang binary files
func CacheDownloadBinary(ver string, f ...bool) error {
	numeric_ver := ver
	ver = GetBinaryVersion(ver)
	expected_sha1, err := Checksum(ver)
	if err != nil {
		return err
	}

	if !Exists(ver) || (len(f) > 0 && f[0]) {
		url := fmt.Sprintf(
			"https://storage.googleapis.com/golang/go%s.tar.gz", ver)
		if version.Compare(version.Normalize(numeric_ver), "1.2.2", "<") {
			url = fmt.Sprintf(
				"https://go.googlecode.com/files/go%s.tar.gz", ver)
		}
		if runtime.GOOS == "windows" {
			url = strings.Replace(url, ".tar.gz", ".zip", -1)
		}
		if err := downloadAndExtract(ver, url, expected_sha1); err != nil {
			return err
		}
		if err := generateManifest(ver); err != nil {
			os.RemoveAll(filepath.Join(CacheDirectory(), ver))
			return err
		}
	}

	return nil
}
