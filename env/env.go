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
	"log"
	"os"
	"os/exec"
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
	if err := e.generateBash(); err != nil {
		return err
	}
	return e.generateFish()
}

func (e *Environment) generateBash() error {
	if os.Getenv("VENGO_HOME") != "" {
		environTemplate = filepath.Join(
			os.Getenv("VENGO_HOME"), "scripts", "tpl", "activate")
	} else {
		if environTemplate == "tpl/activate" {
			_, caller, _, ok := runtime.Caller(1)
			if ok {
				// we are running in a test environment
				environTemplate = filepath.Join(
					path.Dir(caller), "..", "env", environTemplate)
			}
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

func (e *Environment) generateFish() error {
	if os.Getenv("VENGO_HOME") != "" {
		environTemplate = filepath.Join(
			os.Getenv("VENGO_HOME"), "scripts", "tpl", "activate.fish")
	} else {
		if environTemplate == "tpl/activate.fish" {
			_, caller, _, ok := runtime.Caller(1)
			if ok {
				// we are running in a test environment
				environTemplate = filepath.Join(
					path.Dir(caller), "..", "env", environTemplate)
			}
		}
	}
	file, err := e.checkPath(struct{}{})
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
// returns a file value or error if fails
func (e *Environment) checkPath(fish ...struct{}) (*os.File, error) {
	sh := "activate"
	if len(fish) != 0 {
		sh = "activate.fish"
	}
	fileName := filepath.Join(e.VenGO_PATH, "bin", sh)
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
		if err := cache.Compile(ver, false, false); err != nil {
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
					curr, _ := os.Getwd()
					os.Chdir(walkPath)
					defer os.Chdir(curr)
					args := strings.Split(vcs.refCmd, " ")
					out, err := exec.Command(args[0], args[1:]...).CombinedOutput()
					if err != nil {
						_, ok := os.Stat(filepath.Join("."+vcs.name, "test"))
						if ok == nil {
							// we are in the test suite
							out = []byte{}
						} else {
							log.Printf("warning %s skypped: %s", walkPath, string(out))
							return nil
						}
					}

					revision := strings.TrimRight(string(out), "\n")
					options := func(p *Package) {
						splitPaths := strings.Split(walkPath, basePath+"/")
						p.Name = path.Base(splitPaths[1])
						p.Url = splitPaths[1]
						p.Root = path.Dir(p.Url)
						p.Installed = true
						p.Vcs = vcs.name
						p.CodeRevision = revision
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

// activate environment in the current process
func (e *Environment) activate() {
	os.Setenv("VENGO_ENV", cache.ExpandUser(e.VenGO_PATH))
	os.Setenv("GOTOOLDIR", cache.ExpandUser(e.Gotooldir))
	os.Setenv("GORROT", cache.ExpandUser(e.Goroot))
	os.Setenv("GOPATH", cache.ExpandUser(e.Gopath))
}

// deactivate environment in the current process
func (e *Environment) deactivate() {
	os.Setenv("VENGO_ENV", "")
	os.Setenv("GOTOOLDIR", "")
	os.Setenv("GOROOT", "")
	os.Setenv("GOPATH", "")
}
