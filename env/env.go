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
	"path"
	"path/filepath"
	"runtime"
	"strings"
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
		Gopath:     VenGO_PATH,
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
	} else {
		_, caller, _, ok := runtime.Caller(1)
		if ok {
			// we are running in a test environment
			environTemplate = filepath.Join(
				path.Dir(caller), "..", "env", environTemplate)
		}
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
func (e *Environment) Packages(environment ...string) ([]*Package, error) {
	envPath := os.Getenv("VENGO_ENV")
	if len(environment) > 0 {
		envPath = environment[0]
	}
	if envPath == "" {
		return nil, fmt.Errorf("VENGO_ENV environment variable is not set")
	}
	basePath := filepath.Join(envPath, "src")
	packages := []*Package{}
	if err := filepath.Walk(
		basePath,
		func(walkPath string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Fprintf(cache.Output,
					"%s ignored because error: %s\n", walkPath, err)
				return nil
			}
			if !info.IsDir() {
				return nil
			}
			for _, vcs := range vcsTypes {
				_, err := os.Stat(filepath.Join(walkPath, "."+vcs.name))
				if err == nil {
					options := func(p *Package) {
						splitPaths := strings.Split(walkPath, basePath+"/")
						p.Name = path.Base(splitPaths[1])
						p.Url = splitPaths[1]
						p.Installed = true
						p.Vcs = vcs.name
					}
					packages = append(packages, NewPackage(options))
				}
			}
			return nil
		},
	); err != nil {
		return nil, err
	}
	return packages, nil
}

// generates an environment manifest from a configured environment
func (e *Environment) Manifest() (*envManifest, error) {
	general := func(em *envManifest) {
		em.Name = path.Base(e.VenGO_PATH)
		em.Path = e.VenGO_PATH
		em.Packages = []*packageManifest{}
	}
	lib, err := os.Readlink(e.Goroot)
	if err != nil {
		return nil, err
	}
	goVersion := func(em *envManifest) {
		em.GoVersion = path.Base(lib)
	}
	return NewEnvManifest(e, general, goVersion)
}
