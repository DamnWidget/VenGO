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

package main

import (
	"flag"
	"os"

	"github.com/DamnWidget/VenGO/commands"
)

// Main application entry point
func main() {
	flag.Usage = commands.Usage
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		commands.Usage()
	}

	if args[0] == "version" {
		commands.Version(vengo_version)
		return
	}

	if args[0] == "help" {
		commands.Help(args[1:]...)
		return
	}

	cmd, ok := commands.Commands[args[0]]
	if !ok {
		commands.NonCommand(args[0])
		os.Exit(2)
	}
	cmd.Flag.Usage = func() { cmd.DisplayUsageAndExit() }
	cmd.Flag.Parse(args[1:])
	cmd.Execute(cmd, cmd.Flag.Args()...)
	return
}
