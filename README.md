# VenGO

[![Build Status](https://travis-ci.org/DamnWidget/VenGO.png)](https://travis-ci.org/DamnWidget/VenGO)

Create and manage Isolated Virtual Environments for Golang.

## Motivation

Why a tool to generate and manage virtual environments in Go?. Well, sometimes programmers needs to work in or
maintain a project that requires a specific version of Go or use specific versions of 3rd party libraries that
maybe depend themselves on some specific Go version.

There are already tools like `godep` to freeze dependencies and make the programmer able to build a package in
consistent way reproducing the exact package ecosystem that was used when it was developed and Go versions managers
like `gvm` that helps the programmer to install and use different Go versions. But there is no a tool that can do
both and in an easy and familiar way.

VenGO is able to install as many Go versions from as many sources that programmers want and to create as many isolated
environments as they need using one or more Go versions.

## Platforms and Support

VenGO works and is actively maintained in POSIX platforms, it requires go1.2 or higher to be compiled

Platform | Status | Maintainer
-------- | ------ | ----------
GNU/Linux | Stable | [@damnwidget](https://github.com/DamnWidget)
FreeBSD | Stable |
OS X | Stable | [@damnwidget](http://github.com/DamnWidget)
Windows | Garbage |

**note**: Support for Windows is planned

## Installation

VenGO can be installed following three simple steps

### 1 Get dependencies and code

Install the VenGO dependencies

```
$ go get github.com/mcuadros/go-version github.com/ogier/pflag
```

Now get VenGO code itself into the `GOPATH`

```
$ go get -d github.com/DamnWidget/VenGO
```

### 2 Compile and install

VenGO comes with a `Makefile` for your convenience (so the `make` command line tool has to be installed)

```
$ cd $GOPATH/src/github.com/DamnWidget/VenGO && make install
```

### 3 Enable the vengo application in your shell

Finally the command below will enable the `vengo` command in your system

```
$ source $HOME/.VenGO/bin/vengo
```

#### 4 Optional

If you want to enable `vengo` in permanent basis in your system, you can add it to your .bashrc, .zshrc or .profile
files like

```
echo "source $HOME/.VenGO/bin/vengo" >> $HOME/.bashrc
```

## Usage

VenGO is quite similar to Python's virtualenvwrapper tool, if you execute just `vengo` with no arguments you will get
a list of available commands. The most basic usage is install a Go version

**note**: VenGO is not able to use Go installations that has not been made with VenGO itself

The following command will install Go 1.2.2 from the mercurial repository:

```
$ vengo install go1.2.2
```

This install the go1.2.2 version into the VenGO's cache and generates a manifest that guarantee the installation
integrity, now the programmer can create a new environment using the just installed Go version

```
$ vengo mkenv -g go1.2.2 MyEnv
```

This will create a new isolated environment that uses go1.2.2 and uses `$VENGO_HOME/MyEnv` as `GOPATH`

To activate this new environment thw programmer just have to use `vengo activate` with the name of the recently created
environment

```
$ vengo activate MyEnv
```

Now, whatever is installed using `go get` will be installed in the new isolated virtual go environment. It's `GOPATH` bin
will be already added to the programmer `PATH` so new applications should be available in the command line after installation.

To stop using the active environment just execute

```
$ deactivate
```

## Detailed guide on VenGO commands

VenGO comes with eight commands:

    * activate           Activate a Virtual Go Environment
    * install            Installs a new Go version
    * uninstall          Uninstall an installed Go version
    * list               List installed and available Go versions
    * lsenvs             List available Virtual Go Environments
    * mkenv              Create a new Virtual Go Environment
    * rmenv              Remove a Virtual Go Environment
    * vengo-uninstall    Uninstall VenGO and remove all the Virtual Go Environments
