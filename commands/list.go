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
	"path/filepath"
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
	installed, err := cache.GetInstalled(tags, sources, binaries)
	if err != nil {
		fmt.Println("while running List command:", err)
		return "error while running the command", err
	}
	versions["installed"] = append(versions["installed"], installed...)
	versions["available"] = append(
		versions["available"],
		cache.GetNonInstalled(installed, tags, sources, binaries)...,
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
