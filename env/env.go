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

var environTemplate = "tpl/activate"

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

// checks if a environment already exists
func (e *Environment) Exists() bool {
	if _, err := os.Stat(e.VenGO_PATH); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// generate an environment file
func (e *Environment) Generate() error {
	if os.Getenv("VENGO_HOME") != "" {
		environTemplate = filepath.Join(
			os.Getenv("VENGO_HOME"), "scripts", "tpl", "activate")
	}
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

	link := func() error {
		return os.Symlink(path, filepath.Join(e.VenGO_PATH, "lib"))
	}

	if err := link(); err != nil {
		if os.IsExist(err) {
			os.Remove(filepath.Join(e.VenGO_PATH, "lib"))
			if err := link(); err != nil {
				fmt.Println("while creating symlink:", err)
				return err
			}
		} else {
			fmt.Println("while creating symlink:", err)
			return err
		}
	}

	return nil
}

// return back a list of packages installed in the environment
func (e *Environment) Packages() ([]*Package, error) {
	return nil, nil
}

// environment manifest structure
type envManifest struct {
	name      string
	path      string
	goVersion string
	packages  []*packageManifest
}

// creates a new envManifest
func NewEnvManifest(env *Environment, options ...func(em *envManifest)) (*envManifest, error) {
	em := new(envManifest)
	for _, option := range options {
		option(em)
	}
	if err := em.getPackages(env); err != nil {
		return nil, err
	}
	return em, nil
}

// detect all the environment manifest packages and populate its own manifests
func (em *envManifest) getPackages(env *Environment) error {
	packages, err := env.Packages()
	if err != nil {
		return err
	}
	for _, p := range packages {
		name := func(pm *packageManifest) { pm.name = p.Name }
		url := func(pm *packageManifest) { pm.url = p.Url }
		pm, err := NewPackageManifest(env, name, url)
		if err != nil {
			return err
		}
		em.packages = append(em.packages, pm)
	}

	return nil
}
