ARG BASE=ubuntu:22.04
ARG CONVERTER_BUILD_TAG=v0.6.0
ARG CONVERTER_VERSION=1.4.0
ARG BUILDRUN

ARG PORT=8000

# hadolint ignore=DL3006
FROM ${BASE} as converter-base
RUN apt-get update && \
    apt-get install -y --no-install-recommends python3-pip python3-venv libmagic1 curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Stage: build the app
FROM converter-base as app
ARG CONVERTER_BUILD_TAG
ARG SRVDIR=/srv/chemotion

WORKDIR ${SRVDIR}
ADD https://github.com/ComPlat/chemotion-converter-app/archive/refs/tags/${CONVERTER_BUILD_TAG}.tar.gz /tmp/code.tar.gz

RUN tar -xzf /tmp/code.tar.gz --strip-components=1 -C ${SRVDIR} && rm /tmp/code.tar.gz && \
    python3 -m venv env && . env/bin/activate && \
    pip install --no-cache-dir -r ${SRVDIR}/requirements/common.txt

ENV PATH=${SRVDIR}/env/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin \
    VIRTUAL_ENV=${SRVDIR}/env

ARG PORT
EXPOSE ${PORT}

CMD ["gunicorn", "--bind", "0.0.0.0", "converter_app.app:create_app()", "--preload"]

HEALTHCHECK --interval=5s --timeout=3s --start-period=5s --retries=3 \
    CMD curl --fail http://localhost:${PORT-8000}/ || exit 1

FROM app
ARG CONVERTER_VERSION
ARG BUILDRUN
LABEL "org.opencontainers.image.authors"="Chemotion Team" \
    "org.opencontainers.image.title"="Chemotion Converter" \
    "org.opencontainers.image.description"="Image for Chemotion Converter" \
    "org.opencontainers.image.version"="${CONVERTER_VERSION}" \
    "chemotion.internal.buildrun"="${BUILDRUN}" \
    "chemotion.internal.service.id"="converter"
