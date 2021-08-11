# Container files for Visual Studio Code Remote Editing

This folder contains anything you need to use Visual Studio Code's remote developing features to start developing on the Chemotion ELN inside a Docker development container.

Video Tutorial here: [YouTube](https://www.youtube.com/watch?v=HZCAbC6ldzE)

## Requirements

-   Visual Studio Code (non-OSS version. Remote development is not supported by the open source builds.)
-   Visual Studio Code's [Remote Developing extension pack](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.vscode-remote-extensionpack)
-   working [Docker installation](https://docs.docker.com/get-docker/)

For more detailed information check out [Microsoft's knowledge base](https://code.visualstudio.com/docs/remote/remote-overview).

# General Information

The purpose is to set up a VSCode remote development environment in a Docker container. See the knowledge base article above if all this does not sound familiar to you.

The setup is as follows:

`.devcontainer/devcontainer.json` specifies which Dockerfiles/docker-compose files VSCode should use. In addition it contains a postCreation command, that is run as soon as the container finished building. In our case the postCreation command is:

-   creating a database
-   enabling all necessary extensions in the container (for VSCode) and for the database
-   seeding the database

The Dockerfile that is used to build the development environment is called `Dockerfile.vscode`. This file is it's own documentation.

To start an additional container for a database server, `docker-compose.vscode` is used. Again, the file is pretty self-explainatory.

# Usage

:exclamation: **IMPORTANT**: make sure Docker is running before proceeding with the instructions below.

> As of now, the script `createWorkspace.sh <someFolder>` will create a workspace for you in the folder specified.
> This includes cloning the chemotion source repository and copying all files where they belong.
> This should be an easier alternative to the steps below.
> **Optional:** download a suitable version of the file `gems.tar.gz` from [here](https://gems.ptrxyz.de/). It contains precompiled gems and
> speeds up the whole building process by A LOT! Place it in the same folder as `createWorkspace.sh`.

**Step 1:** First, check out the Chemotion ELN source from the [GitHub Repo](https://github.com/ComPlat/chemotion_ELN). In the following `/host/workspace/chemotion` will be used as the source directory:

```
$ git clone https://github.com/ComPlat/chemotion_ELN /host/workspace/chemotion
```

**Step 2:** Do now create the configuration files for Chemotion ELN. For that purpose, remove the ".example" extension from the following files:

-   config/storage.yml.example
-   config/datacollectors.yml.example
-   config/database.yml.example

Only the `database.yml` file should need manual editing: change the database host for all configurations to `db`, as this is the hostname of the sidecar container.

Additionally rename `public/welcome-message-sample.md` to `public/welcome-message.md`, otherwise the ELN will work, yet some tests might fail.

**Step 3:** place the `.devcontainer` folder into Chemotion's source directory:

```
$ cp -R .devcontainer /host/workspace/chemotion
```

**Step 4:** create an empty folder for gems. This folder will be used as gem cache and speeds up the container building process for subsequent builds:

```
$ mkdir -p /host/workspace/chemotion && chown 1000:1000 /host/workspace/chemotion
```

:exclamation: Attention: make sure the folder is writeable by UID 1000 as shown in the snippet above.

**Step 5:** Open VSCode and open the Chemotion folder: `File` -> `Open Folder` -> select the right folder, here `/workspace/chemotion`.

**Step 6:** If the Remote Development Extensions for Docker are installed, you will prompted to reopen the folder in a container. Confirm and the container will be created (this will take a while...).

If the prompt does not show, install the extension pack as mentioned above. Then, you will find an icon similar to `><` in the very bottom left of VSCode's status bar. Click on it and select "Reopen in Container".

**Step 6:** After the container is build, you will be able to access a terminal in the container using VSCode (`Terminal` -> `New Terminal`). If all steps were followed correctly, the container will already be initalised with seed data. This was done by the `postCreateCommand` in the `.devcontainer/devcontainer.json` file. See that file for how it is done.

**Step 7:** You are now ready to do some basic testing: run `bundle exec rails server` in the terminal and open the website in your browser (defaults to `localhost:3000`) when prompted. For the first time, this will take some time as sprites need to be generated and Javascripts need to be transpiled/bundled. You can now log in using the seeded admin user "ADM" with password "PleaseChangeYourPassword".

**Optional:** The container is ready to be used with the rspec testing framework. You can do some basic testing executing this command:

```
$ RAILS_ENV=test bundle exec rspec ./spec/features
```

**Optional:** For a better development experience it's recommended to also install the `solargraph` and the `ruby-debug-ide` gems to enable language server and debugging support in VScode for Ruby. This can be done by executing `gem install solargraph ruby-debug-ide`.
