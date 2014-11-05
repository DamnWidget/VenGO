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
	"runtime"
	"strings"

	"github.com/DamnWidget/VenGO/logger"
	"github.com/mcuadros/go-version"
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

// checks for the existance of the given version in the cache
func Exists(ver string) bool {
	_, err := os.Stat(filepath.Join(CacheDirectory(), ver))
	return err == nil
}

// Download an specific version of Golang binary files
func CacheDownloadBinary(ver string) error {
	numeric_ver := ver
	ver = getBinaryVersion(ver)
	expected_sha1, err := Checksum(ver)
	if err != nil {
		return err
	}

	if !Exists(ver) {
		url := fmt.Sprintf(
			"https://storage.googleapis.com/golang/go%s.tar.gz", ver)
		if version.Compare(version.Normalize(numeric_ver), "1.2.2", "<") {
			url = fmt.Sprintf(
				"https://go.googlecode.com/files/go%s.tar.gz", ver)
		}
		if runtime.GOOS == "windows" {
			url = strings.Replace(url, ".tar.gz", ".zip", -1)
		}
                logger.Println(url)
		if err := downloadAndExtract(ver, url, expected_sha1); err != nil {
			return err
		}
	}

	return nil
}

// Download an specific version of Golang source code
func CacheDownload(ver string) error {
	expected_sha1, err := Checksum(ver)
	if err != nil {
		return err
	}

	if !Exists(ver) {
		url := fmt.Sprintf(
			"https://storage.googleapis.com/golang/go%s.src.tar.gz", ver)
		if version.Compare(version.Normalize(ver), "1.2.2", "<") {
			url = fmt.Sprintf(
				"https://go.googlecode.com/files/go%s.src.tar.gz", ver)
		}
		if err := downloadAndExtract(ver, url, expected_sha1); err != nil {
			return err
		}
	}

	return nil
}

// download and extract the given file checking the given sha1 signature
func downloadAndExtract(ver, url, expected_sha1 string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 400 {
			log.Fatal("Version %s can't be found!\n", ver)
		}
		return fmt.Errorf("%s", resp.Status)
	}
	defer resp.Body.Close()

	logger.Printf("downloading Go%s from %s\n", ver, url)
	buf := new(bytes.Buffer)
	size, err := io.Copy(buf, resp.Body)
	if err != nil {
		return err
	}

	pkg_sha1 := fmt.Sprintf("%x", sha1.Sum(buf.Bytes()))
	if pkg_sha1 != expected_sha1 {
		return fmt.Errorf(
			"Error: SHA1 is different! expected %s got %s",
			expected_sha1, pkg_sha1,
		)
	}
	logger.Printf("%d bytes donwloaded... decompresssing...\n", size)
	prefix := filepath.Join(CacheDirectory(), ver)
	extractTar(prefix, readGzipFile(buf))

	return nil
}

// read the contents of a compressed gzip file
func readGzipFile(data *bytes.Buffer) *bytes.Buffer {
	reader, err := gzip.NewReader(data)
	if err != nil {
		logger.Println("Fatal error reading gzip file contents...")
		logger.Fatal(err)
	}
	defer reader.Close()
	gzipBuf := new(bytes.Buffer)
	if _, err := io.Copy(gzipBuf, reader); err != nil {
		logger.Println(
			"Fatal error while reading gzip file contents into the buffer")
		logger.Fatal(err)
	}

	return gzipBuf
}

// extract the contents of the tar data into the given prefix
func extractTar(prefix string, data *bytes.Buffer) {
	tr := tar.NewReader(data)
	if err := os.MkdirAll(filepath.Join(prefix, "go"), 0766); err != nil {
		logger.Fatal(err)
	}
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err != io.EOF {
				logger.Fatal(err)
			}
			break
		}
		fi := hdr.FileInfo()
		if fi.IsDir() {
			err := os.MkdirAll(filepath.Join(prefix, hdr.Name), 0766)
			if err != nil && os.IsNotExist(err) {
				logger.Fatal(err)
			}
		} else {
			tw, err := os.OpenFile(
				filepath.Join(prefix, hdr.Name), os.O_RDWR|os.O_CREATE|os.O_TRUNC, fi.Mode())
			if err != nil && !os.IsExist(err) {
				logger.Fatal(err)
			}
			if _, err := io.Copy(tw, tr); err != nil {
				logger.Fatal(err)
			}
			tw.Close()
		}
	}
}
