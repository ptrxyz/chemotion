# Chemotion-CLI

 > The following is mostly whishful thinking. Yet feel free to dream with us.

CLI tool to manage chemotion installations.

Following features are planned/thought of:

  - Installation & Deployment: we plan to implement `chemotion instance install` to install a Chemotion instance
  - Upgrade: use `chemotion instance upgrade` to upgrade an existing Chemotion instance
  - Backups: `chemotion backup create|restore` to savely store your data somewhere
  - Instance life cycle commands, such as `chemotion instance start|stop|pause|restart|status`
  - Manage Settings: `chemotion settings import|export` to import/export you settings and `chemotion instance configure` to run configuration wisards that help you to create configuration stubs
  - Frequently asked for features for the Chemotion Administrator: `chemotion user show|add|delete|password-reset`, `chemotion system info|rails-shell|shell`


# Planned concept for CLI

We plan to follow one of the following layouts, depending on which one proves to be more handy in every day use.

```
general: cli-executable  <resource>  <command>  <argument>  <flags>
         └─────┬──────┘  └───┬────┘  └───┬───┘  └───┬────┘  └──┬──┘
example:    chemotion     instance    restart   MyInstance  --force
```

```
general: cli-executable  <command>  <resource>  <argument>  <flags>
         └─────┬──────┘  └───┬───┘  └───┬────┘  └───┬────┘  └──┬──┘
example:    chemotion     restart    instance   MyInstance  --force
```

# Known limitations
- The following flags cannot be specified in the configuration (`chemotion-cli.yml`) file:
  - `--config`: because that creates a circular dependency