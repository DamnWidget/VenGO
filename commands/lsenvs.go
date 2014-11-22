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
	"github.com/DamnWidget/VenGO/logger"
)

type EnvironmentsJSON struct {
	Available []string `json:"available"`
	Invalid   []string `josn:"invalid"`
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
		logger.Println("while running EnvironmentsList command:", err)
		return "error while running the command", err
	}

	return e.display(available, invalid)
}

// generates the output for the lsenvs command
func (e *EnvironmentsList) display(
	available, invalid []string) (string, error) {

	output := []string{}
	if e.DisplayAs == Text {
		output = append(output, Ok("Virtual Go Environments"))
		for _, v := range available {
			output = append(output, fmt.Sprintf("    %s %s", v, Ok("✔")))
		}
		for _, v := range invalid {
			output = append(output, fmt.Sprintf("    %s %s", v, Fail("✖")))
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
		logger.Println("while getting list of environments:", err)
		return nil, nil, err
	}
	available, invalid := []string{}, []string{}
	for _, file := range files {
		stat, err := os.Stat(file)
		if err != nil {
			logger.Println("while getting list of environments:", err)
			return nil, nil, err
		}
		if stat.IsDir() {
			_, err := os.Open(filepath.Join(file, "bin", "activate"))
			if err != nil {
				if os.IsNotExist(err) || os.IsPermission(err) {
					invalid = append(invalid, path.Base(file))
				}
				continue
			}
			available = append(available, path.Base(file))
		}
	}

	return available, invalid, nil
}
