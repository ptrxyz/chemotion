# GLOBALs that should be passed to this file
ARG VERSION=UNSET
ARG BUILD_TAG_SPECTRA=UNSET

# Derived values
ARG BASE=chemotion-build/base:${VERSION}

# Private variables, not passed in from outside, but helpful for this file
# -


# Stage 0: import this for later
FROM chemotion-build/base:${VERSION} as chemotion-build-base

# Stage 1: prepare the base image
# hadolint ignore=DL3006
FROM ${BASE} as spectra-builder
RUN apt-get install -y --no-install-recommends --autoremove --fix-missing git ca-certificates curl
RUN mkdir -p    /builder/tmp && \
    curl -L -o  /builder/tmp/conda.sh https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-x86_64.sh && \
    chmod +x    /builder/tmp/conda.sh

ARG BUILD_TAG_SPECTRA

WORKDIR /builder/app
RUN git init --initial-branch=main . && \
    git remote add origin https://github.com/ComPlat/chem-spectra-app && \
    git fetch --tags --depth 1 origin ${BUILD_TAG_SPECTRA} && \
    git reset --hard FETCH_HEAD && \
    rm -rf .git

COPY ./additives/environment.yml           /builder/app/environment.yml
COPY ./additives/spectra_config.py         /builder/app/instance/config.py
COPY ./additives/fake-docker.py            /builder/bin/docker
COPY --from=chemotion-build-base    /tini  /builder/tini

RUN chmod +x /builder/bin/docker && \
    chmod +x /builder/tini


# Stage 2: add the app
# hadolint ignore=DL3006
FROM ${BASE} as spectra
RUN apt-get install -y --no-install-recommends --autoremove --fix-missing git ca-certificates curl
RUN apt-get install -y --no-install-recommends --autoremove --fix-missing \
    gcc g++ libxrender1 libxext-dev pkg-config \
    libfreetype6-dev `# for matplotlib`

COPY --from=spectra-builder /builder /builder

RUN mv /builder/bin/docker /bin/docker && \
    mv /builder/tini /tini && \
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

WORKDIR /app

# hadolint ignore=SC1008
RUN mkdir -p /shared /app/instance && \
    ln -s /shared /app/chem_spectra/tmp && \
    pip install --no-cache-dir -r requirements.txt

# hadolint ignore=SC1008
RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# FROM ${BASE} as final
# COPY --from=spectra / /
# SHELL ["/anaconda3/condabin/conda", "run", "--no-capture-output", "-n", "chem-spectra", "/bin/bash", "-c"]

# Stage 4: finalize the image
FROM spectra as app
ARG VERSION

ENV FLASK_ENV=production \
    FLASK_DEBUG=0 \
    MSC_HOST=msconvert \
    MSC_PORT=4000 \
    MSC_VALIDATE=true \
    SPECTRA_PORT=4000

EXPOSE 4000

WORKDIR "/app"
ENTRYPOINT ["/tini", "--", "/anaconda3/condabin/conda", "run", "--no-capture-output", "-n", "chem-spectra"]
CMD ["gunicorn", "-w", "4", "-b", "0.0.0.0:4000", "server:app"]

HEALTHCHECK --interval=5s --timeout=3s --start-period=30s --retries=3 \
    CMD curl --fail http://localhost:4000/ping || exit 1

LABEL \
    "org.opencontainers.image.authors"="Chemotion Team" \
    "org.opencontainers.image.title"="Chemotion Spectra" \
    "org.opencontainers.image.version"="${VERSION}" \
    "org.opencontainers.image.description"="Image for Chemotion Spectra" \
    "chemotion.internal.service.id"="spectra"
