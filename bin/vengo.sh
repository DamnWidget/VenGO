
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

if [ "$VENGO_HOME" = "" ]; then
    export VENGO_HOME="$HOME/.VenGO"
    . $VENGO_HOME/bin/includes/output
    . $VENGO_HOME/bin/includes/utils
    . $VENGO_HOME/bin/includes/env
    . $VENGO_HOME/bin/includes/help
fi

alias vengo=${VENGO_HOME}/bin/vengo.sh

# VenGO activate script
function vengo_activate {
    if [ -n "$1" ]; then
        environment="$1"
        if [ "$environment" = "" ]; then
            vengo lsenvs
            return 1
        fi
        if [ "$environment" = "-h" ] || [ "$environment" = "--help" ]; then
            vengo_activate_help
            return 1
        fi
        shift
        for i in "$@"; do
            case $i in
                --pre-activate=*)
                    pre_activate_script=$(echo "$1" | sed 's/[-a-zA-Z0-9]*=//')
                ;;
                --post-activate=*)
                    post_activate_script=$(echo "$i" | sed 's/[-a-zA-Z0-9]*=//')
                ;;
                -h|--help)
                    vengo_activate_help
                    return 0
                ;;
                *)
                    echo "Invalid option $i"
                    vengo_activate_help
                    return 65
                ;;
            esac
        done

        check_environment_exixtance $environment || return 1
        activate="$VENGO_HOME/$environment/bin/activate"
        if [ ! -f "$activate" ]; then
            echo "VenGO: Environment '$VENGO_HOME/$environment' doesn't contains an activate script." >&2
            echo "`Ok`suggestion`Reset`: check the integrity of the environments with 'vengo lsenvs'" >&2
            return 1
        fi

        # call deactivate if we are currently into a virtual environment
        type deactivate >/dev/null 2>&1
        if [ $? -eq 0 ]; then
            deactivate
            unset -f deactivate >/dev/null 2>&1
        fi

        if [ -z ${pre_activate_script+x} ] && [ "$pre_activate_script" != "" ]; then
            pre_activate_script "$environment"
        fi

        source "$activate"

        if [ -z ${post_activate_script+x} ] && [ "$post_activate_script" != "" ]; then
            post_activate_script "$environment"
        fi

        return 0
    else
        echo -e "Usage: vengo_actiate env_name"
    fi
    return 1
}
