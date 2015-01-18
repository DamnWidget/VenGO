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
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// return a list of installed go versions
func GetInstalled(tags, sources, binaries []string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(CacheDirectory(), "*"))
	if err != nil {
		fmt.Println("while getting installed versions:", err)
		return nil, err
	}
	versions := []string{}
	for _, file := range files {
		filename := path.Base(file)
		if filename != "mercurial" && filename != "logs" && filename != "git" {
			stat, err := os.Stat(file)
			if err != nil {
				fmt.Println("while getting installed versions:", err)
				return nil, err
			}
			if stat.IsDir() {
				if isValidVersion(filename, tags, sources, binaries) {
					versions = append(versions, filename)
				}
			}
		}
	}

	return versions, nil
}

// return a list of non installed go versions
func GetNonInstalled(v, tags, sources, binaries []string) []string {
	versions := []string{}
	installed_versions := make([]string, len(v))
	copy(installed_versions, v)
	for _, ver := range append(binaries, append(tags, sources...)...) {
		found := false
		for i, installed := range installed_versions {
			if strings.TrimSpace(installed) == strings.TrimSpace(ver) {
				// skip this element and reduce v
				installed_versions = append(
					installed_versions[:i], installed_versions[i+1:]...)
				found = true
				continue
			}
		}
		if found {
			continue
		}
		versions = append(versions, fmt.Sprintf("    %s", ver))
	}

	return versions
}

// check if a given version is valid in all the possible containers
func isValidVersion(file string, tags, sources, binaries []string) bool {
	// tip is always a valid version
	if file == "tip" || file == "go" {
		return true
	}
	// look on the sources first that is the smaller collection
	for _, ver := range sources {
		if file == ver {
			return true
		}
	}
	// now look on the binaries collection using binary search
	index := sort.SearchStrings(binaries, file)
	if len(binaries) > index && binaries[index] == file {
		return true
	}
	// now look in the git tags using binary search
	// now look in the mercurial tags using binary search
	index = sort.SearchStrings(tags, file)
	if len(tags) > index && tags[index] == file {
		return true
	}

	return false
}
