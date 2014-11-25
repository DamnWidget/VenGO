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

// vcs type structure
type vcsType struct {
	name      string
	refCmd    string
	updateCmd string
}

// Git
var gitVcs = &vcsType{
	name:      "git",
	refCmd:    "git rev-parse --verify HEAD",
	updateCmd: "git checkout {tag}",
}

// Mercurial
var mercurialVcs = &vcsType{
	name:      "hg",
	refCmd:    "hg --debug id -i",
	updateCmd: "hg update -r {tag}",
}

// Bazaar
var bazaarVcs = &vcsType{
	name:      "bzr",
	refCmd:    "bzr revno",
	updateCmd: "hg update -r revno:{tag}",
}

// SubVersion
var svnVcs = &vcsType{
	name:      "svn",
	refCmd:    `svn info | grep "Revision" | awk '{print $2}'`,
	updateCmd: "svn up -r{tag}",
}

// available vcs types
var vcsTypes = []*vscType{
	gitVcs,
	mercurialVcs,
	bazaarVcs,
	svnVcs,
}
