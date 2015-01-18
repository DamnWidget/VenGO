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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/DamnWidget/VenGO/cache"
	"github.com/DamnWidget/VenGO/env"
	"github.com/DamnWidget/VenGO/utils"
)

var cmdExport = &Command{
	Name:  "export",
	Usage: "export [-n] [-f] [-p] environment",
	Short: "Export an environment into a manifest file",
	Long: `Exports a VenGO environment into vengo manifest files using JSON format.
Those manifest files can be used later by the 'vengo import' command to
recreate the exported environment. It generates a JSON file that contains all
the packages that have been installed into the exported VenGO environment
GOPATH using the 'go get' tool (that means that only packages installed using
git, mercurial, bzr or subversion can be exported) capturing the specific
versions used when the package was installed in the exported environment.

If the environment manifest has been already exported, we can force the export
using the -f or --force flags to overwrite it.

A different name can be done to the manifest file passing a parameter to the
-n or --name flag`,
	Execute: runExport,
}

var (
	nameExport     string
	forceExport    bool
	prettifyExport bool
)

// initialize command
func init() {
	cmdExport.Flag.BoolVarP(&forceExport, "force", "f", false, "force export")
	cmdExport.Flag.BoolVarP(&prettifyExport, "prettify", "p", false, "prettify")
	cmdExport.Flag.StringVarP(&nameExport, "name", "n", "", "manifest name")
	cmdExport.register()
}

// run the export command
func runExport(cmd *Command, args ...string) {
	env := func(e *Export) {}
	if len(args) == 0 {
		activeEnv := os.Getenv("VENGO_ENV")
		if activeEnv == "" {
			cmd.DisplayUsageAndExit()
		}
	} else {
		env = func(e *Export) {
			e.Environment = filepath.Join(cache.VenGO_PATH, args[0])
		}
	}
	options := func(e *Export) {
		e.Force = forceExport
		e.Prettify = prettifyExport
		e.Name = nameExport
	}
	e := NewExport(options, env)
	if e.Err() != nil {
		fmt.Println(utils.Fail(fmt.Sprint(e.Err())))
		os.Exit(2)
	}
	if e.Exists() {
		if !e.Force {
			p := filepath.Join(e.Environment, e.Name)
			fmt.Println(utils.Fail(fmt.Sprintf("%s already exists", p)))
			fmt.Printf("%s: use the -f option to overwrite it\n", suggest)
			os.Exit(2)
		}
	}
	_, err := e.Run()
	if err != nil {
		fmt.Println(utils.Fail(fmt.Sprintf("error: %v", err)))
		os.Exit(2)
	}
	os.Exit(0)
}

// export command
type Export struct {
	Environment string
	Name        string
	Force       bool
	Prettify    bool
	err         error
}

// create a new export command and return back it's address
func NewExport(options ...func(e *Export)) *Export {
	export := new(Export)
	for _, option := range options {
		option(export)
	}
	export.normalize()
	return export
}

// implements the Runner interface executing the environment export
func (e *Export) Run() (string, error) {
	return e.envExport()
}

// export the given environment using a VenGO.manifest file
func (e *Export) envExport() (string, error) {
	fmt.Print("Loading environment... ")
	environment, err := e.LoadEnvironment()
	if err != nil {
		fmt.Println(utils.Fail("✖"))
		return "", err
	}
	fmt.Println(utils.Ok("✔"))
	fmt.Print("Generating manifest... ")
	environManifest, err := environment.Manifest()
	if err != nil {
		fmt.Println(utils.Fail("✖"))
		return "", err
	}
	manifest, err := environManifest.Generate()
	if err != nil {
		fmt.Println(utils.Fail("✖"))
		return "", err
	}
	fmt.Println(utils.Ok("✔"))
	if e.Prettify {
		buffer := new(bytes.Buffer)
		json.Indent(buffer, manifest, "", "\t")
		manifest = buffer.Bytes()
	}
	fmt.Printf("Writing manifest into %s... ", environment.VenGO_PATH)
	err = ioutil.WriteFile(
		filepath.Join(environment.VenGO_PATH, e.Name), manifest, 0644)
	if err != nil {
		fmt.Println(utils.Fail("✖"))
	}
	fmt.Println(utils.Ok("✔"))

	return "", err
}

// normalize an export configuration, if there is no environment, try to detect
// if the terminal that called it is in a environment, if so, use it.
// if there is no name use .VenGO.manifest
func (e *Export) normalize() {
	if e.Environment == "" {
		if env := os.Getenv("VENGO_ENV"); env != "" {
			e.Environment = env
		} else {
			e.err = errors.New(
				"there is no active environment and none has been specified")
			return
		}
	}
	if e.Name == "" {
		e.Name = "VenGO.manifest"
	}
}

// expose the internal err property
func (e *Export) Err() error {
	return e.err
}

// load environment using the activate environment script, return an error
// if the operation can't be completed
func (e *Export) LoadEnvironment() (*env.Environment, error) {
	readFile := func(filename string) ([]byte, error) {
		return ioutil.ReadFile(filename)
	}
	filenames := [2]string{
		filepath.Join(e.Environment, "bin", "activate"),
		filepath.Join(cache.VenGO_PATH, e.Environment, "bin", "activate"),
	}
	var err error
	var activateFile []byte
	for _, filename := range filenames {
		activateFile, err = readFile(filename)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}
	byteLines := bytes.Split(activateFile, []byte("\n"))
	re := regexp.MustCompile(`"(.*?) `)
	prompt := strings.TrimRight(strings.TrimLeft(
		re.FindAllString(string(byteLines[86]), 1)[0], `"`), " ")
	environment := env.NewEnvironment(path.Base(e.Environment), prompt)
	return environment, nil
}

// check if a manifest already exists for the given environment
func (e *Export) Exists() bool {
	log.Println(filepath.Join(e.Environment, e.Name))
	_, err := os.Stat(filepath.Join(e.Environment, e.Name))
	return err == nil
}
