function vengo_activate --description 'Activate a VenGO virtual environment'
        if not set -q VENGO_HOME
		set -g VENGO_HOME "$HOME/.VenGO"
		set -x VENGO_HOME $VENGO_HOME
	end
        if not set -q __OK
                set -g __OK (set_color -o green)
        end
        if not set -q __NORMAL
                set -g __NORMAL (set_color normal)
        end

        if count $argv >/dev/null
                set environment "$VENGO_HOME/$argv[1]/bin/activate.fish"
                if test -e "$environment"
                        . $environment
                        return 0
                else
                        echo "VenGO: Environment '$VENGO_HOME/$environment' doesn't  iexists." >&2
            echo -n -s "  " "$__OK" "suggestion" "$__NORMAL" ": check the integrity of the environments with 'vengo lsenvs'" >&2
                end
        else
                vengo lsenvs
                return 2
        end
end

