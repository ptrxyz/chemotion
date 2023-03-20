#!/bin/bash
set -euo pipefail

VERSION=dev

BUILDRUN="$(date +%s)"
BUILDRUN="dev"

# OPTS="--force-rm --no-cache"
OPTS=""

export ARGS_VERSION=${VERSION}
export ARGS_BUILDRUN=${BUILDRUN}
export ARGS_BUILD_TAG_KETCHERSVC="v1.0.0"
export ARGS_BUILD_TAG_CONVERTER="v0.6.0"
export ARGS_BUILD_TAG_SPECTRA="0.11.0"
export ARGS_BUILD_TAG_CHEMOTION="490b86403532106e052710b71de011dbb8e0f9a7"

STATUS_FILE="$(mktemp)"
export STATUS_FILE

trap 'rm -f "${STATUS_FILE}"' EXIT

function listContainers() {
    docker image ls | grep chemotion-build
}

function tagContainer() {
    for i in base converter eln ketchersvc ketchersvc-sc spectra msconvert; do
        docker tag chemotion-build/$i:1.4.1 ptrxyz/internal:$i-1.4.1
    done
}

function buildContainer() {
    # $1: filename of dockerfile relative to this script
    # $2: basename of image
    # $3: optional OPTS
    df=$1
    basename=$2
    myOPTS=${3-$OPTS}

    tag="chemotion-build/${basename}:${VERSION}"
    tag2="chemotion-build/${basename}:latest"
    subdir=$(dirname "$df")
    filename=$(basename "$df")

    IFS=" " read -r -a opts <<<"$myOPTS"
    mapfile -t args < <(set | grep ^ARGS_ | sed 's/^ARGS_//g')

    bargs=()
    for arg in "${args[@]}"; do
        bargs+=("--build-arg")
        bargs+=("$arg")
    done

    if [[ -n "$myOPTS" ]]; then
        echo "Building [$df] as [$tag] in [$subdir] with [$myOPTS]."
    else
        echo "Building [$df] as [$tag] in [$subdir]."
    fi

    if (
        cd "$subdir" || exit 1
        # echo "Build args:"
        # echo "${bargs[@]}"

        systemd-inhibit buildah bud --format docker --layers "${opts[@]}" -f "$filename" \
            "${bargs[@]}" \
            -t "$tag" -t "$tag2" . 2>&1 | tee "build-${basename}.log" >/dev/null
    ); then
        echo "+ Successfully built [$df] as [$tag]." | tee -a "${STATUS_FILE}"
    else
        echo "- Failed to build [$df]. See [${subdir}/build-${basename}.log] for details." | tee -a "${STATUS_FILE}"
    fi
}

export DOCKER_BUILDKIT=0

START_TIME=$(date +%s)
buildContainer "base/build.dockerfile" "base" # "--force-rm --no-cache"
buildContainer "converter/build.dockerfile" "converter" &
buildContainer "ketchersvc/build.dockerfile" "ketchersvc" &
buildContainer "ketchersvc/build2.dockerfile" "ketchersvc-sc" &
buildContainer "spectra/build.dockerfile" "spectra" &
buildContainer "spectra/build2.dockerfile" "msconvert" &
buildContainer "eln/build.dockerfile" "eln" &

while [[ -n "$(jobs -rp)" ]]; do
    sleep 1
    # echo "looping: $(jobs -rp)"
done
END_TIME=$(date +%s)

echo -e "\n\nSummary:"
cat "${STATUS_FILE}"
echo -e "\n--\nBuild took $((END_TIME - START_TIME)) seconds."

listContainers
