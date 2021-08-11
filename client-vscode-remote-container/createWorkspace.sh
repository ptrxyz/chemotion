#!/bin/bash

function confirm() {
	while : ; do
		read -p "$1 (y/n)? " choice
		case "$choice" in
			[yY][eE][sS]|[yY] ) return 0;;
			[nN][oO]|[nN] ) return 1;;
			* ) echo "Invalid choice.";;
		esac
	done
}

dest=$1

[[ -z "${dest}" ]] && {
	echo "No target directory specified."
	exit 1
}

dest=$(realpath ${dest})

echo "Chosen target directory is: ${dest}"
echo "Please make sure the folder is empty!"
confirm "Confirm" || {
	echo "See you around then..."
	exit 0
}

procs=$(nproc 2>/dev/null || echo "4")
echo "Info: ${procs} cores detected."

echo "Creating workspace in ${dest}..."

# clone the repo
mkdir -p ${dest}/
git clone https://github.com/ptrxyz/chemotion_ELN.git ${dest}/ || exit 1

# copy over config files for VSCode
cp -r .devcontainer ${dest}/

# adjust number of CPUs used
sed -i 's/BUNDLE_JOBS=.*/BUNDLE_JOBS='${procs}'/g' ${dest}/.devcontainer/docker/docker-compose.vscode

# adjust config files for chemotion
cp ${dest}/public/welcome-message-sample.md ${dest}/public/welcome-message.md
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

if [[ -f gems.tar.gz ]]; then
	echo "Some gems take ages to compile (even on potent machines). We can put "
	echo "precompiled versions of gems you probably need into your project container "
	echo "upon first start. This will significantly speed up (re-)building the "
	echo "environment at expanse of additional storage (~order or MB) due to unused "
	echo "gems residing on your system."
	# todo: explain why this is not default...
	confirm "Embed gems and precompiled libraries (recommended)?" && {
		echo "Info: precompiled libraries will be embedded."

		( a=$(pwd); cd ${dest}/.devcontainer && tar xfvz ${a}/gems.tar.gz )
		sed -i '/^BUNDLED WITH$/!b;n;c\ \ \ 2.2.25' ${dest}/Gemfile.lock

		# This was the alternate approach. Probably will be removed in the future.
		# cp rdkit_chem.tar.gz ${dest}/rdkit_chem.tar.gz
		# sed -i "s#^gem 'rdkit_chem'.*#gem 'rdkit_chem', git: 'https://github.com/ptrxyz/rdkit_chem', ref: 'b7532a4bbbb154ed2bb7d49d15a79c26eb2c8086'#g" ${dest}/Gemfile
	} || {
		echo "Info: precompiled libraries will be NOT embedded."
		mkdir -p ${dest}/.devcontainer/gems
		chown 1000:1000 ${dest}/.devcontainer/gems || echo "Warning: could not set permissions for gems volume!"
	}
else
	mkdir -p ${dest}/.devcontainer/gems
	chown 1000:1000 ${dest}/.devcontainer/gems || echo "Warning: could not set permissions for gems volume!"
	echo "Info: precompiled libraries not found. skipping..."
fi

[[ -d .devcontainer/gems ]] && echo ".devcontainer/gems" >> ${dest}/.dockerignore
echo ".devcontainer" >> ${dest}/.gitignore

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
