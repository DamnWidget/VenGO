#!/bin/sh

# Copyright (C) 2014  Oscar Campos <oscar.campos@member.fsf.org>

# This program is free software; you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation; either version 2 of the License, or
# (at your option) any later version.

# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.

# You should have received a copy of the GNU General Public License along
# with this program; if not, write to the Free Software Foundation, Inc.,
# 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

# See LICENSE file for more details.

# tools
GO=`which go`
GIT=`which git`

# colors
OK="\033[32m"
FAIL="\033[31m"
RESET="\033[0m"

# variables
CURRDIR=`pwd`
WORKDIR='.vengo_installation'
DESTDIR=$HOME/.VenGO
REPOSITORY='github.com/DamnWidget/VenGO'


[ "$GIT" = "" ] && {
    echo "Git can't be found in your system"
    echo -n "  ${OK}suggestion${RESET}: run '"
    [ x$(uname) = "xLinux" ] && {
        [ x$(which apt-get) != "x" ] && {
            echo "apt-get install git' to install it"
        } || {
            echo "yum install git' to install it"
        }
    } || {
        echo "'brew install git'"
    }
    echo ""
    exit 1
}

echo -n "Getting sources... "
$GIT clone https://$REPOSITORY $WORKDIR 2> /dev/null
echo "${OK}✔${RESET}"

echo -n "Installing tools..."
for tool in "install" "uninstall" "list" "lsenvs" "mkenv" "rmenv"; do
    $GO get $REPOSITORY/applications/$tool
    mv $GOPATH/bin/$tool $WORKDIR/bin/
done
echo "${OK}✔${RESET}"

echo -n "Installing binaries into $DESTDIR..."
mkdir -p $DESTDIR/scripts
mv $WORKDIR/bin $DESTDIR/
mv $WORKDIR/env/tpl $DESTDIR/scripts/
mv $WORKDIR/VERSION $DESTDIR/
rm -Rf $WORKDIR
echo "${OK}✔${RESET}"

echo ""
echo "${OK}VenGO is now installed in your system{RESET}"
echo "add 'source ${HOME}/.VenGO/bin/vengo' to your .bashrc or .profile to activate it"
echo "you can also do '. ${HOME}/.VenGO/bin/vengo' to start using it right now"
