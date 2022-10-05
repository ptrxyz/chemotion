# Evaluated Arguments:
#
# SPECTRA_BUILD_TAG=development
# SPECTRA_VERSION=0.0.0
#
# Files needed:
# ./environment.yml
# ./spectra_config.py
# ./scripts/fake-docker.py
# ./scripts/mscrunner.py
#
# Needed: shared mount point
# spectra:/shared <=> msconvert:/shared
#
# Configuration Spectra:
#
# FLASK_ENV=production
# MSC_HOST=msconvert
# MSC_PORT=8088
# MSC_VALIDATE=true
#
# Exposes Port 3007
#
# Configuration MSConvert:
#
# MSC_PORT=8088
# FLASK_ENV=production
# FLASK_DEBUG=0
# FLASK_APP="mscrunner.py"
#
# Exposes Port 8088

# GLOBALs that should be passed to this file
ARG BASE=chemotion-build/base:1.4.0
ARG SPECTRA_BUILD_TAG=0.10.15
ARG SPECTRA_VERSION=1.4.0
ARG BUILDRUN

# GLOBALs that should not be changed.
ARG MSC_PORT=8088
ARG SPECTRA_PORT=3007


FROM ${BASE} as spectra-builder
RUN apt-get install -y --no-install-recommends git ca-certificates curl
RUN mkdir -p    /builder/tmp && \
    curl -L -o  /builder/tmp/conda.sh https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-x86_64.sh && \
    chmod +x    /builder/tmp/conda.sh

ARG SPECTRA_BUILD_TAG

RUN mkdir -p /builder/app && cd /builder/app && \
    git init --initial-branch=main . && \
    git remote add origin https://github.com/ComPlat/chem-spectra-app && \
    git fetch --tags --depth 1 origin ${SPECTRA_BUILD_TAG} && \
    git reset --hard FETCH_HEAD && \
    rm -rf .git

ADD ./environment.yml        /builder/app/environment.yml
ADD ./spectra_config.py      /builder/app/instance/config.py
ADD ./scripts/fake-docker.py /builder/bin/docker
ADD ./scripts/mscrunner.py   /builder/bin/mscrunner.py

RUN chmod +x /builder/bin/docker


FROM ${BASE} as spectra
RUN apt-get install -y --no-install-recommends git ca-certificates curl
RUN apt-get install -y --no-install-recommends \
    gcc g++ libxrender1 libxext-dev pkg-config \
    libfreetype6-dev # for matplotlib

ARG SPECTRA_PORT

COPY --from=spectra-builder /builder /builder

RUN mv /builder/bin/docker /bin/docker && \
    mv /builder/app/ / && \
    mv /builder/tmp/ / && \
    rm -rf /builder

RUN bash /tmp/conda.sh -p /anaconda3 -b && rm /tmp/conda.sh && \
    echo "PATH=/anaconda3/condabin/:$PATH" >> ~/.profile && \
    /anaconda3/condabin/conda update -y -n base -c defaults conda && \
    /anaconda3/condabin/conda env create -f /app/environment.yml && \
    /anaconda3/condabin/conda install -c anaconda -n chem-spectra setuptools==58.0.4

# Make RUN commands use the new environment:
SHELL ["/anaconda3/condabin/conda", "run", "--no-capture-output", "-n", "chem-spectra", "/bin/bash", "-c"]

RUN cd /app && \
    mkdir -p /shared /app/instance && \
    ln -s /shared /app/chem_spectra/tmp && \
    pip install -r requirements.txt

EXPOSE ${SPECTRA_PORT}

ENV FLASK_ENV=production \
    FLASK_DEBUG=0 \
    MSC_HOST=msconvert \
    MSC_PORT=${MSC_PORT} \
    MSC_VALIDATE=true \
    SPECTRA_PORT=${SPECTRA_PORT}

WORKDIR /app

ENTRYPOINT ["/tini", "--", "/anaconda3/condabin/conda", "run", "--no-capture-output", "-n", "chem-spectra"]
CMD gunicorn -w 4 -b 0.0.0.0:${SPECTRA_PORT} server:app

HEALTHCHECK --interval=5s --timeout=3s --start-period=30s --retries=3 \
    CMD curl --fail http://localhost:${SPECTRA_PORT}/ping || exit 1

FROM spectra
ARG SPECTRA_VERSION
ARG BUILDRUN
LABEL "org.opencontainers.image.authors"="Chemotion Team" \
    "org.opencontainers.image.title"="Chemotion Spectra" \
    "org.opencontainers.image.description"="Image for Chemotion Spectra" \
    "org.opencontainers.image.version"="${SPECTRA_VERSION}" \
    "chemotion.internal.buildrun"="${BUILDRUN}" \
    "chemotion.internal.service.id"="spectra"
