#!/bin/bash
# shellcheck disable=SC2154  # irrelevant cause we inherit env of Makefile
# shellcheck disable=SC2295  # irrelevant cause we want exactly this.

# Global defaults
BUILDER="docker build"
# BUILDER="buildah bud --format docker --layers"
DOCKERFILE="build.dockerfile"
DOCKER_BUILDKIT=0
REPO="chemotion-build"
VERBOSE=0

export DOCKER_BUILDKIT
export BUILDER
export DOCKERFILE
export REPO
export VERBOSE

function strip() {
	local var="$1"
	local remove="${2:-[:space:]}"

	# remove leading characters
	var="${var#"${var%%[!${remove}]*}"}"

	# remove trailing characters
	var="${var%"${var##*[!${remove}]}"}"

	printf '%s' "$var"
}

function info() {
	if [[ -t 1 ]]; then
		txt="\033[34m•\033[0m [$1]: $2"
		echo -e "$txt"
	else
		txt="· [$1]: $2"
		echo "$txt"
	fi
}

function error() {
	if [[ -t 1 ]]; then
		txt="\033[31m✘\033[0m [$1]: $2"
		echo -e "$txt"
	else
		txt="- [$1]: $2"
		echo "$txt"
	fi
}

function ok() {
	if [[ -t 1 ]]; then
		txt="\033[32m✔\033[0m [$1]: $2"
		echo -e "$txt"
	else
		txt="+ [$1]: $2"
		echo "+ "
	fi
}

function nothing() {
	if [[ ${VERBOSE:-"false"} == "false" || ${VERBOSE:-0} -eq 0 ]]; then return 0; fi
	if [[ -t 1 ]]; then
		txt="\033[90m  [$1]: $2\033[0m"
		echo -e "$txt"
	else
		txt="  [$1]: $2"
		echo "+ "
	fi
}

function buildContainer() {
	# Task-local variables passed by make:
	#  - localName: name of the current task as named in the Makefile (e.g. ketchersvc, eln, etc)
	#  - localPath: path to the directory containing the Dockerfile
	#  - localImageTag: image tag base name (e.g. ketchersvc, eln, etc.)
	#  - localAdditionalTags: raw tags to add to the image in addition to "REPO/TAG:VERSION" and "REPO/TAG:latest"
	#  - localHash: commit hash or other identifier
	#  - localDockerfile: name of the Dockerfile (rel. to path)

	#  - localOpts: options passed to `docker build`
	#  - localArgs: build-args passed to `docker build`
	#  - localVersion: version tag (e.g. 1.0.0)

	# Global variables that override/merge with local variables:
	#  - OPTS: options passed to `docker build`
	#  - ARGS: build-args passed to `docker build`
	#  - REPO: repository name (e.g. chemotion-build)
	#  - VERSION: version tag (e.g. 1.0.0)
	#  - INHIBITOR: command to inhibit systemd sleep
	#  - VERBOSE: if set to anything except `false` or `0`, more messages (class: nothing) will be logged.

	local taskname=${localName:-"unnamed"}
	local path=${localPath}/
	local imagetag=${localImageTag}
	local additionalTags=${localAdditionalTags}
	local hash=${localHash} # for future use. as of now, the dockerfiles do not use this
	local dockerfile=${localDockerfile:-${DOCKERFILE}}

	path=$(realpath "${path}")
	path=${path}/
	fulldf=$(realpath "${path}/${dockerfile}")

	if [[ -z "${imagetag}" ]]; then
		error "$taskname" "Failed to build. No image tag provided."
		return 1
	fi

	if ! cd "${path}"; then
		error "$taskname" "Failed to build. Could not change directory to [${path}]."
		return 1
	fi

	if [[ ! -r "${fulldf}" || ! -f "${fulldf}" ]]; then
		error "$taskname" "Failed to build. Can not read dockerfile at [${fulldf}]."
		return 1
	fi

	# version: use localVersion if set, otherwise use VERSION, otherwise use "dev"
	myVERSION=${localVersion:-"${VERSION}"}
	myVERSION=${myVERSION:-"dev"}
	myVERSION=$(strip "${myVERSION}")

	# opts: add global opts, then add local opts
	myOPTS="${OPTS} ${localOpts}"
	myOPTS=$(strip "${myOPTS}")

	# args: add mandatory args, then add global args, then add local args
	myARGS="VERSION=${myVERSION};TASK=${taskname};"
	myARGS="${myARGS};${ARGS};${localArgs}"
	myARGS=$(strip "${myARGS}")

	# repo: use REPO if set, otherwise use "chemotion-build"
	myREPO=${REPO:-"chemotion-build"}
	myREPO=$(strip "${myREPO}" "/[:space:]")

	# inhibitor: use INHIBITOR if set, otherwise use "systemd-inhibit".
	# Check if it exists and is executable.
	myINHIBITOR=${INHIBITOR:-"systemd-inhibit --what=idle"}
	if ! which "$myINHIBITOR" >/dev/null; then
		myINHIBITOR=""
	fi

	# build tags
	fulltag1=${myREPO}/${imagetag}:${myVERSION}
	fulltag2=${myREPO}/${imagetag}:latest

	btags=()
	# shellcheck disable=SC2162  # we want to mangle backslashes...
	while read -d';' tag; do
		if [[ -z "$tag" ]]; then continue; fi
		tag=$(strip "$tag")
		btags+=("--tag")
		btags+=("$tag")
	done <<<"$fulltag1;$fulltag2;$additionalTags;"

	# build log
	logfile=build-${imagetag}.log
	fulllogfile=${path}/${logfile}
	fulllogfile=$(realpath "${fulllogfile}")

	# read OPTS into an array
	IFS=" " read -r -a bopts <<<"$myOPTS"

	# read args into an array and prefix them
	bargs=()
	# shellcheck disable=SC2162  # we want to mangle backslashes...
	while read -d';' arg; do
		if [[ -z "$arg" ]]; then continue; fi
		bargs+=("--build-arg")
		bargs+=("$arg")
	done <<<"$myARGS;"

	info "$taskname" "Building [${fulldf}] in [${path}]."
	if [[ -n "$myOPTS" ]]; then
		nothing "$taskname" "Opts: ${bopts[*]}"
	fi
	if [[ -n "$myARGS" ]]; then
		nothing "$taskname" "Args: ${bargs[*]}"
	fi
	nothing "$taskname" "Tags: ${btags[*]}"

	START_TIME=$(date +%s)
	if (
		# shellcheck disable=SC2086  # we want to expand BUILDER
		$myINHIBITOR $BUILDER "${bopts[@]}" -f "$fulldf" \
			"${bargs[@]}" \
			"${btags[@]}" \
			"$path" 1>"$fulllogfile" 2>&1
	); then
		ok "$taskname" "Build sucessful. Logfile: [${fulllogfile}]"
	else
		error "$taskname" "Failed to build. See [${fulllogfile}] for details."
	fi
	END_TIME=$(date +%s)
	nothing "$taskname" "Build took $((END_TIME - START_TIME)) seconds."
}
