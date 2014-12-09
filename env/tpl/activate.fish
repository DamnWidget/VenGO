# This file can't be executed directly, it has to be
# loaded with 'vengo_activate <ebvironment_name>' from fish

# This script is inspired by virtualenv for Python written by
# Jannis Leidel, Carl Meyer and Brian Rosner

function deactivate --description "Deactivate a VenGO active environment"
    if not set -q VENGO_ENV
        return 0
    end
    # reset environment variables
    set -x PATH $PATH[3..(count $PATH)]
    if test -n "$_VENGO_PREV_PATH"
        set -g PATH "$_VENGO_PREV_PATH"
        set -e _VENGO_PREV_PATH
    end
    if test -n "$_VENGO_PREV_GOROOT"
        set -g GOROOT "$_VENGO_PREV_GOROOT"
        set -e _VENGO_PREV_GOROOT
    end
    if test -n "$_VENGO_PREV_GOTOOLDIR"
        set -g GOTOOLDIR "$_VENGO_PREV_GOTOOLDIR"
        set -e _VENGO_PREV_GOTOOLDIR
    end
    if test -n "$_VENGO_PREV_GOPATH"
        set -g GOPATH "$_VENGO_PREV_GOPATH"
        set -e _VENGO_PREV_GOPATH
    end

    set -e VENGO_ENV
    functions -e deactivate
end

# set paths
set -g VENGO_ENV "{{ .VenGO_PATH }}"
# unset and backup old configuration
set -g _VENGO_PREV_GOROOT (go env GOROOT)
set -e GOROOT

set -g _VENGO_PREV_GOTOOLDIR (go env GOTOOLDIR)
set -e GOTOOLDIR

set -g _VENGO_PREV_GOPATH (go env GOPATH)
set -e GOPATH

# set new environment variables
set -g GOROOT "{{ .Goroot }}"
set -g GOTOOLDIR "{{ .Gotooldir }}"
set -g GOPATH "{{ .Gopath }}"

# set the PATH
set -g PATH "$GOROOT/bin" "$GOPATH/bin" $PATH
