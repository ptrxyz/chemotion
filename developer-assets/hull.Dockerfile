####################################################################################################
# BASE IMAGE: TIMEZONE, LOCALES, YQ, TINI
####################################################################################################
FROM ubuntu:focal
ARG DEBIAN_FRONTEND=noninteractive

# all ARGS are set as meaningless specifically because we want them to specified by user
ARG TZ=region/city
ARG LANG=language_territory
ARG VERSION_PANDOC=0
ARG VERSION_ASDF=0
ARG VERSION_NODE=0
ARG VERSION_RUBY=0
ARG VERSION_BUNDLER=0

# set timezone
RUN ln -s /usr/share/zoneinfo/${TZ} /etc/localtime

# set locale
ENV LANG=${LANG}.UTF-8
ENV LANGUAGE=${LANG}
ENV LC_ALL=${LANG}
RUN echo -e "LANG=${LANG}\nLC_ALL=${LANG}" > /etc/locale.conf && \
    echo "${LANG} UTF-8" > /etc/locale.gen


# install basic packages
RUN apt-get -y update && apt-get -y upgrade && \
    apt-get install locales 

# locale
RUN locale-gen ${LANG}

# include binaries and fontfix
ADD https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 /bin/yq
ADD https://github.com/krallin/tini/releases/latest/download/tini           /bin/tini
ADD https://github.com/jgm/pandoc/releases/download/${VERSION_PANDOC}/pandoc-${VERSION_PANDOC}-1-amd64.deb                               /tmp/pandoc.deb
ADD https://gist.githubusercontent.com/ptrxyz/d32479c72fa73fcc1b86ed3219d41b63/raw/aee4a46ce84d0e18ac958feb84015c779c3664aa/fontfix.conf /etc/fonts/conf.d/99-chemotion-fontfix.conf
RUN chmod +x /bin/tini /bin/yq

####################################################################################################
# INSTALL PACKAGES
####################################################################################################

# X-mas tree of packages we need
RUN apt-get -y --autoremove --fix-missing install \
    `# --> utilitites` \
    vim   \
    wget   \
    bash    \
    nano     \
    sudo      \
    iproute2   \
    `# ---------> for asdf` \
    git          \
    curl          \ 
    `# ------------> for ruby` \
    libssl-dev      \
    zlib1g-dev       \
    libzmq3-dev       \
    libreadline-dev    \
    build-essential     \
    `# ------------------> for gems` \
    swig                  \
    cmake                  \
    libpq-dev               \
    python3-dev              \
    libeigen3-dev             \
    libsqlite3-dev             \
    libmagickcore-dev           \
    libboost-system-dev          \
    libboost-iostreams-dev        \
    libboost-serialization-dev     \
    `# -----------------------------> for chemotion` \
    imagemagick                      \
    libboost-iostreams1.71.0          \
    libboost-serialization1.71.0       \
    curl              `# for pandoc`    \
    inkscape          `# adds python3`   \
    postgresql-client `# adds pg_isready`

RUN dpkg -i /tmp/pandoc.deb && rm /tmp/pandoc.deb

####################################################################################################
# INSTALL ASDF + RUBY + NODE
####################################################################################################

# ASDF
# make asdf and its shims available everywhere i.e. add to PATH and to /etc/environment
ENV ASDF_DIR=/asdf
ENV ASDF_DATA_DIR=/asdf
ENV ASDF_DEFAULT_TOOL_VERSIONS_FILENAME=/asdf/tool-versions
ENV PATH=${ASDF_DIR}/shims:${ASDF_DIR}/bin:${PATH}

RUN git clone https://github.com/asdf-vm/asdf.git ${ASDF_DIR} --branch ${VERSION_ASDF}
RUN chmod a+rw ${ASDF_DIR} && \
    sed -i -E 's#(PATH=)("?)(.*)("?)#\1\2'${ASDF_DIR}/shims:${ASDF_DIR}/bin:'\3\4#' /etc/environment && \
    env | grep ^ASDF_ >> /etc/environment

# NodeJS
RUN asdf plugin add nodejs https://github.com/asdf-vm/asdf-nodejs.git && \
    asdf install nodejs ${VERSION_NODE} && \
    asdf global nodejs $(asdf list nodejs) && \
    npm install -g yarn

# Ruby
RUN asdf plugin add ruby https://github.com/asdf-vm/asdf-ruby.git && \
    asdf install ruby ${VERSION_RUBY} && \
    asdf global ruby $(asdf list ruby) && \
    gem install bundler -v ${VERSION_BUNDLER}

SHELL ["/bin/bash", "-c"]
CMD ["/bin/bash"]
