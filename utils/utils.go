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

package utils

import (
	"bufio"
	"fmt"
	"os/exec"
)

// adds the \x1b[32m prefix and the \x1b[0m suffix to the given string
func Ok(buf string) string {
	return fmt.Sprintf("\x1b[32m%s\x1b[0m", buf)
}

// adds the \x1b[31m prefix and the \x1b[0m suffix to the given string
func Fail(buf string) string {
	return fmt.Sprintf("\x1b[31m%s\x1b[0m", buf)
}

// read logs from external commands
func LogOperation(readers []*bufio.Reader) {
	for _, r := range readers {
		go func(rd *bufio.Reader) {
			for {
				str, err := rd.ReadString('\n')
				if err != nil {
					break
				}
				fmt.Printf("%s", str)
			}
		}(r)
	}
}

// execute an external command and update the output as it's being written
func Exec(verbose bool, args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	outErr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	rd, erd := bufio.NewReader(out), bufio.NewReader(outErr)
	if err := cmd.Start(); err != nil {
		return err
	}
	if verbose {
		LogOperation([]*bufio.Reader{rd, erd})
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
