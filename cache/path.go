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
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/DamnWidget/VenGO/utils"
)

var Output io.Writer = os.Stdout

var VenGO_PATH = filepath.Join(ExpandUser("~"), ".VenGO")

// Expand the user home tilde to the right user home path
func ExpandUser(path string) string {
	u, err := user.Current()
	if err != nil {
		fmt.Fprintln(Output, "Can't get current user:", err)
		return path
	}
	return strings.Replace(path, "~", u.HomeDir, -1)
}

// checks for the existance of the given version in the cache
func Exists(ver string) bool {
	_, err := os.Stat(filepath.Join(CacheDirectory(), ver))
	return err == nil
}

// download and extract the given file checking the given sha1 signature
func downloadAndExtract(ver, url, expected_sha1 string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 400 {
			log.Fatalf("Version %s can't be found!\n", ver)
		}
		return fmt.Errorf("%s", resp.Status)
	}
	defer resp.Body.Close()

	fmt.Fprintf(Output, "downloading Go%s from %s ", ver, url)
	buf := new(bytes.Buffer)
	size, err := io.Copy(buf, resp.Body)
	if err != nil {
		fmt.Fprintln(Output, utils.Fail("✖"))
		return err
	}
	fmt.Fprintln(Output, utils.Ok("✔"))

	pkg_sha1 := fmt.Sprintf("%x", sha1.Sum(buf.Bytes()))
	if pkg_sha1 != expected_sha1 {
		return fmt.Errorf(
			"Error: SHA1 is different! expected %s got %s",
			expected_sha1, pkg_sha1,
		)
	}
	fmt.Fprintf(Output, "%d bytes donwloaded... decompresssing... ", size)
	prefix := filepath.Join(CacheDirectory(), ver)
	extractTar(prefix, readGzipFile(buf))
	buf.Reset()
	buf = nil
	fmt.Fprintln(Output, utils.Ok("✔"))

	return nil
}

// read the contents of a compressed gzip file
func readGzipFile(data *bytes.Buffer) *bytes.Buffer {
	reader, err := gzip.NewReader(data)
	if err != nil {
		fmt.Fprintln(Output, "Fatal error reading gzip file contents...")
		log.Fatal(err)
	}
	defer reader.Close()
	gzipBuf := new(bytes.Buffer)
	if _, err := io.Copy(gzipBuf, reader); err != nil {
		fmt.Fprintln(Output,
			"Fatal error while reading gzip file contents into the buffer")
		log.Fatal(err)
	}

	return gzipBuf
}

// extract the contents of the tar data into the given prefix
func extractTar(prefix string, data *bytes.Buffer) {
	tr := tar.NewReader(data)
	if err := os.MkdirAll(filepath.Join(prefix, "go"), 0766); err != nil {
		log.Fatal(err)
	}
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}
		fi := hdr.FileInfo()
		if fi.IsDir() {
			err := os.MkdirAll(filepath.Join(prefix, hdr.Name), 0766)
			if err != nil && os.IsNotExist(err) {
				log.Fatal(err)
			}
		} else {
			tw, err := os.OpenFile(
				filepath.Join(prefix, hdr.Name), os.O_RDWR|os.O_CREATE|os.O_TRUNC, fi.Mode())
			if err != nil && !os.IsExist(err) {
				log.Fatal(err)
			}
			if _, err := io.Copy(tw, tr); err != nil {
				log.Fatal(err)
			}
			tw.Close()
		}
	}
	data.Reset()
	data = nil
}
