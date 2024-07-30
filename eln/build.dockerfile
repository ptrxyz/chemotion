# GLOBALs that should be passed to this file
ARG VERSION=UNSET
ARG BUILD_TAG_CHEMOTION=UNSET

# Derived values
ARG BASE=chemotion-build/base:${VERSION}
ARG CHEMOTION_RELEASE=${VERSION}

# Private variables, not passed in from outside, but helpful for this file
ARG RUBY_VERSION=latest:2.7
ARG NODE_VERSION=latest:18
# ARG BUNDLER_VERSION=1.17.3
ARG BUNDLER_VERSION=2.2.33
ARG ASDF_VERSION=v0.10.2
ARG PANDOC_VERSION=2.10.1
ARG BUILD_REPO_CHEMOTION=https://github.com/ComPlat/chemotion_ELN

FROM scratch as cache-vol
ADD ./tmp/cache.tar.gz /


########################################################################################################################
# Add ASDF, Ruby, Node, bundler and yarn
########################################################################################################################
# hadolint ignore=DL3006
FROM ${BASE} as asdf-enabled
RUN apt-get install -y --no-install-recommends --autoremove --fix-missing git ca-certificates curl
SHELL ["/bin/bash", "-e", "-o", "pipefail", "-c", "--"]

RUN apt-get install -y --no-install-recommends --autoremove --fix-missing build-essential \
    `# for asdf: ruby` \
    zlib1g-dev libreadline-dev patchelf

ARG ASDF_VERSION
ARG RUBY_VERSION
ARG NODE_VERSION
ARG BUNDLER_VERSION

# make asdf and its shims available everywhere (set ASDF_DIR to where you cloned asdf)
ENV ASDF_DIR=/asdf                                      \
    ASDF_DATA_DIR=/asdf                                 \
    ASDF_DEFAULT_TOOL_VERSIONS_FILENAME=.tool-versions  \
    GEM_HOME=/cache/gems
ENV PATH="${GEM_HOME}/bin:${GEM_HOME}/gems/bin:${ASDF_DIR}/shims:${ASDF_DIR}/bin:${PATH}"

# Ways to get OpenSSL1.1 (required for Ruby 2.7) going on newer Ubuntus:
# - add the .deb packages from a previous Ubuntu version. Needs to be done in all stages.
# - adsf will install OpenSSL if missing, then patchelf can be used to adjust rpaths
# - use asdf's OpenSSL installation and copy libcrypto.so.1.1 to the system directory

# ADD https://launchpad.net/~ubuntu-security/+archive/ubuntu/ppa/+build/24138247/+files/libssl-dev_1.1.1l-1ubuntu1.6_amd64.deb /tmp
# ADD https://launchpad.net/~ubuntu-security/+archive/ubuntu/ppa/+build/24138247/+files/libssl1.1_1.1.1l-1ubuntu1.6_amd64.deb  /tmp
# RUN dpkg -i /tmp/*.deb

# ASDF
RUN git clone https://github.com/asdf-vm/asdf.git /asdf --branch "${ASDF_VERSION}"

COPY --from=cache-vol /cache /cache

# NodeJS
RUN MAKEFLAGS="-j$(getconf _NPROCESSORS_ONLN)" && export MAKEFLAGS && \
    asdf plugin add nodejs https://github.com/asdf-vm/asdf-nodejs.git && \
    asdf install nodejs "${NODE_VERSION}" && \
    asdf global nodejs "${NODE_VERSION}" && \
    npm install -g yarn

# Ruby
RUN MAKEFLAGS="-j$(getconf _NPROCESSORS_ONLN)" && export MAKEFLAGS && \
    asdf plugin add ruby https://github.com/asdf-vm/asdf-ruby.git && \
    asdf install ruby "${RUBY_VERSION}" && \
    asdf global ruby "${RUBY_VERSION}"  && \
    rm -r /asdf/installs/ruby/2.7.8/lib/ruby/gems/2.7.0/gems/bundler-* && \
    rm /asdf/installs/ruby/**/lib/ruby/gems/**/specifications/default/bundler-*.gemspec && \
    gem install bundler -v "${BUNDLER_VERSION}"

# fix the openSSL thing explained above.
RUN SSL_PATH="$(asdf where ruby)" && \
    ln -s "${SSL_PATH}/openssl/lib/libcrypto.so.1.1" /lib/x86_64-linux-gnu/





FROM ${BASE} as raw-eln
RUN apt-get install -y --no-install-recommends --autoremove --fix-missing git ca-certificates curl
SHELL ["/bin/bash", "-e", "-o", "pipefail", "-c", "--"]

ARG BUILD_TAG_CHEMOTION
ARG BUILD_REPO_CHEMOTION
ARG CHEMOTION_RELEASE
ARG VERSION

# checkout Chemotion
RUN test -n "${BUILD_TAG_CHEMOTION}" && echo "BUILD FROM: ${BUILD_REPO_CHEMOTION}" && \
    echo "CHECKOUT: ${BUILD_TAG_CHEMOTION}" && \
    git init --initial-branch=main /chemotion/app && \
    git -C /chemotion/app remote add origin "${BUILD_REPO_CHEMOTION}" && \
    git -C /chemotion/app fetch --tags --depth 1 origin "${BUILD_TAG_CHEMOTION}" && \
    git -C /chemotion/app reset --hard FETCH_HEAD && \
    git -C /chemotion/app describe --abbrev=0 --tags || true

# make sure checkout (kind of) worked
RUN ls -la /chemotion/app | grep -i gemfile
RUN stat /chemotion/app/Gemfile.lock && rm -f /chemotion/app/.tool-versions

WORKDIR /chemotion/app

# Create release file
RUN echo "CHEMOTION_REF=$(git rev-parse --short HEAD || echo unknown)" >> /chemotion/app/.version      && \
    echo "CHEMOTION_TAG=$(git describe --abbrev=0 --tags || echo untagged)" >> /chemotion/app/.version && \
    echo "RELEASE=${CHEMOTION_RELEASE}" >> /chemotion/app/.version                                     && \
    echo "VERSION=${VERSION}" >> /chemotion/app/.version                                               && \
    cat /chemotion/app/.version

# Clean up the configs
RUN find ./config -type f -name '*.yml.example' -print | while read -r f; do \
    echo "$f to ${f%.example}"; \
    yq --input-format yaml -ojson < "$f" | \
    yq -oyaml -P 'with(select(.production != null); . = {"production":.production}) | with(select(.production == null); {})' \
    > "${f%.example}"; \
    done

# Clean up unused files
RUN rm -f config/deploy.rb \
    config/environments/development.rb \
    config/environments/test.rb \
    ./**/*github* \
    ./**/*gitlab* \
    ./**/*travis* \
    ./**/*.bak    \
    ./**/*git* || true

# deal with the weird favicon situation
RUN if [ -f /chemotion/app/public/favicon.ico ] && [ ! -e /chemotion/app/public/favicon.ico.example ]; then \
        cp /chemotion/app/public/favicon.ico /chemotion/app/public/favicon.ico.example; \
    fi

# misc. tweaks: configure yarn
RUN echo -e "--modules-folder ${NODE_PATH}\n--ignore-engines" > /chemotion/app/.yarnrc && \
    if [[ ! -f /chemotion/app/config/klasses.json ]]; then echo '[]' > /chemotion/app/config/klasses.json; fi

# Add application additives as-is. We do this at the end to avoid alterations from previous actions.
COPY ./additives/chemotion/ /chemotion/app

# Create persistent volume directories:
# - move [app]/uploads to /chemotion/data
# - move [app]/public/images to /data
RUN mkdir -p /chemotion/data && \
    mv /chemotion/app/uploads /chemotion/data/ || mkdir -p /chemotion/data/uploads && \
    ln -s /chemotion/data/uploads/ /chemotion/app/uploads && \
    mkdir -p /chemotion/data/public/ && \
    mv /chemotion/app/public/images/ /chemotion/data/public/ || mkdir -p /chemotion/data/public/images && \
    ln -s /chemotion/data/public/images/ /chemotion/app/public/images && \
    mv /chemotion/app/public/safety_sheets/ /chemotion/data/public/ || mkdir -p /chemotion/data/public/safety_sheets && \
    ln -s /chemotion/data/public/safety_sheets/ /chemotion/app/public/safety_sheets && \
    cp /chemotion/app/.version /chemotion/data/.version && \
    mkdir -p /chemotion/data/public/images/thumbnail

RUN sed -i "/gem 'rdkit_chem'/d"                /chemotion/app/Gemfile && \
    sed -i "/gem 'openbabel'/d"                 /chemotion/app/Gemfile && \
    echo "gem 'rdkit_chem', git: 'https://github.com/ptrxyz/rdkit_chem.git', branch: 'pk01'" >> /chemotion/app/Gemfile && \
    echo "gem 'openbabel', '2.4.90.3', git: 'https://github.com/ptrxyz/openbabel-gem.git', branch: 'ptrxyz-ctime-fix'" >> /chemotion/app/Gemfile && \
    echo "gem 'tzinfo-data'"                 >> /chemotion/app/Gemfile && \
    echo "gem 'activerecord-nulldb-adapter'" >> /chemotion/app/Gemfile

RUN touch /raw-eln.done





FROM asdf-enabled as yarn-installed
SHELL ["/bin/bash", "-e", "-o", "pipefail", "-c", "--"]

ENV NODE_ENV=production                                 \
    NODE_PATH=/cache/node_modules/                      \
    NODE_MODULES_PATH=/cache/node_modules/              \
    NODE_OPTIONS=--max_old_space_size=4096              \
    YARN_CACHE_FOLDER=/cache/yarn

WORKDIR /chemotion/app

COPY --from=raw-eln /chemotion/app/package.json /chemotion/app/
COPY --from=raw-eln /chemotion/app/yarn.lock    /chemotion/app/

# RUN --mount=type=bind,source=/tmp/modcache/yarn,destination=/cache/yarn,rw \
RUN MAKEFLAGS="-j$(getconf _NPROCESSORS_ONLN)" && export MAKEFLAGS && \
    yarn install --modules-folder ${NODE_PATH} --ignore-engines --ignore-scripts 2>&1 | grep -v ^warning

RUN touch /yarn.done





FROM asdf-enabled as bundle-installed
SHELL ["/bin/bash", "-e", "-o", "pipefail", "-c", "--"]

RUN apt-get install -y --no-install-recommends --autoremove --fix-missing build-essential \
    `# for the gems` \
    cmake libpq-dev swig \
    libboost-serialization-dev \
    libboost-iostreams-dev \
    libboost-system-dev \
    libeigen3-dev \
    libmagickcore-dev \
    python3-dev libsqlite3-dev

ENV BUNDLE_PATH=/cache/bundle                      \
    BUNDLE_CACHE_PATH=/cache/bundle/package-cache  \
    BUNDLE_USER_HOME=/cache/bundle                 \
    BUNDLE_APP_CONFIG=/cache/bundle                \
    BUNDLE_WITHOUT=development:test                \
    BUNDLE_CACHE_ALL=1                             \
    BUNDLE_SILENCE_ROOT_WARNING=1                  \
    GEM_HOME=/cache/gems                           \
    RAILS_ENV=production                           \
    RAKE_ENV=production                            \
    THOR_SILENCE_DEPRECATION=1

WORKDIR /chemotion/app

COPY --from=raw-eln /chemotion/app/Gemfile      /chemotion/app/
COPY --from=raw-eln /chemotion/app/Gemfile.lock /chemotion/app/

# RUN --mount=type=bind,source=/tmp/modcache,destination=/cache,rw \
RUN MAKEFLAGS="-j$(getconf _NPROCESSORS_ONLN)" && export MAKEFLAGS && \
    `# bundle package --no-install` && \
    bundle install --jobs="$(getconf _NPROCESSORS_ONLN)" --retry=3

# RUN --mount=type=bind,source=/tmp/modcache,destination=/build-cache,rw \
#     mkdir -p /cache/bundle && \
#     cp -ar /build-cache/bundle /cache

RUN bundle clean --force
RUN touch /bundle.done





FROM ${BASE} as eln-ruby-node
RUN apt-get install -y --no-install-recommends --autoremove --fix-missing git ca-certificates curl
SHELL ["/bin/bash", "-e", "-o", "pipefail", "-c", "--"]

# runtime requirements for Chemotion
RUN apt-get install -y --no-install-recommends --autoremove --fix-missing \
    `# for chemotion` \
    libboost-serialization1.74.0 \
    libboost-iostreams1.74.0 \
    postgresql-client `# also adds pg_isready` \
    inkscape `# this installs python3` \
    imagemagick \
    librsvg2-bin `# for thumbnail generation of reasearch plans` \
    locales \
    ghostscript \
    `# utilitites` \
    vim iproute2 sudo make

COPY --from=asdf-enabled /root/.tool-versions /root/.tool-versions
COPY --from=asdf-enabled /asdf/ /asdf/
COPY --from=asdf-enabled /lib/x86_64-linux-gnu/libcrypto.so.1.1 /lib/x86_64-linux-gnu/libcrypto.so.1.1
COPY --from=raw-eln /chemotion /chemotion
# COPY --from=bundle-installed /cache/bundle/ /cache/bundle/
COPY --from=bundle-installed /cache/bundle /cache/bundle
COPY --from=bundle-installed /cache/gems /cache/gems
COPY --from=yarn-installed /cache/node_modules/ /cache/node_modules/

# Add other additives
COPY ./additives/various/fontfix.conf /etc/fonts/conf.d/99-chemotion-fontfix.conf
COPY ./additives/various/policy.xml   /etc/ImageMagick-6/policy.xml
COPY ./additives/embed/               /embed/

# more misc tweaks
ARG PANDOC_VERSION
RUN curl -L -o /tmp/pandoc.deb "https://github.com/jgm/pandoc/releases/download/${PANDOC_VERSION}/pandoc-${PANDOC_VERSION}-1-amd64.deb" && \
    dpkg -i /tmp/pandoc.deb && rm /tmp/pandoc.deb && \
    ln -s /embed/bin/chemotion /bin/chemotion && \
    cp /chemotion/app/.version /.version && \
    chmod +x /embed/bin/*

# CONFIG_PIDFILE: the ELN uses this file to communicate that it is alive
# and the bgworker can start working. (Needs to be on a shared volume.)
ENV ASDF_DIR=/asdf                                      \
    ASDF_DATA_DIR=/asdf                                 \
    ASDF_DEFAULT_TOOL_VERSIONS_FILENAME=.tool-versions  \
    BUNDLE_PATH=/cache/bundle                           \
    BUNDLE_CACHE_PATH=/cache/bundle/package-cache       \
    BUNDLE_USER_HOME=/cache/bundle                      \
    BUNDLE_APP_CONFIG=/cache/bundle                     \
    BUNDLE_WITHOUT=development:test                     \
    BUNDLE_CACHE_ALL=1                                  \
    BUNDLE_SILENCE_ROOT_WARNING=1                       \
    GEM_HOME=/cache/gems                                \
    RAILS_ENV=production                                \
    RAKE_ENV=production                                 \
    NODE_ENV=production                                 \
    NODE_PATH=/cache/node_modules/                      \
    NODE_MODULES_PATH=/cache/node_modules/              \
    NODE_OPTIONS=--max_old_space_size=4096              \
    YARN_CACHE_FOLDER=/cache/yarn                       \
    TERM=xterm-256color                                 \
    THOR_SILENCE_DEPRECATION=1                          \
    RAILS_LOG_TO_STDOUT=1                               \
    PASSENGER_DOWNLOAD_NATIVE_SUPPORT_BINARY=0          \
    CONFIG_PIDFILE=/chemotion/app/tmp/eln.pid
ENV PATH="${GEM_HOME}/bin:${GEM_HOME}/gems/bin:${ASDF_DIR}/shims:${ASDF_DIR}/bin:${PATH}"

WORKDIR /chemotion/app

# # (re)install node packages from cache
# RUN MAKEFLAGS="-j$(getconf _NPROCESSORS_ONLN)" && export MAKEFLAGS && \
#     yarn install --modules-folder ${NODE_PATH} --ignore-engines 2>&1 | grep -v ^warning && \
#     yarn cache clean 2>&1 | grep -v ^warning && \
RUN ln -s ${NODE_PATH} /chemotion/app/node_modules && bash package_postinstall.sh

# (re)install gems from cache
# RUN --mount=type=cache,id=dev-gem-cache,sharing=locked,target=/build-cache \
RUN MAKEFLAGS="-j$(getconf _NPROCESSORS_ONLN)" && export MAKEFLAGS && \
    bundle install --jobs="$(getconf _NPROCESSORS_ONLN)" --retry=3

# Pierre's fix for the style sheet issue
RUN cp -f -l /cache/node_modules/ag-grid-community/dist/styles/ag-grid.css /chemotion/app/app/assets/stylesheets/.  || true && \
    cp -f -l /cache/node_modules/ag-grid-community/dist/styles/ag-theme-alpine.css /chemotion/app/app/assets/stylesheets/.  || true && \
    cp -f -l /cache/node_modules/antd/dist/antd.css /chemotion/app/app/assets/stylesheets/. || true && \
    cp -f -l /cache/node_modules/react-datepicker/dist/react-datepicker.css /chemotion/app/app/assets/stylesheets/. || true && \
    cp -f -l /cache/node_modules/react-select/dist/react-select.css /chemotion/app/app/assets/stylesheets/. || true && \
    cp -f -l /cache/node_modules/react-treeview/react-treeview.css /chemotion/app/app/assets/stylesheets/.  || true && \
    cp -f -l /cache/node_modules/react-vis/dist/style.css /chemotion/app/app/assets/stylesheets/react-vis-styles.css || true && \
    cp -f -l /cache/node_modules/react-virtualized/styles.css /chemotion/app/app/assets/stylesheets/react-virtualized-styles.css || true && \
    cp -f -l /cache/node_modules/react-virtualized-select/styles.css /chemotion/app/app/assets/stylesheets/react-virtualized-select-styles.css || true

# misc. tweaks: configure yarn
RUN echo -e "--modules-folder ${NODE_PATH}\n--ignore-engines" > /chemotion/app/.yarnrc && \
    if [[ ! -f /chemotion/app/config/klasses.json ]] || [[ ! -f /chemotion/app/node_modules/klasses.json ]]; then \
    echo '[]' > /chemotion/app/config/klasses.json; \
    echo '[]' > /chemotion/app/node_modules/klasses.json; \
    echo '[]' > /cache/node_modules/klasses.json; \
    fi

RUN cat /chemotion/app/node_modules/klasses.json

# precompile
RUN export SECRET_KEY_BASE="build"                                                                          && \
    KETCHER_PATH=$(bundle info --path ketcherails)                                                          && \
    UIFILE_PATH="${KETCHER_PATH}/app/assets/javascripts/ketcherails/ui/ui.js.erb"                           && \
    rm /chemotion/app/config/scifinder_n.yml && \
    cp "${UIFILE_PATH}" /tmp/ui.js.bak                                                                      && \
    sed -i 's/Ketcherails::TemplateCategory.with_approved_templates.pluck(:id)/[]/g' "${UIFILE_PATH}"       && \
    bundle exec rake DISABLE_DATABASE_ENVIRONMENT_CHECK=1 DATABASE_URL=nulldb://user:pass@127.0.0.1/dbname webpacker:compile 2>&1 | grep -v ^warning  && \
    bundle exec rake DISABLE_DATABASE_ENVIRONMENT_CHECK=1 DATABASE_URL=nulldb://user:pass@127.0.0.1/dbname assets:precompile 2>&1 | grep -v ^warning  && \
    mv /tmp/ui.js.bak "${UIFILE_PATH}"

# Safety valve to make sure things seem to be fine ...
RUN env ; touch /chemotion/app/.env && bundle exec dotenv erb /chemotion/app/config/secrets.yml

# cleanup
RUN rm -rf /tmp/* /var/tmp/*





# Stage 5: finalize the image
FROM eln-ruby-node as app
ARG VERSION

EXPOSE 4000

WORKDIR /chemotion/app
CMD ["/embed/run.sh"]

HEALTHCHECK --interval=30s --timeout=10s --start-period=300s --retries=3 \
    CMD /embed/health.sh || exit 1

VOLUME [ "/chemotion/app", "/chemotion/data" ]

LABEL \
    "org.opencontainers.image.authors"="Chemotion Team" \
    "org.opencontainers.image.title"="Chemotion ELN" \
    "org.opencontainers.image.description"="Image for Chemotion ELN" \
    "org.opencontainers.image.version"="${VERSION}" \
    "chemotion.internal.service.id"="eln"
