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
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/DamnWidget/VenGO/cache"
	"github.com/DamnWidget/VenGO/utils"
)

var cmdList = &Command{
	Name:  "list",
	Usage: "list list [-a, --all] [-n, --non-installed] [-j, --json]",
	Short: "List installed and available Go versions",
	Long: fmt.Sprintf(`
Shows a list of installed Go versions, available non installed Go versions or
both. If the list command detects that a installed Go version integrity is
compromised, the legend is as shown below:

  Mark  Description
  ----  ------------------------------------
   %s    if the integrity is not compromised
   %s    if the integrity is compromised

If the -n or --non-installed flag is passed to the list command, a complete
list of available sources is shown to the user ordered by binary, mercurial
and source.tar.gz packed versions.

The flag -a or --all is used to show all the available to install and installed
Go versions.

JSON output:
  One can pass the -j or --json option to display the output as a JSON
  structure with the following format:
    {
        "installed": [
             "1.3rc2",
             "go1.2.2",
             "go1.3.2","go1.3.3",
             "go1.4",
             "go1.4rc1",
             "go1.4rc2"
        ]
    }
`, utils.Ok("✔"), utils.Fail("✖")),
	Execute: runList,
}

var (
	all          bool
	nonInstalled bool
	asJson       bool
)

// initialize the command
func init() {
	cmdList.Flag.BoolVarP(&all, "all", "a", false, "")
	cmdList.Flag.BoolVarP(&nonInstalled, "non-installed", "n", false, "")
	cmdList.Flag.BoolVarP(&asJson, "json", "j", false, "")
}

// run the list command
func runList(cmd *Command, args ...string) {
	options := func(l *List) {
		l.DisplayAs = Text
		if asJson {
			l.DisplayAs = Json
		}
		l.ShowBoth = all
		l.ShowInstalled = true
		if nonInstalled {
			l.ShowNotInstalled = true
			if !l.ShowBoth {
				l.ShowInstalled = false
			}
		}
	}
	nl := NewList(options)
	data, err := nl.Run()
	if err == nil {
		fmt.Println(data)
		return
	}
	log.Println(err)
	os.Exit(2)
}

// json brief output structure
type BriefJSON struct {
	Installed []string `json:"installed,omitempty"`
	Available []string `json:"available,omitempty"`
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
