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

package env

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"text/template"

	"github.com/DamnWidget/VenGO/cache"
)

const environTemplate = "tpl/activate"

type Environment struct {
	Goroot     string
	Gotooldir  string
	Gopath     string
	PS1        string
	VenGO_PATH string
}

// Create a new Environment struct and return it addrees back
func NewEnvironment(name, prompt string) *Environment {
	VenGO_PATH := filepath.Join(cache.VenGO_PATH, name)
	osArch := fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)
	return &Environment{
		Goroot:     filepath.Join(VenGO_PATH, "lib"),
		Gotooldir:  filepath.Join(VenGO_PATH, "lib", "pkg", "tool", osArch),
		Gopath:     filepath.Join(VenGO_PATH),
		PS1:        prompt,
		VenGO_PATH: VenGO_PATH,
	}
}

// Generate an environment file
func (e *Environment) Generate() error {
	file, err := e.checkPath()
	if err != nil {
		return err
	}
	defer file.Close()
	activateTpl, err := ioutil.ReadFile(environTemplate)
	if err != nil {
		fmt.Println("while reading activate script template file:", err)
		return err
	}
	tpl := template.Must(template.New("environment").Parse(string(activateTpl)))
	err = tpl.Execute(file, e)
	if err != nil {
		fmt.Println("while generating environment template:", err)
		return err
	}

	return nil
}

// checks if the environment path exists, and create it if doesn't
// returns a file object or error if fails
func (e *Environment) checkPath() (*os.File, error) {
	fileName := filepath.Join(e.VenGO_PATH, "bin", "activate")
	return e.createFile(fileName)
}

func (e *Environment) createFile(filename string) (*os.File, error) {
	file, err := os.OpenFile(
		filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(filepath.Join(e.VenGO_PATH, "bin"), 0755)
			if err != nil {
				return nil, err
			}
			return e.createFile(filename)
		}
		return nil, err
	}
	return file, nil
}

// install the given version into the environment creating a Symlink to it
func (e *Environment) Install(ver string) error {
	if !cache.AlreadyCompiled(ver) {
		if err := cache.Compile(ver, false); err != nil {
			fmt.Println("while installing:", err)
			return err
		}
	}

	path := filepath.Join(cache.CacheDirectory(), ver)
	if _, err := os.Stat(path); err != nil {
		path = filepath.Join(cache.CacheDirectory(), fmt.Sprintf("go%s", ver))
	}
	if err := os.Symlink(path, filepath.Join(e.VenGO_PATH, "lib")); err != nil {
		fmt.Println("while creating symlink:", err)
		return err
	}

	return nil
}
