#!/bin/bash

CACHE="--no-cache"

build_ruby() {
	docker build $CACHE -t chemotion-build:latest-ruby -f dckr-ruby .
	return $?
}

build_node() {
	docker build $CACHE -t chemotion-build:latest-node -f dckr-node .
	return $?
}

build_eln() {
	docker build $CACHE -t chemotion-build:latest-eln -f Dockerfile .
	return $?
}

function parseFlags() {
	while [ -n "$1" ]; do 
		case "$1" in
			--cache)
				CACHE=""				
				;;
			--no-cache)
				CACHE="--no-cache"
				;;
		esac
		shift
	done
}

# parse flags before executing commands
parseFlags $@

echo "Cache is turned: "$([[ -n $CACHE ]] && echo "off" || echo "on")

while [ -n "$1" ]; do 
	case "$1" in
		ruby)
			build_ruby
			;;
		node)
			build_node
			;;
		eln)
			build_eln
			;;
		all)
			build_ruby && build_node && build_eln && echo "Success."
			;;
		--*)
			# ignore all flags
			;;
		*)
			echo "Ignoring: $1"
	esac
	shift
done


