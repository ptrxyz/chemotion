# GLOBALs that should be passed to this file
ARG BASE=chemotion-build/base:1.4.0
ARG CHEMOTION_BUILD_TAG=67d363fa018ff448e0c83919c3e07535022e4978
ARG CHEMOTION_RELEASE=1.4.0p221001
ARG CHEMOTION_VERSION=1.4.0
ARG BUILDRUN

ARG RUBY_VERSION=latest:2.6
ARG NODE_VERSION=14.20.0
ARG BUNDLER_VERSION=1.17.3

# GLOBALs that should not be changed.
ARG ASDF_VERSION=v0.10.2
ARG PANDOC_VERSION=2.10.1

ARG BUILD_PATH=/builder
ARG PORT=4000


########################################################################################################################
# Add ASDF, Ruby, Node, bundler and yarn
########################################################################################################################
# hadolint ignore=DL3006
FROM ${BASE} as eln-ruby-node
RUN apt-get install -y --no-install-recommends git ca-certificates curl
ARG BUILD_PATH

ARG ASDF_VERSION
ARG RUBY_VERSION
ARG NODE_VERSION
ARG BUNDLER_VERSION

# make asdf and its shims available everywhere (set ASDF_DIR to where you cloned asdf)
ENV ASDF_DIR=/asdf                                      \
    ASDF_DATA_DIR=/asdf                                 \
    ASDF_DEFAULT_TOOL_VERSIONS_FILENAME=.tool-versions  \
    GEM_HOME=${BUILD_PATH}/cache/gems
ENV PATH="${GEM_HOME}/bin:${GEM_HOME}/gems/bin:${ASDF_DIR}/shims:${ASDF_DIR}/bin:${PATH}"

RUN apt-get install -y --no-install-recommends build-essential \
    `# for asdf: ruby` \
    zlib1g-dev libssl-dev libreadline-dev

RUN git clone https://github.com/asdf-vm/asdf.git /asdf --branch "${ASDF_VERSION}"

# NodeJS
RUN asdf plugin add nodejs https://github.com/asdf-vm/asdf-nodejs.git && \
    asdf install nodejs "${NODE_VERSION}" && \
    asdf global nodejs "${NODE_VERSION}" && \
    npm install -g yarn

# Ruby
RUN asdf plugin add ruby https://github.com/asdf-vm/asdf-ruby.git && \
    asdf install ruby "${RUBY_VERSION}" && \
    asdf global ruby "${RUBY_VERSION}" && \
    gem install bundler -v "${BUNDLER_VERSION}"

########################################################################################################################
# Runtime
########################################################################################################################
# hadolint ignore=DL3006
FROM "${BASE}" as eln-runtimes
RUN apt-get install -y --no-install-recommends git ca-certificates curl
ARG BUILD_PATH

COPY --from=eln-ruby-node /asdf /asdf
COPY --from=eln-ruby-node /root/.tool-versions /root/.tool-versions
COPY --from=eln-ruby-node ${BUILD_PATH} ${BUILD_PATH}

# make asdf and its shims available everywhere (set ASDF_DIR to where you cloned asdf)
ENV ASDF_DIR=/asdf      \
    ASDF_DATA_DIR=/asdf \
    ASDF_DEFAULT_TOOL_VERSIONS_FILENAME=.tool-versions \
    GEM_HOME=${BUILD_PATH}/cache/gems
ENV PATH="${GEM_HOME}/bin:${GEM_HOME}/gems/bin:${ASDF_DIR}/shims:${ASDF_DIR}/bin:${PATH}"

# Sanity-Test: Make sure this stage is sane.
RUN asdf --version  && \
    asdf list       && \
    node -v         && \
    npm -v          && \
    npx -v          && \
    yarn -v         && \
    ruby -v         && \
    gem -v          && \
    bundle -v

SHELL ["/bin/bash", "-e", "-o", "pipefail", "-c", "--"]

ARG CHEMOTION_RELEASE
ARG CHEMOTION_VERSION
ARG CHEMOTION_BUILD_TAG
ARG PANDOC_VERSION

RUN apt-get install -y --no-install-recommends build-essential \
    `# for the gems` \
    cmake libpq-dev swig \
    libboost-serialization-dev \
    libboost-iostreams-dev \
    libboost-system-dev \
    libeigen3-dev \
    libmagickcore-dev \
    python3-dev libsqlite3-dev
# shared-mime-info `# for the MIMEmagic gem`

RUN mkdir -p ${BUILD_PATH}/tmp     && curl -L -o ${BUILD_PATH}/tmp/pandoc.deb "https://github.com/jgm/pandoc/releases/download/${PANDOC_VERSION}/pandoc-${PANDOC_VERSION}-1-amd64.deb" && \
    mkdir -p ${BUILD_PATH}/usr/bin && curl -L -o ${BUILD_PATH}/usr/bin/yq     "https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64" && chmod +x ${BUILD_PATH}/usr/bin/yq

# checkout Chemotion
RUN [ -n "${CHEMOTION_BUILD_TAG}" ] && echo "CHECKOUT: ${CHEMOTION_BUILD_TAG}" && \
    git init --initial-branch=main ${BUILD_PATH}/chemotion/app && \
    git -C ${BUILD_PATH}/chemotion/app remote add origin https://github.com/ComPlat/chemotion_ELN && \
    git -C ${BUILD_PATH}/chemotion/app fetch --tags --depth 1 origin "${CHEMOTION_BUILD_TAG}" && \
    git -C ${BUILD_PATH}/chemotion/app reset --hard FETCH_HEAD && \
    git -C ${BUILD_PATH}/chemotion/app describe --abbrev=0 --tags || echo ${CHEMOTION_BUILD_TAG} && \
    cat ${BUILD_PATH}/chemotion/app/VERSION

WORKDIR ${BUILD_PATH}/chemotion/app

# Create release file
RUN echo "CHEMOTION_REF=$(git rev-parse --short HEAD || echo unknown)" >> ${BUILD_PATH}/chemotion/app/.version      && \
    echo "CHEMOTION_TAG=$(git describe --abbrev=0 --tags || echo untagged)" >> ${BUILD_PATH}/chemotion/app/.version && \
    echo "RELEASE=${CHEMOTION_RELEASE}" >> ${BUILD_PATH}/chemotion/app/.version                                     && \
    echo "VERSION=${CHEMOTION_VERSION}" >> ${BUILD_PATH}/chemotion/app/.version                                     && \
    cat ${BUILD_PATH}/chemotion/app/.version

# Sanity-Test: Make sure this stage is sane.
RUN asdf --version  && \
    asdf list       && \
    node -v         && \
    npm -v          && \
    npx -v          && \
    yarn -v         && \
    ruby -v         && \
    gem -v          && \
    bundle -v       && \
    cat ${BUILD_PATH}/chemotion/app/VERSION && \
    cat ${BUILD_PATH}/chemotion/app/.version && \
    test -f ${BUILD_PATH}/usr/bin/yq && \
    test -f ${BUILD_PATH}/tmp/pandoc.deb

# Clean up the configs
RUN find ./config -type f -name '*.yml.example' -print | while read -r f; do \
    echo "$f to ${f%.example}"; \
    ${BUILD_PATH}/usr/bin/yq -ojson < "$f" | \
    ${BUILD_PATH}/usr/bin/yq -oyaml -P 'with(select(.production != null); . = {"production":.production}) | with(select(.production == null); {})' > "${f%.example}"; \
    done

# Clean up unused files
RUN rm -f config/deploy.rb \
    config/environments/development.rb \
    config/environments/test.rb \
    ./**/*.example \
    ./**/*github* \
    ./**/*gitlab* \
    ./**/*travis* \
    ./**/*.bak    \
    ./**/*git* || true

# deal with the weird favicon situation
RUN if [ -f ${BUILD_PATH}/chemotion/app/public/favicon.ico ] && [ ! -e ${BUILD_PATH}/chemotion/app/public/favicon.ico.example ]; then \
    cp ${BUILD_PATH}/chemotion/app/public/favicon.ico ${BUILD_PATH}/chemotion/app/public/favicon.ico.example; \
    fi

# Add application additives as-is. we do this at the end to avoid alterations.
COPY ./additives/chemotion/ ${BUILD_PATH}/chemotion/app

# move [app]/uploads to /chemotion/data
RUN mkdir -p ${BUILD_PATH}/chemotion/data && \
    mv ${BUILD_PATH}/chemotion/app/uploads ${BUILD_PATH}/chemotion/data/ || mkdir -p ${BUILD_PATH}/chemotion/data/uploads && \
    ln -s /chemotion/data/uploads/ ${BUILD_PATH}/chemotion/app/uploads && \
    test "$(readlink ${BUILD_PATH}/chemotion/app/uploads)" == "/chemotion/data/uploads/" && \
    test -d ${BUILD_PATH}/chemotion/data/uploads

# move [app]/public/images to /data
RUN mkdir -p ${BUILD_PATH}/chemotion/data/public/ && \
    mv ${BUILD_PATH}/chemotion/app/public/images/ ${BUILD_PATH}/chemotion/data/public/ || mkdir -p ${BUILD_PATH}/chemotion/data/public/images && \
    ln -s /chemotion/data/public/images/ ${BUILD_PATH}/chemotion/app/public/images && \
    test "$(readlink ${BUILD_PATH}/chemotion/app/public/images)" == "/chemotion/data/public/images/" && \
    test -d ${BUILD_PATH}/chemotion/data/public/images

# add gems from cache
# ADD ./cache/rubycache.tar.gz ${BUILD_PATH}/
# ADD ./cache/nodecache.tar.gz ${BUILD_PATH}/chemotion/app/
# ADD ./cache/appcache.tar.gz  ${BUILD_PATH}/chemotion/app/

ADD ./cache/bundle.tar.gz           /cache/
ADD ./cache/gems.tar.gz             ${BUILD_PATH}/cache/
ADD ./cache/node-modules.tar.gz     ${BUILD_PATH}/chemotion/app/

# Sanity-Test: Make sure this stage is sane.
RUN asdf --version  && \
    asdf list       && \
    node -v         && \
    npm -v          && \
    npx -v          && \
    yarn -v         && \
    ruby -v         && \
    gem -v          && \
    bundle -v       && \
    cat /builder/chemotion/app/VERSION && \
    cat /builder/chemotion/app/.version && \
    test -f /builder/usr/bin/yq && \
    test -f /builder/tmp/pandoc.deb && \
    test -L /builder/chemotion/app/public/images && \
    test -d /builder/chemotion/data/public/images && \
    [[ "$(readlink /builder/chemotion/app/public/images)" == "/chemotion/data/public/images/" ]] && \
    test -L /builder/chemotion/app/uploads && \
    test -d /builder/chemotion/data/uploads && \
    [[ "$(readlink /builder/chemotion/app/uploads)" == "/chemotion/data/uploads/" ]]

# setup environment
ENV BUNDLE_PATH=/cache/bundle                           \
    BUNDLE_USER_HOME=/cache/bundle                      \
    BUNDLE_CACHE_ALL=1                                  \
    RAILS_ENV=production                                \
    RAKE_ENV=production                                 \
    NODE_ENV=production                                 \
    NODE_OPTIONS=--max_old_space_size=4096              \
    THOR_SILENCE_DEPRECATION=1

RUN MAKEFLAGS="-j$(getconf _NPROCESSORS_ONLN)" && export MAKEFLAGS && \
    bundle config --global silence_root_warning 1 && \
    bundle install --jobs="$(getconf _NPROCESSORS_ONLN)" --retry=3 --standalone production --without development test

RUN MAKEFLAGS="-j$(getconf _NPROCESSORS_ONLN)" && export MAKEFLAGS && \
    yarn install --ignore-engines 2>&1 | grep -v ^warning


########################################################################################################################
# Builder
########################################################################################################################
# hadolint ignore=DL3006
FROM "${BASE}" as eln-builder
RUN apt-get install -y --no-install-recommends git ca-certificates curl

RUN apt-get install -y --no-install-recommends --autoremove --fix-missing \
    `# for chemotion` \
    libboost-serialization1.74.0 \
    libboost-iostreams1.74.0 \
    postgresql-client `# also adds pg_isready` \
    inkscape `# this installs python3` \
    imagemagick \
    locales \
    `# utilitites` \
    vim iproute2 sudo make patchelf

COPY --from=eln-runtimes /asdf /asdf
COPY --from=eln-runtimes /root/.tool-versions /root/.tool-versions
COPY --from=eln-runtimes /builder /
COPY --from=eln-runtimes /cache /cache

# CONFIG_PIDFILE: the ELN uses this file to communicate that it is alive
# and the bgworker can start working. (Needs to be on a shared volume.)
ENV ASDF_DIR=/asdf                                      \
    ASDF_DATA_DIR=/asdf                                 \
    ASDF_DEFAULT_TOOL_VERSIONS_FILENAME=.tool-versions  \
    BUNDLE_PATH=/cache/bundle                           \
    BUNDLE_USER_HOME=/cache/bundle                      \
    BUNDLE_CACHE_ALL=1                                  \
    GEM_HOME=/cache/gems                                \
    RAILS_ENV=production                                \
    RAKE_ENV=production                                 \
    NODE_ENV=production                                 \
    NODE_PATH=/chemotion/app/node_modules/              \
    NODE_OPTIONS=--max_old_space_size=4096              \
    TERM=xterm-256color                                 \
    THOR_SILENCE_DEPRECATION=1                          \
    RAILS_LOG_TO_STDOUT=1                               \
    PASSENGER_DOWNLOAD_NATIVE_SUPPORT_BINARY=0          \
    CONFIG_PIDFILE=/chemotion/app/tmp/eln.pid
ENV PATH="${GEM_HOME}/bin:${GEM_HOME}/gems/bin:${ASDF_DIR}/shims:${ASDF_DIR}/bin:${PATH}"

SHELL ["/bin/bash", "-e", "-o", "pipefail", "-c", "--"]

RUN cat /root/.tool-versions && echo "tool-version file present." && \
    test -d /chemotion/app && test -d /chemotion/data && echo "Chemotion directories present." && \
    test -d /cache/bundle && test -d /cache/gems && \
    test -d /chemotion/app/node_modules/ && echo "Cache directories present." && \
    test -f /tmp/pandoc.deb && echo "Pandoc deb present."

# Sanity-Test: Make sure this stage is sane.
# hadolint ignore=DL3003
RUN asdf --version  && \
    asdf list       && \
    node -v         && \
    npm -v          && \
    npx -v          && \
    yarn -v         && \
    ruby -v         && \
    gem -v          && \
    bundle -v       && \
    cat /chemotion/app/VERSION && \
    cat /chemotion/app/.version && \
    test -L /chemotion/app/public/images && \
    test -d /chemotion/data/public/images && \
    [[ "$(readlink /chemotion/app/public/images)" == "/chemotion/data/public/images/" ]] && \
    test -L /chemotion/app/uploads && \
    test -d /chemotion/data/uploads && \
    [[ "$(readlink /chemotion/app/uploads)" == "/chemotion/data/uploads/" ]]


COPY ./additives/various/fontfix.conf /etc/fonts/conf.d/99-chemotion-fontfix.conf
COPY ./additives/embed/               /embed/

RUN ln -s /embed/bin/chemotion /bin/chemotion && \
    cp /chemotion/app/.version /.version && \
    cp /chemotion/app/.version /chemotion/data/.version && \
    mkdir -p /shared && \
    dpkg -i  /tmp/pandoc.deb && rm /tmp/pandoc.deb && \
    chmod +x /embed/bin/*

WORKDIR /chemotion/app

RUN MAKEFLAGS="-j$(getconf _NPROCESSORS_ONLN)" && export MAKEFLAGS && \
    bundle config --global silence_root_warning 1 && \
    bundle install --jobs="$(getconf _NPROCESSORS_ONLN)" --retry=3

RUN MAKEFLAGS="-j$(getconf _NPROCESSORS_ONLN)" && export MAKEFLAGS && \
    yarn install --ignore-engines 2>&1 | grep -v ^warning && \
    yarn cache clean 2>&1 | grep -v ^warning

RUN export SECRET_KEY_BASE=build && \
    export RAILS_ENV=production && \
    [ ! -f ./config/klasses.json ] && echo '[]' > ./config/klasses.json || true && \
    rm -f /chemotion/app/public/sprite.png && \
    cp "$(bundle show ketcherails)"/app/assets/javascripts/ketcherails/ui/ui.js.erb /tmp/ui.js.bak && \
    sed -i 's/Ketcherails::TemplateCategory.with_approved_templates.pluck(:id)/[]/g' "$(bundle show ketcherails)"/app/assets/javascripts/ketcherails/ui/ui.js.erb && \
    npx browserslist@latest --update-db 2>&1 | grep -v ^warning && \
    bundle exec rake assets:precompile && \
    bundle exec rake webpacker:compile && \
    mv /tmp/ui.js.bak "$(bundle show ketcherails)"/app/assets/javascripts/ketcherails/ui/ui.js.erb && \
    rm -f /chemotion/app/public/sprite.png

# Sanity-Test: Make sure this stage is sane.
# hadolint ignore=DL3003
RUN asdf --version  && \
    asdf list       && \
    node -v         && \
    npm -v          && \
    npx -v          && \
    yarn -v         && \
    ruby -v         && \
    gem -v          && \
    bundle -v       && \
    (cd /chemotion/app && bundle exec rails -v) && \
    cat /chemotion/app/VERSION && \
    cat /.version   && \
    cat /chemotion/app/.version && \
    cat /chemotion/data/.version && \
    test -L /chemotion/app/public/images && \
    test -d /chemotion/data/public/images && \
    [[ "$(readlink /chemotion/app/public/images)" == "/chemotion/data/public/images/" ]] && \
    test -L /chemotion/app/uploads && \
    test -d /chemotion/data/uploads && \
    [[ "$(readlink /chemotion/app/uploads)" == "/chemotion/data/uploads/" ]] && \
    test -d /shared && \
    test -d /embed && \
    test -f /chemotion/app/public/favicon.ico.example && \
    test ! -e /chemotion/app/public/sprite.png && \
    ls -la /chemotion/data/uploads/ && \
    ls -la /chemotion/app/public/images/

# export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/cache/bundle/bundler/gems/openbabel-gem-3e25548fd95c/openbabel/lib/
# check dependencies: this should be done in testing step
# RUN find /root /asdf /cache /chemotion/app/node_modules -iname '*.so' -type f -exec /embed/bin/lddcheck \{\} \; | tee /lddlog.txt | grep not | grep -v libRD | sort | uniq
RUN find /root /asdf /cache /chemotion/app/node_modules -iname '*.so' -type f -print0 | xargs -0 ldd | grep -i "not found" || true | sort | uniq

ARG ${PORT}
EXPOSE ${PORT}
CMD ["/embed/run.sh"]

HEALTHCHECK --interval=30s --timeout=3s --start-period=60s --retries=3 \
    CMD /embed/health.sh ${PORT} || exit 1

###########################################################################
# Tagging
###########################################################################
FROM eln-builder as eln
ARG CHEMOTION_VERSION
ARG BUILDRUN
LABEL "org.opencontainers.image.authors"="Chemotion Team" \
    "org.opencontainers.image.title"="Chemotion ELN" \
    "org.opencontainers.image.description"="Image for Chemotion ELN" \
    "org.opencontainers.image.version"="${CHEMOTION_VERSION}" \
    "chemotion.internal.buildrun"="${BUILDRUN}" \
    "chemotion.internal.service.id"="eln"
