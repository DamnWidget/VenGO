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

	"github.com/DamnWidget/VenGO/cache"
	"github.com/DamnWidget/VenGO/env"
	"github.com/DamnWidget/VenGO/utils"
)

// mkenv command
type Mkenv struct {
	Force   bool
	Name    string
	Prompt  string
	Version string
}

// Create a new mkenv command and return back it's address
func NewMkenv(options ...func(i *Mkenv)) *Mkenv {
	mkenv := &Mkenv{Name: "", Prompt: "", Version: ""}
	for _, option := range options {
		option(mkenv)
	}
	if mkenv.Prompt == "" {
		mkenv.Prompt = mkenv.Name
	}
	return mkenv
}

// implements the Runner interface creating the new virtual environment
func (m *Mkenv) Run() (string, error) {
	fmt.Fprint(cache.Output, "Checking intalled go versions... ")
	if err := m.checkInstalled(); err != nil {
		fmt.Fprintln(cache.Output, utils.Fail("✖"))
		return "", err
	}
	fmt.Fprintln(cache.Output, utils.Ok("✔"))

	newEnv := env.NewEnvironment(m.Name, m.Prompt)
	if newEnv.Exists() && !m.Force {
		suggest := fmt.Sprintf(
			"  %s: use --force to force reinstallation", utils.Ok("suggestion"))
		return "", fmt.Errorf("error: %s already exists\n%s", m.Name, suggest)
	}
	if err := newEnv.Generate(); err != nil {
		os.RemoveAll(filepath.Join(os.Getenv("VENGO_HOME"), m.Name))
		return "", err
	}
	if err := newEnv.Install(m.Version); err != nil {
		os.RemoveAll(filepath.Join(os.Getenv("VENGO_HOME"), m.Name))
		return "", err
	}

	return fmt.Sprintf(
		"%s", utils.Ok(fmt.Sprintf(
			"Go %s environment created using %s", m.Name, m.Version))), nil
}

// check if the Go version used to generate the virtual environment is
// installed or not, if is not, return a NotIntalled error type
func (m *Mkenv) checkInstalled() error {
	l := NewList(func(l *List) {
		l.DisplayAs = Json
	})
	result, err := l.Run()
	if err != nil {
		return err
	}
	jsonData := new(BriefJSON)
	if err := json.Unmarshal([]byte(result), jsonData); err != nil {
		return err
	}
	for _, v := range jsonData.Installed {
		if v == m.Version {
			return nil
		}
	}

	return ErrNotInstalled
}
