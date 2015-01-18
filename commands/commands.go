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
	"log"
	"strings"
	"text/template"

	"github.com/DamnWidget/VenGO/utils"
	flag "github.com/ogier/pflag"

	"os"
)

// type of output format for commands
const (
	Text = iota
	Json
)

// template constants
const (
	// command template
	commandTpl = `
Usage: vengo command [arguments]

Where command can be one of the list below:
{{ range . }}
   {{ .Name | Ok | printf "%-26s" }} {{ .Short }}{{ end }}

Use "vengo help command" for detailed information about any command
`
	// help template
	helpTpl = `
Usage: vengo {{ .Usage }}

{{ .Long }}
`
)

var suggest = utils.Ok("suggestion")

// execute command function type
type commandFunc func(cmd *Command, args ...string)

// command structure
type Command struct {
	Name    string       // command name
	Usage   string       // short line that contains the usage help
	Short   string       // short description
	Long    string       // long description
	Execute commandFunc  // run the command
	Flag    flag.FlagSet // set of flags for this command
}

// return a string representation of the command
func (cmd *Command) String() string {
	return cmd.Name
}

// display the usage and exit
func (cmd *Command) DisplayUsageAndExit() {
	fmt.Printf("Usage: vengo %s\n\n", cmd.Usage)
	fmt.Printf("%s: execute 'vengo' with no arguments to get a list of valid commands\n", suggest)
	os.Exit(2)
}

// register the command in the commands list
func (cmd *Command) register() {
	Commands[cmd.Name] = cmd
}

// used when the commans is not found
func NonCommand(cmd string) {
	fmt.Printf("Command '%s' doesn't looks like a valid VenGO command...\n", cmd)
	fmt.Printf("%s: execute 'vengo' with no arguments to get a list of valid commands\n", suggest)
}

// version function
func Version(version string) {
	version = utils.Ok(version)
	fmt.Printf("VenGO, Virtual Golang Environment builder %s\n", version)
}

// help function
func Help(args ...string) {
	if len(args) == 0 {
		usage()
		return
	}
	cmd, ok := Commands[args[0]]
	if !ok {
		NonCommand(args[0])
		os.Exit(2)
	}
	t := template.New("help")
	template.Must(t.Parse(fmt.Sprintf("%s\n\n", strings.TrimSpace(helpTpl))))
	if err := t.Execute(os.Stdout, cmd); err != nil {
		log.Fatal(err)
	}
}

// calls usage to prints usage information and exits
func Usage() {
	usage()
	os.Exit(2)
}

// print usage information
func usage() {
	t := template.New("usage")
	t.Funcs(template.FuncMap{"Ok": utils.Ok})
	template.Must(t.Parse(fmt.Sprintf("%s\n\n", strings.TrimSpace(commandTpl))))
	if err := t.Execute(os.Stdout, Commands); err != nil {
		log.Fatal(err)
	}
}

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

// commands map, useful for the help command
var Commands map[string]*Command = map[string]*Command{}
