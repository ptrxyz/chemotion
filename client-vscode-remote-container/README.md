# Contianer files for VS-Code Remote Editing

This folder contains anything you need to use Visual Studio Code's remote developing features to start developing on the Chemotion ELN inside a Docker development container.

## Requirements

-   Visual Studio Code (non-OSS version. Remote development is not supported by the open source builds.)
-   working Docker installation
-   the [Remote Developing extension pack](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.vscode-remote-extensionpack) installed.

For more information check out [Microsoft's knowledge base](https://code.visualstudio.com/docs/remote/remote-overview) for more detailed information.

# General Information

The purpose is to setup a VSCode remote development environment in a Docker container. See the knowledge base article above if all this does not sound familiar to you.

The setup is as follows:

`.devcontainer/devcontainer.json` contains the definition of which Dockerfiles/docker-compose files to use and is read by VSCode. In addition it contains some postCreation command, that is run as soon as the container finished building. In our case this is:

-   creating a database
-   enabling all necessary extensions in the container (for VSCode) and for the database
-   seeding the database

The used Dockerfile to build an development environment is `Dockerfile.vscode`. The file is it's own documentation.

To start an additional container for a database server, `docker-compose.vscode` is used. Again, the file is pretty self-explainatory

# Usage

**Step 1:** First, check out the Chemotion ELN source from the [GitHub Repo](https://github.com/ComPlat/chemotion_ELN). In the following `/workspace/chemotion` will be used as the target directory:

```
$ git clone https://github.com/ComPlat/chemotion_ELN /workspace/chemotion
```

**Step 2:** place all files from this folder into Chemotion's source directory:

```
$ cp * .* /workspace/chemotion
$ ls /workspace/chemotion
.              scripts                    .rubocop.yml                Rakefile
..             spec                       .rubocop_todo.yml           VERSION
.bundle        tmp                        .ruby-gemset.example        config.ru
.devcontainer  uploads                    .ruby-version.example       createDB.mjs
.git           vendor                     .simplecov                  dbinit.sh
.github        .babelrc.bak               .travis.yml                 docker-compose.test.yml
app            .dockerignore              CHANGELOG.md                docker-compose.vscode
backup         .env.development           Capfile                     fontcustom.yml
bin            .env.production.example    Dockerfile.vscode           output.json
config         .env.test                  Dockerfile.focal.gitlab-ci  package.json
data           .eslintrc                  Gemfile                     run.sh
db             .fontcustom-manifest.json  Gemfile.lock                secret_key.conf
lib            .gitignore                 Gemfile.plugin.example      yarn.lock
log            .gitlab-ci.yml             INSTALL.md                  yarn-error.log
node_modules   .nvmrc                     LICENSE
public         .rspec                     README.md
```

(You want to make sure that `.devcontainer`, `docker-compose.vscode` and `Dockerfile.vscode` reside in the Chemotion source's root directory as shown here)

**Step 3:** Open VSCode and open the Chemotion folder: `File` -> `Open Folder` -> select the right folder, here `/workspace/chemotion`.

**Step 4:** If the Remote Development Extensions for Docker are installed, you will prompted to reopen the folder in a container. Confirm and the container will be created (this will take a while...).

If the prompt does not show, install the extension pack as mentioned above. Then, you will find an icon similar to `><` in the very bottom left of VSCode's status bar. Click on it and select "Reopen in Container".

**Step 5:** After the container is build, you will be able to access a terminal in the container using VSCode (`Terminal` -> `New Terminal`). If all steps were followed correclty, the container will already be initalised with seed data. This was done by the `postCreateCommand` in the `.devcontainer/devcontainer.json` file. See that file for how it is done.

You are now ready to do some basic testing: run `bundle exec rails server` in the terminal and open the website in your browser (defaults to `localhost:3000`) when prompted. For the first time, this will take some time as sprites need to be generated and Javascripts need to be transpiled/bundled. You can now log in using the seeded admin user "ADM" with password "PleaseChangeYourPassword".

**Optional**: the container is ready to be used with the rspec testing framework. You can do some basic testing executing this command:

```
$ RAILS_ENV=test bundle exec rspec ./spec/features
```
