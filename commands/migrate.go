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
	"fmt"
	"os"
	"path/filepath"

	"github.com/DamnWidget/VenGO/cache"
	"github.com/DamnWidget/VenGO/env"
	"github.com/DamnWidget/VenGO/utils"
)

var cmdMigrate = &Command{
	Name:    "migrate",
	Usage:   "migrate environment go_version",
	Short:   "Migrate an environment version of Go",
	Long:    fmt.Sprintf("Migrate an already existent VenGO environment from a Go version to another."),
	Execute: runMigrate,
}

// initialize the command
func init() {
	cmdMigrate.register()
}

// run the migrate command
func runMigrate(cmd *Command, args ...string) {
	if len(args) < 2 {
		cmd.DisplayUsageAndExit()
	}
	environName := args[0]
	goVersion := args[1]
	if os.Getenv("VENGO_ENV") == environName {
		fmt.Println(
			"error:", fmt.Sprintf("%s is currently in use as the active environment"))
		fmt.Printf("%s: execute 'deactivate' before call this command\n", suggest)
		os.Exit(2)
	}
	envPath := filepath.Join(os.Getenv("VENGO_HOME"), environName)
	if _, err := os.Stat(envPath); err != nil {
		fmt.Printf("%s is no a VenGO environment...\n", environName)
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Print("Checking installed Go versions...")
	if !env.LookupInstalledVersion(goVersion) {
		fmt.Println(utils.Fail("✖"))
		fmt.Println(fmt.Sprintf(
			"sorry vengo can't perform the operation because %s is %s",
			goVersion, utils.Fail("not installed")))
		fmt.Printf("%s: run 'vengo install %s'\n", suggest, goVersion)
		os.Exit(2)
	}
	fmt.Println(utils.Ok("✔"))

	library := filepath.Join(envPath, "lib")
	installed, _ := os.Readlink(library)
	if installed == goVersion {
		fmt.Printf("%s in use Go version is already %s, skipping...", environName, installed)
		os.Exit(0)
	}
	fmt.Printf("unlinking previous %s Go version from %s...", installed, environName)
	if err := os.Remove(library); err != nil {
		fmt.Println(utils.Fail("✖"))
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Println(utils.Ok("✔"))
	fmt.Printf("Linking Go version %s...", goVersion)
	path := filepath.Join(cache.CacheDirectory(), goVersion)
	err := os.Symlink(path, filepath.Join(envPath, "lib"))
	if err != nil {
		fmt.Println(utils.Fail("✖"))
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Println(utils.Ok("✔"))
	fmt.Println("Done. You may want to run 'go build -a ./... in $GOPATH")
	os.Exit(0)
}
