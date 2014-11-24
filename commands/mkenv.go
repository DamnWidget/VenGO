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
	"fmt"

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
	newEnv := env.NewEnvironment(m.Name, m.Prompt)
	if newEnv.Exists() && !m.Force {
		suggest := fmt.Sprintf(
			"  %s: use --force to force reinstallation", utils.Ok("suggestion"))
		return "", fmt.Errorf("error: %s already exists\n%s", m.Name, suggest)
	}
	if err := newEnv.Generate(); err != nil {
		return "", err
	}
	if err := newEnv.Install(m.Version); err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"%s", utils.Ok(fmt.Sprintf(
			"Go %s environment created using %s", m.Name, m.Version))), nil
}
