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
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/user"
	"strings"

	"github.com/DamnWidget/VenGO/logger"
)

// Expand the user home tilde to the right user home path
func ExpandUser(path string) string {
	u, err := user.Current()
	if err != nil {
		log.Println("Can't get current user:", err)
		return path
	}
	return strings.Replace(path, "~", u.HomeDir, -1)
}

// Download an specific version of Golang
func CacheDownload(version string) error {
	url := fmt.Sprintf(
		"https://storage.googleapis.com/golang/go%s.src.tar.gz", version)
	resp, err := http.Get(url)
	if err != nil {
		logger.Fatal(err)
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 400 {
			log.Fatal("Version %s can't be found!\n", version)
		}
		logger.Fatal(resp.Status)
	}
	defer resp.Body.Close()
	out, err := ioutil.TempFile("", "vengo-")
	if err != nil {
		logger.Fatal(err)
	}

	logger.Printf("downloading go%s.src.tar.gz...\n", version)
	buf := new(bytes.Buffer)
	size, err := io.Copy(buf, resp.Body)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Printf("%s bytes donwloaded... decompresssing...")

}
