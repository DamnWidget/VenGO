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

package env

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/DamnWidget/VenGO/cache"
	"github.com/DamnWidget/VenGO/logger"
)

var environTemplate = `#!/bin/bash

export PREV_GOROOT=$GOROOT
export PREV_GOTOOLDIR=$GOTOOLDIR
export PREV_GOPATH=$GOPATH
export PREV_PS1=$PS1

export GOROOT={{.Goroot}}
export GOTOOLDIR={{.Gotooldir}}
export GOPATH={{.Gopath}}

export PS1="{{.PS1}} $PS1"
`

type Environment struct {
	Goroot, Gotooldir, Gopath, PS1 string
}

// Create a new Environment struct and return it addrees back
func NewEnvironment(root, tooldir, path, name string) *Environment {
	return &Environment{root, tooldir, path, name}
}

// Generate an environment file
func (e *Environment) Generate() {
	file, err := os.OpenFile(
		filepath.Join(cache.VenGO_PATH, e.PS1),
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644,
	)
	if err != nil {
		logger.Fatal(err)
	}
	defer file.Close()
	tpl := template.Must(template.New("environment").Parse(environTemplate))
	err = tpl.Execute(file, e)
	if err != nil {
		logger.Println("while generating environment template:", err)
	}
}
