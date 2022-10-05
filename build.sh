#!/bin/bash
set -euo pipefail

BUILDRUN="dev"  #$(date +%s)
VERSION=1.4.0
ARGZ=""  #--force-rm --squash
LOS=$(date)

# function cleanup() {
#     local pids
#     mapfile -t pids < <(jobs -rp)
#     if [[ ${#pids[@]} -gt 0 ]]; then
#         kill -9 "${pids[@]}"
#     fi
#     echo $LOS
# }

# # trap cleanup EXIT
# trap cleanup SIGTERM
# trap cleanup SIGINT

getImageIdByServiceID() {
    docker images --filter "label=chemotion.internal.service.id=$1" \
    --format "{{.CreatedAt}}\t{{.ID}}" | sort -nr | head -n 1 | cut -f2
}

(
    echo "Building BASE ..."
    cd base

    docker build $ARGZ -f build.dockerfile \
        --build-arg BASE_VERSION="$VERSION" \
        --build-arg BUILDRUN="$BUILDRUN" \
        --target base \
        -t chemotion-build/base:"$VERSION" . 2>&1 | tee build_base.log

    echo ">> BASE built."

    # IMAGE_ID=$(getImageIdByServiceID "base")
    # [[ -n "$IMAGE_ID" ]] && docker tag "$IMAGE_ID" chemotion-build/base:1.4.0
    # echo -e "\n\nTagged [$IMAGE_ID] as BASE.\n\n"
    # sleep 2
)

(
    echo "Building CONVERTER ..."
    cd converter

    docker build $ARGZ -f build.dockerfile \
        --build-arg CONVERTER_BUILD_TAG="v0.6.0" \
        --build-arg CONVERTER_VERSION="$VERSION" \
        --build-arg BUILDRUN="$BUILDRUN" \
        --target app \
        -t chemotion-build/converter:"$VERSION" . 2>&1 | tee build_converter.log

    echo ">> CONVERTER built."

    # [[ -n "$IMAGE_ID" ]] && docker tag "$IMAGE_ID" chemotion-build/converter:1.4.0
    # IMAGE_ID=$(getImageIdByServiceID "converter")
    # echo -e "\n\nTagged [$IMAGE_ID] as CONVERTER.\n\n"
    # sleep 2
)

(
    echo "Building KETCHERSVC ..."
    cd ketchersvc

    docker build $ARGZ -f build.dockerfile \
        --build-arg KETCHERSVC_BUILD_TAG="ba28832" \
        --build-arg KETCHERSVC_VERSION="$VERSION" \
        --build-arg BUILDRUN="$BUILDRUN" \
        --target app \
        -t chemotion-build/ketchersvc:"$VERSION" . 2>&1 | tee build_ketchersvc.log

    echo ">> KETCHERSVC built."

    # IMAGE_ID=$(getImageIdByServiceID "ketchersvc")
    # [[ -n "$IMAGE_ID" ]] && docker tag "$IMAGE_ID" chemotion-build/ketchersvc:1.4.0
    # echo -e "\n\nTagged [$IMAGE_ID] as KETCHERSVC.\n\n"
    # sleep 2
)

(
    echo "Building SPECTRA ..."
    cd spectra

    docker build $ARGZ -f build.dockerfile \
        --build-arg SPECTRA_BUILD_TAG="0.10.15" \
        --build-arg SPECTRA_VERSION="$VERSION" \
        --build-arg BUILDRUN="$BUILDRUN" \
        --target spectra \
        -t chemotion-build/spectra:"$VERSION" . 2>&1 | tee build_spectra.log

    echo ">> SPECTRA built."

    # IMAGE_ID=$(getImageIdByServiceID "spectra")
    # [[ -n "$IMAGE_ID" ]] && docker tag "$IMAGE_ID" chemotion-build/spectra:1.4.0
    # echo -e "\n\nTagged [$IMAGE_ID] as SPECTRA.\n\n"
    # sleep 2

    # IMAGE_ID=$(getImageIdByServiceID "msconvert")
    # [[ -n "$IMAGE_ID" ]] && docker tag "$IMAGE_ID" chemotion-build/msconvert:1.4.0
    # echo -e "\n\nTagged [$IMAGE_ID] as MSCONVERT.\n\n"
    # sleep 2
)

(
    echo "Building MSCONVERT ..."
    cd spectra

    docker build $ARGZ -f build2.dockerfile \
        --build-arg SPECTRA_VERSION="$VERSION" \
        --build-arg BUILDRUN="$BUILDRUN" \
        --target msconvert \
        -t chemotion-build/msconvert:"$VERSION" . 2>&1 | tee build_msconvert.log

    echo ">> MSCONVERT built."
)

(
    echo "Building ELN ..."
    cd eln2

    docker build $ARGZ -f build.dockerfile \
        --build-arg CHEMOTION_BUILD_TAG=67d363fa018ff448e0c83919c3e07535022e4978 \
        --build-arg CHEMOTION_VERSION=1.4.0 \
        --build-arg BUILDRUN="$BUILDRUN" \
        --target eln \
        -t chemotion-build/eln:"$VERSION" . 2>&1 | tee build_eln.log

    echo ">> ELN built."

    # IMAGE_ID=$(getImageIdByServiceID "eln")
    # [[ -n "$IMAGE_ID" ]] && docker tag "$IMAGE_ID" chemotion-build/eln:1.4.0
    # echo -e "\n\nTagged [$IMAGE_ID] as ELN.\n\n"
    # sleep 2
)

while [[ -n "$(jobs -rp)" ]]; do
    sleep 1
    echo "looping: $(jobs -rp)"
done
date
echo "DONE: $LOS"

