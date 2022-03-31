#!/bin/bash

UPLOAD_FOLDER=${UPLOAD_FOLDER-/old/eln/uploads}
IMAGE_FOLDER=${IMAGE_FOLDER-/old/eln/public/images}
OLDPATH=${OLDPATH-/old/eln}
NEWPATH=${NEWPATH-/shared/pullin}

echo "INFO:"
echo "  Upload folder  : ${UPLOAD_FOLDER}"
echo "  Image folder   : ${IMAGE_FOLDER}"
echo "  Old config path: ${OLDPATH}"
echo "  New config path: ${NEWPATH}"
echo "--"
echo

[[ -d "${UPLOAD_FOLDER}" ]] && {
    rm -rf /chemotion/data/uploads
    mkdir -p /chemotion/data/uploads
    tar cH posix --directory="${UPLOAD_FOLDER}" . | tar xh --one-top-level=/chemotion/data/uploads
    echo "✔ Uploads copied from [${UPLOAD_FOLDER}] to [/chemotion/data/uploads]."
} || {
    echo "Can not find upload folder at [${UPLOAD_FOLDER}]"
}

[[ -d "${IMAGE_FOLDER}" ]] && {
    rm -rf /chemotion/data/public/images
    mkdir -p /chemotion/data/public/images
    tar cH posix --directory="${IMAGE_FOLDER}" . | tar xh --one-top-level=/chemotion/data/public/images
    echo "✔ Images copied from [${IMAGE_FOLDER}] to [/chemotion/data/public/images]."
} || {
    echo "Can not find image folder at [${IMAGE_FOLDER}]"
}

toImport=("config/database.yml:config/database.yml" ".env:.env" "config/editors.yml:config/editors.yml" "public/editors:public/editors" "public/welcome-message.md:public/welcome-message.md")
for p2p in ${toImport[@]}; do
    p2p=(${p2p/:/ })
    from=$OLDPATH/${p2p[0]}
    to=$NEWPATH/${p2p[1]}
    [[ -e ${from} ]] && {
        mkdir -p $(dirname ${to})
        cp ${from} ${to}
        [[ "${to}" == *.yml ]] && yq -i '.production | {"production": . }' ${to}
        echo "✔ Copied [${from}] to [$to]."
    } || {
        echo "- No source file found for [$from]. This is not a problem, just information."
    }
done