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
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/DamnWidget/VenGO/cache"
	"github.com/DamnWidget/VenGO/utils"
)

var cmdUninstall = &Command{
	Name:  "uninstall",
	Usage: "uninstall version",
	Short: "Uninstall an installed Go version",
	Long: `Uninstalls a Go installed version, it doesn't remove any virtual Go
environment that has been created using the deleted version but it will be
shown as compromised by the 'lsenvs' command. Remember that environments can
be migrated to other Go versions using the 'vengo migrate' command.
`,
	Execute: runUninstall,
}

// initialize command
func init() {
	cmdUninstall.register()
}

// run the uninstall command
func runUninstall(cmd *Command, args ...string) {
	if len(args) == 0 {
		cmd.DisplayUsageAndExit()
	}
	version := args[0]
	activeEnv := os.Getenv("VENGO_ENV")
	if activeEnv != "" {
		if err := checkEnvironment(version, activeEnv); err != nil {
			fmt.Println("error:", err)
			fmt.Printf("%s: execute 'deactivate' before call this command\n", suggest)
			os.Exit(2)
		}
	}

	versionPath := filepath.Join(cache.CacheDirectory(), version)
	if _, err := os.Stat(versionPath); err != nil {
		log.Println(err)
		if os.IsNotExist(err) {
			fmt.Printf("%s is not a Go installed version...\n", version)
			fmt.Printf("%s: try with 'vengo list'\n", suggest)
			os.Exit(2)
		}
		log.Fatal(err)
	}
	err := os.RemoveAll(versionPath)
	if err == nil {
		fmt.Printf("%s has been uninstalled\n", utils.Ok(version))
	}
	log.Fatal(err)
}

// checks if the in use environment has relation with a specific Go version
func checkEnvironment(version, env string) error {
	versionLink, err := os.Readlink(filepath.Join(env, "lib"))
	if err != nil {
		log.Fatal(err)
	}
	if path.Base(versionLink) == version {
		return fmt.Errorf(
			"%s is currently in use by the active environment", version, path.Base(env))
	}
	return nil
}
