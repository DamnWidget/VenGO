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
	"errors"
	"fmt"

	"github.com/DamnWidget/VenGO/cache"
	"github.com/DamnWidget/VenGO/utils"
)

// possible installation sources
const (
	Mercurial = iota
	Source
	Binary
)

// instal command
type Install struct {
	Force     bool
	Source    int
	Version   string
	DisplayAs int
	Verbose   bool
}

// Create a new install command and return back it's address
func NewInstall(options ...func(i *Install)) *Install {
	install := new(Install)
	for _, option := range options {
		option(install)
	}

	return install
}

// implements the Runner interface executing the required installation
func (i *Install) Run() (string, error) {
	switch i.Source {
	case Mercurial:
		return i.fromMercurial()
	case Source:
		return i.fromSource()
	case Binary:
		return i.fromBinary()
	default:
		return "", errors.New("Intall.Source is not a valid source")
	}
}

// install from mercurial source
func (i *Install) fromMercurial() (string, error) {
	if err := cache.CacheDonwloadMercurial(i.Version, i.Force); err != nil {
		return "error while installing from mercurial", err
	}
	if err := cache.Compile(i.Version, i.Verbose); err != nil {
		return "error while compiling from mercurial", err
	}

	result := fmt.Sprintf(
		"%s", utils.Ok(fmt.Sprintf("Go %s installed", i.Version)))
	return result, nil
}

// install from tar.gz source
func (i *Install) fromSource() (string, error) {
	if err := cache.CacheDownload(i.Version, i.Force); err != nil {
		return "error while installing from tar.gz source", err
	}
	if err := cache.Compile(i.Version, i.Verbose); err != nil {
		return "error while installing from tar.gz source", err
	}

	result := fmt.Sprintf(
		"%s", utils.Ok(fmt.Sprintf("Go %s installed", i.Version)))
	return result, nil
}

// install from binary source
func (i *Install) fromBinary() (string, error) {
	if err := cache.CacheDownloadBinary(i.Version, i.Force); err != nil {
		return "wrror while installing from binary", err
	}

	result := fmt.Sprintf(
		"%s", utils.Ok(fmt.Sprintf("Go %s installed", i.Version)))
	return result, nil
}
