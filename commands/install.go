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
	"os"

	"github.com/DamnWidget/VenGO/cache"
	"github.com/DamnWidget/VenGO/utils"
)

var cmdInstall = &Command{
	Name:  "install",
	Usage: "install [-s] [-b] [-v] [-f] [-n] version",
	Short: "Installs a new Go version",
	Long: `Install a new version of Go, it can be installed directly from the official
mercurial or git repositories, from a tarball packaed source or directly in
binary format in case that you don't want to compile it.

If the -s or --source flag is passed, a compressed tarball is downloaded and
compiled. If the -b or --binary flag is used, an already compiled tarball is
used.

If the given version is already installed, we can force it's reinstallation
using the -f or --force flags, to compile the newly downloaded Go version
with CGO_ENABLED=0 the -n or --ncgo flag should be passed.

Use the -v or --verbose flags to run the command with verbose output, this
is useful to debug in case of errors during the compilation phase.
`,
	Execute: runInstall,
}

var (
	forceInstall   bool
	binaryInstall  bool
	sourceInstall  bool
	verboseInstall bool
	nocgoInstall   bool
)

// possible installation sources
const (
	Mercurial = iota
	Source
	Binary
)

// install command
type Install struct {
	Force     bool
	Source    int
	Version   string
	DisplayAs int
	Verbose   bool
	NoCGO     bool
}

// initialize the command
func init() {
	cmdInstall.Flag.BoolVarP(&forceInstall, "force", "f", false, "Force operation")
	cmdInstall.Flag.BoolVarP(&binaryInstall, "binary", "b", false, "Download binaries")
	cmdInstall.Flag.BoolVarP(&sourceInstall, "source", "s", false, "Download tarball")
	cmdInstall.Flag.BoolVarP(&verboseInstall, "verbose", "v", false, "verbose output")
	cmdInstall.Flag.BoolVarP(&nocgoInstall, "ncgo", "n", false, "CGO_ENABLE=0")
	cmdInstall.register()
}

// fun the install command
func runInstall(cmd *Command, args ...string) {
	if len(args) == 0 {
		cmd.DisplayUsageAndExit()
	}
	options := func(i *Install) {
		i.Verbose = verboseInstall
		i.Force = forceInstall
		i.NoCGO = nocgoInstall
		if binaryInstall {
			i.Source = Binary
		} else {
			if sourceInstall {
				i.Source = Source
			}
		}
		i.Version = args[0]
	}
	i := NewInstall(options)
	data, err := i.Run()
	if err != nil {
		fmt.Println(utils.Fail(fmt.Sprintf("error: %v", err)))
		if !verboseInstall {
			fmt.Printf(
				"%s: run the install command with the '-v' option\n", suggest)
		}
		os.Exit(2)
	}
	fmt.Println(data)
	os.Exit(0)
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
		return i.fromGit()
	case Source:
		return i.fromSource()
	case Binary:
		return i.fromBinary()
	default:
		return "", errors.New("Install.Source is not a valid source")
	}
}

// install from mercurial source
func (i *Install) fromGit() (string, error) {
	if err := cache.CacheDownloadGit(i.Version, i.Force); err != nil {
		return "error while installing from mercurial", err
	}
	if err := cache.Compile(i.Version, i.Verbose, i.NoCGO); err != nil {
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
	if err := cache.Compile(i.Version, i.Verbose, i.NoCGO); err != nil {
		return "error while installing from tar.gz source", err
	}

	result := fmt.Sprintf(
		"%s", utils.Ok(fmt.Sprintf("Go %s installed", i.Version)))
	return result, nil
}

// install from binary source
func (i *Install) fromBinary() (string, error) {
	if err := cache.CacheDownloadBinary(i.Version, i.Force); err != nil {
		return "error while installing from binary", err
	}

	result := fmt.Sprintf(
		"%s", utils.Ok(fmt.Sprintf("Go %s installed", i.Version)))
	return result, nil
}
