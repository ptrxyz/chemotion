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

# adjust config files for chemotion
cp ${dest}/public/welcome-message-sample.md ${dest}/welcome-message.md
cp ${dest}/config/datacollectors.yml.example ${dest}/config/datacollectors.yml
cp ${dest}/config/storage.yml.example ${dest}/config/storage.yml
cp ${dest}/config/database.yml.example ${dest}/config/database.yml
sed -i 's/host: .*/host: db/g' ${dest}/config/database.yml

# Install VSCode extension
command -v code &>/dev/null && {
	code --install-extension ms-vscode-remote.vscode-remote-extensionpack
} || {
	echo "VSCode not detected. Please make sure the Remote Development Extension Pack is installed."
	echo "Get it here: https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.vscode-remote-extensionpack"
}

echo "done."
echo ""
command -v code &>/dev/null && {
	echo "Opening folder [${dest}] in VSCode for you. Please confirm"
	echo "to change to the container environment when prompted."
	echo ""
	code ${dest}
} || {
	echo "Please open the folder [${dest}] in VSCode and confirm"
	echo "to change to the container environment when prompted."
	echo ""
	echo "Please Note: this requires the non-OSS build of VSCode "
	echo "             with the Remote Development Extension Pack"
	echo "             installed and working."
	echo ""
}
