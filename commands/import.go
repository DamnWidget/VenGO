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

type Import struct {
	Manifest string
	Prompt   string
	Verbose  bool
	Force    bool
}

// create a new import command and return back it's address
func NewImport(options ...func(i *Import)) *Import {
	imp := new(Import)
	for _, option := range options {
		option(imp)
	}
	if imp.Prompt == "" {
		imp.Prompt = fmt.Sprintf("[%s]", imp.Manifest)
	}
	return imp
}

// implements the Runner interface importing and recreating an exported env
func (i *Import) Run() (string, error) {
	return i.envImport()
}

// import the given manifest and create a new environment based on it
func (i *Import) envImport() (string, error) {
	fmt.Printf("Loading manifest file %s... ", i.Manifest)
	manifest, err := env.LoadManifest(i.Manifest)
	if err != nil {
		fmt.Println(utils.Fail("✖"))
		return "", err
	}
	fmt.Println(utils.Ok("✔"))
	fmt.Printf(
		"Creating %s environment (this may take a while) ...", manifest.Path)
	err = manifest.GenerateEnvironment(i.Verbose, i.Prompt)
	if err != nil {
		fmt.Println(utils.Fail("✖"))
		return "", err
	}
	fmt.Println(utils.Ok("✔"))
	return fmt.Sprintf(
		"%s has been created into %s use vengo activate %s to active it",
		manifest.Name, manifest.Path,
	), nil
}
