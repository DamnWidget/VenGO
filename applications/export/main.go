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
	"path/filepath"

	"github.com/DamnWidget/VenGO/commands"
	"github.com/DamnWidget/VenGO/utils"
	flag "github.com/ogier/pflag"
)

var name string
var force, prettify, help bool

func init() {
	flag.StringVarP(&name, "name", "n", "", "Manifest name")
	flag.BoolVarP(&force, "force", "f", false, "Force operation")
	flag.BoolVarP(&prettify, "prettify", "p", false, "Prettify output")
	flag.BoolVarP(&help, "help", "h", false, "display help message")
	flag.Parse()
}

// main function entry point
func main() {
	if help {
		displayHelp()
		os.Exit(1)
	}

	// build the list value based on the given options
	options := buildCommandListOptions()
	// set the environment if any
	env := flag.Args()
	if len(env) != 0 {
		options = append(options, func(e *commands.Export) {
			e.Environment = env[0]
		})
	}

	e := commands.NewExport(options...)
	if e.Err() != nil {
		fmt.Println(utils.Fail(fmt.Sprint(e.Err())))
		os.Exit(1)
	}
	if e.Exists() {
		if !force {
			fmt.Println(utils.Fail(fmt.Sprintf(
				"%s already exists", filepath.Join(e.Environment, e.Name))))
			fmt.Printf(
				"  %s: use the --force option to overwrite it\n",
				utils.Ok("suggestion"),
			)
			os.Exit(1)
		}
	}
	_, err := e.Run()
	if err != nil {
		fmt.Println(utils.Fail(fmt.Sprintf("error: %v", err)))
		os.Exit(1)
	}
	os.Exit(0)
}

// build the command list options based in the passed flags
func buildCommandListOptions() []func(*commands.Export) {
	options := []func(*commands.Export){}
	if force {
		options = append(options, func(e *commands.Export) {
			e.Force = true
		})
	}
	if prettify {
		options = append(options, func(e *commands.Export) {
			e.Prettify = true
		})
	}
	options = append(options, func(e *commands.Export) {
		e.Name = name
	})

	return options
}

// display help message
func displayHelp() {
	fmt.Println(fmt.Sprintf(`%s: vengo export [options] (environmet)
	-n, --name                  The name for the manifest file that will be created
	-f, --force                 Force Will overwrite other exports already present
	-p, --prettify              Write prettify JSON output

	-h, --help                  Display this message

The environment to export is optional, if nothing is passed, VenGO just try
to export the in use environment, if there is no environment being used and
no environment is specified in the command invocation, it will fail
	`, utils.Ok("Usage")))
}
