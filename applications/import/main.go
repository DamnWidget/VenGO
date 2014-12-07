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

var prompt string
var force, verbose, help bool

func init() {
	flag.StringVarP(&prompt, "prompt", "p", "", "Use it as env prompt")
	flag.BoolVarP(&force, "force", "f", false, "Force operation")
	flag.BoolVarP(&verbose, "verbose", "v", false, "Verbosity")
	flag.BoolVarP(&help, "help", "h", false, "display help message")
	flag.Parse()
}

// main function entry point
func main() {
	manifest := flag.Args()
	if len(manifest) == 0 || help {
		displayHelp()
		os.Exit(1)
	}
	// make sure that the manifest file exists
	if _, err := os.Stat(manifest[0]); err != nil {
		fmt.Println(utils.Fail(
			fmt.Sprintf("error: can't open %s manifest file", manifest[0])))
		os.Exit(1)
	}

	// build the list value based on the given options
	options := buildCommandListOptions()
	options = append(options, func(i *commands.Import) {
		i.Manifest = manifest[0]
	})
	i := commands.NewImport(options...)
	out, err := i.Run()
	if err != nil {
		fmt.Println(utils.Fail(fmt.Sprintf("error: %v", err)))
		os.Exit(1)
	}
	fmt.Printf("%s\n", utils.Ok(out))
	os.Exit(0)
}

// build the command list options based in the passed flags
func buildCommandListOptions() []func(*commands.Import) {
	options := []func(*commands.Import){}
	if force {
		options = append(options, func(i *commands.Import) {
			i.Force = true
		})
	}
	if verbose {
		options = append(options, func(i *commands.Import) {
			i.Verbose = true
		})
	}
	options = append(options, func(i *commands.Import) {
		i.Prompt = prompt
	})

	return options
}

// display help message
func displayHelp() {
	fmt.Println(fmt.Sprintf(`%s: vengo import [options] (manifest_file)
	-p, --prompt			The prompt of the environment (force it)
	-f, --force			Force Will overwrite any already present environment
	-v, --verbose			Verbose output

	-h, --help			Display this message
	`, utils.Ok("Usage")))
}
