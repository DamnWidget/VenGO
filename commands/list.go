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
	"strings"

	"github.com/DamnWidget/VenGO/cache"
)

const (
	Brief = iota
	Extended
	JsonBrief
	JsonExtended
)

// list command
type List struct {
	ShowInstalled    bool
	ShowNotInstalled bool
	ShowBoth         bool
	Details          int
}

// Create a new list and return back it's address
func NewList(options ...func(*List)) *List {
	list := new(List)
	list.ShowInstalled = true

	for _, option := range options {
		option(list)
	}

	return list
}

// implements the Runner interface returning back a list of installed
// go versions, a list of not installed versions or all versions depending
// on the list options
func (l *List) Run() (string, error) {
	tags := cache.Tags()
	return strings.Join(tags, "\n"), nil
}
