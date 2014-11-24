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
	"path/filepath"

	"github.com/DamnWidget/VenGO/utils"
)

func main() {
	if len(os.Args) < 2 {
		displayHelp()
		os.Exit(1)
	}
	env := os.Args[1]
	if os.Getenv("VENGO_ENV") == env {
		fmt.Println("error:", fmt.Sprintf(
			"%s is currently in use as the active environment", env))
		fmt.Printf(
			"  %s: execute 'deactivate' before call this command\n",
			utils.Ok("suggestion"),
		)
		os.Exit(1)
	}

	envPath := filepath.Join(os.Getenv("VENGO_HOME"), env)
	if _, err := os.Stat(envPath); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%s is not a VenGO environment...\n", env)
			fmt.Println(err)
			os.Exit(1)
		}
	}
	err := os.RemoveAll(envPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s has been removed\n", utils.Ok(env))
}

// display help message
func displayHelp() {
	fmt.Println(fmt.Sprintf("%s: vengo rmenv environment", utils.Ok("Usage")))
}
