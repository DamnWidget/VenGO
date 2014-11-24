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

import "errors"

// type of output format for commands
const (
	Text = iota
	Json
)

// Runner is a interface that wraps the execution of a command
//
// Runner returns a string (that can be empty) with the results of the
// executed command and an error, that should be nil if no error occurred
type Runner interface {
	Run() (string, error)
}

// error used when Go version used for mkenv is not installed yet
var ErrNotInstalled = errors.New("Go version not installed")

// determine if the given error is of ErrNotInstalled type
func IsNotInstalledError(err error) bool {
	return err == ErrNotInstalled
}
