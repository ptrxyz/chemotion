#!/bin/bash
set -e

WORKDIR=$(pwd)
BASEDIR=$(realpath $(dirname $0))

# Sanity check as we do some "rm -rf" later, we want to be sure that this is not messed up.
[[ -z "${WORKDIR}" || "${WORKDIR}" == "/" ]] && (echo "Working directory not set to a proper value. Will not continue."; exit 1)

[[ -z "$1" ]] && LOGFILE="build.log" || LOGFILE=$1
[[ -z "$2" ]] && TARGET="<default>"  || TARGET=$2

[[ -d "${BASEDIR}/.git" ]] && BLDGITDIR=${BASEDIR}/.git
[[ -d "${BASEDIR}/../.git" ]] && BLDGITDIR=$(realpath ${BASEDIR}/../.git)

REPO="$WORKDIR/src"
GIT="git -c advice.detachedHead=false --git-dir $REPO/.git "
BLDGIT="git -c advice.detachedHead=false --git-dir ${BLDGITDIR}/.git "
YMLPARSE="python3 ${BASEDIR}/scripts/parseYML.py"
LOGFILE=$(realpath $LOGFILE)
LOG="tee -a ${LOGFILE}"

echo "Logfile: $LOGFILE"
echo "Target: $TARGET"

echo "---- Preparation script ----"
# TODO some kind of version file output, e.g. in .version -> needs to be parsed with the CLI "info" option
# TODO tar ball the src folder as chemotion/app and move it to the container via ADD
# TODO Do we need to move the template folder as well or can it stay outside of the container?

rm -rf ${REPO}
rm -rf ${WORKDIR}/shared/

${BASEDIR}/createSetupScript.sh > ${WORKDIR}/setup.sh && chmod +x ${WORKDIR}/setup.sh
${BASEDIR}/setup.sh

[[ -n "$TARGET" && "$TARGET" != "<default>" ]] && BRANCH="--branch $TARGET"
$GIT clone ${BRANCH} https://github.com/ComPlat/chemotion_ELN $REPO

ELNREF=$($GIT rev-parse --short HEAD)
ELNTAG=$($GIT describe --abbrev=0 --tags)""
BLDREF=$(git --git-dir=${BLDGITDIR} rev-parse --short HEAD 2>/dev/null || echo "D0.1")
BLDTAG=$(git --git-dir=${BLDGITDIR} describe --abbrev=0 --tags 2>/dev/null || echo "<no tag>")
echo "Versions: " | $LOG
echo -e "CHEMOTION_REF=${ELNREF}\nCHEMOTION_TAG=${ELNTAG}\nBUILDSYSTEM_REF=${BLDREF}\nBUILDSYSTEM_TAG=${BLDTAG}" | tee $REPO/.version | sed 's/^/  /g' | $LOG

[[ -d ${WORKDIR}/fixes ]] && (
    cd $REPO; $GIT apply ${WORKDIR}/fixes/*.patch
)

for foldername in $($YMLPARSE read --collect ${BASEDIR}/configFileStructure.yml folders.item); do
    echo "Exposing [${foldername}] ..."; \
    rm -rf $REPO/${foldername}; \
    ln -s /shared/eln/${foldername} $REPO/${foldername}; \
done

for filename in $($YMLPARSE read --collect ${BASEDIR}/configFileStructure.yml files.item); do
    echo "Exposing [${filename}] ..."; \
    rm -f $REPO/${filename}; \
    ln -s /shared/eln/${filename} $REPO/${filename}; \
done
