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

package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var logFile *os.File
var stdLogFile *log.Logger

func init() {
	log.SetOutput(os.Stdout)
	logFile, err := os.Create(filepath.Join(os.TempDir(), "VenGO.log"))
	if err != nil {
		log.Fatal(err)
	}
	stdLogFile = log.New(logFile, "", log.LstdFlags)
}

// Printf calls Logger.Printf to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
	go stdLogFile.Printf(format, v...)
}

// Print calls Logger.Print to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	log.Print(v...)
	go stdLogFile.Print(v...)
}

// Println calls Logger.Println to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) {
	log.Println(v...)
	go stdLogFile.Println(v...)
}

// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	log.Fatal(v...)
	stdLogFile.Output(2, fmt.Sprint(v...))
	logFile.Close()
	os.Exit(1)
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
	s := fmt.Sprintf(format, v...)
	stdLogFile.Output(2, s)
	logFile.Close()
	os.Exit(1)
}

// Fatalln is equivalent to l.Println() followed by a call to os.Exit(1).
func Fatalln(v ...interface{}) {
	log.Fatalln(v...)
	s := fmt.Sprintln(v...)
	stdLogFile.Output(2, s)
	logFile.Close()
	os.Exit(1)
}

// Panic is equivalent to l.Print() followed by a call to panic().
func Panic(v ...interface{}) {
	log.Panic(v...)
	s := fmt.Sprint(v...)
	stdLogFile.Output(2, s)
	panic(s)
}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func Panicf(format string, v ...interface{}) {
	log.Panicf(format, v...)
	s := fmt.Sprintf(format, v...)
	stdLogFile.Output(2, s)
	panic(s)
}

// Panicln is equivalent to l.Println() followed by a call to panic().
func Panicln(v ...interface{}) {
	log.Panicln(v...)
	s := fmt.Sprintln(v...)
	stdLogFile.Output(2, s)
	panic(s)
}
