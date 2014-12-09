function vengo --description "Genrate and manage isolated virtual Go environments"
	bash -c 'for inc in ~/.VenGO/bin/includes/*; do . $inc; done && . ~/.VenGO/bin/vengo; vengo "$@"' vengo $argv
end
