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

package cache

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/DamnWidget/VenGO/utils"
)

const REPO = "https://github.com/golang/go.git"

var TARGET = filepath.Join(CacheDirectory(), "git")
var logFile *os.File

// get git tags from git repo
func Tags() []string {
	return getVersionTags()
}

// Download git repository and clone the given version
func CacheDownloadGit(ver string, f ...bool) error {
	logFile = openGitLogs()
	availableVersions := getVersionTags()
	if availableVersions == nil {
		log.Fatal("Fatal error, exiting...")
	}
	ver = normalizeVersion(ver)

	index := lookupVersion(ver, availableVersions)
	if index == -1 {
		return fmt.Errorf("%s doesn't seems to be a valid Go release\n", ver)
	}
	if err := cloneSource(); err != nil {
		return err
	}

	force := false
	if len(f) != 0 && f[0] {
		force = true
	}
	if exists, err := SourceExists(ver); !force && err != nil {
		log.Fatal(err)
	} else if !exists || force {
		if err := copySource(ver); err != nil {
			return err
		}
	}
	return nil
}

func normalizeVersion(ver string) string {
	if !strings.HasPrefix(ver, "go") {
		if strings.HasPrefix(ver, "1") {
			return fmt.Sprintf("go%s", ver)
		}
		if !strings.HasPrefix(ver, "release") {
			if strings.HasPrefix(ver, "5") || strings.HasPrefix(ver, "6") {
				return fmt.Sprintf("release.r%s", ver)
			}
		}
	}

	return ver
}

func checkSource(tag string) error {
	fmt.Fprintf(Output, "Checking %s... ", tag)
	out, err := exec.Command("hg", "pull", "-R", TARGET).CombinedOutput()
	if err != nil {
		fmt.Fprintln(Output, utils.Fail("✖"))
		return err
	}
	fmt.Fprintln(Output, utils.Ok("✔"))
	logOutput(out)
	return nil
}

func cloneSource() error {
	// check if git command line is installed
	if _, err := exec.LookPath("git"); err != nil {
		log.Fatal("Git is not installed on this system.")
	}

	if GitExists() {
		return pull()
	}
	fmt.Fprint(Output, "Cloning Go sources from Github... ")

	out, err := exec.Command("git", "clone", REPO, TARGET).CombinedOutput()
	if err != nil {
		fmt.Println(Output, utils.Fail("✖"))
		return err
	}
	fmt.Fprintln(Output, utils.Ok("✔"))
	logOutput(out)
	return nil
}

func copySource(ver string) error {
	var out []byte
	fmt.Fprint(Output, "Copying source... ")
	destination := filepath.Join(CacheDirectory(), ver)
	os.RemoveAll(destination)
	curr, err := os.Getwd()
	if err != nil {
		return err
	}
	defer func() {
		exec.Command("git", "checkout", "master").Run()
		os.Chdir(curr)
	}()
	os.Chdir(TARGET)
	if ver != "go" && ver != "tip" {
		out, err := exec.Command("git", "checkout", ver).CombinedOutput()
		log.Println(string(out), err)
		if err != nil {
			fmt.Fprintln(Output, utils.Fail("✖"))
			return fmt.Errorf("%s", out)
		}
	}
	out, err = exec.Command("cp", "-R", "../git", destination).CombinedOutput()
	if err != nil {
		fmt.Fprintln(Output, utils.Fail("✖"))
		return err
	}
	fmt.Fprintln(Output, utils.Ok("✔"))
	logOutput(out)
	return nil
}

func pull() error {
	curr, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(curr)
	os.Chdir(TARGET)
	fmt.Fprintf(Output, "Pulling Go sources from Github... ")
	out, err := exec.Command("git", "pull").CombinedOutput()
	if err != nil {
		fmt.Println(Output, utils.Fail("✖"))
		return err
	}
	fmt.Fprintln(Output, utils.Ok("✔"))
	logOutput(out)
	return nil
}

func lookupVersion(ver string, availableVersions []string) (index int) {
	if ver == "go" || ver == "tip" {
		return 0xBEDEAD
	}

	if !strings.HasPrefix(ver, "go") && !strings.HasPrefix(ver, "release") {
		return -1
	}

	for i, v := range availableVersions {
		if v == ver {
			return i
		}
	}
	return -1
}

func getVersionTagsFromGitRepo() ([]string, error) {
	tags := []string{}
	curr, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	os.Chdir(TARGET)
	defer os.Chdir(curr)

	out, err := exec.Command("git", "tag").Output()
	if err != nil {
		return nil, err
	}
	for _, line := range strings.Split(string(out), "\n") {
		if line != "" && !strings.Contains(line, "weekly") {
			tags = append(tags, strings.TrimRight(line, "\n"))
		}
	}
	return tags, nil
}

func getVersionTags() (tags []string) {
	tags = append(tags, "go")
	if err := cloneSource(); err != nil {
		log.Fatal(err)
	}

	newTags, err := getVersionTagsFromGitRepo()
	if err != nil {
		log.Fatal(err)
	}
	tags = append(tags, newTags...)
	sort.Strings(tags)
	return tags
}

func openGitLogs() *os.File {
	logsDir := filepath.Join(CacheDirectory(), "logs")
	openFlags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(filepath.Join(logsDir, "git-go.log"), openFlags, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(logsDir, 0755)
			file, err = os.OpenFile(filepath.Join(logsDir, "git-go.log"), openFlags, 0644)
			if err != nil {
				fmt.Fprintf(Output, "error: can't open log file to write: %s\n", err)
				fmt.Fprintf(Output, "this is a non fatal error, ignoring...")
				return nil
			}
		}
	}
	return file
}

func logOutput(out []byte) {
	logFile.Write(out)
}
