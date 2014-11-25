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

import "github.com/DamnWidget/VenGO/utils"

// vcs type structure
type vcsType struct {
	name      string
	refCmd    string
	updateCmd string
	cloneCmd  func(repo, tag string) error
}

// Git
var gitVcs = &vcsType{
	name:      "git",
	refCmd:    "git rev-parse --verify HEAD",
	updateCmd: "git checkout {tag}",
	cloneCmd: func(repo, tag string) error {
		if err := utils.Exec(true, "git", "clone", repo); err != nil {
			return err
		}
		err := utils.Exec(true, "git", "checkout", tag)
		return err
	},
}

// Mercurial
var mercurialVcs = &vcsType{
	name:      "hg",
	refCmd:    "hg --debug id -i",
	updateCmd: "hg update -r {tag}",
	cloneCmd: func(repo, tag string) error {
		return utils.Exec(true, "hg", "clone", "-r", tag, repo)
	},
}

// Bazaar
var bazaarVcs = &vcsType{
	name:      "bzr",
	refCmd:    "bzr revno",
	updateCmd: "bzr update -r revno:{tag}",
	cloneCmd: func(branch, rev string) error {
		return utils.Exec(true, "bzr", "branch", branch, "-r", rev)
	},
}

// SubVersion
var svnVcs = &vcsType{
	name:      "svn",
	refCmd:    `svn info | grep "Revision" | awk '{print $2}'`,
	updateCmd: "svn up -r{tag}",
	cloneCmd: func(repo, rev string) error {
		return utils.Exec(true, "svn", "checkout", "-r", rev, repo)
	},
}

// available vcs types
var vcsTypes = []*vcsType{
	gitVcs,
	mercurialVcs,
	bazaarVcs,
	svnVcs,
}

// clone the repo in an scpecific revision, tag or commit
func (vcs *vcsType) Clone(repo, tag string) error {
	return vcs.cloneCmd(repo, tag)
}
