# Changelog

## [1.5.4] - 2023-05-05
-   See Changelog for [ELN](https://github.com/ComPlat/chemotion_ELN/blob/main/CHANGELOG.md)

## [1.5.3] - 2023-04-25
-   See Changelog for [ELN](https://github.com/ComPlat/chemotion_ELN/blob/main/CHANGELOG.md#v153)

## [1.5.2] - 2023-04-11
-   See Changelog for [ELN](https://github.com/ComPlat/chemotion_ELN/blob/main/CHANGELOG.md#v152)

## [1.5.1] - 2023-03-24
-   Fixed bugs in embeded scripts.

## [1.5.0] - 2023-03-14

-   Chore: new version of spectra supported
-   Feature: added support for NMRium
-   Feature: added healthchecks to all containers
-   Feature: containers now show internal versions on startup
-   Feature: it's now possible to dump the current config with the `chemotion dumpConfig` command

-   Change: Recreation of Ketcher sprite sheet is now triggered on boot up (during [dbcheck] step)
-   Change: Container startup is now cancelled if subtasks fail
-   Change: added config for Shrine

-   Fixed a bug with `chemotion railsc` command.
-   Fixed a bug where thumbnails were not generated for research plans. (rsvg-convert is now a dependency)
-   Fixed an issue where convert did not allow to convert PDF/EPS/PS files.
-   Fixed an issue with `chemotion restore`-command not being found under certain circumstances
-   Fixed an issue where the ELN application's secret key was not generated correctly.

## [1.4.1-3] - 2022-12-19

-   Fixed a bug where the ELN looses some stylesheets under certain circumstances.
-   Added a default secret key. This is not secure, users are encouraged to change it.
-   Added a Startup-Delay to ketcherSVC to prevent it from continuously restarting.

## [1.4.1-2] - 2022-12-02

### Changed

-   Fixed a bug where the ELN application's secret key was not generated correctly.
