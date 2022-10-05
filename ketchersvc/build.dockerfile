ARG BASE=ubuntu:22.04
ARG KETCHERSVC_BUILD_TAG=ba28832
ARG KETCHERSVC_VERSION=1.4.0
ARG BUILDRUN

FROM ${BASE} as ketcher-base

RUN apt-get update; \
    apt-get install -y --no-install-recommends gnupg wget curl unzip ca-certificates;
RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | \
    gpg --no-default-keyring --keyring gnupg-ring:/etc/apt/trusted.gpg.d/google.gpg --import; \
    chmod 644 /etc/apt/trusted.gpg.d/google.gpg; \
    echo "deb https://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google.list; \
    apt-get update -y;
RUN apt-get install -y --no-install-recommends google-chrome-stable; \
    CHROME_VERSION=$(google-chrome --product-version | grep -o "[^\.]*\.[^\.]*\.[^\.]*"); \
    CHROMEDRIVER_VERSION=$(curl -s "https://chromedriver.storage.googleapis.com/LATEST_RELEASE_$CHROME_VERSION"); \
    wget -q --continue -P /chromedriver "https://chromedriver.storage.googleapis.com/$CHROMEDRIVER_VERSION/chromedriver_linux64.zip"; \
    unzip /chromedriver/chromedriver* -d /usr/local/bin/
RUN curl -sL https://deb.nodesource.com/setup_16.x | bash - && apt-get update -y
RUN apt-get install -y --no-install-recommends nodejs


# Phase: build the app
FROM ketcher-base as builder
ARG KETCHERSVC_BUILD_TAG

RUN apt-get install -y --no-install-recommends git
RUN git clone --depth 1 https://github.com/ptrxyz/chemotion-ketchersvc.git /src && \
    git -C /src checkout ${KETCHERSVC_BUILD_TAG}

WORKDIR /src

RUN npm install && \
    npm run build && \
    cp -R package.json package-lock.json src/.env.example dist/src/


# Phase: add the app
FROM ketcher-base as app
COPY --from=builder /src/dist/src /app

WORKDIR /app
RUN npm install --omit=dev

EXPOSE 9000
ENV CONFIG_PORT=9000 \
    CONFIG_KETCHER_URL=https://chemotion/ketcher \
    CONFIG_MIN_WORKERS=1 \
    CONFIG_MAX_WORKERS=4

CMD [ "node", "app.js" ]

HEALTHCHECK --interval=5s --timeout=1s --start-period=30s --retries=3 \
    CMD /usr/bin/pidof node || exit 1

FROM app
ARG KETCHERSVC_VERSION
ARG BUILDRUN
LABEL "org.opencontainers.image.authors"="Chemotion Team" \
    "org.opencontainers.image.title"="Chemotion KetcherSVC" \
    "org.opencontainers.image.description"="Image for Chemotion KetcherSVC" \
    "org.opencontainers.image.version"="${KETCHERSVC_VERSION}" \
    "chemotion.internal.buildrun"="${BUILDRUN}" \
    "chemotion.internal.service.id"="ketchersvc"
