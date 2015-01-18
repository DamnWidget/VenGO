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
	"os"
	"path/filepath"
)

func MercurialExists() bool {
	_, err := os.Stat(filepath.Join(CacheDirectory(), "mercurial"))
	if err != nil {
		return !os.IsNotExist(err)
	}

	return true
}

func GitExists() bool {
	_, err := os.Stat(filepath.Join(CacheDirectory(), "git"))
	if err != nil {
		return !os.IsNotExist(err)
	}

	return true
}

func SourceExists(ver string) (bool, error) {
	_, err := os.Stat(filepath.Join(CacheDirectory(), ver))
	if err != nil {
		if !os.IsNotExist(err) {
			return false, err
		}
	} else {
		return true, nil
	}

	return false, nil
}
