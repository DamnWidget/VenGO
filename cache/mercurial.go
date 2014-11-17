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
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/DamnWidget/VenGO/logger"
)

const REPO = "https://go.googlecode.com/hg/"

var TARGET = filepath.Join(CacheDirectory(), "mercurial")

var logFile *os.File

func CacheDonwloadMercurial(ver string, forceDownload ...bool) error {
	logFile = openMercurialLogs()
	availableVersions := getVersionTags()
	if availableVersions == nil {
		logger.Fatal("Fatal error, exiting...")
	}
	ver = normalizeVersion(ver)

	index := lookupVersion(ver, availableVersions)
	if index == -1 {
		return fmt.Errorf("%s doesn't seems to be a valid Go release\n", ver)
	}
	if err := cloneSource(); err != nil {
		return err
	}
	if err := checkSource(availableVersions[index]); err != nil {
		return err
	}

	force := false
	if len(forceDownload) != 0 && forceDownload[0] {
		force = true
	}
	if exists, err := SourceExists(ver); err != nil {
		logger.Fatal(err)
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
	logger.Printf("Checking %s...\n", tag)
	out, err := exec.Command("hg", "pull", "-R", TARGET).CombinedOutput()
	if err != nil {
		return err
	}
	logOutput(out)
	return nil
}

func cloneSource() error {
	if MercurialExists() {
		return nil
	}

	logger.Println("Downloading Go source from mercurial...")
	// check if mercurial command line is installed
	if _, err := exec.LookPath("hg"); err != nil {
		logger.Fatal("Mercurial is not installed on your machine.")
	}
	out, err := exec.Command("hg", "clone", REPO, TARGET).CombinedOutput()
	if err != nil {
		return err
	}
	logOutput(out)
	return nil
}

func copySource(ver string) error {
	logger.Println("Copying source...")
	destination := filepath.Join(CacheDirectory(), ver)
	os.RemoveAll(destination)
	out, err := exec.Command(
		"hg", "clone", "-u", ver, TARGET, destination).CombinedOutput()
	if err != nil {
		logger.Println(err)
		return err
	}
	logOutput(out)
	return nil
}

func lookupVersion(ver string, availableVersions []string) (index int) {
	if !strings.HasPrefix(ver, "go") && !strings.HasPrefix(ver, "release") {
		return -1
	}

	return sort.SearchStrings(availableVersions, ver)
}

func getVersionTags() (tags []string) {
	resp, err := http.Get("https://go.googlecode.com/hg/.hgtags")
	if err != nil {
		log.Println(err)
		return nil
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 400 {
			logger.Fatal("Cant't get go versions list from Google servers")
			logger.Println(fmt.Errorf("%s", resp.Status))
		}
		return nil
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		logger.Println(err)
		return nil
	}

	// return releases only
	for _, tag := range strings.Split(buf.String(), "\n") {
		if tag == "" {
			continue
		}

		ver := strings.Split(tag, " ")[1]
		if strings.HasPrefix(ver, "release") || strings.HasPrefix(ver, "go1") {
			tags = append(tags, ver)
		}
	}

	// sort tags in increasing order
	sort.Strings(tags)
	return tags
}

func logOutput(out []byte) {
	logFile.Write(out)
}

func openMercurialLogs() *os.File {
	logsDir := filepath.Join(CacheDirectory(), "logs")
	openFlags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(
		filepath.Join(logsDir, "mercurial-go.log"), openFlags, 0644,
	)
	if err != nil {
		logger.Printf("error: can't open log file to write: %s\n", err)
		logger.Println("this is a non fatal error, ignoring...")
		return nil
	}
	return file
}
