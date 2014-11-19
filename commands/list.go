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

package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/DamnWidget/VenGO/cache"
	"github.com/DamnWidget/VenGO/logger"
)

const (
	Text = iota
	Json
)

// json brief output structure
type BriefJSON struct {
	Installed []string `json:"installed"`
	Available []string `json:"available"`
}

// list command
type List struct {
	ShowInstalled    bool
	ShowNotInstalled bool
	ShowBoth         bool
	DisplayAs        int
}

// Create a new list and return back it's address
func NewList(options ...func(*List)) *List {
	list := new(List)
	list.ShowInstalled = true

	for _, option := range options {
		option(list)
	}

	return list
}

// implements the Runner interface returning back a list of installed
// go versions, a list of not installed versions or all versions depending
// on the list options
func (l *List) Run() (string, error) {
	tags := cache.Tags()
	sources := cache.AvailableSources()
	binaries := cache.AvailableBinaries()

	versions := []string{}
	installed, err := l.getInstalled(tags, sources, binaries)
	if err != nil {
		logger.Println("while running List command:", err)
		return "error while running the command", err
	}
	if l.ShowBoth {
		versions = []string{"Installed"}
		versions = append(versions, installed...)
		versions = append(versions, "\nAvailable for Installation")
		versions = append(
			versions, l.getNonInstalled(installed, tags, sources, binaries)...)
	} else {
		if l.ShowInstalled {
			versions = []string{"Installed"}
			versions = append(versions, installed...)
		}
		if l.ShowNotInstalled {
			versions = append(versions, "\nAvailable for Installation")
			versions = append(
				versions, l.getNonInstalled(installed, tags, sources, binaries)...)
		}
	}

	return l.display(versions)
}

// generates the output for the list command
func (l *List) display(versions []string) (string, error) {
	if l.DisplayAs == Text {
		return strings.Join(versions, "\n"), nil
	}

	if l.DisplayAs == Json {
		jsonData := &BriefJSON{[]string{}, []string{}}
		doneInstalled := false
		for _, v := range versions {
			if v != "Installed" {
				if !doneInstalled {
					if v == "\nAvailable for Installation" {
						doneInstalled = true
						continue
					}
					v := strings.TrimLeft(v, "    ")
					jsonData.Installed = append(jsonData.Installed, v)
				} else {
					v := strings.TrimLeft(v, "    ")
					jsonData.Available = append(jsonData.Available, v)
				}
			}

		}
		data, err := json.Marshal(jsonData)
		return string(data), err
	}

	return "", fmt.Errorf("List.DisplayAs is not set to a valid value!")
}

// return a list of installed go versions
func (l *List) getInstalled(tags, sources, binaries []string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(cache.CacheDirectory(), "*"))
	if err != nil {
		logger.Println("while getting installed versions:", err)
		return nil, err
	}
	versions := []string{}
	for _, file := range files {
		filename := path.Base(file)
		if filename != "mercurial" && filename != "logs" {
			stat, err := os.Stat(file)
			if err != nil {
				logger.Println("while getting installed versions:", err)
				return nil, err
			}
			if stat.IsDir() {
				if l.isValidVersion(filename, tags, sources, binaries) {
					versions = append(versions, fmt.Sprintf("    %s", filename))
				}
			}
		}
	}

	return versions, nil
}

// return a list of non installed go versions
func (l *List) getNonInstalled(v, tags, sources, binaries []string) []string {
	var versions = make([]string, len(tags)+len(sources)+len(binaries))
	installed_versions := make([]string, len(v))
	copy(installed_versions, v)
	c := 0
	for _, ver := range append(binaries, append(tags, sources...)...) {
		for i, installed := range installed_versions {
			if installed == ver {
				// skip this element and reduce v
				installed_versions = append(
					installed_versions[:i], installed_versions[i+1:]...)
				continue
			}
		}
		versions[c] = fmt.Sprintf("    %s", ver)
		c++
	}

	return versions
}

// check if a given version is valid in all the possible containers
func (l *List) isValidVersion(file string, tags, sources, binaries []string) bool {
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
	// now look in the mercurial tags using binary search
	index = sort.SearchStrings(tags, file)
	if len(tags) > index && tags[index] == file {
		return true
	}

	return false
}
