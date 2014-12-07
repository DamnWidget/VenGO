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

	"github.com/DamnWidget/VenGO/utils"
)

// Package struct
type Package struct {
	Name         string
	Url          string
	Root         string
	Installed    bool
	Vcs          string
	CodeRevision string
}

// create a new package and returns it's address back
func NewPackage(options ...func(p *Package)) *Package {
	p := new(Package)
	for _, option := range options {
		option(p)
	}
	return p
}

// create a string representation of a package
func (p *Package) String() string {
	check := utils.Ok("✔")
	if !p.Installed {
		check = utils.Fail("✖")
	}
	return fmt.Sprintf(`    %s(%s) %s`, p.Name, p.Url, check)
}
