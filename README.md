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

VenGO comes with eight different commands that will be used trough the vengo command line application
![VenGO no arguments](https://raw.githubusercontent.com/DamnWidget/VenGO/images/vengo.png)

### VenGO install

Vengo install is used to install new versions of Go, it can install them directly from the official mercurial repository, from a `tar.gz` packed source or directly in binary format in case that the user doesn't want to compile it.
![VenGO Install](https://raw.githubusercontent.com/DamnWidget/VenGO/images/install00.png)

Install will download the sources from the official mercurial repository by default, then check and copy the specific version into a directory (in the VenGO cache directory) named as the version itself, compile it and generate a `manifest` for the installation.  To download from a packaged `tar.gz` source use the `-s` or `--source` flag like in:
```
$ vengo install -s 1.3.3
```

In similar way, to download from a binary source use the `-b` or `--binary` flag like in:
```
$ vengo install -b 1.3.3
```

### VenGO list

Vengo list is used to show a list of installed Go versions, available Go versions or both. If the list command detects that a installed Go version integrity is compromised, it will display a red ✖ mark, a green ✔ mark if not
![VenGO List](https://raw.githubusercontent.com/DamnWidget/VenGO/images/list00.png)

If the `-n` or `--non-installed` flag is passed to the list command, a complete list of available sources is returned back to the user ordered by binary, mercurial and `tar.gz` packed versions.

#### How do I know from which source is each version?

Versions that are **prefixed** like `1.2.2.<platform>-<arch>` are binaries, note that is not neccesary to add the platform and architecture to the install command to donwload the version so for example if the list command return to us that the version `1.3.3.darwin-amd64-osx10.8` is available, we will write just:
```
$ vengo install --binary 1.3.3
```
The install command is smart enough to know that we are using a 64bits OS X and it's version, it will work in the exact same way on GNU/Linux and Windows

**note**: Windows support is not complete yet

Versions **prefixed** with `go` or `release` like `go1.1` or `release.r56` come from the official mercurial repository, the install command doesn't need any special flag to use it as it's the default download option, note that is not needed to add the `go` prefix neither but is a good practice to use it just to avoid confusion.

Finally, all the versions that doesn't have any prefix or suffix are `tar.gz` packaged versions of the source, just pass the `--source` flag to the install command in other to download them.

### VenGO uninstall

Vengo uninstall is used to uninstall a Go installed version, it doesn't remove any Virtual Go Environment that has been created using the deleted version but it will be shown by the `lsenvs` command as integrity compromised.

### VenGO mkenv

Vengo mkenv is used to create new Isolated Virtual Go Environments, the Go version to use must can be specified as argument for the parameter `-g` or `--go`, if no version is pased, `tip` is tried to be used automatically.
![VenGO Install](https://raw.githubusercontent.com/DamnWidget/VenGO/images/mkenv.png)

Vengo mkenv will use the name of the environment as prefix of the terminal prompt when the user switch to an environment using `vengo activate` but the users can specify whatever other prompt that they like passing a string to the parameter `-p` or `--prompt` so for exmample:

```
$ vengo mkenv -p "(VenGO [go1.4rc1])" -g go1.4rc1 vengo_go14rc1
```

Will give you a prompt like this one when you switch to it, `(VenGO [go1.4rc1]) damnwidget@iMacStation ~ $`

You can also force the environment reinstallation passing the flag `-f` or `--force` in case that the environment already exists

### VenGO lsenvs

Vengo lsenvs is used to list Isolated Virtual Go Environments in your system. Integrity compromised environments will be shown with a red ✖ mark, a green ✔ mark will be shown otherwise
![VenGO Lsenvs](https://raw.githubusercontent.com/DamnWidget/VenGO/images/lsenvs00.png)

### VenGO rmenv

Vengo rmenv is used to delete Virtual Go Environments, delete an environment doesn't affect the Go version used to install the environment or other environments using that Go version

### VenGO vengo-uninstall

Vengo vengo-uninstall will delete all the environments, Go versions and VenGO installation itself.
