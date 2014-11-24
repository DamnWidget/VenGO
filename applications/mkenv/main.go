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

var force, help bool
var prompt, version string

func init() {
	flag.BoolVarP(&force, "force", "f", false, "force re-installs")
	flag.BoolVarP(&help, "help", "h", false, "display help message")
	flag.StringVarP(&prompt, "prompt", "p", "", "environment prompt")
	flag.StringVarP(&version, "go", "g", "tip", "Go version to use")
	flag.Parse()
}

// main function entry point
func main() {
	name := flag.Args()
	if len(name) == 0 || help {
		displayHelp()
		os.Exit(1)
	}

	// build the list object based on the given options
	options := buildCommandListOptions()
	options = append(options, func(m *commands.Mkenv) {
		m.Name = name[0]
	})

	mkenv := commands.NewMkenv(options...)
	data, err := mkenv.Run()
	if err != nil {
		if commands.IsNotInstalledError(err) {
			fmt.Println(fmt.Sprintf(
				"sorry vengo can't perform the operation because %s is %s",
				mkenv.Version, utils.Fail("not installed")),
			)
			fmt.Printf(
				"  %s: run 'vengo install %s'\n",
				utils.Ok("suggestion"), mkenv.Version,
			)
			os.Exit(1)
		}
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(data)
}

// build the command list options based in the passed flags
func buildCommandListOptions() []func(*commands.Mkenv) {
	options := []func(m *commands.Mkenv){}
	if prompt != "" {
		options = append(options, func(m *commands.Mkenv) {
			m.Prompt = prompt
		})
	}
	if version != "" {
		options = append(options, func(m *commands.Mkenv) {
			m.Version = version
		})
	}
	if force {
		options = append(options, func(m *commands.Mkenv) {
			m.Force = true
		})
	}

	return options
}

// display help message
func displayHelp() {
	fmt.Println(fmt.Sprintf(`%s: vengo mkenv [options] env_name
    -f, --force             Force Reinstallation if environment exists
    -p, --prompt            Environment prompt to be used
    -g, --go                Go version to be used (tip by default)

    -h, --help              Display this message
    `, utils.Ok("Usage")))
}
