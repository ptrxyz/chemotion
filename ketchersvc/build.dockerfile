# GLOBALs that should be passed to this file
ARG VERSION=UNSET
ARG BUILD_TAG_KETCHERSVC=UNSET

# Derived values
ARG BASE=chemotion-build/base:${VERSION}

# Private variables, not passed in from outside, but helpful for this file
# -


# Stage 1: prepare the base image
# hadolint ignore=DL3006
FROM ${BASE} as ketcher-base

# hadolint ignore=DL3009
RUN apt-get -y update && apt-get -y upgrade && \
    apt-get install -y --no-install-recommends --autoremove --fix-missing gnupg curl unzip ca-certificates

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

RUN LATEST_CD_RELEASE=$(curl -s https://chromedriver.storage.googleapis.com/LATEST_RELEASE) && echo "LATEST_CD_RELEASE: ${LATEST_CD_RELEASE}" && \
    LATEST_MILESTONE=$(echo $LATEST_CD_RELEASE | egrep -io '^[0-9]+') && echo "LASTET_MILESTONE: ${LATEST_MILESTONE}" && \
    LATEST_C_VERSION=$(curl -s 'https://chromiumdash.appspot.com/fetch_releases?channel=Stable&platform=linux&num=1&milestone='${LATEST_MILESTONE} | yq '.[0].version') && echo "LATEST_C_VERSION: ${LATEST_C_VERSION}" && \
    curl -L -o /tmp/chrome.deb https://dl.google.com/linux/chrome/deb/pool/main/g/google-chrome-stable/google-chrome-stable_${LATEST_C_VERSION}-1_amd64.deb && \
    dpkg -i /tmp/chrome.deb || true && apt install --fix-broken --yes --no-install-recommends --autoremove --fix-missing && \
    curl -L -o /tmp/chromedriver.zip "https://chromedriver.storage.googleapis.com/${LATEST_CD_RELEASE}/chromedriver_linux64.zip" && \
    unzip /tmp/chromedriver.zip -d /usr/local/bin/

RUN NODE_MAJOR=16; apt-get update && apt-get install -y \
    ca-certificates \
    curl \
    gnupg \
    && mkdir -p /etc/apt/keyrings \
    && curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg \
    && echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_${NODE_MAJOR}.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list \
    && apt update && apt install --yes --no-install-recommends --autoremove nodejs && \
    (npm -v || apt-get install -y --no-install-recommends --autoremove --fix-missing --fix-broken npm) && \
    node -v && npm -v

RUN exit 0 ; curl -sL https://deb.nodesource.com/setup_16.x | bash - && apt-get update -y && \
    apt-get install -y --no-install-recommends --autoremove --fix-missing nodejs && \
    (npm -v || apt-get install -y --no-install-recommends --autoremove --fix-missing --fix-broken npm) && \
    node -v && npm -v

# Stage 2: add the app
FROM ketcher-base as builder
RUN apt-get install -y --no-install-recommends --autoremove --fix-missing git

ARG BUILD_TAG_KETCHERSVC

RUN git clone --depth 1 https://github.com/ptrxyz/chemotion-ketchersvc.git /src && \
    git -C /src checkout ${BUILD_TAG_KETCHERSVC}

WORKDIR /src

RUN npm install && \
    npm run build && \
    cp -R package.json package-lock.json src/.env.example dist/src/

# Stage 3: combine
FROM ketcher-base as ketcher
COPY --from=builder /src/dist/src /app

WORKDIR /app
RUN npm install --omit=dev

COPY --from=builder /tini /tini
ENTRYPOINT ["/tini", "-s", "--"]

RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Stage 4: prepare the final image
# FROM ${BASE} as final
# COPY --from=ketcher / /
# SHELL ["/bin/bash", "-e", "-o", "pipefail", "-c", "--"]

# Stage 4: finalize the image
FROM ketcher as app
ARG VERSION

ENV CONFIG_PORT=4000 \
    CONFIG_KETCHER_URL=https://chemotion/ketcher \
    CONFIG_MIN_WORKERS=1 \
    CONFIG_MAX_WORKERS=4

EXPOSE 4000

WORKDIR /app
CMD [ "bash", "-c", "sleep ${STARTUPDELAY:-0}; exec node app.js" ]

HEALTHCHECK --interval=5s --timeout=1s --start-period=30s --retries=3 \
    CMD /usr/bin/pidof node

LABEL \
    "org.opencontainers.image.authors"="Chemotion Team" \
    "org.opencontainers.image.title"="Chemotion KetcherSVC" \
    "org.opencontainers.image.description"="Image for Chemotion KetcherSVC" \
    "org.opencontainers.image.version"="${VERSION}" \    
    "chemotion.internal.service.id"="ketchersvc"
