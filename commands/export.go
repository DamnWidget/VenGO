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
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/DamnWidget/VenGO/env"
	"github.com/DamnWidget/VenGO/utils"
)

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
	filename := filepath.Join(e.Environment, "bin", "activate")
	activateFile, err := ioutil.ReadFile(filename)
	byteLines := bytes.Split(activateFile, []byte("\n"))
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`"(.*?) `)
	prompt := strings.TrimRight(strings.TrimLeft(
		re.FindAllString(string(byteLines[86]), 1)[0], `"`), " ")
	environment := env.NewEnvironment(path.Base(e.Environment), prompt)
	return environment, nil
}

// check if a manifest already exists for the given environment
func (e *Export) Exists() bool {
	_, err := os.Stat(filepath.Join(e.Environment, e.Name))
	return err == nil
}
