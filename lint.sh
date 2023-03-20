#!/bin/bash
set -euo pipefail

# Chemotion container image linter to check if the build is at least "sane".
# This script should be run in the final container.
cd /

# asdf
asdf --version
asdf list
asdf list | grep nodejs > /dev/null
asdf list | grep ruby > /dev/null

# Node
[[ "${NODE_ENV:-}" == "production" ]]
[[ "${NODE_PATH:-}" == "/cache/node_modules/" ]]
[[ "${NODE_MODULES_PATH:-}" == "/cache/node_modules/" ]]
[[ "${NODE_OPTIONS:-}" == *"max_old_space_size"* ]]
test -d "${NODE_PATH:-}"
test -d "${NODE_MODULES_PATH:-}"
test -L "/chemotion/app/node_modules"
[[ "$(readlink /chemotion/app/node_modules)" == "${NODE_PATH:-}" ]]
(
    cd /chemotion/app
    node -v
    npm -v
    npx -v
    yarn -v
    yarn check --integrity
    node -e 'require("webpack")'
)

# Ruby + Rails
test -n "${BUNDLE_PATH:-}"
[[ "${BUNDLE_PATH:-}" == "/cache/bundle" && -d "${BUNDLE_PATH:-}" ]]

test -n "${BUNDLE_USER_HOME:-}"
[[ "${BUNDLE_USER_HOME:-}" == "${BUNDLE_PATH:-}" && -d "${BUNDLE_USER_HOME:-}" ]]

test -n "${GEM_HOME:-}"
[[ "${GEM_HOME:-}" == "/cache/gems" && -d "${GEM_HOME:-}" ]]

test -n "${BUNDLE_CACHE_ALL:-}"
test -n "${RAILS_ENV:-}"

[[ "${RAILS_ENV:-}" == "production" ]]
[[ "${RAKE_ENV:-}" == "${RAILS_ENV:-}" ]]
test -n "${RAILS_LOG_TO_STDOUT:-}"
[[ "${PASSENGER_DOWNLOAD_NATIVE_SUPPORT_BINARY:-}" == 0 ]]

(
    cd /chemotion/app
    ruby -v
    gem -v
    bundle -v
    bundle exec rails -v
    test -d "$(bundle info --path rdkit_chem)"
)

# Chemotion + Container specific
[[ "${CONFIG_PIDFILE:-}" == "/chemotion/app/tmp/eln.pid" ]]
test -d "$(dirname "${CONFIG_PIDFILE:-}")"
test -d /shared
test -d /embed

test -d /chemotion/app
test -d /chemotion/data

test -d /chemotion/data/public/images
test -L /chemotion/app/public/images
[[ "$(readlink /chemotion/app/public/images)" == "/chemotion/data/public/images/" ]]

test -d /chemotion/data/uploads
test -L /chemotion/app/uploads
[[ "$(readlink /chemotion/app/uploads)" == "/chemotion/data/uploads/" ]]

test -d /chemotion/data/public/images/thumbnail
test -d /chemotion/app/public/images/thumbnail

test -f /chemotion/app/public/favicon.ico.example
test -f /chemotion/data/public/images/transparent.png
test ! -e /chemotion/app/public/sprite.png

# check for some additives
test -f /etc/ImageMagick-6/policy.xml
grep -i -E "read" -E "write" -E "pdf" /etc/ImageMagick-6/policy.xml
test -f /etc/fonts/conf.d/99-chemotion-fontfix.conf


# check versions
cat /chemotion/app/VERSION
cat /.version
v1="$(md5sum /.version)"
v2="$(md5sum /chemotion/app/.version)"
v3="$(md5sum /chemotion/data/.version)"
[[ "$v1" == "$v2" && "$v2" == "$v3" ]]

# check dynamically linked dependencies
find / -iname '*.so' -type f -print0 | xargs -0 ldd | grep -i "not found" || true | sort | uniq
