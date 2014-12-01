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
	"encoding/json"
	"log"
)

// environment manifest structure
type envManifest struct {
	Name      string             `json:"environment_name"`
	Path      string             `json:"environment_path"`
	GoVersion string             `json:"environment_go_version"`
	Packages  []*packageManifest `json:"environment_packages"`
}

// creates a new envManifest
func NewEnvManifest(env *Environment, options ...func(em *envManifest)) (*envManifest, error) {
	em := new(envManifest)
	for _, option := range options {
		option(em)
	}
	if err := em.getPackages(env); err != nil {
		log.Println(err)
		return nil, err
	}
	return em, nil
}

// detect all the environment manifest packages and populate its own manifests
func (em *envManifest) getPackages(env *Environment) error {
	packages, err := env.Packages()
	if err != nil {
		return err
	}
	for _, p := range packages {
		name := func(pm *packageManifest) { pm.Name = p.Name }
		url := func(pm *packageManifest) { pm.Url = p.Url }
		rev := func(pm *packageManifest) { pm.CodeRevision = p.CodeRevision }
		pm, err := NewPackageManifest(env, name, url, rev)
		if err != nil {
			return err
		}
		em.Packages = append(em.Packages, pm)
	}

	return nil
}

// generate the environment manifest for exports/import
func (em *envManifest) Generate() ([]byte, error) {
	b, err := json.Marshal(em)
	if err != nil {
		return nil, err
	}
	return b, nil
}
