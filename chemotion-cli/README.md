# Chemotion-CLI

Chemotion CLI tool is there to help you manage installation(s) of Chemotion on a machine. The goal is to make installation, maintenance and upgradation of Chemotion as easy as possible.

> :information_source: [Link to quick intro video](https://youtu.be/10fk2C6qku0)

## Download

### Get the binary

The Chemotion CLI tool is a binary file and needs no installation. The only prerequisite is that you install [Docker Desktop](https://www.docker.com/products/docker-desktop/) (and, on Windows, [WSL](https://docs.microsoft.com/en-us/windows/wsl/install)). Depending on your OS, you can download the lastest release of the CLI from [here](https://github.com/harivyasi/chemotion/releases/latest). Builds for the following systems are available:

- Linux, amd64
- Windows, amd64; remember to turn on [Docker integration with WSL](https://docs.docker.com/desktop/windows/wsl/)
- macOS, apple-silicon
- macOS, amd64

Please be sure that you have both, `docker` and `docker compose` commands. This should be the case if you install Docker Desktop following the instructions [here](https://docs.docker.com/desktop/#download-and-install). If you choose to install only Docker Engine, then please make sure that you _also_ have `docker compose` as a command (as opposed to `docker-compose`).

### Make it an executable

On Linux, make this file executable by doing: `chmod u+x chemotion`.

On Windows, the file should be executable by default, i.e. do nothing.

On macOS, make this file executable by doing: `chmod u+x chemotion.amd.osx` or `chmod u+x chemotion.arm.osx`. If the there is a security pop-up when running the command, please also `Allow` the executable in `System Preferences > Security & Privacy`.

### Important Note:

All commands here, and all the documentation of the tool, use term `chemotion` to refer to the executable. Depending on your configuration, you may have to use any one of the following:

- `./chemotion`
- `.\chemotion.exe`
- `./chemotion.arm.osx`
- `./chemotion.amd.osx`

## First run

### Make a dedicated folder

Make a folder where you want to store installation(s) of Chemotion. Ideally this folder should be in the largest drive (in terms of free space) of your system. Remember that Chemotion also uses space via Docker (docker containers, volumes etc.) and therefore you need to make sure that your system partition has abundant free space.

### Install

To begin with installation, execute: `chemotion install` and follow the prompt. The first installation can take really long time (15-30 minutes depending on your download and processor speeds).

This will create the first (production-grade) `instance` of Chemotion on your system. Generally, this is suffice if you want to use Chemotion in a single scientific group/lab. By default

- this first instance will be available on port 4000
- this first instance will be the `selected` instance.

> :warning: **chemotion-cli.yml**: Installation also creates a file called `chemotion-cli.yml`. This file is critical as it contains information regarding existing installations. Removing the file will render the CLI clueless about existing installations and it will behave as if Chemotion was never installed. Please do not remove the file. Ideally there should be no need for you to modify it manually.

### The `selected` instance

Once you install multiple instances of Chemotion, the actions of CLI will pertain to only one of them i.e. you will be managing only one of them. This instance is referred to as the `selected` instance and it's name is stored in a local file (`chemotion-cli.yml`). You can do `chemotion instance switch` to switch to another instance.

You can also select an instance _temporarily_ by giving its name to the CLI as a flag when you start it e.g. `chemotion instance status --instance the-other-one`.

### Start and Stop Chemotion

To turn on, and off, the `chosen` instance, issue the commands:

- `chemotion on`, and
- `chemotion off`.

### Upgrading an instance (for versions 1.3 and above)

As long as you installed an instance of Chemotion using this tool, the upgrade process is quite straightforward:

- First make sure that you have the latest version of this tool. You can check the version of your chemotion binary by doing `chemotion --version`. If necessary, follow the instructions in the [download](#download) section again. Feel free to replace the existing `chemotion` file. DO NOT remove/replace the `chemotion-cli.yml` file.
- Prepare for update by running `chemotion advanced pull-image`. This will download the latest chemotion image from the internet if not already present on the system. Downloading the image outside of downtime saves you time later on.
- Schedule a downtime of at least 15 minutes; more if you have a lot of data that needs to backed up. During the downtime, run `chemotion instance backup` to backup your data followed by `chemotion instance upgrade` to update the instance.

### Updating this CLI tool

Starting from version `0.2.0-alpha`, the tool itself can be updated to the latest version by running `chemotion advanced update`.

> :info: If you are updating from version `0.1.x-alpha` to a newer version, please get in touch with your `chemotion-cli.yml` file. We will update this file manually (< 5 minutes) to get your installation up and working again! If you want, you can do it yourself following the example at the bottom of this page.

### Uninstallation

> :warning: be sure about what you want to do!

You can uninstall everything created by the CLI tool by running: `chemotion advanced uninstall`. Last you can simply delete the downloaded binary itself.

## Silent and Debug Use

Almost all features of the CLI can be used in silent mode i.e. without any input from user as long as all required pieces of information have been provided using flags. In silent mode, most of the output from the CLI (but not that of docker) is logged only in the log file, and not put on screen.

To use the CLI in silent mode, add the flag `-q`/`--quiet` to your command. The CLI will then use default values and other flags to try and accomplish the action. Examples:

```bash
./chemotion install -q --name first-instance --address https://myuni.de:3000
./chemotion instance switch --name switch-to-this-instance -q
```

Similarly, the CLI can be run in Debug mode when you encounter an error. This produces a very detailed log file containing a trace of actions you undertake. Telling us about the error and sending us the log file can help us a lot when it comes to helping you.

## Known limitations and bugs

- The following flags cannot be specified in the configuration (`chemotion-cli.yml`) file:
  - `--config`: because that creates a circular dependency
  - `chemotion off`: does not lead to exit of containers with exit code 0.
- Everything happens in the folder (and subfolders) where `chemotion` is executed. All files and folders are expected to be there; otherwise failures can happen.

# Planned concept for CLI

The commands have the following general layout:

```
general: cli-executable  <resource>  <command>  <flags>
         └─────┬──────┘  └───┬────┘  └───┬───┘  └──┬──┘
example:    chemotion     instance    restart   --force
```

Following features are exist:

- ✔ Installation & Deployment: `chemotion install` installs a production instance that is ready to use.
- ✔ Instance life cycle commands: `chemotion on|off` and `chemotion instance status|stats|list|restart`.
- ✔ Multiple instances: `chemotion instance add|switch|remove` can be used to manage multiple instances.
- ✔ Upgrade: use `chemotion instance upgrade` to upgrade an existing Chemotion instance.
- ✔ Backups: use `chemotion instance backup` to save the data associated with an instance.
- ✔ Shell access: using `chemotion instance console` to access shell/rails/SQL console of an instance

Following features are planned:

- Restore backup: `chemotion restore`
- Manage Settings: `chemotion instance settings --import|--export` to import/export settings and to run auto-configuring wizards.
- Features for the Chemotion Administrator: `chemotion user show|add|delete|password-reset`
- Command to manage underlying docker installation i.e. free up space and prune network

## Moving from CLI version 0.1.x-alpha to 0.2.0-alpha

The only difference is formatting of the `chemotion-cli.yml` file.

- The global keys that handle state of the tool i.e. `selected`, `quiet` and `debug` have now been moved to `cli_state:selected`, `cli_state:quiet` and `cli_state:debug` respectively.
- The `instances:<instance_name>:address` and `instances:<instance_name>:protocol` keys have been removed. Instead, we have `instances:<instance_name>:accessaddress` which stores the full URL that is used to access the ELN instance.
- A new key called `instances:<instance_name>:environment` has been introduced. This is now used to create the `shared/pullin/.env` file **everytime** the instance is (re)started. **Please** move all your `key=value` pairs from this `.env` file to the `chemotion-cli.yml` in `key: value` format as sub-keys of the `instances:<instance_name>:environment` key.
- With these changes, the version of this YAML file has been changed from `"1.0"` to `"1.1"`.

Therefore, if your file looked as follows:

```yaml
instances:
  main:
    address: mynotebook.kit.edu
    debug: false
    kind: Production
    name: main-ee5e5424
    port: 4000
    protocol: http
    quiet: false
  second:
    address: localhost
    debug: false
    kind: Production
    name: second-ff6f6535
    port: 4100
    protocol: http
    quiet: false
selected: main
version: "1.0"
```

Please change it to look as follows:

```yaml
cli_state:
  debug: false
  quiet: false
  selected: main
instances:
  main:
    accessaddress: http://mynotebook.kit.edu
    environment:
      url_host: ifgs6.ifg.kit.edu
      url_protocol: http
    kind: Production
    name: main-ee5e5044
    port: 4000
  second:
    accessaddress: http://localhost:4100
    environment:
      url_host: localhost:4100
      url_protocol: http
      smtp_port: ...<key-value pairs from /shared/pullin/.env file>...
    kind: Production
    name: second-ff6f6535
    port: 4100
version: "1.1"
```

Once you have made the changes, please downloaded the latest version of the CLI from [here](https://github.com/harivyasi/chemotion/releases/latest).
