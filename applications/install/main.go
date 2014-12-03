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
	"fmt"
	"os"

	"github.com/DamnWidget/VenGO/commands"
	"github.com/DamnWidget/VenGO/utils"
	flag "github.com/ogier/pflag"
)

var force, binary, source, help, verbose bool

func init() {
	flag.BoolVarP(&force, "force", "f", false, "Force operation")
	flag.BoolVarP(&binary, "binary", "b", false, "Download binaries")
	flag.BoolVarP(&source, "source", "s", false, "Download tar.gz source")
	flag.BoolVarP(&help, "help", "h", false, "display help message")
	flag.BoolVarP(&verbose, "vebose", "v", false, "verbose output")
	flag.Parse()
}

// main function entry point
func main() {
	ver := flag.Args()
	if len(ver) == 0 || help {
		displayHelp()
		os.Exit(1)
	}

	// build the list value based on the given options
	options := buildCommandListOptions()
	// set the version
	options = append(options, func(i *commands.Install) {
		i.Version = ver[0]
	})

	i := commands.NewInstall(options...)
	data, err := i.Run()
	if err != nil {
		fmt.Println(utils.Fail(fmt.Sprintf("error: %v", err)))
		if !verbose {
			fmt.Printf(
				"  %s: run the install commadn with the '-v' option\n",
				utils.Ok("suggestion"),
			)
		}
		os.Exit(1)
	}
	fmt.Println(data)
	os.Exit(0)
}

// build the command list options based in the passed flags
func buildCommandListOptions() []func(*commands.Install) {
	options := []func(*commands.Install){}
	if verbose {
		options = append(options, func(i *commands.Install) {
			i.Verbose = true
		})
	}
	if force {
		options = append(options, func(i *commands.Install) {
			i.Force = true
		})
	}
	if binary {
		options = append(options, func(i *commands.Install) {
			i.Source = commands.Binary
		})
	} else {
		if source {
			options = append(options, func(i *commands.Install) {
				i.Source = commands.Source
			})
		}
	}

	return options
}

// display help message
func displayHelp() {
	fmt.Println(fmt.Sprintf(`%s: vengo install [options] version
	-f, --force		Force Download
	-s, --source		Download from tar.gz source instead of mercurial
	-b, --binary		Download from a binary file (doesn't compile)
	-v, --verbose		Verbose output

	-h, --help              Display this message

%s.
	`, utils.Ok("Usage"),
		utils.Fail("note: if binary is set, source is ignored")))
}
