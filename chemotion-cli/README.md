# Chemotion-CLI

 > The following is mostly whishful thinking. Yet feel free to dream with us.

CLI tool to manage chemotion installations.

Following features are planned/thought of:

 - Installation & Deployment: we plan to implement `chemotion instance create <name>` to create a Chemotion instance called `<name>`.
 - Upgrade: use `chemotion instance upgrade <name>` to upgrade an existing Chemotion instance.
 - Backup data: `chemotion instance backup|restore <name> --data -f <my_backup.tgz>` to safely store your data somewhere
 - Manage Settings: `chemotion instance backup|restore <name> --settings -f <my_settings.dat>` to import/export your settings and `chemotion configure <name>` to run configuration wizards that help you to create configuration stubs.
 - Instance life cycle commands, such as `chemotion instance status|create|upgrade|start|pause|stop|restart|delete <name>`
 - Frequently required for features as a Chemotion administrator: `chemotion user show|add|remove|passwd <user_name>`, `chemotion system info|rails-shell|shell`


# Planned concept for CLI

We plan to follow the following syntax:

```
general: cli-executable  <resource>  <command>  <argument>  <flags>
         └─────┬──────┘  └───┬────┘  └───┬───┘  └───┬────┘  └──┬──┘
example:    chemotion     instance    restart   MyInstance  --force
```

One should be able to execute the following combinations:

resource | alias | supported actions
---------|-------|------------------
instance | i     | status, create, upgrade,     start,  pause, stop, restart, delete
user     | u     | show,   add,    remove,      passwd
system   | s     | info,   shell,  rails-shell
