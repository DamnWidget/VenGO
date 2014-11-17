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
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/DamnWidget/VenGO/logger"
	"github.com/mcuadros/go-version"
)

// determine if a Go version has been already compiled in the cache
func alreadyCompiled(ver string) bool {
	_, err := os.Stat(filepath.Join(CacheDirectory(), ver, "go", "bin", "go"))
	return err == nil
}

// compile a given version of go in the cache
func Compile(ver string) error {
	currdir, _ := os.Getwd()
	err := os.Chdir(filepath.Join(CacheDirectory(), ver, "go", "src"))
	if err != nil {
		return err
	}
	defer func() { os.Chdir(currdir) }()

	cmd := "./make.bash"
	if runtime.GOOS == "windows" {
		cmd = "./make.bat"
	}
	p := exec.Command(cmd)
	out, err := p.StdoutPipe()
	if err != nil {
		return err
	}
	rd := bufio.NewReader(out)
	if err := p.Start(); err != nil {
		return err
	}

	// read the command output and update the terminal
	for {
		str, err := rd.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				logger.Fatal(err)
			}
			break
		}
		logger.Printf("%s", str)
	}

	if _, err := os.Stat(
		filepath.Join(CacheDirectory(), ver, "go", "bin", "go")); err != nil {
		return fmt.Errorf("Go %s wasn't compiled properly!", ver)
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
