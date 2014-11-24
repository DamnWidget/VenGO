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
	"github.com/DamnWidget/VenGO/utils"
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

	versions := map[string][]string{
		"installed": []string{},
		"available": []string{},
	}
	installed, err := l.getInstalled(tags, sources, binaries)
	if err != nil {
		fmt.Println("while running List command:", err)
		return "error while running the command", err
	}
	versions["installed"] = append(versions["installed"], installed...)
	versions["available"] = append(
		versions["available"],
		l.getNonInstalled(installed, tags, sources, binaries)...,
	)

	return l.display(versions)
}

// generates the output for the list command
func (l *List) display(versions map[string][]string) (string, error) {
	output := []string{}
	if l.DisplayAs == Text {
		if l.ShowBoth || l.ShowInstalled {
			output = append(output, utils.Ok("Installed"))
			for _, v := range versions["installed"] {
				_, err := os.Stat(
					filepath.Join(cache.CacheDirectory(), v, ".vengo-manifest"))
				check := utils.Ok("✔")
				if err != nil {
					check = utils.Fail("✖")
				}
				output = append(output, fmt.Sprintf("    %s %s", v, check))
			}
		}
		if l.ShowBoth || l.ShowNotInstalled {
			output = append(output, utils.Ok("Available for Installation"))
			output = append(output, versions["available"]...)
		}
		return strings.Join(output, "\n"), nil
	}

	if l.DisplayAs == Json {
		jsonData := &BriefJSON{[]string{}, []string{}}
		if l.ShowBoth || l.ShowInstalled {
			for _, v := range versions["installed"] {
				v := strings.TrimLeft(v, "    ")
				jsonData.Installed = append(jsonData.Installed, v)
			}
		}
		if l.ShowBoth || l.ShowNotInstalled {
			for _, v := range versions["available"] {
				v := strings.TrimLeft(v, "    ")
				jsonData.Available = append(jsonData.Available, v)
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
		fmt.Println("while getting installed versions:", err)
		return nil, err
	}
	versions := []string{}
	for _, file := range files {
		filename := path.Base(file)
		if filename != "mercurial" && filename != "logs" {
			stat, err := os.Stat(file)
			if err != nil {
				fmt.Println("while getting installed versions:", err)
				return nil, err
			}
			if stat.IsDir() {
				if l.isValidVersion(filename, tags, sources, binaries) {
					versions = append(versions, filename)
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
		versions[c] = fmt.Sprintf("    %s", ver)
		c++
	}

	return versions
}

// check if a given version is valid in all the possible containers
func (l *List) isValidVersion(file string, tags, sources, binaries []string) bool {
	// tip is always a valid version
	if file == "tip" {
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
	// now look in the mercurial tags using binary search
	index = sort.SearchStrings(tags, file)
	if len(tags) > index && tags[index] == file {
		return true
	}

	return false
}
