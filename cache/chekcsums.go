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

import "fmt"

var checksums map[string]string = map[string]string{
	// stable
	"1.2.2": "3ce0ac4db434fc1546fec074841ff40dc48c1167",
	"1.3":   "9f9dfcbcb4fa126b2b66c0830dc733215f2f056e",
	"1.3.1": "bc296c9c305bacfbd7bff9e1b54f6f66ae421e6e",
	"1.3.2": "67d3a692588c259f9fe9dca5b80109e5b99271df",
	"1.3.3": "b54b7deb7b7afe9f5d9a3f5dd830c7dede35393a",

	// unstable
	"1.4beta1": "f2fece0c9f9cdc6e8a85ab56b7f1ffcb57c3e7cd",
	"1.3rc2":   "53a5b75c8bb2399c36ed8fe14f64bd2df34ca4d9",
	"1.3rc1":   "6a9dac2e65c07627fe51899e0031e298560b0097",
}

// check if a given version is supported by VenGO to auto donwload/compile
// if the version is valid, it returns it's SHA1 fingerprint, error is
// returned otherwise
func Checksum(version string) (string, error) {
	if sha1, ok := checksums[version]; ok {
		return sha1, nil
	}
	return "", fmt.Errorf("%s is not a VenGO supported version you must donwload and compile it yourself", version)
}
