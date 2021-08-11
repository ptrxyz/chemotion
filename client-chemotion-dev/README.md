# Container for development

Containerized dev environment for programmers.
Please run `./firstun.sh` after checkout.

Source files will be placed in `./src`, ready to be edited ;)
To get access to the container, run `docker-compose run eln (shell|user-shell)`.
See `docker-compose run eln help` for more info.

## Usage:

```
./firstrun.sh                          # necessary before first run to initialize the container
docker-compose up                      # to start the development environment
docker-compose down --remove-orphans   # tear things down
docker-compose logs -f eln             # show logs of a specific service (-f for follow)
```
