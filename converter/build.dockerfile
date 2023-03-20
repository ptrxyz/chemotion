# GLOBALs that should be passed to this file
ARG VERSION=UNSET
ARG BUILD_TAG_CONVERTER=UNSET

# Derived values
ARG BASE=chemotion-build/base:${VERSION}

# Private variables, not passed in from outside, but helpful for this file
# -


# Stage 1: prepare the base image
# hadolint ignore=DL3006
FROM ${BASE} as converter-base
RUN apt-get update && \
    apt-get install -y --no-install-recommends --autoremove --fix-missing python3-pip python3-venv libmagic1 curl

# Stage 2: build the app
FROM converter-base as converter
ARG BUILD_TAG_CONVERTER

WORKDIR /srv/chemotion
ADD https://github.com/ComPlat/chemotion-converter-app/archive/refs/tags/${BUILD_TAG_CONVERTER}.tar.gz /tmp/code.tar.gz

RUN tar -xzf /tmp/code.tar.gz --strip-components=1 -C /srv/chemotion && rm /tmp/code.tar.gz && \
    python3 -m venv env && . env/bin/activate && \
    pip install --no-cache-dir -r /srv/chemotion/requirements/common.txt

COPY pass /bin/genpass
RUN chmod +x /bin/genpass && echo "$(/bin/genpass)" > /srv/chemotion/htpasswd   # use echo to append newline.

ENV PATH=/srv/chemotion/env/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin \
    VIRTUAL_ENV=/srv/chemotion/env       \
    MAX_CONTENT_LENGTH=10M               \
    PROFILES_DIR=/srv/chemotion/profiles \
    DATASETS_DIR=/srv/chemotion/datasets \
    HTPASSWD_PATH=/srv/chemotion/htpasswd

RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# FROM ${BASE} as final
# COPY --from=converter / /

# SHELL ["/bin/bash", "-e", "-o", "pipefail", "-c", "--"]

# ENV PATH=/srv/chemotion/env/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin \
#     VIRTUAL_ENV=/srv/chemotion/env       \
#     MAX_CONTENT_LENGTH=10M               \
#     PROFILES_DIR=/srv/chemotion/profiles \
#     DATASETS_DIR=/srv/chemotion/datasets \
#     HTPASSWD_PATH=/srv/chemotion/htpasswd

# Stage 4: finalize the image
FROM converter as app
ARG VERSION

EXPOSE 4000

WORKDIR /srv/chemotion
CMD ["gunicorn", "--bind", "0.0.0.0:4000", "converter_app.app:create_app()", "--preload"]

HEALTHCHECK --interval=5s --timeout=3s --start-period=5s --retries=3 \
    CMD curl --fail http://localhost:4000/

LABEL \
    "org.opencontainers.image.authors"="Chemotion Team" \
    "org.opencontainers.image.title"="Chemotion Converter" \
    "org.opencontainers.image.description"="Image for Chemotion Converter" \
    "org.opencontainers.image.version"="${VERSION}" \
    "chemotion.internal.service.id"="converter"
