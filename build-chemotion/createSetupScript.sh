#!/bin/bash

BASEDIR=$(realpath $(dirname $0))
YMLPARSE="python3 ${BASEDIR}/scripts/parseYML.py"

echo "#!/bin/bash"
echo "set -e"
for foldername in $($YMLPARSE read --collect ${BASEDIR}/configFileStructure.yml folders.item); do
    echo "mkdir -p shared/eln/${foldername}";
done
echo "mkdir -p shared/eln/config"
echo "mkdir -p db-data"
