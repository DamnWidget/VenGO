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
	"fmt"
	"os"
	"path/filepath"
)

// package manifest structure
type packageManifest struct {
	Name         string   `json:"package_name"`
	Url          string   `json:"package_url"`
	Vcs          *vcsType `json:"package_vcs,omitempty"`
	CodeRevision string   `json:"package_vcs_revision,omitempty"`
}

type funcOpts func(*packageManifest)

// creates a new packageManifest
func NewPackageManifest(
	env *Environment, options ...funcOpts) (*packageManifest, error) {

	pm := new(packageManifest)
	for _, option := range options {
		option(pm)
	}
	if err := pm.getVcs(env); err != nil {
		return nil, err
	}
	return pm, nil
}

// detect the version control system used for a go package and assign it
func (pm *packageManifest) getVcs(env *Environment) error {
	packagePath := filepath.Join(env.Gopath, "src", pm.Url)
	currdir, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.Chdir(packagePath); err != nil {
		return err
	}
	defer func() { os.Chdir(currdir) }()
	for _, vcs := range vcsTypes {
		vcsdir := fmt.Sprintf(".%s", vcs.name)
		if fi, err := os.Stat(filepath.Join(packagePath, vcsdir)); err == nil {
			if fi.IsDir() {
				pm.Vcs = vcs
				return nil
			}
		}
	}
	return fmt.Errorf("%s is using an unknown version control system", pm.Name)
}
