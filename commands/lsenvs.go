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
	"strings"

	"github.com/DamnWidget/VenGO/cache"
	"github.com/DamnWidget/VenGO/utils"
)

var cmdLsenvs = &Command{
	Name:  "lsenvs",
	Usage: "lsenvs [-j]",
	Short: "Lists available virtual Go environments",
	Long: fmt.Sprintf(`Lists isolated virtual Go environments in your system. Integrity compromised
environments are shown as the legend shown below:

	Mark  Description
  ----  ------------------------------------
   %s    if the integrity is not compromised
   %s    if the integrity is compromised

If the -j or --json option is passed, the command resturn a JSON string instead.
`, utils.Ok("✔"), utils.Fail("✖")),
	Execute: runLsenvs,
}

// initialize the command
func init() {
	cmdLsenvs.Flag.BoolVarP(&asJsonList, "json", "j", false, "display JSON")
	cmdLsenvs.register()
}

// run the lsenvs command
func runLsenvs(cmd *Command, args ...string) {
	options := []func(el *EnvironmentsList){}

	if asJsonList {
		options = append(options, func(el *EnvironmentsList) {
			el.DisplayAs = Json
		})
	}
	nel := NewEnvironmentsList(options...)
	data, err := nel.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Println(data)
	os.Exit(0)
}

type EnvironmentsJSON struct {
	Available []string `json:"available"`
	Invalid   []string `json:"invalid,omitempty"`
}

// EnvironmentsList command
type EnvironmentsList struct {
	DisplayAs int
}

// Create a new lsenv adn returns back it's address
func NewEnvironmentsList(options ...func(*EnvironmentsList)) *EnvironmentsList {
	lsenv := new(EnvironmentsList)

	for _, option := range options {
		option(lsenv)
	}

	return lsenv
}

// implements the Runner interface returning back a list of available
// environments in text or json format depending on the given options
func (e *EnvironmentsList) Run() (string, error) {
	available, invalid, err := e.getEnvironments()
	if err != nil {
		fmt.Println("while running EnvironmentsList command:", err)
		return "error while running the command", err
	}

	return e.display(available, invalid)
}

// generates the output for the lsenvs command
func (e *EnvironmentsList) display(
	available, invalid []string) (string, error) {

	output := []string{}
	if e.DisplayAs == Text {
		output = append(output, utils.Ok("Virtual Go Environments"))
		for _, v := range available {
			output = append(output, fmt.Sprintf("    %s %s", v, utils.Ok("✔")))
		}
		for _, v := range invalid {
			output = append(
				output, fmt.Sprintf("    %s %s", v, utils.Fail("✖")))
		}

		return strings.Join(output, "\n"), nil
	}

	if e.DisplayAs == Json {
		data, err := json.Marshal(&EnvironmentsJSON{available, invalid})
		return string(data), err
	}

	return "", fmt.Errorf("EnvironmentsList.DisplayAs is not a valid value!")
}

// return a list of available virtual go environments for the user
func (e *EnvironmentsList) getEnvironments() ([]string, []string, error) {
	envs_path := filepath.Join("~", ".VenGO", "*")
	files, err := filepath.Glob(cache.ExpandUser(envs_path))
	if err != nil {
		fmt.Println("while getting list of environments:", err)
		return nil, nil, err
	}
	available, invalid := []string{}, []string{}
	for _, file := range files {
		var vengoenv string
		filename := path.Base(file)
		stat, err := os.Stat(file)
		if err != nil {
			fmt.Println("while getting list of environments:", err)
			return nil, nil, err
		}
		if stat.IsDir() && filename != "bin" && filename != "scripts" {
			_, err := os.Open(filepath.Join(file, "bin", "activate"))
			if err != nil {
				if os.IsNotExist(err) || os.IsPermission(err) {
					invalid = append(invalid, filename)
				}
				continue
			}
			if r, err := os.Readlink(filepath.Join(file, "lib")); err != nil {
				if os.IsNotExist(err) || os.IsPermission(err) {
					invalid = append(invalid, filename)
				}
				continue
			} else {
				if e.DisplayAs == Text {
					vengoenv = fmt.Sprintf("%-22s%-8s", filename, path.Base(r))
				} else {
					vengoenv = filename
				}
				if _, err := os.Stat(r); err != nil {
					if os.IsNotExist(err) || os.IsPermission(err) {
						invalid = append(invalid, vengoenv)
					}
					continue
				}
			}

			available = append(available, vengoenv)
		}
	}

	return available, invalid, nil
}
