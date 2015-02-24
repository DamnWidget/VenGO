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
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/DamnWidget/VenGO/utils"
	"github.com/mcuadros/go-version"
)

// determine if a Go version has been already compiled in the cache
func AlreadyCompiled(ver string) bool {
	manifest := filepath.Join(CacheDirectory(), ver, ".vengo-manifest")
	if _, err := os.Stat(manifest); err != nil {
		return false
	}
	fmt.Fprint(Output, "Checking manifest integrity... ")
	if err := CheckManifestIntegrity(manifest); err != nil {
		fmt.Println(utils.Fail("✖"))
		log.Println(err)
		return false
	}
	fmt.Fprintln(Output, utils.Ok("✔"))
	return true
}

// compile a given version of go in the cache
func Compile(ver string, verbose, nocgo bool, boostrap ...string) error {
	fmt.Fprint(Output, "Compiling... ")
	if verbose {
		fmt.Fprint(Output, "\n")
	}

	bs := ""
	if len(boostrap) > 0 {
		bs = boostrap[0]
	}

	currdir, _ := os.Getwd()
	prefixed := false
	err := os.Chdir(filepath.Join(CacheDirectory(), ver, "go", "src"))
	if err != nil {
		if !strings.HasPrefix(ver, "go") && ver != "tip" {
			ver = fmt.Sprintf("go%s", ver)
		}
		prefixed = true
		if err := os.Chdir(
			filepath.Join(CacheDirectory(), ver, "src")); err != nil {
			if !verbose {
				fmt.Fprintln(Output, utils.Fail("✖"))
			}
			return err
		}
	}
	defer func() { os.Chdir(currdir) }()

	cmd := "./make.bash"
	if runtime.GOOS == "windows" {
		cmd = "./make.bat"
	}
	if nocgo {
		os.Setenv("CGO_ENABLED", "0")
	}
	if bs != "" {
		os.Setenv("GOROOT_BOOTSTRAP", bs)
	}
	if err := utils.Exec(verbose, cmd); err != nil {
		return err
	}
	goBin := filepath.Join(CacheDirectory(), ver, "go", "bin", "go")
	if prefixed {
		goBin = filepath.Join(CacheDirectory(), ver, "bin", "go")
	}
	if _, err := os.Stat(goBin); err != nil {
		if !verbose {
			fmt.Fprintln(Output, utils.Fail("✖"))
		}
		fmt.Fprintln(Output, err)
		return fmt.Errorf("Go %s wasn't compiled properly! %v", ver, err)
	}
	if !verbose {
		fmt.Fprintln(Output, utils.Ok("✔"))
	}
	fmt.Fprintf(Output, "Generating manifest... ")
	if err := generateManifest(ver); err != nil {
		os.RemoveAll(filepath.Join(CacheDirectory(), ver))
		fmt.Fprintln(Output, utils.Fail("✖"))
		return err
	}
	fmt.Fprintln(Output, utils.Ok("✔"))

	return nil
}

// log compilation process
func logCompilation(rd, erd *bufio.Reader) {
	go func() {
		for {
			str, err := rd.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					log.Fatal(err)
				}
				break
			}
			fmt.Printf("%s", str)
		}
	}()

	// read the command error output and update the terminal
	go func() {
		for {
			str, err := erd.ReadString('\n')
			if err != nil {
				break
			}
			fmt.Printf("%s", str)
		}
	}()
}

// Download an specific version of Golang source code
func CacheDownload(ver string, f ...bool) error {
	expected_sha1, err := Checksum(ver)
	if err != nil {
		return err
	}
	force := (len(f) != 0 && f[0] == true)

	if !Exists(ver) || force {
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
