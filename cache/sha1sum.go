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
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// generate installation manifest
func generateManifest(ver string) error {
	manifest := []string{}
	versionPath := filepath.Join(CacheDirectory(), ver)
	if err := filepath.Walk(
		versionPath,
		func(path string, info os.FileInfo, err error) error {
			if info.Name() == ".vengo-manifest" {
				return nil
			}
			data := []byte(path)
			if !info.IsDir() {
				var e error
				data, e = ioutil.ReadFile(path)
				if e != nil {
					return fmt.Errorf("while generating manifest: %s", e)
				}
			}

			fileSha1 := fmt.Sprintf("%x %s", sha1.Sum(data), path)
			manifest = append(manifest, fileSha1)
			return nil
		},
	); err != nil {
		return err
	}

	fileName := filepath.Join(versionPath, ".vengo-manifest")
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("while generating manifest:", err)
	}
	file.WriteString(fmt.Sprintf("%s\n", strings.Join(manifest, "\n")))
	file.Close()

	return nil
}

// checks a manifest integrity
func CheckManifestIntegrity(manifestName string) error {
	file, err := os.Open(manifestName)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if line != "" {
			splitData := strings.Split(line, " ")
			f := strings.TrimRight(splitData[1], "\n\r")
			fi, statErr := os.Stat(f)
			if statErr != nil {
				return fmt.Errorf("Integirty check failed! %s", statErr)
			}
			data := []byte(f)
			if !fi.IsDir() {
				data, _ = ioutil.ReadFile(f)
			}
			fileSha1 := fmt.Sprintf("%x", sha1.Sum(data))
			if splitData[0] != fileSha1 {
				fmt.Printf("%s -> %s\n", f, splitData[0])
				fmt.Printf("%s -> %s\n", f, fileSha1)
				return fmt.Errorf("Integrity check failed!")
			}
		}

		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
	}
	return nil
}
