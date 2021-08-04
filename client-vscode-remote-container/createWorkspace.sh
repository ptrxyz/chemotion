#!/bin/bash

dest=$1

[[ -z "${dest}" ]] && {
	echo "No target directory specified."
	exit 1
}

dest=$(realpath ${dest})

echo "Chosen target directory is: ${dest}"
echo "Please make sure the folder ist empty!"
while : ; do
	read -p "Continue (Y/n)? " choice
	case "$choice" in
		[yY][eE][sS]|[yY] ) break;;
		[nN][oO]|[nN] ) echo "See you around then..."; exit 0;;
		* ) echo "Invalid choice.";;
	esac
done

procs=$(nproc 2>/dev/null || echo "4")
echo "Info: ${procs} cores detected."


echo "Creating workspace in ${dest}..."

# clone the repo
mkdir -p ${dest}/
git clone https://github.com/ptrxyz/chemotion_ELN.git ${dest}/ || exit 1

# copy over config files for VSCode
cp -r .devcontainer ${dest}/
cp Dockerfile.vscode ${dest}/
cp dbinit.sh ${dest}/
cp docker-compose.vscode ${dest}/

# adjust number of CPUs used
sed -i 's/BUNDLE_JOBS=.*/BUNDLE_JOBS='${procs}'/g' ${dest}/docker-compose.vscode

echo "done."
echo "Please open the folder [${dest}] in VSCode and confirm"
echo "to change to the container environment when prompted."
echo ""
echo "Please Note: this requires the non-OSS build of VSCode "
echo "             with the Remote Development Extension Pack"
echo "             installed and working."
echo ""

