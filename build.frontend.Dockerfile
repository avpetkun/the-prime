FROM golang:alpine

RUN apk add nodejs yarn

WORKDIR /srv

COPY package.json yarn.lock frontend/
RUN cd frontend && \
    yarn install --no-progress --no-interactive && \
    yarn cache clean

COPY . frontend
RUN cd frontend && yarn build

RUN mkdir dist && \
    mv frontend/root dist && \
    mv frontend/static dist && \
    mv dist/static/*.html dist/root && \
    mv dist/root/tonconnect-manifest.* . && \
    cp -r dist/root dist/root_dev && \
    mv ./tonconnect-manifest.prod.json dist/root/tonconnect-manifest.json && \
    mv ./tonconnect-manifest.dev.json dist/root_dev/tonconnect-manifest.json
