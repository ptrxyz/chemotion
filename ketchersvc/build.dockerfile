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

# hadolint ignore=DL3009
RUN curl -L https://dl-ssl.google.com/linux/linux_signing_key.pub | \
    gpg --no-default-keyring --keyring gnupg-ring:/etc/apt/trusted.gpg.d/google.gpg --import && \
    chmod 644 /etc/apt/trusted.gpg.d/google.gpg && \
    echo "deb https://dl.google.com/linux/chrome/deb/ stable main" > /etc/apt/sources.list.d/google-chrome.list && \
    apt-get update -y && apt-get install -y --no-install-recommends --autoremove --fix-missing google-chrome-stable && \
    `# next line: we only go for the latest major release.` && \
    CHROME_VERSION=$(google-chrome --product-version | grep -oE '[0-9]+' | head -n1) && \
    CHROME_VERSION=114 && \
    echo "CURLING: https://chromedriver.storage.googleapis.com/LATEST_RELEASE_$CHROME_VERSION" && \
    CHROMEDRIVER_VERSION=$(curl -s "https://chromedriver.storage.googleapis.com/LATEST_RELEASE_$CHROME_VERSION") && \
    mkdir -p /chromedriver && \
    curl -L -o /chromedriver/chromedriver_linux64.zip "https://chromedriver.storage.googleapis.com/$CHROMEDRIVER_VERSION/chromedriver_linux64.zip" && \
    unzip /chromedriver/chromedriver_linux64.zip -d /usr/local/bin/ && \
    curl -sL https://deb.nodesource.com/setup_16.x | bash - && apt-get update -y && \
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
