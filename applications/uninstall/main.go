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
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/DamnWidget/VenGO/cache"
	"github.com/DamnWidget/VenGO/utils"
)

func main() {
	if len(os.Args) < 2 {
		displayHelp()
		os.Exit(1)
	}
	version := os.Args[1]
	activeEnv := os.Getenv("VENGO_ENV")
	if activeEnv != "" {
		if err := checkEnvironment(version, activeEnv); err != nil {
			fmt.Println("error:", err)
			fmt.Printf(
				"  %s: execute 'deactivate' before call this command\n",
				utils.Ok("suggestion"),
			)
			os.Exit(1)
		}
	}

	versionPath := filepath.Join(cache.CacheDirectory(), version)
	if _, err := os.Stat(versionPath); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%s is not a Go installed version...\n", version)
			os.Exit(1)
		}
		log.Fatal(err)
	}
	err := os.RemoveAll(filepath.Join(cache.CacheDirectory(), version))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s has been uninstalled\n", utils.Ok(version))
}

// display help message
func displayHelp() {
	fmt.Println(fmt.Sprintf("%s: vengo uninstall version", utils.Ok("Usage")))
}

// checks if the in use environment has relation with a specific Go version
func checkEnvironment(version, env string) error {
	versionLink, err := os.Readlink(filepath.Join(env, "lib"))
	if err != nil {
		log.Fatal(err)
	}
	if path.Base(versionLink) == version {
		return fmt.Errorf(
			"%s is currently in use by active environment '%s'",
			version, path.Base(env),
		)
	}
	return nil
}
